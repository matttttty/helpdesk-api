package model

import "time"

type Role string

const (
	RoleClient Role = "client"
	RoleAgent  Role = "agent"
	RoleAdmin  Role = "admin"
)

type User struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Password  string    `json:"-" db:"password"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created" db:"created_at"`
	Role      Role      `json:"role" db:"role"`
}
