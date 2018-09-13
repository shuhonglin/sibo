package component

import (
	log "github.com/sirupsen/logrus"
	"reflect"
	"sibo/entity"
)

type UserComponent struct {
	BaseComponent
	userEntity *entity.User
}

func (u *UserComponent) InitComponent(playerId int64) {
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

func (u *UserComponent) Save2DB() error {
	u.userEntity.GetStructMap()
	log.Info("save user: ", u.userEntity)
	return nil
}
func (u *UserComponent) InitFromDB() error {
	if u.userEntity == nil {
		u.userEntity = &entity.User{
			Players:make([]int64, 5),
		}
	}
	return nil
}

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

func (u *UserComponent) Players() []int64 {
	return u.userEntity.Players
}

func (u *UserComponent) AddPlayer(playerId int64) {
	if len(u.userEntity.Players) > 0 {
		for _, p := range u.userEntity.Players {
			if p == playerId {
				return
			}
		}
	}
	u.userEntity.Players = append(u.userEntity.Players, playerId)
}
