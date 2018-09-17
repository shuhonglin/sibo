package component

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
	"reflect"
	"sibo/entity"
	"strings"
)

type PlayerComponent struct {
	BaseComponent
	playerEntity *entity.Player
}

func (p *PlayerComponent) InitComponent(playerId int64) {
	p.dbSaveProxy = p
	p.playerId = playerId
	p.keyPrefix = "player_"
	if p.init == false {
		p.playerEntity = &entity.Player{}
		p.init = true
	}
}

func (p PlayerComponent) GetType() reflect.Type {
	return reflect.TypeOf(p)
}

/*func (p PlayerComponent) Key() string {
	return p.prefix+strconv.FormatInt(p.playerEntity.PlayerId, 10)
}*/

/*func (p PlayerComponent) ID() int64 {
	return p.playerEntity.PlayerId
}*/


func (p *PlayerComponent) save2SqlDB() error {
	log.Info("save to sql database")
	//insertSql := "REPLACE INTO tb_skill(skillId, playerId, hole) VALUES (:skillId, :playerId, :hole)"
	insertSql := "REPLACE INTO tb_player("
	nameValues := p.playerEntity.GetStructMap()
	names := make([]string, len(nameValues))
	values := make([]string, len(nameValues))
	i:=0
	for k := range nameValues {
		names[i] = strings.ToLower(k)
		values[i] = ":"+names[i]
		i++
	}
	insertSql += strings.Join(names, ",")+") VALUES (" + strings.Join(values,",")+")"
	log.Info(insertSql)
	tx := SQL_DB.MustBegin()
	_,err:= tx.NamedExec(insertSql, p.playerEntity)
	if err!=nil {
		log.Error(err)
	}
	tx.Commit()
	return nil
}

func (p *PlayerComponent) save2NoSqlDB() error {
	log.Info("save to nosql database")
	//m := p.playerEntity.GetStructMap()
	jsonData, _ := json.Marshal(p.playerEntity)
	reply, err := REDIS_DB.Get().Do("SET", p.Key(), jsonData)
	log.Println(reply, err)
	log.Info("save player: ", p.playerEntity)
	return nil
}

func (p *PlayerComponent) initFromSqlDB() error {
	if p.playerEntity == nil {
		p.playerEntity = &entity.Player{}
	}
	log.Info("init player from sql db")
	selectSql := "SELECT * FROM tb_player where playerId = ?"
	err := SQL_DB.Select(p.playerEntity, selectSql, p.Key())
	if err != nil {
		log.Error(err)
	}
	return nil
}
func (p *PlayerComponent) initFromNoSqlDB() error {
	if p.playerEntity == nil {
		p.playerEntity = &entity.Player{}
	}
	jsonData, err := redis.Bytes(REDIS_DB.Get().Do("GET", p.Key()))
	if err != nil {
		return err
	}
	json.Unmarshal(jsonData, p.playerEntity)
	return nil
}

/*func (p *PlayerComponent) InitFromDB() error {
	if p.playerEntity == nil {
		p.playerEntity = &entity.Player{}
	}
	jsonData, err := redis.Bytes(REDIS_DB.Get().Do("GET", p.Key()))
	if err != nil {
		return err
	}
	json.Unmarshal(jsonData, p.playerEntity)
	return nil
}*/

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
	pos := make([]int, 3, 3)
	pos[0] = x
	pos[1] = y
	pos[2] = z
	jsonData,_ := json.Marshal(pos)
	p.playerEntity.Pos = string(jsonData)
}

func (p PlayerComponent) Position() [3]int {
	pos := [3]int{}
	json.Unmarshal([]byte(p.playerEntity.Pos), pos)
	return pos
}
