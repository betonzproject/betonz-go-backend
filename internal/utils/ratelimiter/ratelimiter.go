package ratelimiter

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	redis *redis.Client
}

type LimiterOptions struct {
	Tokens int
	Window time.Duration
}

var RateLimited = errors.New("Rate limit reached")

func NewLimiter(redis *redis.Client) *Limiter {
	return &Limiter{redis: redis}
}

// Consumes a token based on fixed window rate limiting.
//
// If the limit is reached, it returns `ratelimiter.RateLimited`.
func (l *Limiter) Consume(ctx context.Context, key string, opts LimiterOptions) error {
	result, err := l.redis.Get(ctx, key).Int()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Panicln("Can't get redis: " + err.Error())
		}

		err := l.redis.Set(ctx, key, 1, opts.Window).Err()
		if err != nil {
			log.Panicln("Can't set redis: " + err.Error())
		}
		return nil
	}

	if result >= opts.Tokens {
		return RateLimited
	}

	err = l.redis.Incr(ctx, key).Err()
	if err != nil {
		log.Panicln("Can't incr redis: " + err.Error())
	}
	return nil
}

func (l *Limiter) Reset(ctx context.Context, key string) {
	err := l.redis.Del(ctx, key).Err()
	if err != nil {
		log.Panicln("Can't del redis: " + err.Error())
	}
}
