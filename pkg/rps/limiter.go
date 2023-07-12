package rps

import (
	"github.com/kuzxnia/mongoload/pkg/config"
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

func NewLimiter(cfg *config.Job) Limiter {
	if cfg.Pace == 0 {
		return Limiter(NewNoLimitLimiter())
	} else {
		return Limiter(NewSimpleLimiter(cfg.Pace))
	}
}

type NoLimitLimiter struct{}

func (*NoLimitLimiter) Take() {}

func NewNoLimitLimiter() *NoLimitLimiter {
	return &NoLimitLimiter{}
}

type SimpleLimiter struct {
	rateLimit ratelimit.Limiter
}

func NewSimpleLimiter(rps uint64) *SimpleLimiter {
	return &SimpleLimiter{
		rateLimit: ratelimit.New(int(rps)),
	}
}

func (limiter *SimpleLimiter) Take() {
	limiter.rateLimit.Take()
}
