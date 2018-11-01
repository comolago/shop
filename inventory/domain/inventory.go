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

// Delete an Item by its ID from the persistence backend
func (i Inventory) DelItemById(id int) (string, *ErrHandler) {
   err := i.Db.DelItemById(id)
   if err != nil {
      return "", err
   }
   return "item deleted", nil
}


// Add an item to the persistence backend
func (i Inventory) AddItem(item Item) (string, *ErrHandler) {
   if item.Id == 0 {
      return "", &ErrHandler{10, "func (Inventory)", "AddItem", ""}
   }
   if item.Name == "" {
      return "", &ErrHandler{11, "func (Inventory)", "AddItem", ""}
   }
   if item.Quantity == 0 {
      return "", &ErrHandler{12, "func (Inventory)", "AddItem", ""}
   }
   err := i.Db.AddItem(item)
   if err != nil {
      return "", err
   }
   return "item added", nil
}

