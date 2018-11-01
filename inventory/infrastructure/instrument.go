// infrastructure implementation details
package infrastructure

import (
   "time"
   "github.com/comolago/shop/inventory/domain"
   "github.com/go-kit/kit/metrics"
)

// Define metrics
type metricsMiddleware struct {
   domain.InventoryHandler
   requestCount   metrics.Counter
   requestLatency metrics.Histogram
}

// Define the metrics function
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

// Define AddItem helper function to handle metrics
func (mw metricsMiddleware) AddItem(item domain.Item) (output string, err *domain.ErrHandler) {
   defer func(begin time.Time) {
      lvs := []string{"method", "AddItem"}
      mw.requestCount.With(lvs...).Add(1)
      mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
   }(time.Now())
   output, err = mw.InventoryHandler.AddItem(item)
   return output, err
}

// Define GetItemById helper function to handle metrics
func (mw metricsMiddleware) GetItemById(id int) (output domain.Item, err *domain.ErrHandler) {
   defer func(begin time.Time) {
      lvs := []string{"method", "GetItemById"}
      mw.requestCount.With(lvs...).Add(1)
      mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
   }(time.Now())
   output, err = mw.InventoryHandler.GetItemById(id)
   return output, err
}

// Define DelItemById helper function to handle metrics
func (mw metricsMiddleware) DelItemById(id int) (output string, err *domain.ErrHandler)  {
   defer func(begin time.Time) {
      lvs := []string{"method", "DelItem"}
      mw.requestCount.With(lvs...).Add(1)
      mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
   }(time.Now())
   output, err  = mw.InventoryHandler.DelItemById(id)
   return output, err
}
