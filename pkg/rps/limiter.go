package rps

import (
	"go.uber.org/ratelimit"
)

// type token uint64

// const (
// 	limiterDone token = iota
// 	limiterContinue
// )

type Limiter interface {
	Take()
}

type NoLimitLimiter struct{}

func (*NoLimitLimiter) Take() {}

func NewNoLimitLimiter() *NoLimitLimiter {
	return &NoLimitLimiter{}
}

type SimpleLimiter struct {
	rateLimit ratelimit.Limiter
}

func NewSimpleLimiter(rps int) *SimpleLimiter {
	return &SimpleLimiter{
		rateLimit: ratelimit.New(rps),
	}
}

func (limiter *SimpleLimiter) Take() {
	limiter.rateLimit.Take()
}
