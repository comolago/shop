// infrastructure implementation details
package infrastructure

import (
   "github.com/juju/ratelimit"
   "github.com/go-kit/kit/endpoint"
   "context"
   "github.com/comolago/shop/inventory/domain"
)

// create a rate limit toklen bucket
func NewTokenBucketLimiter(tb *ratelimit.Bucket) endpoint.Middleware {
   return func(next endpoint.Endpoint) endpoint.Endpoint {
      return func(ctx context.Context, request interface{}) (interface{}, error) {
         if tb.TakeAvailable(1) == 0 {
            return nil, &domain.ErrHandler{8, "func ", "NewTokenBucketLimiter(tb *ratelimit.Bucket)", ""}
         }
         return next(ctx, request)
      }
   }
}
