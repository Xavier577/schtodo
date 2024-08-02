package handlers

import (
	"net/http"
	"schtodo/http/dtos"
	"schtodo/http/middlewares"
	"schtodo/internal"
	"schtodo/internal/repositories"
	"schtodo/pkg/passwd"
	"schtodo/pkg/token"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidUsernameOrPassword = internal.NewHttpReponse(http.StatusBadRequest, "Invalid Username/Password", nil)
	ErrUsernameTaken             = internal.NewHttpReponse(http.StatusConflict, "Username is taken", nil)
)

var Login = internal.POSTController("/auth/login", middlewares.ValidateReq(&dtos.LoginPayload{}, middlewares.Body), func(appCnt *internal.AppContainer) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		payload := ctx.MustGet("payload").(*dtos.LoginPayload)

		user, errGetUser := appCnt.UserRepo.GetByUsername(payload.Username)

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

		accessToken, errAccessToken := token.GenerateHS256Token(token.TokenGenOptions{Payload: token.JWTPayload{"id": user.ID},
			ExpiryDate: time.Now().Add(time.Duration(time.Second) * 10)})

		if errAccessToken != nil {
			panic(errAccessToken)
		}

		internal.NewHttpReponse(http.StatusOK, "Success", gin.H{"token": accessToken}).Send(ctx)

	}
})

var SignUp = internal.POSTController("/auth/signup", middlewares.ValidateReq(&dtos.SignupPayload{}, middlewares.Body), func(appCnt *internal.AppContainer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload := ctx.MustGet("payload").(*dtos.SignupPayload)

		userWithUsername, errGetUser := appCnt.UserRepo.GetByUsername(payload.Username)

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

		user, errCreateUser := appCnt.UserRepo.Create(&repositories.CreateUserData{Username: payload.Username, Password: passwdHash})

		if errCreateUser != nil {
			panic(errCreateUser)
		}

		internal.NewHttpReponse(http.StatusCreated, "Success", user).Send(ctx)
	}
})
