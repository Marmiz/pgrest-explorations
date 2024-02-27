package data

import (
	"time"

	"github.com/marmiz/pgrest-explorations/internal/db"
)

type User struct {
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	PasswordHash []byte    `json:"-"`
	Active       bool      `json:"active"`
	Role         string    `json:"role"`
}

func NewFromDB(a *db.AuthUser) User {
	return User{
		Email:        a.Email,
		CreatedAt:    a.CreatedAt.Time,
		PasswordHash: a.PasswordHash,
		Active:       a.Active,
		Role:         a.Role,
	}
}
