package infrastructure

import (
	"time"

	"github.com/comolago/shop/inventory/domain"
	"github.com/go-kit/kit/metrics"
)

type metricsMiddleware struct {
	domain.InventoryInt
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

func Metrics(requestCount metrics.Counter,
	requestLatency metrics.Histogram) domain.InventoryMiddleware {
	return func(next domain.InventoryInt) domain.InventoryInt {
		return metricsMiddleware{
			next,
			requestCount,
			requestLatency,
		}
	}
}

func (mw metricsMiddleware) AddItem(id int, name string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Word"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output, err = mw.InventoryInt.AddItem(id, name)
	return output, err
}
