package infrastructure

import (
	"time"

	"github.com/comolago/shop/inventory/domain"
	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	domain.InventoryInt
	logger log.Logger
}


func LoggingMiddleware(logger log.Logger) domain.InventoryMiddleware {
	return func(next domain.InventoryInt) domain.InventoryInt {
		return loggingMiddleware{next, logger}
	}
}

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
