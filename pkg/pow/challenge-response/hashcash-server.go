package challenge_response

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"time"

	"alexander.rassanov/pow-tcp-server/pkg/cache"
	"alexander.rassanov/pow-tcp-server/pkg/pow"
	"alexander.rassanov/pow-tcp-server/pkg/protocol"
)

// All possible headers can be used during Client <-> Server interactions.
const (
	RequestChallenge = iota
	ResponseChallenge
	RequestService
	ResponseService
	Quit
)

// ClientExpiration indicates how soon data of a client will be expired.
const ClientExpiration = time.Hour * 24 * 2

// HashExpiration indicates how soon data of a hash will be expired.
const HashExpiration = time.Hour * 24 * 365 * 10

// RandomForHashCash is used to init a hashCash.
const RandomForHashCash = 10000

// Service represents a functions that should be invoked when a service can be accessed.
type Service func() interface{}

// ErrBadRequest is used when a request is bad or received unexpected format.
var ErrBadRequest = errors.New("bad request")

// ErrHashAlreadyExist hash is already exist and cannot be used for the verification.
var ErrHashAlreadyExist = errors.New("hash is already exist")

// ErrBadHash hash is not resolved.
var ErrBadHash = errors.New("bad hash")

// ErrRequestExpired the request cannot be served any more due to expiration.
var ErrRequestExpired = errors.New("request is expired")

// StreamWithHashCash represents a structure to provide word of wisdom to a stream.
// It uses hash cash as a protection of spam requests.
type StreamWithHashCash struct {
	// cache is used to store necessary information for hash cash algorithm.
	cache cache.Cache
	// stream is used to exchange protocol data.
	stream io.ReadWriter
	// clientID is a client id.
	clientID string
	// zeroCount represents the complexity of hash cash challenge.
	zeroCount int
	// service is a function that will return any value that will be sent to clients.
	service Service
}

// NewStreamWithHashCash returns a new StreamWithHashCash.
func NewStreamWithHashCash(cache cache.Cache, clientID string, zeroCount int, stream io.ReadWriter, service Service) StreamWithHashCash {
	return StreamWithHashCash{
		cache:     cache,
		clientID:  clientID,
		zeroCount: zeroCount,
		stream:    stream,
		service:   service,
	}
}

// handleRequestChallenge handles request challenge step.
func (ww StreamWithHashCash) handleRequestChallenge() protocol.Message {
	random := rand.Intn(RandomForHashCash)
	key := fmt.Sprintf("%s:%d", ww.clientID, random)
	// ClientExpiration window compensates for clock skew and network routing time between different systems.
	ww.cache.Set(key, random, ClientExpiration)
	challenge := pow.NewHashCashDataChallenge(ww.clientID, ww.zeroCount, random)
	return protocol.NewMessage(ResponseChallenge, challenge)
}

// handleRequestChallenge handles request service step.
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
	return protocol.NewMessage(ResponseService, ww.service()), nil
}

// ProcessMessage proceeds a message of a hash cash protocol.
func (ww StreamWithHashCash) ProcessMessage(m protocol.Message) (protocol.Message, error) {
	switch m.Header {
	case RequestChallenge:
		m := ww.handleRequestChallenge()
		return m, nil
	case RequestService:
		return ww.handleRequestService(m)
	case Quit:
		return protocol.Message{}, nil
	default:
		return protocol.Message{}, ErrBadRequest
	}
}

// ProcessStream processes a stream.
// It receives incoming traffic, handles it and sends results back.
func (ww StreamWithHashCash) ProcessStream(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		m, err := protocol.ReadPackage(ww.stream)
		if err != nil {
			return err
		}
		mOut, err := ww.ProcessMessage(m)
		if err != nil {
			return err
		}
		if err := protocol.SendPackage(mOut, ww.stream); err != nil {
			return err
		}
	}
}
