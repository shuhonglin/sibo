package server

import (
	"net"
	"sync"
	"github.com/satori/go.uuid"
	"time"
	"sibo/protocol"
	"log"
	"io"
)

var PlayerId2PlayerMap = &PlayerId2Player{
	playerMap:make(map[int64]IPlayer),
}

type PlayerId2Player struct {
	mu sync.RWMutex
	playerMap map[int64]IPlayer
}

func (p PlayerId2Player) Get(playerId int64) (IPlayer, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	player,ok := p.playerMap[playerId]
	return player, ok
}

func (p *PlayerId2Player) Put(playerId int64, player IPlayer) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.playerMap[playerId] = player
}

func (p *PlayerId2Player) Remove(playerId int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.playerMap, playerId)
}

func (p *PlayerId2Player) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for k := range p.playerMap {
		delete(p.playerMap, k)
	}
}

func (p PlayerId2Player) Keys() []int64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	keys := make([]int64, len(p.playerMap))
	i:=0
	for k := range p.playerMap  {
		keys[i] = k
		i++
	}
	return keys
}

func (p PlayerId2Player) Values() []IPlayer {
	p.mu.RLock()
	defer p.mu.RUnlock()
	values := make([]IPlayer, len(p.playerMap))
	i:=0
	for _,v := range p.playerMap  {
		values[i] = v
		i++
	}
	return values
}

func (p PlayerId2Player) SendMessage(playerId int64, msg *protocol.Message) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if player,ok := p.playerMap[playerId];ok {
		player.Session().SendResponse(msg)
		protocol.FreeMsg(msg)
	} else {
		log.Println("player not exist with playerId ", playerId)
	}
}

func(p *PlayerId2Player) AutoSave2DB() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,v:=range p.playerMap  {
		v.SaveAll()
		if v.Session().Status() == NOT_CONNECTED {
			delete(p.playerMap, v.PlayerId())
		}
	}
}

// 定时任务存储玩家数据
func AutoSaveDisConnectedPlayer() {
	PlayerId2PlayerMap.AutoSave2DB()
}


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
	Close(onClose func()) error  // 不要在onClose执行耗时任务
	Status() Status
	SessionId() string
	Createtime() int64
	SendResponse(msg *protocol.Message) error
	ReadRequest(r io.Reader) (*protocol.Message, error)
	Conn() net.Conn
}

type IPlayer interface {
	SetUserId(userId int64)
	SetPlayerId(playerId int64)
	SetPlayerName(playerName string)
	UserId() int64
	PlayerId() int64
	PlayerName() string
	Session() ISession
	Login(token string)
	Logout(onlogout func()) error
	DispatchMsg(msgId uint32, payload interface{}) (interface{}, error)
	SaveAll()
}

type Session struct {
	status         Status
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
	s.status = status
}

func (s Session) Status() Status {
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
	s.UpdateStatus(NOT_CONNECTED)
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

type Player struct {
	session    ISession
	userId     int64
	playerId   int64
	playerName string
	components sync.Map
}

func NewPlayer(session ISession) IPlayer {
	session.UpdateStatus(CONNECTED)
	return &Player{
		session: session,
	}
}

func (p *Player) SetUserId(userId int64) {
	p.userId = userId
}
func (p *Player) SetPlayerId(playerId int64) {
	p.playerId = playerId
}
func (p *Player) SetPlayerName(playerName string) {
	p.playerName = playerName
}
func (p Player) UserId() int64 {
	return p.userId
}
func (p Player) PlayerId() int64 {
	return p.playerId
}
func (p Player) PlayerName() string {
	return p.playerName
}
func (p *Player) Session() ISession {
	return p.session
}
func (p *Player) Login(token string) {
	// decode userId and playerId
	p.userId = 1
	p.playerId = 1
}
func (p *Player) Logout(onlogout func()) error {
	onlogout()
	p.session.UpdateStatus(NOT_CONNECTED)

	p.SaveAll()
	PlayerId2PlayerMap.Remove(p.playerId)  // todo 定时任务删除

	log.Println("player log out")
	return nil
}

func (p *Player) DispatchMsg(msgId uint32, payload interface{}) (interface{}, error) {
	log.Println("receive msg -> ", msgId, " ", payload)
	return nil, nil
}
func (p Player) SaveAll() {
	log.Println("save all data to db. -> ", p.playerId)
}
