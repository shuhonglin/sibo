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

func (u UserComponent) GetType() reflect.Type {
	return reflect.TypeOf(u)
}
func (u *UserComponent) Save2DB() error {
	u.userEntity.GetStructMap()
	log.Info("save user: ", u.userEntity)
	return nil
}
func (u *UserComponent) InitFromDB(playerId int64) error {
	u.userEntity = &entity.User{}
	u.userEntity.Players = make([]int64, 5)
	return nil
}

func (u UserComponent) IsInit() bool {
	return u.init
}

func (u *UserComponent) SetUserToken(token string) {
	u.userEntity.UserToken = token
}

func (u *UserComponent) SetUserId(userId int64) {
	u.userEntity.UserId = userId
	u.userEntity.KeyId = userId
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
