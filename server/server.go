package server

import (
	"net"
	"time"
	"sync"
	"crypto/tls"
	"log"
	"runtime"
	"bufio"
	"strings"
	"io"
	"sibo/protocol"
	"sibo/share"
	"errors"
	"fmt"
	"sibo/proto"
	"github.com/robfig/cron"
)

var ErrServerClosed = errors.New("sibo: Server closed")

var shutdownPollInterval = 500 * time.Millisecond

const (
	ReaderBufferSize = 1024
	WriterBufferSize = 1024
)

type Server struct {
	ln           net.Listener
	readTimeout  time.Duration
	writeTimeout time.Duration

	mu            sync.RWMutex
	activeSession map[ISession]struct{}
	doneChan      chan struct{}
	inShutdown    int32 // accessed atomically (non-zero means we're in Shutdown)

	// TLSConfig for creating tls tcp connection.
	tlsConfig *tls.Config
	waitGroup *sync.WaitGroup

	cron *cron.Cron
}

func NewServer() *Server {
	return &Server{
		waitGroup: &sync.WaitGroup{},
		cron:      cron.New(),
	}
}

func (s *Server) Address() net.Addr {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.ln == nil {
		return nil
	}
	return s.ln.Addr()
}

func (s *Server) Serve(network, address string) (err error) {
	var ln net.Listener
	if s.tlsConfig == nil {
		ln, err = net.Listen("tcp", address)
	} else {
		ln, err = tls.Listen("tcp", address, s.tlsConfig)
	}
	s.cron.Start()
	s.scheduleSavePlayer2DB() // start save job
	return s.serveListener(ln)
}

func (s *Server) serveListener(ln net.Listener) error {
	var tempDelay time.Duration

	s.mu.Lock()
	s.ln = ln
	if s.activeSession == nil {
		s.activeSession = make(map[ISession]struct{})
	}
	s.mu.Unlock()

	for {
		conn, e := ln.Accept()
		if e != nil {
			select {
			case <-s.doneChan:
				return ErrServerClosed
			default:
			}
			if ne, ok := e.(net.Error); ok && ne.Temporary() { // temporary error, retry in one second
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("sibo: Accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			log.Println("closed listener")
			return e
		}
		tempDelay = 0

		if tc, ok := conn.(*net.TCPConn); ok {
			tc.SetKeepAlive(true)
			tc.SetKeepAlivePeriod(3 * time.Minute)
		}

		s.mu.Lock()
		session := NewSession(conn)
		s.activeSession[session] = struct{}{}
		s.mu.Unlock()
		go s.serveSession(session)
	}
}

func (s *Server) serveSession(session ISession) {
	log.Println("serve session -> ", session.SessionId())
	player := NewPlayer(session)
	//atomic.AddInt32(&s.inShutdown, 5)
	defer func() { // execute when client disconnect from server
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			ss := runtime.Stack(buf, false)
			if ss > size {
				ss = size
			}
			buf = buf[:ss]
			log.Printf("serving %s panic error: %s, stack:\n %s", session.Conn().RemoteAddr(), err, buf)
		}
		s.mu.Lock()
		delete(s.activeSession, session)
		s.mu.Unlock()
		session.Close(func() {
			log.Println("exec on session close.")
			/*player.Logout(func() {
				log.Println("exec on player logout.")
			})*/
		})
	}()

	if tlsConn, ok := session.Conn().(*tls.Conn); ok {
		if d := s.readTimeout; d != 0 {
			session.Conn().SetReadDeadline(time.Now().Add(d))
		}
		if d := s.writeTimeout; d != 0 {
			session.Conn().SetWriteDeadline(time.Now().Add(d))
		}
		if err := tlsConn.Handshake(); err != nil {
			log.Printf("rpcx: TLS handshake error from %s: %v", session.Conn().RemoteAddr(), err)
			return
		}
	}
	r := bufio.NewReaderSize(session.Conn(), ReaderBufferSize)

	for {
		t0 := time.Now()
		if s.readTimeout != 0 {
			session.Conn().SetReadDeadline(t0.Add(s.readTimeout))
		}
		req, err := session.ReadRequest(r)
		if err != nil {
			if err == io.EOF {
				//player.SaveAll()
				log.Printf("client has closed this connection: %s", session.Conn().RemoteAddr().String())
			} else if strings.Contains(err.Error(), "use of closed network connection") {
				log.Printf("sibo: connection %s is closed", session.Conn().RemoteAddr().String())
			} else {
				log.Printf("sibo: failed to read request: %v", err)
			}
			session.Close(func() {
				log.Println("session close : " + session.SessionId())
			})
			return
		}
		if s.writeTimeout != 0 {
			session.Conn().SetWriteDeadline(t0.Add(s.writeTimeout))
		}

		go func() {
			if req.IsHeartbeat() {
				req.SetMessageType(protocol.Response)
				session.SendResponse(req)
				return
			}

			res, err := s.handleRequest(player, req)
			if err != nil {
				log.Printf("sibo: failed to handle request: %v", err)
			}
			if !req.IsOneway() {
				session.SendResponse(res)
			}
			protocol.FreeMsg(req)
			protocol.FreeMsg(res)
		}()
	}
}

func (s *Server) handleRequest(player IPlayer, req *protocol.Message) (res *protocol.Message, err error) {
	res = req.Clone()
	res.SetMessageType(protocol.Response)

	codec, ok := share.Codecs[req.SerializeType()]
	if !ok {
		err = fmt.Errorf("can not find codec for %d", req.SerializeType())
		handleError(res, err)
		return res, err
	}
	compression, ok := share.Compression[req.CompressType()]
	if !ok {
		err = fmt.Errorf("can not find compress type for %d", req.CompressType())
		handleError(res, err)
		return res, err
	}

	payload, err := compression.Decompress(req.Payload)
	if err != nil {
		err = fmt.Errorf("UnZip error for compress type %d", req.CompressType())
		handleError(res, err)
		return res, err
	}

	mmId := req.ModuleMessageID()
	if _, ok := proto.RequestMap[mmId]; !ok {
		handleError(res, errors.New("messageID not exist"))
	}
	request := proto.RequestMap[mmId]()
	err = codec.Decode(payload, request)
	if err != nil {
		return handleError(res, err)
	}
	payloadRes, err := ProcessorMap[mmId].Process(player, request)

	if !req.IsOneway() {
		data, err := codec.Encode(payloadRes)
		if err != nil {
			handleError(res, err)
			return res, err
		}
		payload, err := compression.Compress(data)
		if err != nil {
			handleError(res, err)
			return res, err
		}
		res.Payload = payload
	}
	return res, nil
}

func (s *Server) Close() error {
	s.mu.Lock()
	s.closeDoneChan()
	var err error
	if s.ln != nil {
		err = s.ln.Close()
	}
	s.mu.Unlock()

	s.cron.Stop() // stop schedule job
	s.saveAllPlayer2DB()
	s.closeSessions()
	s.waitGroup.Wait()
	return err
}

func (s *Server) closeSessions() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for session := range s.activeSession {
		session.Close(func() {
			log.Println("session close on by server: " + session.SessionId())
		})
		delete(s.activeSession, session)
	}
}

// 定时任务存储玩家数据
func (s *Server) scheduleSavePlayer2DB() {
	//spec := "0 */2 * * * ?"
	s.cron.Schedule(cron.Every(2*time.Minute), new(SavePlayerJob))
}

func (s *Server) saveAllPlayer2DB() {
	if PlayerId2PlayerMap.Len() > 0 {
		for _, player := range PlayerId2PlayerMap.Values() {
			//player, ok := PlayerId2PlayerMap.Get(k)
			s.waitGroup.Add(1)
			go func(player IPlayer) {
				defer s.waitGroup.Done()
				log.Println("task start")
				time.Sleep(2 * time.Second)
				player.SaveAll()
				log.Println("task end")
			}(player)
		}
		PlayerId2PlayerMap.Clear()
	}
}

func (s *Server) closeDoneChan() {
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}
	select {
	case <-s.doneChan:
		// Already closed. Don't close again.
	default:
		// Safe to close here. We're the only closer, guarded
		// by s.mu.
		close(s.doneChan)
	}
}

func handleError(res *protocol.Message, err error) (*protocol.Message, error) {
	res.SetMessageStatusType(protocol.Error)
	return res, err
}
