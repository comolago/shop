// Business Rules
package domain

// Inventory is an array of items along with the
// DB interface to store and retrieve them
type Inventory struct{
   Items []Item
   Db DbHandler
}

// Open a connection to the persistence backend
func (i Inventory) Open() *ErrHandler  {
   return i.Db.Open()
}

// Retrive an Item by its ID from the persistence backend
func (i Inventory) GetItemById(id int) (Item, *ErrHandler) {
   var item Item
   err := i.Db.GetItemById(id,&item)
   return item, err
}

// Add an item to the persistence backend
func (i Inventory) AddItem(id int, name string) (string, *ErrHandler) {
   if id == 0 {
      return "", &ErrHandler{10, "func (Inventory)", "AddItem", ""}
   }
   if name == "" {
      return "", &ErrHandler{11, "func (Inventory)", "AddItem", ""}
   }
   return "item added", nil
}

