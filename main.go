package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MovieDescription struct {
	Budget            int64   `bson:"budget"`
	Original_language string  `bson:"original_language"`
	Description       string  `bson:"overview"`
	Popularity        float64 `bson:"popularity"`
	Release_date      string  `bson:"release_date"`
	Revenue           int32   `bson:"revenue"`
	Runtime           int32   `bson:"runtime"`
	Title             string  `bson:"title"`
	Vote_average      float32 `bson:"vote_average"`
	Vote_count        int32   `bson:"vote_count"`
}

type MovieVector = [521]float32

const db_uri = "mongodb://wassi-algo:poney@138.195.138.30:27017/wassidb?authMechanism=SCRAM-SHA-256"

func main() {
}

func FromTitle(title string) *MovieDescription {

	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(db_uri).SetServerAPIOptions(serverApi)

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("wassidb").Collection("movies")
	filter := bson.M{"title": title}

	var movie MovieDescription
	coll.FindOne(context.TODO(), filter).Decode(&movie)

	return &movie
}
