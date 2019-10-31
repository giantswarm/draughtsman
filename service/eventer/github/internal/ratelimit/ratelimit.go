package ratelimit

import (
	"sync"
	"time"

	jujuratelimit "github.com/juju/ratelimit"
)

const (
	// defaultFillInterval is the duration of time in which rate limiting token
	// bucket is refilled.
	defaultFillInterval = 1 * time.Minute

	// defaultCapacity is the default number of tokens initially filled in
	// token bucket and which gets refilled after defaultFillInterval.
	defaultCapacity = 60
)

// RateLimiter is wrapper type to encapsulate a 3rd party rate limiting
// implementation and provide functionality to reconfigure it.
type RateLimiter struct {
	bucket *jujuratelimit.Bucket
	mutex  *sync.Mutex
}

// New initializes RateLimiter with default values.
func New() *RateLimiter {
	rl := &RateLimiter{
		mutex: &sync.Mutex{},
	}

	rl.Update(defaultFillInterval, defaultCapacity)

	return rl
}

// Wait takes a value from token bucket when available. If bucket is empty,
// call blocks until token is available.
func (rl *RateLimiter) Wait() {
	rl.bucket.Wait(1)
}

// Update reconfigures RateLimiter's token bucket with given fill interval and
// capacity. This is useful when dealing with APIs that return current rate
// limit values in response headers and where same API rate limiting token
// bucket is shared between many distinct instances.
func (rl *RateLimiter) Update(fillInterval time.Duration, capacity int64) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	rl.bucket = jujuratelimit.NewBucket(fillInterval, capacity)
}
