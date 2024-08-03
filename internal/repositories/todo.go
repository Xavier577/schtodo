package repositories

import "time"

type Todo struct {
	ID        string    `db:"id" json:"id"`
	Task      string    `db:"task" json:"task"`
	IsTimed   bool      `db:"is_timed" json:"is_timed"`
	Deadline  time.Time `db:"deadline" json:"deadline"`
	UserID    string    `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateTodoFields struct {
	Task     string    `db:"task" json:"task"`
	IsTimed  bool      `db:"is_timed" json:"is_timed"`
	UserID   string    `db:"user_id" json:"user_id"`
	Deadline time.Time `db:"deadline" json:"deadline"`
}

type UpdatableTodoFields struct {
	CreateTodoFields
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type TodoRepository interface {
	Create(data *CreateTodoFields) (*Todo, error)
	GetById(id string) (*Todo, error)
	GetUserOwnTodo(id, userID string) (*Todo, error)
	GetUserTodos(userID string) ([]*Todo, error)
	Update(id string, data *UpdatableTodoFields) (*Todo, error)
	Delete(id string) error
}
