package server

import log "github.com/sirupsen/logrus"

type SavePlayerJob int

func(job *SavePlayerJob) Run() {
	log.Println("schedule save...")
	PlayerId2PlayerMap.AutoSave2DB()
}
