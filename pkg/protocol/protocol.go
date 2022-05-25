package protocol

import (
	"bytes"
	"encoding/gob"
)

type MessageHeader int

const (
	RequestChallenge MessageHeader = iota
	ResponseChallenge
	RequestService
	ResponseService
	Quit
)

type Message struct {
	Header  MessageHeader
	Payload interface{}
}

func NewMessage(h MessageHeader, p interface{}) *Message {
	return &Message{
		h, p,
	}
}

// Encode encodes Message struct to bytes that can be send via Network.
// It uses Gob as an encoder. Read https://go.dev/blog/gob.
func (m Message) Encode() []byte {
	var ret bytes.Buffer
	enc := gob.NewEncoder(&ret)
	enc.Encode(m)
	return ret.Bytes()
}

// Encode decodes bytes to the Message struct.
// It uses Gob as an encoder. Read https://go.dev/blog/gob.
func ParseMessage(m []byte) (*Message, error) {
	var ret bytes.Buffer
	_, err := ret.Write(m)
	if err != nil {
		return nil, err
	}
	enc := gob.NewDecoder(&ret)
	message := &Message{}
	err = enc.Decode(message)
	return message, err
}
