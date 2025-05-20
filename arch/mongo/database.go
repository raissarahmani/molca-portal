package mongo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DbConfig struct {
	User        string
	Pwd         string
	Host        string
	Port        uint16
	Name        string
	MinPoolSize uint16
	MaxPoolSize uint16
	Timeout     time.Duration
}

type database struct {
	*mongo.Database
	context context.Context
	config  DbConfig
}

type Database interface {
	GetInstance() *database
	Connect()
	Disconnect()
}

type Document[T any] interface {
	EnsureIndexes(Database)
	GetValue() *T
	Validate() error
}

func NewDatabase(ctx context.Context, config DbConfig) Database {
	db := database{
		context: ctx,
		config:  config,
	}

	return &db
}

func (d *database) GetInstance() *database {
	return d
}

func (d *database) Connect() {
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		d.config.User, d.config.Pwd, d.config.Host, d.config.Port,
	)

	clientOptions := options.Client().ApplyURI(uri)

	clientOptions.SetMaxPoolSize(uint64(d.config.MaxPoolSize))
	clientOptions.SetMinPoolSize(uint64(d.config.MinPoolSize))
	clientOptions.SetTimeout(d.config.Timeout)

	fmt.Println("Connecting to MongoDB...")

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("Connection failed: ", err)
	}

	err = client.Ping(d.context, nil)
	if err != nil {
		log.Fatal("Ping failed: ", err)
	}

	fmt.Println("Connected to MongoDB")

	d.Database = client.Database(d.config.Name)
}

func (d *database) Disconnect() {
	fmt.Println("Disconnecting from MongoDB...")

	err := d.Client().Disconnect(d.context)
	if err != nil {
		log.Panic("Disconnection failed: ", err)
	}

	fmt.Println("Disconnected from MongoDB")
}

func NewObjectID(id string) (bson.ObjectID, error) {
	i, err := bson.ObjectIDFromHex(id)
	if err != nil {
		err = errors.New(id + " is not a valid ObjectID")
	}

	return i, err
}
