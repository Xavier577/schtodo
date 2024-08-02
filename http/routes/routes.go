package routes

import (
	"schtodo/internal"

	_ "schtodo/http/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(cnt *internal.AppContainer, r *gin.Engine) {

	for _, v := range internal.AppControllerList() {

		routeHandlers := []gin.HandlerFunc{}

		for _, controllerHandler := range v.Handlers {
			routeHandlers = append(routeHandlers, controllerHandler(cnt))
		}

		switch v.Method {
		case internal.GET:
			r.GET(v.Path, routeHandlers...)
		case internal.POST:
			r.POST(v.Path, routeHandlers...)
		case internal.PUT:
			r.PUT(v.Path, routeHandlers...)
		case internal.PATCH:
			r.PATCH(v.Path, routeHandlers...)
		case internal.DELETE:
			r.PATCH(v.Path, routeHandlers...)
		}

	}

}
