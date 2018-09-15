package login

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
