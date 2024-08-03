package app

import (
	"fmt"
	"github.com/Xavier577/schtodo/config/env"
	"github.com/Xavier577/schtodo/http/middlewares"
	"github.com/Xavier577/schtodo/http/routes"
	"github.com/Xavier577/schtodo/internal"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func NewApp(cnt *internal.AppContainer) *gin.Engine {

	appEnv := os.Getenv("APP_ENV")

	if appEnv == env.Production {
		gin.SetMode(gin.ReleaseMode)

	}

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gin.CustomRecovery(middlewares.ErrRecovery))
	r.Use(middlewares.GlobalErrHandler)

	r.GET("/api/health", func(c *gin.Context) {

		var result time.Time

		err := cnt.DB.Get(&result, "select now()")

		if err != nil {
			panic(internal.NewHttpReponse(http.StatusInternalServerError, fmt.Sprintf("%s state not healthy", appEnv), nil))
		}

		log.Printf("Health check at %v", result)

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s running 🚀", appEnv)})
	})

	routes.SetupRoutes(cnt, r)

	return r

}
