package main

import (
	"context"
	"fmt"
	"wassimovie-algo/internal/database"

	// "math/rand"

	"go.mongodb.org/mongo-driver/bson"
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

func main() {
	// t := annoyindex.NewAnnoyIndexAngular(f)
	// for i := 0; i < 1000000; i++ {
	// 	item := make([]float32, 0, f)
	// 	for x := 0; x < f; x++ {
	// 		item = append(item, rand.Float32())
	// 	}
	// 	t.AddItem(i, item)
	// }
	// t.Build(40)
	// t.Save("test.ann")

	// annoyindex.DeleteAnnoyIndexAngular(t)

	FromTitle("The Godfather")
}

func FromTitle(title string) {

	client, err := database.MongoConnect("wassidb")

	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("wassidb").Collection("movies")
	filter := bson.M{}

	var results []bson.M
	cursor, _ := coll.Find(context.TODO(), filter)
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(results))

}
