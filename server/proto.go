package server

import (
	"errors"
	"log"
	"time"
	"math/rand"
)

var (
	CREATE_PLAYER, _ = EncodeMessageID(1, 1000)
	LOGIN, _         = EncodeMessageID(1, 1001)
	CHANGENAME, _    = EncodeMessageID(1, 1002)
)

var (
	RequestMap = map[uint32]func() interface{}{
		CREATE_PLAYER: func() interface{} {
			return &CreatePlayerRequest{
			}
		},
		LOGIN: func() interface{} {
			return &LoginRequest{
			}
		},
	}
	/*RequestMap = map[uint32]interface{} {
		CREATE_PLAYER: &CreatePlayerRequest{},
		LOGIN: &LoginRequest{},
	}*/

	ProcessorMap = map[uint32]Processor{
		CREATE_PLAYER: new(CreatePlayerProcessor),
		LOGIN:         new(LoginProcessor),
	}
)

type Processor interface {
	Process(player IPlayer, req interface{}) (interface{}, error)
}

type LoginProcessor uint32
type CreatePlayerProcessor uint32

func (p *LoginProcessor) Process(player IPlayer, req interface{}) (interface{}, error) {
	loginMsg, ok := req.(*LoginRequest)
	if !ok {
		return nil, errors.New("request type transform error!")
	}
	log.Println(loginMsg.Token)
	// todo
	loginResponse := &LoginResponse{
		Token: time.Now().String(),
	}
	return loginResponse, nil
}

func (p *CreatePlayerProcessor) Process(player IPlayer, req interface{}) (interface{}, error) {
	createPlayerRequest, ok := req.(*CreatePlayerRequest)
	if !ok {
		return nil, errors.New("request type transform error!")
	}
	log.Println(createPlayerRequest.UserId, createPlayerRequest.PlayerName, createPlayerRequest.Sex)
	player.(*Player).PlayerId = int64(rand.Int())
	player.(*Player).PlayerName = createPlayerRequest.PlayerName
	player.(*Player).UserId = createPlayerRequest.UserId
	player.(*Player).Sex = createPlayerRequest.Sex
	// todo
	createPlayerResponse := &CreatePlayerResponse{
		PlayerId: player.(*Player).PlayerId,
		Token:    player.(*Player).Token,
	}
	log.Println(player)
	return createPlayerResponse, nil
}

type CreatePlayerRequest struct {
	UserId     int64
	PlayerName string
	Sex        byte
}

type LoginRequest struct {
	Token string
}

type LoginResponse struct {
	Token string
}

type CreatePlayerResponse struct {
	PlayerId int64
	Token    string
}

func EncodeMessageID(moduleId byte, messageID uint32) (uint32, error) {
	if (messageID >> 24) > 0 {
		return 0, errors.New("messageId out of bound")
	}
	msgId := uint32(moduleId)
	return (msgId << 24) | messageID, nil
}
