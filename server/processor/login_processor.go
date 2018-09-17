package processor

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"sibo/component"
	"sibo/proto"
	"sibo/proto/errcode"
	"sibo/proto/login"
	"time"
)

type LoginProcessor uint32
type CreatePlayerProcessor uint32
type ReconnectProcessor uint32
type EntergameProcessor uint32

func (p *ReconnectProcessor) Process(player IPlayer, req interface{}) (interface{}, error) {
	reconnectMsg, ok := req.(*login.ReconnectRequest)
	if !ok {
		return nil, errors.New("reconnect request transform error")
	}
	log.Debug("reconnect msg with token: ", reconnectMsg.Token)

	if reconnectMsg.Token == "" {
		return &errcode.ErrorResponse{ErrCode: proto.PROCESS_ERROR}, nil
	}
	// todo decode playerId from Token
	playerId := time.Now().Unix()

	if p, ok := PlayerId2PlayerMap.Get(playerId); ok { // player在内存中
		p.Session().UpdateStatus(CONNECTED) // 更新内存中session的状态，及时阻止保存到数据库的操作被执行
		p.SetSession(player.Session())      // 设置新的session
		player = p
	} else { // player不在内存中
		player = NewPlayer(player.Session())
		playerSession := player.(*PlayerSession)
		playerSession.PlayerId = playerId
		playerSession.Token = reconnectMsg.Token

		userComponent, err := player.CreateIfNotExist(component.UserComponent{}.GetType())
		if err != nil {
			return nil, err
		}
		uComponent := userComponent.(*component.UserComponent)
		playerSession.UserId = uComponent.UserId()
		player.InitFromDB(playerId)
	}
	reconnectResponse := &login.ReconnectResponse{}
	return reconnectResponse, nil
}

func (p *LoginProcessor) Process(player IPlayer, req interface{}) (interface{}, error) {
	loginMsg, ok := req.(*login.LoginRequest)
	if !ok {
		return nil, errors.New("login request cast error")
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
		playerSession := player.(*PlayerSession)
		playerSession.PlayerId = playerId
		playerSession.Token = loginMsg.Token
		player.InitFromDB(playerId)
		userComponent, err := player.CreateIfNotExist(component.UserComponent{}.GetType())
		if err != nil {
			return nil, err
		}
		uComponent := userComponent.(*component.UserComponent)
		playerSession.UserId = uComponent.UserId()
	}
	loginResponse := &login.LoginResponse{
		Token: time.Now().String(),
	}
	return loginResponse, nil
}

func (p *CreatePlayerProcessor) Process(player IPlayer, req interface{}) (interface{}, error) {
	createPlayerRequest, ok := req.(*login.CreatePlayerRequest)
	if !ok {
		return nil, errors.New("createPlayer request type cast error")
	}
	log.Println(createPlayerRequest.UserToken, createPlayerRequest.PlayerName, createPlayerRequest.Sex)

	// todo decode userId from userToken and generate playerId,playerToken
	playerId := int64(10001000004)
	userId := int64(2018)
	playerToken := createPlayerRequest.UserToken

	playerSession := player.(*PlayerSession)
	playerSession.PlayerId = playerId
	playerSession.UserId = userId
	playerSession.Token = playerToken

	userComponent, err := player.CreateIfNotExist(component.UserComponent{}.GetType())
	if err != nil {
		return nil, err
	}
	uComponent := userComponent.(*component.UserComponent)
	/*if uComponent.AddPlayer() {

	}*/

	uComponent.SetUserToken(createPlayerRequest.UserToken)
	uComponent.SetUserId(userId)
	uComponent.SetPlayerId(playerId)
	uComponent.Save2DB()

	//player.(*PlayerSession).PlayerId = playerId
	//player.(*PlayerSession).UserId = userId
	//player.(*PlayerSession).Token = playerToken
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
	createPlayerResponse := &login.CreatePlayerResponse{
		Token: playerToken,
	}
	log.Println(player)
	return createPlayerResponse, nil
}

func (p *EntergameProcessor) Process(player IPlayer, req interface{}) (interface{}, error) {
	entergameReq, ok := req.(*login.EntergameRequest)
	if !ok {
		return nil, errors.New("entergame request cast error")
	}
	playerComponent, _ := player.CreateIfNotExist(component.PlayerComponent{}.GetType())
	pComponent := playerComponent.(*component.PlayerComponent)
	if entergameReq.Token != pComponent.Token() {
		return &errcode.ErrorResponse{ErrCode: proto.PROCESS_ERROR}, nil
	}

	entergameResponse := &login.EntergameResponse{
		Token:      pComponent.Token(),
		PlayerName: pComponent.PlayerName(),
		Sex:        pComponent.Sex(),
		Position:   pComponent.Position(),
	}
	return entergameResponse, nil

}