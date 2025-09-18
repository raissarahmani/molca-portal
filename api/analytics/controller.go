package analytics

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/molca-id/portal-app-api/arch/network"
)

type controller struct {
	network.BaseController
	service Service
}

func NewController(
	authMFunc network.AuthenticationProvider,
	authorizeMFunc network.AuthorizationProvider,
	service Service,
) network.Controller {
	return &controller{
		BaseController: network.NewBaseController("/analytics", authMFunc, authorizeMFunc),
		service:        service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.Use(c.Authentication())
	group.GET("/user-visits", c.getVisitsTimeSeriesHandler)
	group.GET("/view-project/title", c.getRankingByTitleHandler)
	group.GET("/view-project/type", c.getProjectVisitByTypeHandler)
	group.GET("/project/type", c.getProjectTypeHandler)
}

func (c *controller) getVisitsTimeSeriesHandler(ctx *gin.Context) {
	rangeType := ctx.DefaultQuery("range", "daily")
	debug, _ := strconv.ParseBool(ctx.DefaultQuery("debug", "false"))

	series, err := c.service.GetVisitsTimeSeries(rangeType, debug)
	if err != nil {
		c.Send(ctx).InternalServerError("Error fetching visits by range", err)
		return
	}
	c.Send(ctx).SuccessDataResponse("Visits by range fetched successfully", series)
}

func (c *controller) getRankingByTitleHandler(ctx *gin.Context) {
	limit, order := c.parseRankingParams(ctx)
	rangeType := ctx.DefaultQuery("range", "daily")
	debug, _ := strconv.ParseBool(ctx.DefaultQuery("debug", "false"))

	ranking, err := c.service.GetRankingByTitle(order, limit, rangeType, debug)
	if err != nil {
		c.Send(ctx).InternalServerError("Error fetching ranking by title", err)
		return
	}
	c.Send(ctx).SuccessDataResponse("Ranking by title fetched successfully", ranking)
}

func (c *controller) getProjectVisitByTypeHandler(ctx *gin.Context) {
	limit, order := c.parseRankingParams(ctx)
	rangeType := ctx.DefaultQuery("range", "daily")
	debug, _ := strconv.ParseBool(ctx.DefaultQuery("debug", "false"))

	ranking, err := c.service.GetProjectVisitByType(order, limit, rangeType, debug)
	if err != nil {
		c.Send(ctx).InternalServerError("Error fetching project visit by type", err)
		return
	}
	c.Send(ctx).SuccessDataResponse("Project visit by type fetched successfully", ranking)
}

func (c *controller) getProjectTypeHandler(ctx *gin.Context) {
	debug, _ := strconv.ParseBool(ctx.DefaultQuery("debug", "false"))
	data, err := c.service.GetProjectType(debug)
	if err != nil {
		c.Send(ctx).InternalServerError("Error fetching project type data", err)
		return
	}
	c.Send(ctx).SuccessDataResponse("Project type fetched successfully", data)
}

func (c *controller) parseRankingParams(ctx *gin.Context) (int64, string) {
	limitStr := ctx.DefaultQuery("limit", "6")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 6
	}
	order := ctx.DefaultQuery("order", "desc")
	return limit, order
}
