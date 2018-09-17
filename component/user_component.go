package component

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
	"reflect"
	"sibo/entity"
	"strings"
)

type UserComponent struct {
	BaseComponent
	userEntity *entity.User
}

func (u *UserComponent) InitComponent(playerId int64) {
	u.dbSaveProxy = u
	u.playerId = playerId
	u.keyPrefix = "user_"
	if u.init == false {
		u.userEntity = &entity.User{}
		u.init = true
	}
}

func (u UserComponent) GetType() reflect.Type {
	return reflect.TypeOf(u)
}

/*func (u UserComponent) ID() int64 {
	return u.userEntity.UserId
}*/

/*func (u *UserComponent) Save2DB() error {
	u.userEntity.GetStructMap()
	log.Info("save user: ", u.userEntity)
	return nil
}*/

func (u *UserComponent) save2SqlDB() error {
	log.Info("save usercomponent to sql db")
	insertSql := "REPLACE INTO tb_user("
	nameValues := u.userEntity.GetStructMap()
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
	_, err := tx.NamedExec(insertSql, u.userEntity)
	if err!=nil {
		log.Error(err)
	}
	tx.Commit()
	return nil
}

func (u *UserComponent) save2NoSqlDB() error {
	log.Info("save to nosql database")
	//m := p.playerEntity.GetStructMap()
	jsonData, _ := json.Marshal(u.userEntity)
	reply, err := REDIS_DB.Get().Do("SET", u.Key(), jsonData)
	log.Println(reply, err)
	log.Info("save user: ", u.userEntity)
	return nil
}

func (u *UserComponent) initFromSqlDB() error {
	if u.userEntity == nil {
		u.userEntity = &entity.User{}
	}
	selectSql := "SELECT * FROM tb_user where playerId = ?"
	err := SQL_DB.Select(u.userEntity, selectSql, u.Key())
	if err != nil {
		log.Error(err)
	}
	return nil
}
func (u *UserComponent) initFromNoSqlDB() error {
	if u.userEntity == nil {
		u.userEntity = &entity.User{}
	}
	jsonData, err := redis.Bytes(REDIS_DB.Get().Do("GET", u.Key()))
	if err != nil {
		return err
	}
	json.Unmarshal(jsonData, u.userEntity)
	return nil
}

/*func (u *UserComponent) InitFromDB() error {
	u.dbSaveProxy = u
	if u.userEntity == nil {
		u.userEntity = &entity.User{
			Players: make([]int64, 5),
		}
	}
	return nil
}*/

func (u UserComponent) IsInit() bool {
	return u.init
}

func (u *UserComponent) SetInit(init bool) {
	u.init = init
}

func (u *UserComponent) SetUserToken(token string) {
	u.userEntity.UserToken = token
}

func (u *UserComponent) SetUserId(userId int64) {
	u.userEntity.UserId = userId
}

func (u UserComponent) UserId() int64 {
	return u.userEntity.UserId
}

func (u *UserComponent) SetPlayerId(playerId int64) {
	u.userEntity.PlayerId = playerId
}

func (u UserComponent) PlayerId() int64 {
	return u.userEntity.PlayerId
}

/*func (u *UserComponent) Players() []int64 {
	players := make([]int64, 0, 10)
	json.Unmarshal([]byte(u.userEntity.Players), players)
	return players
}

func (u *UserComponent) AddPlayer(playerId int64) {
	players := make([]int64, 0, 10)
	json.Unmarshal([]byte(u.userEntity.Players), players)
	if len(players) > 0 {
		for _, p := range players {
			if p == playerId {
				return
			}
		}
	}
	jsonData,_ := json.Marshal(append(players, playerId))
	u.userEntity.Players = string(jsonData)
}*/
