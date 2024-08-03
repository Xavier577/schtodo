package handlers

import (
	"net/http"
	"time"

	"github.com/Xavier577/schtodo/http/middlewares"
	"github.com/Xavier577/schtodo/internal"
	"github.com/Xavier577/schtodo/internal/repositories"
	"github.com/Xavier577/schtodo/pkg/passwd"
	"github.com/Xavier577/schtodo/pkg/token"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidUsernameOrPassword = internal.NewHttpReponse(http.StatusBadRequest, "Invalid Username/Password", nil)
	ErrUsernameTaken             = internal.NewHttpReponse(http.StatusConflict, "Username is taken", nil)
)

var Login = internal.POSTController("/auth/login", middlewares.ValidateReq(&LoginPayload{}, middlewares.Body),
	func(cnt *internal.AppContainer) gin.HandlerFunc {
		return func(ctx *gin.Context) {

			payload := ctx.MustGet(middlewares.Body).(*LoginPayload)

			user, errGetUser := cnt.UserRepo.GetByUsername(payload.Username)

			if errGetUser != nil {
				panic(errGetUser)
			}

			if user == nil {
				ErrInvalidUsernameOrPassword.Send(ctx)
				return
			}

			passwdCorrect := passwd.Compare(payload.Password, user.Password)

			if !passwdCorrect {
				ErrInvalidUsernameOrPassword.Send(ctx)
				return
			}

			accessToken, errAccessToken := token.GenerateHS256Token(&token.TokenGenOptions{Payload: token.JWTPayload{"id": user.ID},
				ExpiryDate: time.Now().Add(14 * 24 * time.Hour)})

			if errAccessToken != nil {
				panic(errAccessToken)
			}

			internal.NewHttpReponse(http.StatusOK, "Success", gin.H{"token": accessToken}).Send(ctx)

		}
	})

var Signup = internal.POSTController("/auth/signup", middlewares.ValidateReq(&SignupPayload{}, middlewares.Body),
	func(cnt *internal.AppContainer) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			payload := ctx.MustGet(middlewares.Body).(*SignupPayload)

			userWithUsername, errGetUser := cnt.UserRepo.GetByUsername(payload.Username)

			if errGetUser != nil {
				panic(errGetUser)
			}

			if userWithUsername != nil {
				ErrUsernameTaken.Send(ctx)
				return
			}

			passwdHash, errHashPasswd := passwd.Hash(payload.Password)

			if errHashPasswd != nil {
				panic(errHashPasswd)
			}

			user, errCreateUser := cnt.UserRepo.Create(&repositories.CreateUserData{Username: payload.Username, Password: passwdHash})

			if errCreateUser != nil {
				panic(errCreateUser)
			}

			internal.NewHttpReponse(http.StatusCreated, "Success", user).Send(ctx)
		}
	})
