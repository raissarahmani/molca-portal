package projects

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/molca-id/portal-app-api/api/projects/dto"
	"github.com/molca-id/portal-app-api/arch/network"
)

type controller struct {
	network.BaseController
	service Service
}

func NewController(
	service Service,
) network.Controller {
	return &controller{
		BaseController: network.NewBaseController("/projects"),
		service:        service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/latest", c.getPaginatedLatestProjectsHandler)
}

func (c *controller) getPaginatedLatestProjectsHandler(ctx *gin.Context) {
	fp, err := network.ReqQuery(ctx, &dto.FilterPaginatedProject{})

	fmt.Println(fp)
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	projects, err := c.service.GetPaginatedLatestProjects(fp)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("Projects fetched successfully", projects)
}
