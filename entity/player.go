package entity


type Player struct {
	Base
	Token      string `json:"token"`
	UserId 	   int64
	PlayerId   int64
	PlayerName string
	Sex        byte
	Pos        [3]int
}

func (p Player) GetStructMap()(map[string]interface{})  {
	return p.Base.getStructMap(p)
}
