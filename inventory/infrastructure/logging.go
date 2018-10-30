package infrastructure

import (
   "time"
   "github.com/comolago/shop/inventory/domain"
   "github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
   domain.InventoryHandler
   logger log.Logger
}


func LoggingMiddleware(logger log.Logger) InventoryMiddleware {
   return func(next domain.InventoryHandler) domain.InventoryHandler {
      return loggingMiddleware{next, logger}
   }
}

func (mw loggingMiddleware) AddItem(id int, name string) (output string, err *domain.ErrHandler) {
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
   output, err = mw.InventoryHandler.AddItem(id, name)
   return output, err
}

func (mw loggingMiddleware) GetItemById(id int) (output domain.Item, err *domain.ErrHandler) {
   defer func(begin time.Time) {
      mw.logger.Log(
         "object", "item",
         "action", "get",
         "type", "id",
         "id", id,
         "result", output,
         "took", time.Since(begin),
      )
   }(time.Now())
   output, err = mw.InventoryHandler.GetItemById(id)
   return output, err
}
