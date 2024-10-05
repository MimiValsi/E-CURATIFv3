package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time

	DB *pgxpool.Pool
}

func (m *User) Insert(conn *pgxpool.Conn) error {
	ctx := context.Background()

	hashedPasswd, err := bcrypt.GenerateFromPassword(m.HashedPassword, 12)
	if err != nil {
		return err
	}

	query := `
INSERT INTO users (name, email, hashed_password, created)
VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
RETURNING id;
	`

	id := 0
	err = conn.QueryRow(ctx, query, m.Name, m.Email, string(hashedPasswd)).Scan(&id)
	if err != nil {
		return err
	}
	fmt.Println("user sent")

	return nil

}

func (m *User) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *User) Exists(id int) (bool, error) {
	return false, nil
}
