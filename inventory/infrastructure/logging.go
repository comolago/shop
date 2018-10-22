package infrastructure

import (
	"time"

	"github.com/comolago/shop/inventory/domain"
	"github.com/go-kit/kit/log"
)

// implement function to return domain.InventoryMiddleware
func LoggingMiddleware(logger log.Logger) domain.InventoryMiddleware {
	return func(next domain.InventoryInt) domain.InventoryInt {
		return loggingMiddleware{next, logger}
	}
}

// Make a new type and wrap into domain.InventoryInt interface
// Add logger property to this type
type loggingMiddleware struct {
	domain.InventoryInt
	logger log.Logger
}

// Implement domain.InventoryInt Interface for LoggingMiddleware
func (mw loggingMiddleware) AddItem(id, name string) (output string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"object", "item",
			"action", "add",
			"id", id,
			"name", name,
			"result", output,
			"took", time.Since(begin),
		)
	}(time.Now())
	output, err = mw.InventoryInt.AddItem(id, name)
	return output, err
}

/*func (mw loggingMiddleware) Sentence(min, max int) (output string) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "Sentence",
			"min", min,
			"max", max,
			"result", output,
			"took", time.Since(begin),
		)
	}(time.Now())
	output = mw.InventoryInt.Sentence(min, max)
	return
}

func (mw loggingMiddleware) Paragraph(min, max int) (output string) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "Paragraph",
			"min", min,
			"max", max,
			"result", output,
			"took", time.Since(begin),
		)
	}(time.Now())
	output = mw.InventoryInt.Paragraph(min, max)
	return
}*/
