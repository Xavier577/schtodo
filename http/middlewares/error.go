package middlewares

import (
	"github.com/Xavier577/schtodo/internal"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrRecovery(ctx *gin.Context, err any) {
	httpErr, isHttpErr := err.(*internal.HttpResponse)

	if isHttpErr {

		if httpErr.StatusCode > http.StatusInternalServerError {
			log.Println(err)
		}

		httpErr.Send(ctx)

	} else {

		// send 500 for unknown err type
		log.Println(err)
		internal.ErrInternalServer.Send(ctx)
	}

}

func GlobalErrHandler(ctx *gin.Context) {

	ctx.Next()

	for _, err := range ctx.Errors {
		ErrRecovery(ctx, err.Err)

	}

}
