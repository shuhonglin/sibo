package component

import "reflect"

type IComponent interface {
	GetType() reflect.Type
	Save2DB() error
}