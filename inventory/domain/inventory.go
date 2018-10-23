package domain

type InventoryInt interface {
	AddItem(id, name string) (string, error)
	Count(string) int
}

type Inventory struct{}

func (Inventory) AddItem(id, name string) (string, error) {
	if id == "" {
		return "", ErrEmpty
	}
	if name == "" {
		return "", ErrEmpty
	}
	return "item added", nil
}

type InventoryMiddleware func(InventoryInt) InventoryInt
