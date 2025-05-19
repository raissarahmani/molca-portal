package startup

import (
	"context"

	"github.com/molca-id/portal-app-api/arch/mongo"
	"github.com/molca-id/portal-app-api/arch/network"
	"github.com/molca-id/portal-app-api/config"
)

type Module network.Module[module]

type module struct {
	Context context.Context
	Env     *config.Env
	DB      mongo.Database
}

func (m *module) GetInstance() *module {
	return m
}

func (m *module) Controllers() []network.Controller {
	return []network.Controller{}
}

func (m *module) RootMiddlewares() []network.RootMiddleware {
	return []network.RootMiddleware{}
}

func NewModule(ctx context.Context, env *config.Env, db mongo.Database) Module {
	return &module{
		Context: ctx,
		Env:     env,
		DB:      db,
	}
}
