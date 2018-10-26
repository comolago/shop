package domain

type InventoryInt interface {
	AddItem(id int, name string) (string, error)
}

type Inventory struct{
	Items []Item
}

func (Inventory) AddItem(id int, name string) (string, error) {
	if id == 0 {
		return "", ErrNoId
	}
	if name == "" {
		return "", ErrNoName
	}
	return "item added", nil
}

type InventoryMiddleware func(InventoryInt) InventoryInt
