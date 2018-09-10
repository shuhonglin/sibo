package protocol

import (
	"errors"
	"encoding/binary"
	"io"
)

// MaxMessageLength is the max length of a message.
// Default is 0 that means does not limit length of messages.
// It is used to validate when read messages from io.Reader.
var MaxMessageLength = 0

const (
	magicNumber byte = 0x09
)

var (
	// message is too long
	ErrMessageToLong = errors.New("message is too long")
	ErrMessageMagicNumberError = errors.New("magic number is wrong")
)

// MessageType is message type of request and response
type MessageType byte

const (
	// message type of request
	Request MessageType = iota
	// message type of response
	Response
)

// MessageStatusType is status of message
type MessageStatusType byte

const (
	// Normal is normal requests and responses
	Normal MessageStatusType = iota
	Error
)

// CompressType defines decompression type.
type CompressType byte

const (
	// Nonde does not compress
	None CompressType = iota
	// Gzip uses gzip compression
	Gzip
	Zlib
)

// SerializeType defines serialization type of message's payload.
type SerializeType byte

const (
	SerializeNone SerializeType = iota
	JSON
	ProtoBuffer
	MsgPack
	Bson
)

// Message is the generic type of Request and Response.
type Message struct {
	*Header
	Payload []byte
	data    []byte
}

func NewMessage() *Message {
	header := Header([8]byte{})
	header[0] = magicNumber
	return &Message{
		Header: &header,
	}
}

// Header is the first part of Message and has fixed size.
// Format:
//   1			  1			1				1           4
//magicNumber|version|MessageType|serializationType|MessageId
type Header [8]byte

// CheckMagicNumber checks whether header starts sibo magic number.
func (h Header) CheckMagicNumber() bool {
	return h[0] == magicNumber
}

// Version returns version of sibo protocol.
func (h Header) Version() byte {
	return h[1]
}

// SetVersion sets version of this header.
func (h *Header) SetVersion(version byte) {
	h[1] = version
}

// MessageType returns the message type.
func (h Header) MessageType() MessageType {
	return MessageType(h[2]&0x80) >> 7
}

// SetMessageType sets message type.
func (h *Header) SetMessageType(messageType MessageType) {
	h[2] = h[2] | (byte(messageType) << 7)
}

// IsHeartbeat returns whether the message is heartbeat message.
func (h Header) IsHeartbeat() bool {
	return h[2]&0x40 == 0x40
}

// SetHeartbeat sets the heartbeat flag.
func (h *Header) SetHeartbeat(hb bool) {
	if hb {
		h[2] = h[2] | 0x40
	} else {
		//h[2] = h[2] &^ 0x40 or h[2] = h[2] & (^0x40); or h[2] = h[2] & 0xbf
		h[2] = h[2] & 0xbf
	}
}

// IsOneway returns whether the message is one-way message.
// If true, server won't send responses.
func (h Header) IsOneway() bool {
	return h[2]&0x20 == 0x20
}

// SetOneway sets the oneway flag.
func (h *Header) SetOneway(oneway bool) {
	if oneway {
		h[2] = h[2] | 0x20
	} else {
		h[2] = h[2] &^ 0x20
	}
}

// CompressType returns compression type of messages.
func (h Header) CompressType() CompressType {
	return CompressType((h[2] & 0x1C) >> 2)
}

// SetCompressType sets the compression type.
func (h *Header) SetCompressType(ct CompressType) {
	h[2] = h[2] | ((byte(ct) << 2) & 0x1C)
}

// MessageStatusType returns the message status type.
func (h Header) MessageStatusType() MessageStatusType {
	return MessageStatusType(h[2] & 0x03)
}

// SetMessageStatusType sets message status type.
func (h *Header) SetMessageStatusType(mt MessageStatusType) {
	h[2] = h[2] | (byte(mt) & 0x03)
}

// SerializeType returns serialization type of payload.
func (h Header) SerializeType() SerializeType {
	return SerializeType((h[3] & 0xF0) >> 4)
}

// SetSerializeType sets the serialization type.
func (h *Header) SetSerializeType(st SerializeType) {
	h[3] = h[3] | (byte(st) << 4)
}

func (h Header) ModuleMessageID() uint32 {
	return binary.BigEndian.Uint32(h[4:8])
}

// Module returns system module
func (h Header) Module() byte {
	return h[4]
}

// SetModule sets the system module
func (h *Header) SetModule(m byte) {
	h[4] = m
}

// MessageID returns message id
func (h Header) MessageID() uint32 {
	msgId := binary.BigEndian.Uint32(h[4:8])
	return (msgId << 8) >> 8
}

// SetMessage sets the message id
func (h *Header) SetMessageID(msgId uint32) {
	realMsgId := (uint32(h[4]) << 24) | ((msgId << 8) >> 8)
	binary.BigEndian.PutUint32(h[4:8], realMsgId)
}

func (m Message) Clone() *Message {
	header := *m.Header
	c := GetPooledMsg()
	c.Header = &header
	return c
}

// Reset clean data of this message but keep allocated data
func (m *Message) Reset() {
	resetHeader(m.Header)
	m.Payload = m.Payload[:0]
	m.data = m.data[:0]
}

var zeroHeaderArray Header
var zeroHeader = zeroHeaderArray[1:]

func resetHeader(h *Header) {
	copy(h[1:], zeroHeader)
}

func (m Message) Encode() []byte {
	l := 12 + len(m.Payload)
	data := make([]byte, l)
	copy(data, m.Header[:])
	binary.BigEndian.PutUint32(data[8:12], uint32(len(m.Payload)))
	copy(data[12:], m.Payload)
	return data
}

func (m Message) WriteTo(w io.Writer) error {
	_, err := w.Write(m.Header[:])
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.BigEndian, uint32(len(m.Payload)))
	if err != nil {
		return err
	}
	_, err = w.Write(m.Payload)
	return err
}

func (m *Message) Decode(r io.Reader) error {
	_, err := io.ReadFull(r, m.Header[:])
	if err != nil {
		return err
	}
	lenData := poolUint32Data.Get().(*[]byte)
	_, err = io.ReadFull(r, *lenData)
	if err != nil {
		poolUint32Data.Put(lenData)
		return err
	}
	l := binary.BigEndian.Uint32(*lenData)
	poolUint32Data.Put(lenData)
	if MaxMessageLength > 0 && int(l) > MaxMessageLength {
		return ErrMessageToLong
	}

	data := make([]byte, int(l))
	_, err = io.ReadFull(r, data)
	if err != nil {
		return err
	}
	m.data = data
	m.Payload = data[:]
	return err
}

func Read(r io.Reader) (*Message, error) {
	msg := NewMessage()
	err := msg.Decode(r)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
