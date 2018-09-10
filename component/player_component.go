package component

import (
	"sibo/entity"
	"reflect"
	log "github.com/sirupsen/logrus"
)

type PlayerComponent struct {
	BaseComponent
	playerEntity *entity.Player
}

func (p PlayerComponent) GetType() reflect.Type {
	return reflect.TypeOf(p)
}
func (p *PlayerComponent) Save2DB() error {
	p.playerEntity.GetStructMap()
	log.Info("save player: ", p.playerEntity)
	return nil
}
func (p *PlayerComponent) InitFromDB(playerId int64) error {
	p.playerEntity = &entity.Player{}
	return nil
}

func (p PlayerComponent) IsInit() bool {
	return p.init
}

func (p *PlayerComponent) SetToken(token string) {
	p.playerEntity.Token = token
}

func (p *PlayerComponent) SetUserId(userId int64) {
	p.playerEntity.UserId = userId
}

func (p *PlayerComponent) SetPlayerId(playerId int64) {
	p.playerEntity.PlayerId = playerId
}

func (p *PlayerComponent) SetPlayerName(playerName string) {
	p.playerEntity.PlayerName = playerName
}

func (p *PlayerComponent) SetSex(sex byte) {
	p.playerEntity.Sex = sex
}

func (p *PlayerComponent) SetPosition(x, y, z int) {
	p.playerEntity.Pos[0] = x
	p.playerEntity.Pos[1] = y
	p.playerEntity.Pos[2] = z
}
