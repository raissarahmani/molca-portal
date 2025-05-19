package network

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type SendResponse interface {
	SuccessMsgResponse(message string)
	SuccessDataResponse(message string, data any)
	BadRequestError(message string, err error)
	ForbiddenError(message string, err error)
	UnauthorizedError(message string, err error)
	NotFoundError(message string, err error)
	InternalServerError(message string, err error)
	MixedError(err error)
}

type BaseController interface {
	ResponseSender
	Path() string
	Authentication() gin.HandlerFunc
	Authorization(role string) gin.HandlerFunc
}

type Controller interface {
	BaseController
	MountRoutes(group *gin.RouterGroup)
}

type ResponseSender interface {
	Debug() bool
	Send(ctx *gin.Context) SendResponse
}

type BaseMiddleware interface {
	ResponseSender
	Debug() bool
}

type RootMiddleware interface {
	BaseMiddleware
	Attach(engine *gin.Engine)
	Handler(ctx *gin.Context)
}

type BaseRouter interface {
	GetEngine() *gin.Engine
	RegisterValidationParsers(tagNameFunc validator.TagNameFunc)
	LoadRootMiddlewares(middlewares []RootMiddleware)
	Start(ip string, port uint16)
}

type Router interface {
	BaseRouter
	LoadControllers(controllers []Controller)
}

type BaseModule[T any] interface {
	GetInstance() *T
	RootMiddlewares() []RootMiddleware
	// AuthenticationProvider() AuthenticationProvider
	// AuthorizationProvider() AuthorizationProvider
}

type Module[T any] interface {
	BaseModule[T]
	Controllers() []Controller
}
