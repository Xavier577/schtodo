package internal

import (
	"github.com/Xavier577/schtodo/internal/repositories"

	"github.com/jmoiron/sqlx"
)

type AppContainer struct {
	DB       *sqlx.DB
	UserRepo repositories.UserRepository
	TodoRepo repositories.TodoRepository
}
