package proto

/*
Token由playerId与expireTime等加密而成
*/

const (
	PROCESS_ERROR int = -1
)

type ErrorResponse struct {
	ErrCode int
}

type ReconnectResponse struct {

}

type LoginResponse struct {
	Token string
}

type CreatePlayerResponse struct {
	Token    string
}

type EntergameResponse struct {
	Token string
	PlayerName string
	Sex byte
	Position [3]int
}
