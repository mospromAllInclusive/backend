package entities

import "time"

type Database struct {
	ID        int64      `db:"id"`
	Name      string     `db:"name"`
	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type UsersDatabase struct {
	UserID     int64      `db:"user_id"`
	DatabaseID int64      `db:"database_id"`
	Role       Role       `db:"role"`
	CreatedAt  time.Time  `db:"created_at"`
	DeletedAt  *time.Time `db:"deleted_at"`

	Name string `db:"name"`
}

type Role string

const (
	RoleReader Role = "reader"
	RoleWriter Role = "writer"
	RoleAdmin  Role = "admin"
)

func (r Role) Priority() int {
	switch r {
	case RoleAdmin:
		return 3
	case RoleWriter:
		return 2
	case RoleReader:
		return 1
	default:
		return 0
	}
}

func (r Role) Authorize(role Role) bool {
	return r.Priority() >= role.Priority()
}
