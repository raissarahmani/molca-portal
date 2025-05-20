package network

// import "github.com/gin-gonic/gin"

type baseController struct {
	ResponseSender
	basePath string
	// authProvider AuthencationProvider
	// authorizeProvider AuthorizationProvider
}

func NewBaseController(
	basePath string,
	// authProvider AuthencationProvider,
	// authorizeProvider AuthorizationProvider,
) BaseController {
	return &baseController{
		ResponseSender: NewResponseSender(),
		basePath:       basePath,
		// authProvider: authProvider,
		// authorizeProvider: authorizeProvider,
	}
}

func (c *baseController) Path() string {
	return c.basePath
}
