package domain

import (
   "fmt"
)

type InventoryHandler interface {
   Open()
   AddItem(id int, name string) (string, *ErrHandler)
   GetItemById(id int)
}

type DbHandler interface {
   Open() *ErrHandler
   GetItemById(int, *Item) *ErrHandler
}

type Inventory struct{
   Items []Item
   Db DbHandler
}

func (i Inventory) AddItem(id int, name string) (string, *ErrHandler) {
   if id == 0 {
      return "", &ErrHandler{1, "func (Inventory)", "AddItem", ""}
   }
   if name == "" {
      return "", &ErrHandler{2, "func (Inventory)", "AddItem", ""}
   }
   return "item added", nil
}

func (i Inventory) Open()  {
   i.Db.Open()
}

func (i Inventory) GetItemById(id int) {
   var item Item
   i.Db.GetItemById(id,&item)
   fmt.Println(item.Name)
}
