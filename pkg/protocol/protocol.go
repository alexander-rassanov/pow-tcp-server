package protocol

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
)

type MessageHeader int

var ErrBadPayload = errors.New("bad payload")

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

func NewMessage(h MessageHeader, p interface{}) Message {
	return Message{
		h, p,
	}
}

// Encode encodes Message struct to bytes that can be sent via Network.
// It uses Gob as an encoder. Read https://go.dev/blog/gob.
func (m Message) Encode() []byte {
	gob.Register(&Message{})
	var ret bytes.Buffer
	enc := gob.NewEncoder(&ret)
	if err := enc.Encode(&m); err != nil {
		fmt.Println(err)
	}
	return ret.Bytes()
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
