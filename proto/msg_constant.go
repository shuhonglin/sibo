package proto

import (
	"sibo/proto/login"
	"sibo/share"
)

const (
	PROCESS_ERROR = -1
)

var (
	// login request constant
	RECONNECT, _     = share.EncodeMessageID(1, 1000)
	CREATE_PLAYER, _ = share.EncodeMessageID(1, 1001)
	LOGIN, _         = share.EncodeMessageID(1, 1002)
	ENTERGAME,_		 = share.EncodeMessageID(1, 1003)

	// other module request constant
	// ...
)

var (
	RequestMap = map[uint32]func() interface{}{
		RECONNECT: func() interface{} {
			return &login.ReconnectRequest{}
		},
		CREATE_PLAYER: func() interface{} {
			return &login.CreatePlayerRequest{
			}
		},
		LOGIN: func() interface{} {
			return &login.LoginRequest{
			}
		},
		ENTERGAME: func() interface{} {
			return &login.EntergameRequest{}
		},
	}
)