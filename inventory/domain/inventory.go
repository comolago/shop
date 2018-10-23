package domain

type InventoryInt interface {
	AddItem(id, name string) (string, error)
}

type Inventory struct{}

func (Inventory) AddItem(id, name string) (string, error) {
	if id == "" {
		return "", ErrNoId
	}
	if name == "" {
		return "", ErrNoName
	}
	return "item added", nil
}

type InventoryMiddleware func(InventoryInt) InventoryInt
