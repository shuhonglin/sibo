package component

import (
	"github.com/deckarep/golang-set"
	"reflect"
	"strconv"
)

type IComponent interface {
	InitComponent(playerId int64)

	Key() string
	PlayerId() int64
	GetType() reflect.Type
	Save2DB() error
	InitFromDB() error

	IsInit() bool
	SetInit(init bool)
}

type IMapComponent interface {
	IComponent
	saveUpdateEntityToDB()
	saveNewEntityToDB()
	deleteEntityFromDB()
}

type BaseComponent struct {
	init bool
	playerId  int64
	keyPrefix string
}

func (b BaseComponent) PlayerId() int64 {
	return b.playerId
}

func (b BaseComponent) Key() string {
	if DB_TYPE == MYSQL_TYPE {
		return strconv.FormatInt(b.playerId, 10)
	} else if DB_TYPE == REDIS_TYPE {
		return b.keyPrefix+strconv.FormatInt(b.playerId, 10)
	} else {
		return strconv.FormatInt(b.playerId, 10)
	}
}


type MapComponent struct {
	BaseComponent
	updateSet mapset.Set
	addSet    mapset.Set
	delSet    mapset.Set
}

func (m *MapComponent) InitComponent(playerId int64) {
	m.addSet = mapset.NewSet()
	m.delSet = mapset.NewSet()
	m.updateSet = mapset.NewSet()
}
