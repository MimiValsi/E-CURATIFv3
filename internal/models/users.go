package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"E-CURATIFv3/database"
)

type User struct {
	ID       int
	Name     string
	NNI      string
	Email    string
	Password string
	Created  time.Time

	DB *pgxpool.Pool
}

func (u *User) Insert(conn *pgxpool.Conn) error {
	ctx := context.Background()

	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return err
	}

	query := `
INSERT INTO users (name, email, hashed_password, created)
VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
RETURNING id;
	`

	id := 0
	err = conn.QueryRow(ctx, query, u.Name, u.Email, string(hashedPasswd)).Scan(&id)
	if err != nil {
		return err
	}

	return nil

}

func (u *User) Authenticate(conn *pgxpool.Conn) (int, error) {
	ctx := context.Background()

	var hashed_passwd []byte

	query := `SELECT id, hashed_password FROM users WHERE email = $1`

	var id int
	err := conn.QueryRow(ctx, query, u.Email).Scan(&id, &hashed_passwd)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, database.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashed_passwd, []byte(u.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, database.ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

// func (u *User) GetUser(conn *pgxpool.Conn) (string, error) {
// 	ctx := context.Background()
//
// 	query := `SELECT nni FROM users WHERE`
// }

func (u *User) Exists(conn *pgxpool.Conn) (bool, error) {
	ctx := context.Background()
	var exists bool

	query := `SELECT EXISTS(SELECT true FROM users WHERE id = $1)`
	err := conn.QueryRow(ctx, query, u.ID).Scan(&exists)
	return exists, err
}
