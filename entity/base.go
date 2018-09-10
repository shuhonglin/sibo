package entity

import (
	"reflect"
	"strings"
	log "github.com/sirupsen/logrus"
)

type Base struct {
	KeyId int64
}

func (b Base)getStructMap(entity interface{}) (fieldMap map[string]interface{}) {
	//fieldMap := make(map[string]interface{})
	t := reflect.TypeOf(entity)
	//v := reflect.ValueOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Errorln("Check type error not Struct")
	}
	fieldNum := t.NumField()
	var tagName, fieldName string
	for i:=0;i<fieldNum;i++  {
		fieldName = t.Field(i).Name
		tags := strings.Split(string(t.Field(i).Tag), "\"")
		if len(tags) > 1 {
			tagName = tags[1]
		} else {
			tagName = fieldName
		}
		//fieldMap[tagName] = v.FieldByName(fieldName)
		fieldMap[tagName] = fieldName
	}
	log.Println(fieldMap)
	return
}