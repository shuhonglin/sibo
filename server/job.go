package server

import "log"

type SavePlayerJob int

func(job *SavePlayerJob) Run() {
	log.Println("schedule save...")
	PlayerId2PlayerMap.AutoSave2DB()
}
