package handlers

import (
	"net/http"

	"github.com/Xavier577/schtodo/http/middlewares"
	"github.com/Xavier577/schtodo/internal"
	"github.com/Xavier577/schtodo/internal/repositories"
	"github.com/Xavier577/schtodo/pkg/objects"
	"github.com/Xavier577/schtodo/pkg/token"

	"github.com/gin-gonic/gin"
)

var (
	ErrTodoNotFound = internal.NewHttpReponse(http.StatusNotFound, "Todo not found", nil)
)

var CreateTodo = internal.POSTController("/todo", middlewares.ValidateReq(&CreateTodoPayload{}, middlewares.Body),
	middlewares.Auth, func(cnt *internal.AppContainer) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			claim := ctx.MustGet("claim").(token.JWTPayload)

			userID := claim["id"].(string)

			payload := ctx.MustGet(middlewares.Body).(*CreateTodoPayload)

			todoField := &repositories.CreateTodoFields{UserID: userID}

			todoData := map[string]any{
				"task":     payload.Task,
				"is_timed": payload.IsTimed,
				"deadline": payload.Deadline.Time,
			}

			objects.MustMarshalStructMerge(todoField, todoData)

			todo, errCreateTodo := cnt.TodoRepo.Create(todoField)

			if errCreateTodo != nil {
				panic(errCreateTodo)
			}

			internal.NewHttpReponse(http.StatusOK, "Success", todo).Send(ctx)

		}
	})

var GetUserTodoList = internal.GETController("/todo", middlewares.Auth, func(cnt *internal.AppContainer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claim := ctx.MustGet("claim").(token.JWTPayload)

		userID := claim["id"].(string)

		todos, errFetchTodos := cnt.TodoRepo.GetUserTodos(userID)

		if errFetchTodos != nil {
			panic(errFetchTodos)
		}

		internal.NewHttpReponse(http.StatusOK, "Success", todos).Send(ctx)

	}
})

var GetTodo = internal.GETController("/todo/:id", middlewares.ValidateReq(&IDParam{}, middlewares.Params), middlewares.Auth,
	func(cnt *internal.AppContainer) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			claim := ctx.MustGet("claim").(token.JWTPayload)
			userID := claim["id"].(string)

			params := ctx.MustGet(middlewares.Params).(*IDParam)

			todo, errFetchTodo := cnt.TodoRepo.GetUserOwnTodo(params.ID, userID)

			if errFetchTodo != nil {
				panic(errFetchTodo)
			}

			internal.NewHttpReponse(http.StatusOK, "Success", todo).Send(ctx)

		}
	})

var DeleteTodo = internal.DELETEController("/todo/:id", middlewares.ValidateReq(&IDParam{}, middlewares.Params), middlewares.Auth,
	func(cnt *internal.AppContainer) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			claim := ctx.MustGet("claim").(token.JWTPayload)
			userID := claim["id"].(string)

			params := ctx.MustGet(middlewares.Params).(*IDParam)

			todo, errFetchTodo := cnt.TodoRepo.GetUserOwnTodo(params.ID, userID)

			if errFetchTodo != nil {
				panic(errFetchTodo)
			}

			if todo == nil {
				ErrTodoNotFound.Send(ctx)
				return
			}

			errDeleteTodo := cnt.TodoRepo.Delete(params.ID)

			if errDeleteTodo != nil {
				panic(errDeleteTodo)
			}

			internal.NewHttpReponse(http.StatusOK, "Success", nil).Send(ctx)

		}
	})
