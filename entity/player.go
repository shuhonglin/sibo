package entity


type Player struct {
	UserId     int64
	PlayerId   int64
	PlayerName string
	Sex        byte
	Pos        [3]int
	Token      string
}
