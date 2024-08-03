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
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}

type UpdatableUserFields struct {
	Username  string    `db:"username" json:"username"`
	Password  string    `db:"password" json:"password"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type UserRepository interface {
	Create(data *CreateUserData) (*User, error)
	GetById(id string) (*User, error)
	GetByUsername(username string) (*User, error)
	Update(id string, data *UpdatableUserFields) (*User, error)
}
