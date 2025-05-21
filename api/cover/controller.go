package cover

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/molca-id/portal-app-api/api/cover/dto"
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
		BaseController: network.NewBaseController("/cover"),
		service:        service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.POST("/upload", c.saveCoverHandler)
}

func (c *controller) saveCoverHandler(ctx *gin.Context) {
	d, err := network.ReqBodyMultipart(ctx, &dto.SaveImage{})
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	fmt.Println(d)

	cover, err := c.service.SaveCover(d)
	if err != nil {
		c.Send(ctx).InternalServerError(err.Error(), err)
		return
	}

	c.Send(ctx).SuccessDataResponse("Cover saved successfully", cover)
}
