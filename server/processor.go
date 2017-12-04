package server

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"time"
	"math/rand"
	"sibo/proto"
)

var (
	/*RequestMap = map[uint32]interface{} {
		CREATE_PLAYER: &CreatePlayerRequest{},
		LOGIN: &LoginRequest{},
	}*/

	ProcessorMap = map[uint32]Processor{
		proto.CREATE_PLAYER: new(CreatePlayerProcessor),
		proto.LOGIN:         new(LoginProcessor),
	}
)

type Processor interface {
	Process(player IPlayer, req interface{}) (interface{}, error)
}

type LoginProcessor uint32
type CreatePlayerProcessor uint32
type ReconnectProcessor uint32

func (p *ReconnectProcessor) Process(player IPlayer, req interface{}) (interface{}, error) {
	reconnectMsg, ok := req.(*proto.ReconnectRequest)
	if !ok {
		return nil, errors.New("reconnect request transform error")
	}
	log.Debug("reconnect msg with token: ", reconnectMsg.Token)

	// todo decode playerId from Token
	playerId := time.Now().Unix()
	if p, ok := PlayerId2PlayerMap.Get(playerId); ok { // player在内存中
		p.Session().UpdateStatus(CONNECTED) // 更新内存中session的状态，及时阻止保存到数据库的操作被执行
		p.SetSession(player.Session()) // 设置新的session
		player = p
	} else { // player不在内存中
		player = NewPlayer(player.Session())
		player.InitFromDB(playerId)
	}
	reconnectResponse := &proto.ReconnectResponse{}
	return reconnectResponse, nil
}

func (p *LoginProcessor) Process(player IPlayer, req interface{}) (interface{}, error) {
	loginMsg, ok := req.(*proto.LoginRequest)
	if !ok {
		return nil, errors.New("login request transform error")
	}
	log.Debug("login msg with token: ", loginMsg.Token)

	// todo decode playerId from Token
	playerId := time.Now().Unix()
	if p, ok := PlayerId2PlayerMap.Get(playerId); ok { // player在内存中
		p.SetSession(player.Session())
		//p.Session().UpdateStatus(CONNECTED)
		player = p
	} else { // player不在内存中,从数据库中初始化
		player = NewPlayer(player.Session())
		player.InitFromDB(playerId)
	}
	loginResponse := &proto.LoginResponse{
		Token: time.Now().String(),
	}
	return loginResponse, nil
}

func (p *CreatePlayerProcessor) Process(player IPlayer, req interface{}) (interface{}, error) {
	createPlayerRequest, ok := req.(*proto.CreatePlayerRequest)
	if !ok {
		return nil, errors.New("request type transform error")
	}
	log.Println(createPlayerRequest.UserId, createPlayerRequest.PlayerName, createPlayerRequest.Sex)
	player.(*PlayerSession).PlayerId = int64(rand.Int())
	player.(*PlayerSession).PlayerName = createPlayerRequest.PlayerName
	player.(*PlayerSession).UserId = createPlayerRequest.UserId
	player.(*PlayerSession).Sex = createPlayerRequest.Sex

	PlayerId2PlayerMap.Put(player.(*PlayerSession).PlayerId, player)
	// todo
	createPlayerResponse := &proto.CreatePlayerResponse{
		PlayerId: player.(*PlayerSession).PlayerId,
		Token:    player.(*PlayerSession).Token,
	}
	log.Println(player)
	return createPlayerResponse, nil
}
