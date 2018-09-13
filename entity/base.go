package entity

import (
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

type Base struct {
}

func (b Base) GetStructMap(entity interface{}) map[string]interface{} {
	fieldMap := make(map[string]interface{})
	t := reflect.TypeOf(entity)
	v := reflect.ValueOf(entity)
	getStructFieldNameValues(t, v, fieldMap)
	for k, v := range fieldMap {
		log.Println(k, v)
	}
	return fieldMap
}

func getStructFieldNameValues(t reflect.Type, v reflect.Value, fieldMap map[string]interface{}) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Errorln("Check type error not Struct")
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	fieldNum := t.NumField()
	var tagName, fieldName string
	var fieldValue reflect.Value
	for i := 0; i < fieldNum; i++ {
		fieldName = t.Field(i).Name
		fieldValue = v.FieldByName(fieldName)
		if fieldValue.Kind() == reflect.Struct {
			getStructFieldNameValues(fieldValue.Type(), fieldValue, fieldMap)
			log.Println("field is struct ", t.Field(i).Name)
		} else {
			tags := strings.Split(string(t.Field(i).Tag), "\"")
			if len(tags) > 1 {
				tagName = tags[1]
			} else {
				tagName = fieldName
			}
			fieldMap[tagName] = v.FieldByName(fieldName)
		}
	}
}
