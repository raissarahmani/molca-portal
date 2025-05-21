package startup

import (
	"context"

	cover "github.com/molca-id/portal-app-api/api/cover"
	project "github.com/molca-id/portal-app-api/api/project"
	projects "github.com/molca-id/portal-app-api/api/projects"
	coreMW "github.com/molca-id/portal-app-api/arch/middleware"
	"github.com/molca-id/portal-app-api/arch/mongo"
	"github.com/molca-id/portal-app-api/arch/network"
	"github.com/molca-id/portal-app-api/arch/redis"
	"github.com/molca-id/portal-app-api/config"
)

type Module network.Module[module]

type module struct {
	Context context.Context
	Env     *config.Env
	DB      mongo.Database
	Store   redis.Store
}

func (m *module) GetInstance() *module {
	return m
}

func (m *module) Controllers() []network.Controller {
	return []network.Controller{
		project.NewController(project.NewService(m.DB)),
		projects.NewController(projects.NewService(m.DB)),
		cover.NewController(cover.NewService(m.Env)),
	}
}

func (m *module) RootMiddlewares() []network.RootMiddleware {
	return []network.RootMiddleware{
		coreMW.NewErrorCatcher(),
		coreMW.NewNotFound(),
	}
}

func NewModule(ctx context.Context, env *config.Env, db mongo.Database, store redis.Store) Module {
	return &module{
		Context: ctx,
		Env:     env,
		DB:      db,
		Store:   store,
	}
}
