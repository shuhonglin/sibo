package component

import (
	"reflect"
	"github.com/deckarep/golang-set"
)

type IComponent interface {
	GetType() reflect.Type
	Save2DB() error
	InitFromDB(playerId int64) error

	IsInit() bool
}

type IMapComponent interface {
	IComponent
	saveUpdateEntityToDB()
	saveNewEntityToDB()
	deleteEntityFromDB()
}

type BaseComponent struct {
	init     bool
}

type MapComponent struct {
	BaseComponent
	playerId int64
	updateSet mapset.Set
	addSet mapset.Set
	delSet mapset.Set
}