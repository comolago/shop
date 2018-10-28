package domain

type InventoryHandler interface {
	AddItem(id int, name string) (string, *ErrHandler)
}

type Inventory struct{
	Items []Item
}

func (Inventory) AddItem(id int, name string) (string, *ErrHandler) {
	if id == 0 {
		return "", &ErrHandler{1, "func (Inventory)", "AddItem", ""}
	}
	if name == "" {
		return "", &ErrHandler{2, "func (Inventory)", "AddItem", ""}
	}
	return "item added", nil
}

