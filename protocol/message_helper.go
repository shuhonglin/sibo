package protocol

import "errors"

var (
	CREATE_PLAYER, _ = EncodeMessageID(1, 1000)
	LOGIN, _         = EncodeMessageID(1, 1001)
	CHANGENAME, _    = EncodeMessageID(1, 1002)
)

type CreatePlayer struct {
	userId     uint64
	playerName string
	sex        byte
}

type Login struct {
	token string
}

var (
	PayloadTypeMap = map[uint32]interface{}{
		CREATE_PLAYER: CreatePlayer{},
		LOGIN:         Login{},
	}
)

func EncodeMessageID(moduleId byte, messageID uint32) (uint32, error) {
	if (messageID >> 24) > 0 {
		return 0, errors.New("messageId out of bound")
	}
	msgId := uint32(moduleId)
	return (msgId << 24) | messageID, nil
}
