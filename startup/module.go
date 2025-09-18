package startup

import (
	"context"

	apiAnalytics "github.com/molca-id/portal-app-api/api/analytics"
	"github.com/molca-id/portal-app-api/api/auth"
	authMW "github.com/molca-id/portal-app-api/api/auth/middleware"
	cover "github.com/molca-id/portal-app-api/api/cover"
	project "github.com/molca-id/portal-app-api/api/project"
	projects "github.com/molca-id/portal-app-api/api/projects"
	coreAnalytics "github.com/molca-id/portal-app-api/arch/analytics"
	coreMW "github.com/molca-id/portal-app-api/arch/middleware"
	"github.com/molca-id/portal-app-api/arch/mongo"
	"github.com/molca-id/portal-app-api/arch/network"
	"github.com/molca-id/portal-app-api/arch/redis"
	"github.com/molca-id/portal-app-api/config"
)

type Module network.Module[module]

type module struct {
	Context              context.Context
	Env                  *config.Env
	DB                   mongo.Database
	Store                redis.Store
	AuthService          auth.Service
	CoreAnalyticsService *coreAnalytics.Service
	ApiAnalyticsService  apiAnalytics.Service
}

func (m *module) GetInstance() *module {
	return m
}

func (m *module) Controllers() []network.Controller {
	return []network.Controller{
		project.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), project.NewService(m.DB, m.Store)),
		projects.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), projects.NewService(m.DB)),
		cover.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), cover.NewService(m.Env)),
		apiAnalytics.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), m.ApiAnalyticsService),
	}
}

func (m *module) RootMiddlewares() []network.RootMiddleware {
	return []network.RootMiddleware{
		coreMW.NewCorsMiddleware(),
		coreMW.NewErrorCatcher(),
		coreMW.NewNotFound(),
	}
}

func (m *module) AuthenticationProvider() network.AuthenticationProvider {
	return authMW.NewAuthenticationProvider(m.Env, m.AuthService)
}

func (m *module) AuthorizationProvider() network.AuthorizationProvider {
	return authMW.NewAuthorizationProvider()
}

func NewModule(ctx context.Context, env *config.Env, db mongo.Database, store redis.Store) Module {
	authService := auth.NewService(env)
	coreAnalyticsService, err := coreAnalytics.NewService(ctx, env.GAPropertyID, env.GoogleApplicationCredentials)
	if err != nil {
		panic(err)
	}

	apiAnalyticsService := apiAnalytics.NewService(coreAnalyticsService, db)

	return &module{
		Context:              ctx,
		Env:                  env,
		DB:                   db,
		Store:                store,
		AuthService:          authService,
		CoreAnalyticsService: coreAnalyticsService,
		ApiAnalyticsService:  apiAnalyticsService,
	}
}
