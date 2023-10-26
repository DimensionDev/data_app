package server

import (
	"context"
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	libredis "github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/redis"
)

var ErrLimitExceed = errors.New(429, "RATELIMIT", "service unavailable due to rate limit exceeded")
var ErrGetLimitKey = errors.New(500, "INTERERROR", "failed to get key for rate limit")

func NewRateLimiter() *limiter.Limiter {
	rate, err := limiter.NewRateFromFormatted("10-M")
	if err != nil {
		fmt.Println("error creating rate limit")
		return nil
	}
	redis_client, err := newRedisClient()
	if err != nil {
		fmt.Println("create redis client failed: ", err)
		return nil
	}
	store, err := redis.NewStoreWithOptions(redis_client, limiter.StoreOptions{
		Prefix: "middle_platform_api_rate_limit",
	})
	if err != nil {
		fmt.Println("create store of limiter failed", err)
		return nil
	}
	instance := limiter.New(store, rate, limiter.WithClientIPHeader("Cf-Connecting-Ip"))
	fmt.Println("create rate limit instance success.")
	return instance
}

// var rate_limiter = NewRateLimiter()

func NewCustomRateLimiter() func(handler middleware.Handler) middleware.Handler {
	rate_limiter := NewRateLimiter()
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if rate_limiter == nil {
				fmt.Println("create rate limiter failed")
				reply, err = handler(ctx, req)
				return
			}

			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				fmt.Println("get server context failed")
				reply, err = handler(ctx, req)
				return
			}
			fmt.Println("Cf-Connecting-Ip:", tr.RequestHeader().Get("Cf-Connecting-Ip"))
			// fmt.Println("uri:", tr.Operation())
			if tr.Operation() == "/api.nft_transfer.v1.NftTransfer/PostReportSpam" && tr.RequestHeader().Get("Cf-Connecting-Ip") != "" {

				ctx1, ok := ctx.(http.Context)
				if !ok {
					fmt.Println("get http context failed")
					reply, err = handler(ctx, req)
					return
				}
				key := rate_limiter.GetIPKey(ctx1.Request())
				context, err := rate_limiter.Get(ctx, key)

				if err != nil {
					return nil, ErrGetLimitKey
				}
				if context.Reached {
					// rejected
					fmt.Println("rate limit reached:", tr.RequestHeader().Get("Cf-Connecting-Ip"))
					return nil, ErrLimitExceed
				}
			}

			// allowed
			reply, err = handler(ctx, req)
			return
		}
	}
}

func newRedisClient() (*libredis.Client, error) {
	uri := "redis://localhost:6379/0"

	addr := os.Getenv("REDIS_ADDR")
	db := os.Getenv("REDIS_DB")

	if addr != "" && db != "" {
		uri = fmt.Sprintf("redis://%s/%s", addr, db)
		fmt.Println("redis uri:", uri)
	}

	opt, err := libredis.ParseURL(uri)
	if err != nil {
		return nil, err
	}

	client := libredis.NewClient(opt)
	return client, nil
}
