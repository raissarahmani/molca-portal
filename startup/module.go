package startup

import (
	"context"

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
	Context         context.Context
	Env             *config.Env
	DB              mongo.Database
	Store           redis.Store
	ProjectService  project.Service
	ProjectsService projects.Service
}

func (m *module) GetInstance() *module {
	return m
}

func (m *module) Controllers() []network.Controller {
	return []network.Controller{
		project.NewController(project.NewService(m.DB)),
		projects.NewController(projects.NewService(m.DB)),
	}
}

func (m *module) RootMiddlewares() []network.RootMiddleware {
	return []network.RootMiddleware{
		coreMW.NewErrorCatcher(),
		coreMW.NewNotFound(),
	}
}

func NewModule(ctx context.Context, env *config.Env, db mongo.Database, store redis.Store) Module {
	projectService := project.NewService(db)
	projectsService := projects.NewService(db)
	return &module{
		Context:         ctx,
		Env:             env,
		DB:              db,
		Store:           store,
		ProjectService:  projectService,
		ProjectsService: projectsService,
	}
}
