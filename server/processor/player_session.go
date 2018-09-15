package processor

import (
	"github.com/kataras/go-errors"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"reflect"
	"sibo/component"
	"sibo/protocol"
	"strconv"
	"sync"
	"time"
)

type Status byte

const (
	NOT_CONNECTED Status = iota
	CONNECTING
	CONNECTED
	RE_CONNECTING
	CLOSED
)

type ISession interface {
	UpdateStatus(status Status)
	Close(onClose func()) error // 不要在onClose执行耗时任务
	Status() Status
	SessionId() string
	Createtime() int64
	SendResponse(msg *protocol.Message) error
	ReadRequest(r io.Reader) (*protocol.Message, error)
	Conn() net.Conn
}

type IPlayer interface {
	Session() ISession
	/*MsgChan() chan *protocol.Message*/
	SetSession(session ISession)
	Login(token string)
	Logout(onlogout func()) error
	/*DispatchMsg(msgId uint32, payload interface{}) (interface{}, error)*/
	SaveAll() error
	InitFromDB(player int64) error
	SaveComponent(t reflect.Type) error
	GetComponent(t reflect.Type) (component.IComponent, error)
	CreateIfNotExist(t reflect.Type) (component.IComponent, error)
}

type Session struct {
	status         Status
	mu             sync.Mutex
	sessionId      string
	createtime     int64
	isShuttingDown bool
	conn           net.Conn
}

func NewSession(conn net.Conn) ISession {
	return &Session{
		status:         CONNECTING,
		sessionId:      uuid.NewV4().String(),
		createtime:     time.Now().Unix(),
		isShuttingDown: false,
		conn:           conn,
	}
}

func (s *Session) Conn() net.Conn {
	return s.conn
}

func (s *Session) UpdateStatus(status Status) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status = status
}

func (s Session) Status() Status {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status
}

func (s Session) SessionId() string {
	return s.sessionId
}

func (s Session) Createtime() int64 {
	return s.createtime
}

func (s *Session) Close(onClose func()) error {
	onClose()
	s.isShuttingDown = true
	s.UpdateStatus(CLOSED)
	err := s.conn.Close()
	if err != nil {
		return err
	}
	log.Println("session close -> ", s.sessionId)
	return nil
}

func (s *Session) SendResponse(msg *protocol.Message) error {
	_, err := s.conn.Write(msg.Encode())
	return err
}
func (s *Session) ReadRequest(r io.Reader) (*protocol.Message, error) {
	req := protocol.GetPooledMsg()
	err := req.Decode(r)
	return req, err
}

type PlayerSession struct {
	session    ISession
	components map[reflect.Type]component.IComponent
	UserId     int64
	PlayerId   int64
	Token      string
	//msgChan chan *protocol.Message
	//components sync.Map
}

func NewPlayer(session ISession) IPlayer {
	session.UpdateStatus(CONNECTED)
	return &PlayerSession{
		session:    session,
		components: make(map[reflect.Type]component.IComponent),
		//msgChan: make(chan *protocol.Message, 1000),
	}
}

func (p *PlayerSession) Session() ISession {
	return p.session
}

/*func (p *PlayerSession) MsgChan() chan *protocol.Message {
	return p.msgChan
}*/

func (p *PlayerSession) SetSession(session ISession) {
	p.session = session
}
func (p *PlayerSession) Login(token string) {
	// decode userId and playerId
	//p.UserId = 1
	//p.PlayerId = 1
}
func (p *PlayerSession) Logout(onlogout func()) error {
	onlogout()

	p.session.UpdateStatus(CLOSED)
	if PlayerId2PlayerMap.ConstainsKey(p.PlayerId) {
		p.SaveAll()
		/*_, ok := <-p.msgChan
		if ok {
			close(p.msgChan) // 关闭消息队列
		}*/
		//PlayerId2PlayerMap.Remove(p.PlayerId) // todo 定时任务删除
		log.Println("player logout")
	}
	return nil
}

/*func (p *Player) DispatchMsg(msgId uint32, payload interface{}) (interface{}, error) {
	log.Println("receive msg -> ", msgId, " ", payload)
	return nil, nil
}*/
func (p PlayerSession) SaveAll() error {
	log.Println("保存所有component到数据库 -> ", p.PlayerId, p.session.SessionId())
	var err error
	for k, cn := range p.components {
		err = cn.Save2DB()
		if err != nil {
			log.Println("无法保存component ", k.String())
		}
	}
	return nil
}

func (p *PlayerSession) InitFromDB(player int64) error {
	log.Info("从数据库中初始化化玩家数据", player)

	// todo 从数据库或其他地方加载玩家数据
	return nil
}

func (p PlayerSession) SaveComponent(t reflect.Type) error {
	log.Println("保存component到数据库. -> ", p.PlayerId, p.session.SessionId())
	return p.components[t].Save2DB()
}

func (p *PlayerSession) GetComponent(t reflect.Type) (component.IComponent, error) {
	if _, ok := p.components[t]; !ok {
		return nil, errors.New(strconv.FormatInt(p.PlayerId, 10) + " component不存在")
	}
	if p.components[t].IsInit() == false {
		return nil, errors.New(strconv.FormatInt(p.PlayerId, 10) + " component未初始化")
	}
	return p.components[t], nil
}
func (p *PlayerSession) CreateIfNotExist(t reflect.Type) (component.IComponent, error) {
	if _, ok := p.components[t]; !ok {
		cn := reflect.New(t).Interface()
		p.components[t], ok = cn.(component.IComponent)
		if !ok {
			log.Errorln("无法将实例t转化为IComponent类型")
			return nil, errors.New(strconv.FormatInt(p.PlayerId, 10) + " 无法将实例t转化为IComponent类型")
		}
		p.components[t].InitComponent(p.PlayerId)
		err := p.components[t].InitFromDB()
		if err != nil {
			log.Errorln("无法从数据库中初始化 ", p.PlayerId, t.String(), " for ", err)
			//return nil, err
		}
	}
	return p.components[t], nil
}
