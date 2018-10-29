package domain

type InventoryHandler interface {
	AddItem(id int, name string) (string, *ErrHandler)
   AddItemDBHandler()
Open()
}

type DbHandler interface {
   Open() *ErrHandler
   getItem(inventory *Inventory) *ErrHandler
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

func (i Inventory) AddItemDBHandler()  {
   
   //i.Db=db
}

func (i Inventory) Open()  {
   i.Db.Open()
}
