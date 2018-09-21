package component

import (
	"github.com/deckarep/golang-set"
	"reflect"
	"strconv"
)

type IDBSaveProxy interface {
	initFromSqlDB() error
	initFromNoSqlDB() error
	save2SqlDB() error
	save2NoSqlDB() error
}

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
	init      bool
	playerId  int64
	keyPrefix string
	dbSaveProxy IDBSaveProxy // 用于实现父类调用子类，即实现abstract class的效果
	insertSql string
	deleteSql string
	selectSql string
}

func (b BaseComponent) PlayerId() int64 {
	return b.playerId
}

func (b *BaseComponent) InitFromDB() error {
	if DB_TYPE == MYSQL_TYPE {
		return b.dbSaveProxy.initFromSqlDB()
	} else if DB_TYPE == REDIS_TYPE {
		return b.dbSaveProxy.initFromNoSqlDB()
	}
	return nil
}

func (b *BaseComponent) Save2DB() error {
	if DB_TYPE == MYSQL_TYPE {
		return b.dbSaveProxy.save2SqlDB()
	} else if DB_TYPE == REDIS_TYPE {
		return b.dbSaveProxy.save2NoSqlDB()
	}
	return nil
}

func (b BaseComponent) Key() string {
	if DB_TYPE&1 == MYSQL_TYPE {
		return strconv.FormatInt(b.playerId, 10)
	} else if DB_TYPE&2 == REDIS_TYPE {
		return b.keyPrefix + strconv.FormatInt(b.playerId, 10)
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
