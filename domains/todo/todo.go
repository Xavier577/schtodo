package todo

import (
	"database/sql"
	"errors"

	"github.com/Xavier577/schtodo/internal/repositories"
	"github.com/Xavier577/schtodo/pkg/objects"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type todoRepo struct {
	DB *sqlx.DB
}

func NewTodoRepo(pg *sqlx.DB) repositories.TodoRepository {
	return &todoRepo{DB: pg}
}

func (t *todoRepo) Create(data *repositories.CreateTodoFields) (*repositories.Todo, error) {
	todo := &repositories.Todo{ID: ulid.Make().String()}

	errStructMerge := objects.MarshalStructMerge(todo, data)

	if errStructMerge != nil {
		return nil, errStructMerge
	}

	query, _, errQueryBuilder := goqu.Insert("todos").Rows(todo).Returning("*").ToSQL()

	if errQueryBuilder != nil {
		return nil, errQueryBuilder
	}

	errQueryResult := t.DB.Get(todo, query)

	if errQueryResult != nil {
		return nil, errQueryResult
	}

	return todo, nil

}

func (t *todoRepo) GetById(id string) (*repositories.Todo, error) {

	todo := &repositories.Todo{}

	query, _, _ := goqu.From("todos").Select("*").Where(goqu.C("id").Eq(id)).ToSQL()

	errQuery := t.DB.Get(todo, query)

	if errQuery != nil {
		if errors.Is(errQuery, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errQuery
	}

	return todo, nil
}

func (t *todoRepo) GetUserOwnTodo(id, userID string) (*repositories.Todo, error) {

	todo := &repositories.Todo{}

	query, _, _ := goqu.From("todos").Select("*").Where(goqu.Ex{"id": id, "user_id": userID}).ToSQL()

	errQuery := t.DB.Get(todo, query)

	if errQuery != nil {
		if errors.Is(errQuery, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errQuery
	}

	return todo, nil
}

func (t *todoRepo) GetUserTodos(userID string) ([]*repositories.Todo, error) {
	todo := []*repositories.Todo{}

	query, _, _ := goqu.From("todos").Select("*").Where(goqu.C("user_id").Eq(userID)).ToSQL()

	errQuery := t.DB.Select(&todo, query)

	if errQuery != nil {
		if errors.Is(errQuery, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errQuery
	}

	return todo, nil
}

func (t *todoRepo) Update(id string, data *repositories.UpdatableTodoFields) (*repositories.Todo, error) {

	query, _, errQueryBuilder := goqu.Update("todos").Where(goqu.C("id").Eq("id")).Returning("*").ToSQL()

	if errQueryBuilder != nil {
		return nil, errQueryBuilder
	}

	todo := &repositories.Todo{}

	errQueryResult := t.DB.Get(todo, query)

	if errQueryResult != nil {
		return nil, errQueryResult
	}

	return todo, nil
}

func (t *todoRepo) Delete(id string) error {

	query, _, _ := goqu.Delete("todos").Where(goqu.C("id").Eq(id)).ToSQL()

	_, errQueryExec := t.DB.Exec(query)

	if errQueryExec != nil {
		return errQueryExec
	}

	return nil
}
