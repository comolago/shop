// Business Rules
package domain

// Inventory interface
// Open() open connection to database backend
// AddItem add one item to the inventory
// GetItemById retrieve an item given its id
type InventoryHandler interface {
   Open() *ErrHandler
   AddItem(item Item) (string, *ErrHandler)
   GetItemById(id int) (Item, *ErrHandler)
   DelItemById(id int) (string, *ErrHandler)
}

// DB interface
// Open() open connection to DB backend
// GetItemById retrieve an item given its id from database
type DbHandler interface {
   Open() *ErrHandler
   GetItemById(int, *Item) *ErrHandler
   AddItem(Item) *ErrHandler
   DelItemById(id int) *ErrHandler
}
