package internal

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var ErrInternalServer = NewHttpReponse(500, "Internal Server Error", nil)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	PATCH  = "PATCH"
	DELETE = "DELETE"
)

type ControllerHandler func(cnt *AppContainer) gin.HandlerFunc

type AppController struct {
	Path     string
	Method   string
	Handlers []ControllerHandler
}

var appControllerList []*AppController

func GETController(path string, handlers ...ControllerHandler) *AppController {
	return Controller(path, GET, handlers...)
}

func POSTController(path string, handlers ...ControllerHandler) *AppController {
	return Controller(path, POST, handlers...)
}

func PUTController(path string, handlers ...ControllerHandler) *AppController {
	return Controller(path, PUT, handlers...)
}

func PATCHController(path string, handlers ...ControllerHandler) *AppController {
	return Controller(path, PATCH, handlers...)
}

func DELETEController(path string, handlers ...ControllerHandler) *AppController {
	return Controller(path, DELETE, handlers...)
}

func Controller(path, method string, handlers ...ControllerHandler) *AppController {

	newController := &AppController{Path: path, Method: method, Handlers: handlers}

	appControllerList = append(appControllerList, newController)

	return newController
}

func AppControllerList() []*AppController {
	return appControllerList
}

type HttpResponse struct {
	StatusCode int
	Response   *HttpResponseBody
}

type HttpResponseBody struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (e HttpResponse) ResponseBody() map[string]any {
	return map[string]any{"message": e.Response.Message, "data": e.Response.Data}
}

func (e HttpResponse) Error() string {
	return fmt.Sprintln(e.Response.Message)
}

func (e HttpResponse) Send(ctx *gin.Context) {
	ctx.JSON(e.StatusCode, e.Response)
}

func NewHttpReponse(code int, message string, data any) *HttpResponse {
	httpErrResBody := &HttpResponseBody{Message: message, Data: data}

	return &HttpResponse{StatusCode: code, Response: httpErrResBody}
}
