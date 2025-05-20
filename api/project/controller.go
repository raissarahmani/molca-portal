package project

import (
	"github.com/gin-gonic/gin"
	"github.com/molca-id/portal-app-api/api/project/dto"
	coredto "github.com/molca-id/portal-app-api/arch/dto"
	"github.com/molca-id/portal-app-api/arch/network"
	"github.com/molca-id/portal-app-api/utils"
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
	group.POST("/", c.createProjectHandler)
	group.GET("/:id", c.getProjectByIdHandler)
	group.GET("/slug/:slug", c.getProjectBySlugHandler)
}

func (c *controller) createProjectHandler(ctx *gin.Context) {
	body, err := network.ReqBody(ctx, &dto.CreateProject{})
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	project, err := c.service.SaveProject(body)
	if err != nil {
		c.Send(ctx).InternalServerError("Something went wrong", err)
		return
	}

	data, err := utils.MapTo[dto.InfoProject](project)

	if err != nil {
		c.Send(ctx).InternalServerError("Something went wrong", err)
		return
	}
	c.Send(ctx).SuccessDataResponse("Project created successfully", data)
}

func (c *controller) getProjectByIdHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, coredto.EmptyMongoId())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	project, err := c.service.GetProjectById(mongoId.ID)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	data, err := utils.MapTo[dto.InfoProject](project)
	if err != nil {
		c.Send(ctx).InternalServerError("Something went wrong", err)
		return
	}

	c.Send(ctx).SuccessDataResponse("Project fetched successfully", data)
}

func (c *controller) getProjectBySlugHandler(ctx *gin.Context) {
	slug, err := network.ReqParams(ctx, coredto.EmptySlug())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	project, err := c.service.GetProjectBySlug(slug.Slug)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	data, err := utils.MapTo[dto.InfoProject](project)
	if err != nil {
		c.Send(ctx).InternalServerError("Something went wrong", err)
		return
	}

	c.Send(ctx).SuccessDataResponse("Project fetched successfully", data)
}
