package app

import (
	"log"
	"net/http"
	"schtodo/http/middlewares"
	"schtodo/http/routes"
	"schtodo/internal"
	"time"

	"github.com/gin-gonic/gin"
)

func NewApp(cnt *internal.AppContainer) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gin.CustomRecovery(middlewares.ErrRecovery))
	r.Use(middlewares.GlobalErrHandler)

	r.GET("/api/health", func(c *gin.Context) {

		var result time.Time

		err := cnt.DB.Get(&result, "select now()")

		if err != nil {
			panic(internal.NewHttpReponse(http.StatusInternalServerError, "Server not healthy", nil))
		}

		log.Printf("Health check at %v", result)

		c.JSON(200, gin.H{"message": "running 🚀"})
	})

	routes.SetupRoutes(cnt, r)

	return r

}
