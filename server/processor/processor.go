package processor

import "sibo/proto"

var (
	ProcessorMap = map[uint32]Processor{
		proto.RECONNECT: new(ReconnectProcessor),
		proto.CREATE_PLAYER: new(CreatePlayerProcessor),
		proto.LOGIN:         new(LoginProcessor),
		proto.ENTERGAME:     new(EntergameProcessor),
	}
)

type Processor interface {
	Process(player IPlayer, req interface{}) (interface{}, error)
}
