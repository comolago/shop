package infrastructure

import (
   "github.com/juju/ratelimit"
   "github.com/go-kit/kit/endpoint"
   "context"
   "github.com/comolago/shop/inventory/domain"
)

//var ErrLimitExceed = errors.New("Rate Limit Exceed")

func NewTokenBucketLimiter(tb *ratelimit.Bucket) endpoint.Middleware {
   return func(next endpoint.Endpoint) endpoint.Endpoint {
      return func(ctx context.Context, request interface{}) (interface{}, error) {
         if tb.TakeAvailable(1) == 0 {
            return nil, &domain.ErrHandler{8, "func ", "NewTokenBucketLimiter(tb *ratelimit.Bucket)", ""}
            //return nil, ErrLimitExceed
         }
         return next(ctx, request)
      }
   }
}
