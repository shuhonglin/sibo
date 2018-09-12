package component

import (
	"github.com/deckarep/golang-set"
	"reflect"
)

type IComponent interface {
	InitComponent()

	ID() int64
	GetType() reflect.Type
	Save2DB() error
	InitFromDB(id int64) error

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
}

type MapComponent struct {
	BaseComponent
	playerId  int64
	updateSet mapset.Set
	addSet    mapset.Set
	delSet    mapset.Set
}

func (m *MapComponent) InitComponent() {
	m.playerId = 0
	m.addSet = mapset.NewSet()
	m.delSet = mapset.NewSet()
	m.updateSet = mapset.NewSet()
}
