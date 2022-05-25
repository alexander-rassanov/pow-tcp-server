package pow

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// ErrLimitExceeded is an error when limit of computations is exceeded.
var ErrLimitExceeded = errors.New("limit exceeded")

// ComputationLimit limits how much times a hash can be computed.
// It is used to prevent infinite loops.
const ComputationLimit = 1000

// TimeLayout represents layout for HashCashData which is should be YYMMDD[hhmm[ss]].
const TimeLayout = "060102"

// HashCashData represents the structure to compute challenges.
// See https://en.wikipedia.org/wiki/Hashcash.
type HashCashData struct {
	// Hashcash format version, 1 (which supersedes version 0).
	Ver int
	// Number of "partial pre-image" (zero) bits in the hashed code.
	Bits int
	// The time that the message was sent, in the format YYMMDD[hhmm[ss]].
	Date string
	// Resource data string being transmitted, e.g., an IP address or email address.
	Resource string
	// String of random characters, encoded in base-64 format.
	Rand int
	// Binary counter, encoded in base-64 format.
	Counter int
}

// NewHashCashDataChallenge returns a new challenge based on HashCashData algorithm.
func NewHashCashDataChallenge(resource string, zeroCount int, random int) *HashCashData {
	return &HashCashData{
		Ver:     1,
		Bits:    zeroCount,
		Date:    getHashCashDate(),
		Rand:    random,
		Counter: rand.Intn(100),
	}
}

func getHashCashDate() string {
	return time.Now().Format(TimeLayout)
}

// ToString returns string representation of HashCashData.
func (h HashCashData) ToString() string {
	return fmt.Sprintf("%d:%d:%s:%s::%s:%s",
		h.Ver,
		h.Bits,
		h.Date,
		h.Resource,
		base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(h.Rand))),
		base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(h.Counter))),
	)
}

// sha1Hash calculates sha1 for HashcahData.
func (h HashCashData) Sha1Hash() string {
	sha1 := sha1.New()
	sha1.Write([]byte(h.ToString()))
	bs := sha1.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// isHashCorrect returns true when hash is correct otherwise false.
func isHashCorrect(hash string, zeroCount int) bool {
	if zeroCount == 0 {
		// what to do?
	}
	if len(hash) < zeroCount {
		return false
	}
	for _, ch := range hash[:zeroCount] {
		if ch != '0' {
			return false
		}
	}
	return true
}

func (h HashCashData) IsCorrect() bool {
	if h.Bits == 0 {
		// what to do?
	}
	hash := h.Sha1Hash()
	if len(hash) < h.Bits {
		return false
	}
	for _, ch := range hash[:h.Bits] {
		if ch != '0' {
			return false
		}
	}
	return true
}

// Resolve calculates correct hashcash.
func (h HashCashData) Resolve() (HashCashData, error) {
	for i := 0; i < ComputationLimit; i, h.Counter = i+1, h.Counter+1 {
		if h.IsCorrect() {
			return h, nil
		}
		h.Counter++
	}
	return h, ErrLimitExceeded
}
