package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MONGOConfig struct {
	User     string
	Password string
	Host     string
	Database string
}

type MongoDatabase struct {
	name string
}

type Databases struct{}

func MongoConfig() *MONGOConfig {
	return &MONGOConfig{
		User:     "wassi-algo",
		Password: "poney",
		Host:     "138.195.138.30",
	}
}

func (db *MongoDatabase) formatURI(conf *MONGOConfig) (string, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:27017/%s", conf.User, conf.Password, conf.Host, db.name)

	return uri, nil
}

func (db *MongoDatabase) make() (*mongo.Client, error) {
	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	uri, err := db.formatURI(MongoConfig())
	if err != nil {
		return nil, err
	}
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverApi)

	client, err := mongo.Connect(context.TODO(), opts)

	return client, err
}

func MongoConnect(db string) (*mongo.Client, error) {

	instance := &MongoDatabase{
		name: db,
	}

	return instance.make()
}
