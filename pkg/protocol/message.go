package protocol

import (
	"bytes"
	"encoding/gob"
	"errors"
)

// RegisterType accepts a value of a type.
// It is required to be invoked for types that will be used in payload field of Message.
func RegisterType(x interface{}) {
	gob.Register(x)
}

// MessageHeader represents the Header of a packet message.
type MessageHeader int

// ErrBadPayload indicates a payload is bad.
var ErrBadPayload = errors.New("bad payload")

// Message is a network packet that will be used to keep proper interaction between Client and Server side.
type Message struct {
	Header  MessageHeader
	Payload interface{}
}

// NewMessage inits Message struct.
func NewMessage(h MessageHeader, p interface{}) Message {
	return Message{
		h, p,
	}
}

// Encode encodes Message struct to bytes that can be sent via Network.
// It uses Gob as an encoder. Read https://go.dev/blog/gob.
func (m Message) Encode() ([]byte, error) {
	gob.Register(&Message{})
	var ret bytes.Buffer
	enc := gob.NewEncoder(&ret)
	if err := enc.Encode(&m); err != nil {
		return nil, err
	}
	return ret.Bytes(), nil
}

// ParseMessage decodes bytes to the Message struct.
// It uses Gob as an encoder. Read https://go.dev/blog/gob.
func ParseMessage(m []byte) (Message, error) {
	inputBuffered := bytes.NewBuffer(m)
	enc := gob.NewDecoder(inputBuffered)
	message := Message{}
	err := enc.Decode(&message)
	return message, err
}
