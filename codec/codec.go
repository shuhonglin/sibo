package codec

import (
	"fmt"
	"reflect"
	"encoding/json"
	proto "github.com/gogo/protobuf/proto"
	pb "github.com/golang/protobuf/proto"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/mgo.v2/bson"
)

type Codec interface {
	Encode(i interface{})([]byte, error)
	Decode(data []byte, i interface{}) error
}

type ByteCodec struct {}

func (c ByteCodec) Encode(i interface{}) ([]byte, error) {
	if data, ok := i.([]byte); ok {
		return data, nil
	}
	return nil, fmt.Errorf("%T is not a []byte", i)
}

func (c ByteCodec) Decode(data []byte, i interface{}) error {
	reflect.ValueOf(i).SetBytes(data)
	return nil
}

type JSONCodec struct{}

func (c JSONCodec) Encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (c JSONCodec) Decode(data []byte, i interface{}) error {
	return json.Unmarshal(data, i)
}

// PBCodec uses protobuf marshaler and unmarshaler.
type PBCodec struct{}

// Encode encodes an object into slice of bytes.
func (c PBCodec) Encode(i interface{}) ([]byte, error) {
	if m, ok := i.(proto.Marshaler); ok {
		return m.Marshal()
	}

	if m, ok := i.(pb.Message); ok {
		return pb.Marshal(m)
	}

	return nil, fmt.Errorf("%T is not a proto.Marshaler", i)
}

// Decode decodes an object from slice of bytes.
func (c PBCodec) Decode(data []byte, i interface{}) error {
	if m, ok := i.(proto.Unmarshaler); ok {
		return m.Unmarshal(data)
	}

	if m, ok := i.(pb.Message); ok {
		return pb.Unmarshal(data, m)
	}

	return fmt.Errorf("%T is not a proto.Unmarshaler", i)
}


// MsgpackCodec uses messagepack marshaler and unmarshaler.
type MsgpackCodec struct{}

// Encode encodes an object into slice of bytes.
func (c MsgpackCodec) Encode(i interface{}) ([]byte, error) {
	return msgpack.Marshal(i)
}

// Decode decodes an object from slice of bytes.
func (c MsgpackCodec) Decode(data []byte, i interface{}) error {
	return msgpack.Unmarshal(data, i)
}

// BsonCodec uses bson marshaler and unmarshaler.
type BsonCodec struct{}

// Encode encodes an object into slice of bytes.
func (c BsonCodec) Encode(i interface{}) ([]byte, error) {
	return bson.Marshal(i)
}

// Decode decodes an object from slice of bytes.
func (c BsonCodec) Decode(data []byte, i interface{}) error {
	return bson.Unmarshal(data, i)
}

