// infrastructure implementation details
package infrastructure

import (
   "github.com/comolago/shop/inventory/domain"
)

// define a type for the middleware helpers
type AuthMiddleware func(domain.AuthHandler) domain.AuthHandler

