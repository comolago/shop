package domain

import (
	"errors"
)

// Inventory provides operations on strings.
type InventoryInt interface {
	AddItem(id, name string) (string, error)
	Count(string) int
}

// Inventory is a concrete implementation of Inventory
type Inventory struct{}

// create type that return function.
// this will be needed in main.go
type InventoryMiddleware func(InventoryInt) InventoryInt

func (Inventory) AddItem(id, name string) (string, error) {
	if id == "" {
		return "", ErrEmpty
	}
	if name == "" {
		return "", ErrEmpty
	}
	return "item added", nil
}

func (Inventory) Count(s string) int {
	return len(s)
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")
