package database

import (
	"errors"
)

var (
	ErrNoRecord = errors.New("Models: No matching record found")
)
