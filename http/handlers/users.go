package handlers

import (
	"net/http"

	"github.com/Xavier577/schtodo/http/middlewares"
	"github.com/Xavier577/schtodo/internal"
	"github.com/Xavier577/schtodo/pkg/token"

	"github.com/gin-gonic/gin"
)

var GetUser = internal.GETController("/user", middlewares.Auth, func(cnt *internal.AppContainer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claim := ctx.MustGet(middlewares.Body).(token.JWTPayload)

		userID := claim["id"].(string)

		user, errGetUser := cnt.UserRepo.GetById(userID)

		if errGetUser != nil {
			panic(errGetUser)
		}

		internal.NewHttpReponse(http.StatusOK, "Success", user).Send(ctx)
	}
})
