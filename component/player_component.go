package component

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
	"reflect"
	"sibo/entity"
)

type PlayerComponent struct {
	BaseComponent
	playerEntity *entity.Player
}

func (p *PlayerComponent) InitComponent() {
	p.playerEntity = &entity.Player{}
}

func (p PlayerComponent) GetType() reflect.Type {
	return reflect.TypeOf(p)
}

func (p PlayerComponent) ID() int64 {
	return p.playerEntity.PlayerId
}

func (p *PlayerComponent) Save2DB() error {
	//m := p.playerEntity.GetStructMap()
	jsonData, _ := json.Marshal(p.playerEntity)
	reply, err := REDIS_DB.Get().Do("SET", p.playerEntity.PlayerId, jsonData)
	log.Println(reply, err)
	log.Info("save player: ", p.playerEntity)
	return nil
}

func (p *PlayerComponent) InitFromDB(playerId int64) error {
	p.playerEntity = &entity.Player{}
	jsonData, err := redis.Bytes(REDIS_DB.Get().Do("GET", playerId))
	if err != nil {
		return err
	}
	json.Unmarshal(jsonData, p.playerEntity)
	return nil
}

func (p PlayerComponent) IsInit() bool {
	return p.init
}

func (p *PlayerComponent) SetInit(init bool) {
	p.init = init
}

func (p *PlayerComponent) SetToken(token string) {
	p.playerEntity.Token = token
}

func (p PlayerComponent) Token() string {
	return p.playerEntity.Token
}

func (p *PlayerComponent) SetUserId(userId int64) {
	p.playerEntity.UserId = userId
}

func (p *PlayerComponent) SetPlayerId(playerId int64) {
	p.playerEntity.PlayerId = playerId
	p.playerEntity.KeyId = playerId
}

func (p *PlayerComponent) SetPlayerName(playerName string) {
	p.playerEntity.PlayerName = playerName
}

func (p PlayerComponent) PlayerName() string {
	return p.playerEntity.PlayerName
}

func (p *PlayerComponent) SetSex(sex byte) {
	p.playerEntity.Sex = sex
}

func (p PlayerComponent) Sex() byte {
	return p.playerEntity.Sex
}

func (p *PlayerComponent) SetPosition(x, y, z int) {
	p.playerEntity.Pos[0] = x
	p.playerEntity.Pos[1] = y
	p.playerEntity.Pos[2] = z
}

func (p PlayerComponent) Position() [3]int {
	return p.playerEntity.Pos
}
