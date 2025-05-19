package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host string
	Port uint16
	Pwd  string
	DB   int
}

type Store interface {
	GetInstance() *store
	Connect()
	Disconnect()
}

type store struct {
	*redis.Client
	context context.Context
}

func NewStore(context context.Context, config *Config) Store {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Pwd,
		DB:       config.DB,
	})

	return &store{
		context: context,
		Client:  client,
	}
}

func (r *store) GetInstance() *store {
	return r
}

func (r *store) Connect() {
	fmt.Println("Connecting to Redis")
	pong, err := r.Ping(r.context).Result()
	if err != nil {
		panic(fmt.Errorf("failed to ping redis: %v", err))
	}
	fmt.Println("Connected to Redis", pong)
}

func (r *store) Disconnect() {
	fmt.Println("Disconnecting from Redis")
	err := r.Close()
	if err != nil {
		panic(fmt.Errorf("failed to close redis: %v", err))
	}
	fmt.Println("Disconnected from Redis")
}
