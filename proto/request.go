package proto

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
)

type CreatePlayerRequest struct {
	UserId     int64
	PlayerName string
	Sex        byte
}

type LoginRequest struct {
	Token string
}
