package driver

import (
	"github.com/kuzxnia/mongoload/pkg/config"
	"go.uber.org/ratelimit"
)

type Limiter interface {
	Take()
}

func NewLimiter(cfg *config.Job) Limiter {
	if cfg.Pace == 0 {
		return Limiter(NewNoLimitLimiter())
	} else {
		return Limiter(NewBucketLeakingLimiter(cfg.Pace))
	}
}

type NoLimitLimiter struct{}

func (*NoLimitLimiter) Take() {}

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
