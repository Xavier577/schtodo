package user

import (
	"database/sql"
	"errors"
	"schtodo/internal/repositories"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type userRepo struct {
	DB *sqlx.DB
}

func NewUserRepo(pg *sqlx.DB) repositories.UserRepository {
	return &userRepo{DB: pg}
}

func (u *userRepo) Create(data *repositories.CreateUserData) (*repositories.User, error) {
	user := &repositories.User{ID: ulid.Make().String(), Username: data.Username, Password: data.Password}

	query, _, errSqlBuilder := goqu.Insert("users").Rows(user).Returning("*").ToSQL()

	if errSqlBuilder != nil {
		return nil, errSqlBuilder
	}

	errQueryResult := u.DB.Get(user, query)

	if errQueryResult != nil {
		return nil, errQueryResult
	}

	return user, nil
}

func (u *userRepo) GetById(id string) (*repositories.User, error) {

	query, _, _ := goqu.From("users").Select("*").Where(goqu.C("id").Eq(id)).ToSQL()

	user := &repositories.User{}

	errQuery := u.DB.Get(user, query)

	if errQuery != nil {
		if errors.Is(errQuery, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errQuery
	}

	return user, nil

}

func (u *userRepo) GetByUsername(username string) (*repositories.User, error) {

	query, _, _ := goqu.From("users").Select("*").Where(goqu.C("username").Eq(username)).ToSQL()

	user := &repositories.User{}

	errQuery := u.DB.Get(user, query)

	if errQuery != nil {
		if errors.Is(errQuery, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errQuery
	}

	return user, nil

}
