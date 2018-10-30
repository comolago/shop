package domain

/*import (
   "fmt"
)*/

type Inventory struct{
   Items []Item
   Db DbHandler
}

func (i Inventory) Open() *ErrHandler  {
   return i.Db.Open()
}

func (i Inventory) GetItemById(id int) (Item, *ErrHandler) {
   var item Item
   err := i.Db.GetItemById(id,&item)
   return item, err
}

func (i Inventory) AddItem(id int, name string) (string, *ErrHandler) {
   if id == 0 {
      return "", &ErrHandler{10, "func (Inventory)", "AddItem", ""}
   }
   if name == "" {
      return "", &ErrHandler{11, "func (Inventory)", "AddItem", ""}
   }
   return "item added", nil
}

