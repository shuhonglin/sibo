package proto

type ReconnectResponse struct {

}

type LoginResponse struct {
	Token string
}

type CreatePlayerResponse struct {
	PlayerId int64
	Token    string
}
