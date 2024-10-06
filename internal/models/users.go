package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"E-CURATIFv3/database"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time

	DB *pgxpool.Pool
}

func (u *User) Insert(conn *pgxpool.Conn) error {
	ctx := context.Background()

	hashedPasswd, err := bcrypt.GenerateFromPassword(u.HashedPassword, 12)
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
	fmt.Println("user sent")

	return nil

}

func (u *User) Authenticate(conn *pgxpool.Conn) (int, error) {
	ctx := context.Background()

	var hashed_passwd []byte

	query := `SELECT id, hashed_password FROM users WHERE email = $1`

	err := conn.QueryRow(ctx, query, u.Email).Scan(&u.ID, &hashed_passwd)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			fmt.Println("pas bon du tout")
			return 0, database.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashed_passwd, u.HashedPassword)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, database.ErrInvalidCredentials
		}
		return 0, err
	}

	return u.ID, nil
}

func (u *User) Exists(id int) (bool, error) {
	return false, nil
}
