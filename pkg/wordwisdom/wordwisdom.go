package wordwisdom

import (
	"bufio"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"time"

	"alexander.rassanov/pow-tcp-server/pkg/cache"
	"alexander.rassanov/pow-tcp-server/pkg/pow"
	"alexander.rassanov/pow-tcp-server/pkg/protocol"
)

const ClientExpiration = time.Hour * 24 * 2
const HashExpiration = time.Hour * 24 * 365 * 10
const RandomForHashCash = 10000

var ErrBadRequest = errors.New("bad request")
var ErrHashAlreadyExist = errors.New("hash is already exist")
var ErrBadHash = errors.New("bad hash")
var ErrRequestExpired = errors.New("request is expired")

type StreamWithHashCash struct {
	cache     cache.Cache
	stream    io.ReadWriter
	clientID  string
	zeroCount int
}

func NewStreamWithHashCash(cache cache.Cache, clientID string, zeroCount int, stream io.ReadWriter) StreamWithHashCash {
	gob.Register(pow.HashCashData{})
	return StreamWithHashCash{
		cache:     cache,
		clientID:  clientID,
		zeroCount: zeroCount,
		stream:    stream,
	}
}

func (ww StreamWithHashCash) getNextMessage() (protocol.Message, error) {
	b, err := bufio.NewReader(ww.stream).ReadBytes('\n')
	if err != nil {
		return protocol.Message{}, err
	}
	return protocol.ParseMessage(b)
}

func (ww StreamWithHashCash) handleRequestChallenge() protocol.Message {
	random := rand.Intn(RandomForHashCash)
	key := fmt.Sprintf("%s:%d", ww.clientID, random)
	// ClientExpiration window compensates for clock skew and network routing time between different systems.
	ww.cache.Set(key, random, ClientExpiration)
	challenge := pow.NewHashCashDataChallenge(ww.clientID, ww.zeroCount, random)
	return protocol.NewMessage(protocol.ResponseChallenge, challenge)
}

func (ww StreamWithHashCash) handleRequestService(m protocol.Message) (protocol.Message, error) {
	hc, ok := m.Payload.(pow.HashCashData)
	if !ok {
		return protocol.Message{}, ErrBadRequest
	}
	key := fmt.Sprintf("%s:%d", ww.clientID, hc.Rand)
	_, ok = ww.cache.Get(key)
	if !ok {
		return protocol.Message{}, ErrRequestExpired
	}
	if !hc.IsCorrect() {
		return protocol.Message{}, ErrBadHash
	}
	_, ok = ww.cache.Get(hc.Sha1Hash())
	if ok {
		return protocol.Message{}, ErrHashAlreadyExist
	}
	ww.cache.Set(hc.Sha1Hash(), "", HashExpiration)
	return protocol.NewMessage(protocol.ResponseService, getRandQuote()), nil
}

func (ww StreamWithHashCash) ProcessMessage(m protocol.Message) (protocol.Message, error) {
	switch m.Header {
	case protocol.RequestChallenge:
		m := ww.handleRequestChallenge()
		return m, nil
	case protocol.RequestService:
		return ww.handleRequestService(m)
	case protocol.Quit:
		return protocol.Message{}, nil
	default:
		return protocol.Message{}, ErrBadRequest
	}
}

func (ww StreamWithHashCash) ProcessStream(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		m, err := ww.getNextMessage()
		if err != nil {
			return err
		}
		mOut, err := ww.ProcessMessage(m)
		if err != nil {
			return err
		}
		if _, err = ww.stream.Write(mOut.Encode()); err != nil {
			return err
		}
		ww.stream.Write([]byte{'\n'})
		fmt.Printf("send bytes %v %v\n", mOut.Encode(), []byte{'\n'})
	}
}
