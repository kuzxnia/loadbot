package worker

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/benbjohnson/clock"
	"go.uber.org/ratelimit"
)

type Limiter interface {
	Take()
	SetRate(uint64)
}

func NewLimiter(rate uint64) Limiter {
	if rate == 0 {
		return Limiter(NewNoLimitLimiter())
	} else {
		return Limiter(NewMutableBucketLeakingLimiter(rate))
	}
}

type NoLimitLimiter struct{}

func (*NoLimitLimiter) Take()          {}
func (*NoLimitLimiter) SetRate(uint64) {}

func NewNoLimitLimiter() *NoLimitLimiter {
	return &NoLimitLimiter{}
}

type BucketLeakingLimiter struct {
	rateLimit ratelimit.Limiter
}

func NewBucketLeakingLimiter(rps uint64) *BucketLeakingLimiter {
	return &BucketLeakingLimiter{
		rateLimit: ratelimit.New(int(rps), ratelimit.WithSlack(1000)),
	}
}

func (limiter *BucketLeakingLimiter) Take() {
	limiter.rateLimit.Take()
}

type MutableBucketLeakingLimiter struct {
	//lint:ignore U1000 Padding is unused but it is crucial to maintain performance
	// of this rate limiter in case of collocation with other frequently accessed memory.
	prepadding [64]byte // cache line size = 64; created to avoid false sharing.
	state      int64    // unix nanoseconds of the next permissions issue.
	//lint:ignore U1000 like prepadding.
	postpadding [56]byte // cache line size - state size = 64 - 8; created to avoid false sharing.

	perRequest time.Duration
	maxSlack   time.Duration
	clock      ratelimit.Clock

	mu sync.RWMutex
}

// newAtomicBased returns a new atomic based limiter.
func NewMutableBucketLeakingLimiter(rate uint64) *MutableBucketLeakingLimiter {
	// TODO consider moving config building to the implementation
	// independent code.
	perRequest := time.Second / time.Duration(rate)
	l := &MutableBucketLeakingLimiter{
		perRequest: perRequest,
		maxSlack:   time.Duration(1000) * perRequest,
		clock:      clock.New(),
	}
	atomic.StoreInt64(&l.state, 0)
	return l
}

// Take blocks to ensure that the time spent between multiple
// Take calls is on average time.Second/rate.
func (l *MutableBucketLeakingLimiter) Take() {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var (
		newTimeOfNextPermissionIssue int64
		now                          int64
	)
	for {
		now = l.clock.Now().UnixNano()
		timeOfNextPermissionIssue := atomic.LoadInt64(&l.state)

		switch {
		case timeOfNextPermissionIssue == 0 || (l.maxSlack == 0 && now-timeOfNextPermissionIssue > int64(l.perRequest)):
			// if this is our first call or t.maxSlack == 0 we need to shrink issue time to now
			newTimeOfNextPermissionIssue = now
		case l.maxSlack > 0 && now-timeOfNextPermissionIssue > int64(l.maxSlack):
			// a lot of nanoseconds passed since the last Take call
			// we will limit max accumulated time to maxSlack
			newTimeOfNextPermissionIssue = now - int64(l.maxSlack)
		default:
			// calculate the time at which our permission was issued
			newTimeOfNextPermissionIssue = timeOfNextPermissionIssue + int64(l.perRequest)
		}

		if atomic.CompareAndSwapInt64(&l.state, timeOfNextPermissionIssue, newTimeOfNextPermissionIssue) {
			break
		}
	}

	sleepDuration := time.Duration(newTimeOfNextPermissionIssue - now)
	if sleepDuration > 0 {
		l.clock.Sleep(sleepDuration)
	}
}

func (l *MutableBucketLeakingLimiter) SetRate(rate uint64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.perRequest = time.Second / time.Duration(rate)
	l.maxSlack = time.Duration(1000) * l.perRequest
	l.clock = clock.New()
	l.state = 0
}
