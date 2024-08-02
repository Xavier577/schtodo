package middlewares

import (
	"io"
	"net/http"
	"schtodo/internal"

	"github.com/gin-gonic/gin"
)

const (
	Params = "params"
	Body   = "param"
	Query  = "query"
)

var (
	ErrInvalidJson = internal.NewHttpReponse(http.StatusBadRequest, "Invalid Json", nil)
)

func ValidateReq(dto any, context string) internal.ControllerHandler {

	return func(cnt *internal.AppContainer) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			var bindFn func(obj any) error

			switch context {
			case Params:
				bindFn = ctx.ShouldBindUri
			case Query:
			case Body:
				bindFn = ctx.ShouldBind
			}

			if bindFn != nil {
				if err := bindFn(dto); err != nil {
					if err.Error() == "EOF" || err == io.ErrUnexpectedEOF {
						ErrInvalidJson.Send(ctx)
						ctx.Abort()
						return
					}
					internal.NewHttpReponse(http.StatusBadRequest, err.Error(), nil).Send(ctx)
					ctx.Abort()
				} else {

					ctx.Set("payload", dto)
					ctx.Next()
				}
			}

		}
	}
}
