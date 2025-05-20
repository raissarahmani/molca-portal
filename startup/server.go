package startup

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/molca-id/portal-app-api/arch/mongo"
	"github.com/molca-id/portal-app-api/arch/network"
	"github.com/molca-id/portal-app-api/arch/redis"
	"github.com/molca-id/portal-app-api/config"
)

type Shutdown = func()

func Server() {
	env := config.NewEnv(".env", true)
	router, _, shutdown := create(env)
	defer shutdown()

	router.Start(env.ServerHost, env.ServerPort)
}

func create(env *config.Env) (network.Router, Module, Shutdown) {
	ctx := context.Background()

	dbConfig := mongo.DbConfig{
		User:        env.DBUser,
		Pwd:         env.DBUserPwd,
		Host:        env.DBHost,
		Port:        env.DBPort,
		Name:        env.DBName,
		MinPoolSize: env.DBMinPoolSize,
		MaxPoolSize: env.DBMaxPoolSize,
		Timeout:     time.Duration(env.DBQueryTimeout) * time.Second * 10,
	}

	db := mongo.NewDatabase(ctx, dbConfig)
	db.Connect()

	if env.GoMode != gin.TestMode {
		EnsureIndexes(db)
	}

	redisConfig := redis.Config{
		Host: env.RedisHost,
		Port: env.RedisPort,
		Pwd:  env.RedisPwd,
		DB:   env.RedisDB,
	}

	store := redis.NewStore(ctx, &redisConfig)
	store.Connect()

	module := NewModule(ctx, env, db, store)

	router := network.NewRouter(env.GoMode)
	router.RegisterValidationParsers(network.CustomTagNameFunc())
	router.LoadRootMiddlewares(module.RootMiddlewares())
	router.LoadControllers(module.Controllers())

	shutdown := func() {
		db.Disconnect()
		store.Disconnect()
	}

	return router, module, shutdown
}
