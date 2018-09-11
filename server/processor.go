package server

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"sibo/component"
	"sibo/proto"
	"time"
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
		p.SetSession(player.Session())      // 设置新的session
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
	log.Println(createPlayerRequest.UserToken, createPlayerRequest.PlayerName, createPlayerRequest.Sex)
	userComponent, err := player.CreateIfNotExist(component.UserComponent{}.GetType())
	if err != nil {
		return nil, err
	}
	// todo decode userToken
	playerId := int64(10001000001)
	userId := int64(rand.Int())
	uComponent := userComponent.(*component.UserComponent)
	uComponent.SetUserToken(createPlayerRequest.UserToken)
	uComponent.SetUserId(userId)
	uComponent.AddPlayer(playerId)
	uComponent.Save2DB()

	//todo generate player token
	playerToken := createPlayerRequest.UserToken

	player.(*PlayerSession).PlayerId = playerId
	player.(*PlayerSession).UserId = userId
	player.(*PlayerSession).Token = playerToken
	// player.(*PlayerSession).Sex = createPlayerRequest.Sex
	//player.(*PlayerSession).PlayerName = createPlayerRequest.PlayerName

	playerComponent, err := player.CreateIfNotExist(component.PlayerComponent{}.GetType())
	pComponent := playerComponent.(*component.PlayerComponent)
	pComponent.SetUserId(userId)
	pComponent.SetPlayerId(playerId)
	pComponent.SetToken(playerToken)
	pComponent.SetPlayerName(createPlayerRequest.PlayerName)
	pComponent.SetSex(createPlayerRequest.Sex)
	pComponent.SetPosition(0, 0, 0)
	pComponent.Save2DB()

	PlayerId2PlayerMap.Put(player.(*PlayerSession).PlayerId, player)
	// todo
	createPlayerResponse := &proto.CreatePlayerResponse{
		PlayerId: playerId,
		Token:    playerToken,
	}
	log.Println(player)
	return createPlayerResponse, nil
}
