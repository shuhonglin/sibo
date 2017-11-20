package share

import (
	"sibo/codec"
	"sibo/protocol"
	"errors"
)

var (
	Codecs = map[protocol.SerializeType]codec.Codec {
		protocol.SerializeNone: &codec.ByteCodec{},
		protocol.JSON: &codec.JSONCodec{},
		protocol.ProtoBuffer: &codec.PBCodec{},
		protocol.MsgPack: &codec.MsgpackCodec{},
		protocol.Bson: &codec.BsonCodec{},
	}
)

func RegisterCodec(t protocol.SerializeType, c codec.Codec) {
	Codecs[t] = c
}

func EncodeMessageID(moduleId byte, messageID uint32) (uint32, error) {
	if (messageID >> 24) > 0 {
		return 0, errors.New("messageId out of bound")
	}
	msgId := uint32(moduleId)
	return (msgId << 24) | messageID, nil
}