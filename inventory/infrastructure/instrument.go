package infrastructure

import (
	"time"

	"github.com/comolago/shop/inventory/domain"
	"github.com/go-kit/kit/metrics"
)

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

// Make a new type and wrap into domain.InventoryInt interface
// Add expected metrics property to this type
type metricsMiddleware struct {
	domain.InventoryInt
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

// Implement service functions and add label method for our metrics
func (mw metricsMiddleware) AddItem(id, name string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Word"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output, err = mw.InventoryInt.AddItem(id, name)
	return output, err
}

/*func (mw metricsMiddleware) Sentence(min, max int) (output string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Sentence"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output = mw.domain.InventoryInt.Sentence(min, max)
	return
}

func (mw metricsMiddleware) Paragraph(min, max int) (output string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Paragraph"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output = mw.domain.InventoryInt.Paragraph(min, max)
	return
}*/
