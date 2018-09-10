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
	}
)

type CreatePlayerRequest struct {
	UserToken string
	PlayerName string
	Sex        byte
}

type LoginRequest struct {
	Token string
}

type ReconnectRequest struct {
	Token string
}
