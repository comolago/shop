package domain

import (
	"errors"
)


var ErrNoId = errors.New("Id is empty - please provie a value")
var ErrNoName = errors.New("Name is empty - please provie a value")
