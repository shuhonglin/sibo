package server

import "log"

type SavePlayerJob int

func(job *SavePlayerJob) Run() {
	log.Println("schedule save...")
	for _, player := range PlayerId2PlayerMap.Values() { //todo 可优化，使用多个goroutine加快保存
		player.SaveAll()
	}
}
