// Business Rules
package domain

// Inventory interface
// Open() open connection to database backend
// AddItem add one item to the inventory
// GetItemById retrieve an item given its id
type InventoryHandler interface {
   Open() *ErrHandler
   AddItem(id int, name string) (string, *ErrHandler)
   GetItemById(id int) (Item, *ErrHandler)
}

// DB interface
// Open() open connection to DB backend
// GetItemById retrieve an item given its id from database
type DbHandler interface {
   Open() *ErrHandler
   GetItemById(int, *Item) *ErrHandler
}
