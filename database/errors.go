package database

import (
	"errors"
)

// variable global pour être utilisée pour chaque demande de connection
// avec PSQL
// global variable(s) to be used for each PSQL connexion
var (
	ErrNoRecord           = errors.New("models: No matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)
