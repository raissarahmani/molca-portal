package project

import (
	"github.com/gin-gonic/gin"
	"github.com/molca-id/portal-app-api/api/project/dto"
	coredto "github.com/molca-id/portal-app-api/arch/dto"
	"github.com/molca-id/portal-app-api/arch/network"
	"github.com/molca-id/portal-app-api/common"
	"github.com/molca-id/portal-app-api/utils"
)

type controller struct {
	network.BaseController
	common.ContextPayload
	service Service
}

func NewController(
	authMFunc network.AuthenticationProvider,
	authorizeMFunc network.AuthorizationProvider,
	service Service,
) network.Controller {
	return &controller{
		BaseController: network.NewBaseController("/project", authMFunc, authorizeMFunc),
		ContextPayload: common.NewContextPayload(),
		service:        service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.Use(c.Authentication())
	group.POST("/", c.createProjectHandler)
	group.GET("/:id", c.getProjectByIdHandler)
	group.PUT("/:id", c.updateProjectHandler)
	group.DELETE("/:id", c.deleteProjectHandler)
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

func (c *controller) updateProjectHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, coredto.EmptyMongoId())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}
	body, err := network.ReqBody(ctx, &dto.UpdateProject{})
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	err = c.service.UpdateProject(mongoId.ID, body)
	if err != nil {
		c.Send(ctx).InternalServerError("Something went wrong", err)
		return
	}

	c.Send(ctx).SuccessDataResponse("Project updated successfully", nil)
}

func (c *controller) getProjectByIdHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, coredto.EmptyMongoId())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	project, err := c.service.GetProjectDtoCacheById(mongoId.ID)
	if err == nil {
		c.Send(ctx).SuccessDataResponse("Project fetched from cache successfully", project)
		return
	}

	project, err = c.service.GetProjectById(mongoId.ID)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("Project fetched successfully", project)
	c.service.SetProjectDtoCacheById(project)
}

func (c *controller) getProjectBySlugHandler(ctx *gin.Context) {
	slug, err := network.ReqParams(ctx, coredto.EmptySlug())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	project, err := c.service.GetProjectDtoCacheBySlug(slug.Slug)
	if err == nil {
		c.Send(ctx).SuccessDataResponse("Project fetched from cache successfully", project)
		return
	}

	project, err = c.service.GetProjectBySlug(slug.Slug)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("Project fetched successfully", project)
	c.service.SetProjectDtoCacheBySlug(project)
}

func (c *controller) deleteProjectHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, coredto.EmptyMongoId())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	err = c.service.DeleteProject(mongoId.ID)
	if err != nil {
		c.Send(ctx).InternalServerError("Something went wrong", err)
		return
	}

	c.Send(ctx).SuccessDataResponse("Project deleted successfully", nil)
}
