package database

import (
	"errors"
)

// global variable(s) to be used for each PSQL connexion
var (
	ErrNoRecord = errors.New("Models: No matching record found")
)
