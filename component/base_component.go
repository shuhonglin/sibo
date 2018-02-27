package component

import (
	"reflect"
	"github.com/deckarep/golang-set"
)

type IComponent interface {
	GetType() reflect.Type
	Save2DB() error
	InitFromDB(playerId int64) error

	saveUpdateEntityToDB()
	saveNewEntityToDB()
	deleteEntityFromDB()

	IsInit() bool
}

type BaseComponent struct {
	playerId int64
	updateSet mapset.Set
	addSet mapset.Set
	delSet mapset.Set
}