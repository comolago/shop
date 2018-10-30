package infrastructure

import (
	"github.com/juju/ratelimit"
	"github.com/go-kit/kit/endpoint"
	"context"
)

//var ErrLimitExceed = errors.New("Rate Limit Exceed")

func NewTokenBucketLimiter(tb *ratelimit.Bucket) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			if tb.TakeAvailable(1) == 0 {
				//return nil, ErrLimitExceed
				return nil, nil
			}
			return next(ctx, request)
		}
	}
}
