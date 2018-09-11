package server

import (
	log "github.com/sirupsen/logrus"
	"sibo/protocol"
	"sync"
)

var PlayerId2PlayerMap = &PlayerId2Player{
	playerMap: make(map[int64]IPlayer),
}

type PlayerId2Player struct {
	mu        sync.RWMutex
	playerMap map[int64]IPlayer
}

func (p PlayerId2Player) Get(playerId int64) (IPlayer, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	player, ok := p.playerMap[playerId]
	return player, ok
}

func (p *PlayerId2Player) Put(playerId int64, player IPlayer) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.playerMap[playerId] = player
}

func (p *PlayerId2Player) Remove(playerId int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.playerMap, playerId)
}

func (p *PlayerId2Player) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for k := range p.playerMap {
		delete(p.playerMap, k)
	}
}

func (p PlayerId2Player) Keys() []int64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	keys := make([]int64, len(p.playerMap))
	i := 0
	for k := range p.playerMap {
		keys[i] = k
		i++
	}
	return keys
}

func (p PlayerId2Player) Len() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.playerMap)
}

func (p PlayerId2Player) ConstainsKey(playerId int64) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	_, ok := p.playerMap[playerId]
	return ok
}

func (p PlayerId2Player) Values() []IPlayer {
	p.mu.RLock()
	defer p.mu.RUnlock()
	values := make([]IPlayer, len(p.playerMap))
	i := 0
	for _, v := range p.playerMap {
		values[i] = v
		i++
	}
	return values
}

func (p PlayerId2Player) SendMessage(playerId int64, msg *protocol.Message) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if player, ok := p.playerMap[playerId]; ok {
		if player.Session().Status() == CONNECTED {
			player.Session().SendResponse(msg)
			protocol.FreeMsg(msg)
		} else {
			log.Println("player session status error, need CONNECTED but current is ", player.Session().Status())
		}
	} else {
		log.Println("player not exist with playerId ", playerId)
	}
}

func (p *PlayerId2Player) AutoSave2DB() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, player := range p.playerMap {
		player.SaveAll()
		if player.Session().Status() == NOT_CONNECTED || player.Session().Status() == CLOSED {
			delete(p.playerMap, player.(*PlayerSession).PlayerId)
		}
	}
}
