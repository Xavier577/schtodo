package middlewares

import (
	"github.com/Xavier577/schtodo/internal"
	"github.com/Xavier577/schtodo/pkg/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrUnauthorized = internal.NewHttpReponse(http.StatusUnauthorized, "Unauthorized", nil)
)

func Auth(cnt *internal.AppContainer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken := token.ExtractFromHeader(ctx)

		if accessToken == "" {
			ErrUnauthorized.Send(ctx)
			return
		}

		tokenValid, claim, errTokenVerify := token.Verify(&token.TokenVerifyOptions{SignedToken: accessToken})

		if errTokenVerify != nil {
			switch errTokenVerify {
			case token.ErrExpiredToken:
				internal.NewHttpReponse(http.StatusBadRequest, token.ErrExpiredToken.Error(), nil).Send(ctx)
				ctx.Abort()
				return
			case token.ErrParseToken:
				internal.NewHttpReponse(http.StatusBadRequest, token.ErrParseToken.Error(), nil).Send(ctx)
				ctx.Abort()
				return
			default:
				panic(errTokenVerify)

			}
		}

		if !tokenValid {
			ErrUnauthorized.Send(ctx)
			return
		}

		ctx.Set("claim", claim.Payload)
		ctx.Next()

	}
}
