package proto

type LoginResponse struct {
	Token string
}

type CreatePlayerResponse struct {
	PlayerId int64
	Token    string
}
