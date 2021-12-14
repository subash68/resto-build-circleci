package database

import (
	"context"
	"sync"

	"github.com/subash68/ate/ate_onboard_service/configuration"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientInstance *mongo.Client
var clientInstanceError error
var mongoOnce sync.Once

type Config struct {
	DatabaseInstance string
	DatabaseName     string
	Collection       string
}

func LoadDatabase() Config {
	var cfg Config

	configuration.LoadConfig()

	// cfg.DatabaseInstance = configuration.DbConfig().DatabaseInstance
	// cfg.DatabaseName = configuration.DbConfig().DatabaseName
	// cfg.Collection = configuration.DbConfig().Collection

	return cfg
}

//GetMongoClient - Return mongodb connection to work with
func GetMongoClient() (*mongo.Client, error) {
	//Perform connection creation operation only once.

	mongoOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(LoadDatabase().DatabaseInstance)
		// Connect to MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
		}
		// Check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
		}
		clientInstance = client
	})
	return clientInstance, clientInstanceError
}
