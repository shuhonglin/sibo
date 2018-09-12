package proto

var (
	RequestMap = map[uint32]func() interface{}{
		RECONNECT: func() interface{} {
			return &ReconnectRequest{}
		},
		CREATE_PLAYER: func() interface{} {
			return &CreatePlayerRequest{
			}
		},
		LOGIN: func() interface{} {
			return &LoginRequest{
			}
		},
		ENTERGAME: func() interface{} {
			return &EntergameRequest{}
		},
	}
)

type CreatePlayerRequest struct {
	UserToken string
	PlayerName string
	Sex        byte
}

type LoginRequest struct {
	Token string // 解码可得到playerId,时效等信息
}

type ReconnectRequest struct {
	Token string
}

type EntergameRequest struct {
	Token string
}
