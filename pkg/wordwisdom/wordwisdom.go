package wordwisdom

import (
	"io"
	"time"
	"bufio"
	"log"
	"context"
	"math/rand"
	"fmt"
	"errors"

	"alexander.rassanov/pow-tcp-server/pkg/protocol"
	"alexander.rassanov/pow-tcp-server/pkg/pow"
	"alexander.rassanov/pow-tcp-server/pkg/cache"
)


const ClientExpiration = time.Hour*24*2

var ErrBadRequest = errors.New("bad request")
var ErrBadHash = errors.New("bad hash")
var ErrRequestExpired = errors.New("request is expired")

type WordWisdom struct {
	cache cache.Cache
}

func NewWordWisdom(cache cache.Cache) WordWisdom {
	return WordWisdom{
		cache: cache,
	}
}

func (ww WordWisdom) ProcessStream(ctx context.Context, stream io.ReadWriteCloser, clientID string, zeroCount int) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		b, err := bufio.NewReader(stream).ReadBytes('\n')
		if err != nil {
			return err
		}
		m, err := protocol.ParseMessage(b)
		if err != nil {
			return err
		}
		switch m.Header {
		case protocol.RequestChallenge:
			random := rand.Intn(10000)
			key := fmt.Sprintf("%s:%d", clientID, random)
			ww.cache.Set(key, random, ClientExpiration)
			// if we will loose connections or the communication will be interrupted
			defer ww.cache.Delete(key)
			challenge := pow.NewHashCashDataChallenge(clientID, zeroCount, random)
			_, err := stream.Write(protocol.NewMessage(protocol.ResponseChallenge, challenge).Encode())
			if err != nil {
				return err
			}
		case protocol.RequestService:
			hc, ok := m.Payload.(*pow.HashCashData)
			if !ok {
				return ErrBadRequest
			}
			key := fmt.Sprintf("%s:%d", clientID, hc.Rand)
			_, ok = ww.cache.Get(key)
			if !ok {
				return ErrRequestExpired
			}
			if !hc.IsCorrect() {
				return ErrBadHash
			}
			
		case protocol.Quit:
			return nil
		default:
			return ErrBadRequest
		}
	}
}
