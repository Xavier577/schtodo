package internal

import (
	"schtodo/internal/repositories"

	"github.com/jmoiron/sqlx"
)

type AppContainer struct {
	DB       *sqlx.DB
	UserRepo repositories.UserRepository
}
