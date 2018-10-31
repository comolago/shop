// infrastructure implementation details
package infrastructure

import (
   "github.com/comolago/shop/inventory/domain"
)

// define a type for the middleware helpers
type InventoryMiddleware func(domain.InventoryHandler) domain.InventoryHandler
