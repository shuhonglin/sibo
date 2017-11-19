package share

import (
	"sibo/codec"
	"sibo/protocol"
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