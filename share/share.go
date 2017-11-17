package share

import (
	"sibo/protocol"
	"sibo/codec"
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