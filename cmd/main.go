package main

import (
	"context"
	"time"

	"github.com/molca-id/portal-app-api/arch/mongo"
)

func main() {
	dbConfig := mongo.DbConfig{
		User:        "root",
		Pwd:         "root",
		Host:        "localhost",
		Port:        27017,
		Database:    "portal-app",
		MinPoolSize: 10,
		MaxPoolSize: 100,
		Timeout:     10 * time.Second,
	}

	db := mongo.NewDatabase(context.Background(), dbConfig)
	db.Connect()
	defer db.Disconnect()
}
