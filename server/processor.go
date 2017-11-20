package server

import (
	"errors"
	"log"
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

func (p *LoginProcessor) Process(player IPlayer, req interface{}) (interface{}, error) {
	loginMsg, ok := req.(*proto.LoginRequest)
	if !ok {
		return nil, errors.New("request type transform error!")
	}
	log.Println(loginMsg.Token)
	// todo
	loginResponse := &proto.LoginResponse{
		Token: time.Now().String(),
	}
	return loginResponse, nil
}

func (p *CreatePlayerProcessor) Process(player IPlayer, req interface{}) (interface{}, error) {
	createPlayerRequest, ok := req.(*proto.CreatePlayerRequest)
	if !ok {
		return nil, errors.New("request type transform error!")
	}
	log.Println(createPlayerRequest.UserId, createPlayerRequest.PlayerName, createPlayerRequest.Sex)
	player.(*Player).PlayerId = int64(rand.Int())
	player.(*Player).PlayerName = createPlayerRequest.PlayerName
	player.(*Player).UserId = createPlayerRequest.UserId
	player.(*Player).Sex = createPlayerRequest.Sex

	PlayerId2PlayerMap.Put(player.(*Player).PlayerId, player)
	// todo
	createPlayerResponse := &proto.CreatePlayerResponse{
		PlayerId: player.(*Player).PlayerId,
		Token:    player.(*Player).Token,
	}
	log.Println(player)
	return createPlayerResponse, nil
}
