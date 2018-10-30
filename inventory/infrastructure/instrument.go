package infrastructure

import (
   "time"
   "github.com/comolago/shop/inventory/domain"
   "github.com/go-kit/kit/metrics"
)

type metricsMiddleware struct {
   domain.InventoryHandler
   requestCount   metrics.Counter
   requestLatency metrics.Histogram
}

func Metrics(requestCount metrics.Counter,
   requestLatency metrics.Histogram) InventoryMiddleware {
   return func(next domain.InventoryHandler) domain.InventoryHandler {
      return metricsMiddleware{
         next,
         requestCount,
         requestLatency,
      }
   }
}

func (mw metricsMiddleware) AddItem(id int, name string) (output string, err *domain.ErrHandler) {
   defer func(begin time.Time) {
      lvs := []string{"method", "AddItem"}
      mw.requestCount.With(lvs...).Add(1)
      mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
   }(time.Now())
   output, err = mw.InventoryHandler.AddItem(id, name)
   return output, err
}

func (mw metricsMiddleware) GetItemById(id int) (output domain.Item, err *domain.ErrHandler) {
   defer func(begin time.Time) {
      lvs := []string{"method", "GetItem"}
      mw.requestCount.With(lvs...).Add(1)
      mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
   }(time.Now())
   output, err = mw.InventoryHandler.GetItemById(id)
   return output, err
}
