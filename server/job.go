package server

import (
	log "github.com/sirupsen/logrus"
	"sibo/server/processor"
)

type SavePlayerJob int

func (job *SavePlayerJob) Run() {
	log.Println("schedule save...")
	processor.PlayerId2PlayerMap.AutoSave2DB()
}
