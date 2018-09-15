package login

/*
Token由playerId与expireTime等加密而成
*/

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
