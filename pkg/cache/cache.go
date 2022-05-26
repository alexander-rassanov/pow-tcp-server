package cache

import (
	"time"
)

// Cache is an interface which should be used across all places.
// It allows us to use different cache providers.
type Cache interface {
	// Set sets a key value pair with an expiration.
	Set(key string, value interface{}, d time.Duration)
	// Get returns the value of the key.
	Get(key string) (interface{}, bool)
	// Delete deleted the key from the store.
	Delete(key string)
}
