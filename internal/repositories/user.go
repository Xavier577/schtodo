package repositories

import "time"

type User struct {
	ID        string    `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Password  string    `db:"password" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateUserData struct {
	Username string `db:"username"`
	Password string `db:"password"`
}

type UserRepository interface {
	Create(data *CreateUserData) (*User, error)
	GetById(id string) (*User, error)
	GetByUsername(username string) (*User, error)
}
