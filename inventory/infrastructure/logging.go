// infrastructure implementation details
package infrastructure

import (
   "time"
   "github.com/comolago/shop/inventory/domain"
   "github.com/go-kit/kit/log"
)

// Define logger
type loggingMiddleware struct {
   domain.InventoryHandler
   logger log.Logger
}


// Define logging function
func LoggingMiddleware(logger log.Logger) InventoryMiddleware {
   return func(next domain.InventoryHandler) domain.InventoryHandler {
      return loggingMiddleware{next, logger}
   }
}

// Define AddItem helper function to handle logging
func (mw loggingMiddleware) AddItem(item domain.Item) (output string, err *domain.ErrHandler) {
   defer func(begin time.Time) {
      mw.logger.Log(
         "object", "item",
         "action", "add",
         "id", item.Id,
         "name", item.Name,
         "quantity", item.Quantity,
         "result", output,
         "took", time.Since(begin),
      )
   }(time.Now())
   output, err = mw.InventoryHandler.AddItem(item)
   return output, err
}

// Define GetItemById helper function to handle logging
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

// Define DelItemById helper function to handle logging
func (mw loggingMiddleware) DelItemById(id int) (output string, err *domain.ErrHandler) {
   defer func(begin time.Time) {
      mw.logger.Log(
         "object", "item",
         "action", "del",
         "type", "id",
         "id", id,
         "result", "deleted",
         "took", time.Since(begin),
      )
   }(time.Now())
   output, err  = mw.InventoryHandler.DelItemById(id)
   return output, err 
}
