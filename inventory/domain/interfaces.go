package domain

type InventoryHandler interface {
   Open() *ErrHandler
   AddItem(id int, name string) (string, *ErrHandler)
   GetItemById(id int) (Item, *ErrHandler)
}

type DbHandler interface {
   Open() *ErrHandler
   GetItemById(int, *Item) *ErrHandler
}
