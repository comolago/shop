package infrastructure

import (
   "github.com/comolago/shop/inventory/domain"
)

type InventoryMiddleware func(domain.InventoryHandler) domain.InventoryHandler
