package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
)

type Genre struct {
	Id int64 `bson:"id"`
	Name string `bson:"name"`
}

type MovieDescription struct {
	Budget            int64   `bson:"budget"`
	Original_language string  `bson:"original_language"`
	Description       string  `bson:"overview"`
	Popularity        float64 `bson:"popularity"`
	Release_date      string  `bson:"release_date"`
	Revenue           int32   `bson:"revenue"`
	Runtime           int32   `bson:"runtime"`
	Title             string  `bson:"title"`
	Vote_average      float64 `bson:"vote_average"`
	Vote_count        int32   `bson:"vote_count"`
	//Description_vector [384]float64 `bson:"description_vector"`
	Genres			  bson.A `bson:"bson:genre"`
}


type MovieVector = [406]float64

const db_uri = "mongodb://wassi-algo:poney@138.195.138.30:27017/wassidb?authMechanism=SCRAM-SHA-256"

func main() {
	fmt.Println(Cosine(*BuildMovieVector("Tarzan"),*BuildMovieVector("Cars")))
}

func BuildMovieVector(title string) *MovieVector {

	movie := FromTitle(title)
	var movie_vec [406]float64
	movie_vec[0] = movie["popularity"].(float64)/100
	movie_vec[1] = float64(movie["runtime"].(int32))/95 
	movie_vec[2] = movie["vote_average"].(float64)/10
	for _, s := range movie["genres"].(bson.A) {
		if s.(bson.M)["id"] == int32(12) { 
			movie_vec[4] = float64(1)
		} else if s.(bson.M)["id"] == int32(14) {
			movie_vec[5] = float64(1)
		} else if s.(bson.M)["id"] == int32(16) {
			movie_vec[6] = float64(1)
		} else if s.(bson.M)["id"] == int32(18) {
			movie_vec[7] = float64(1)
		} else if s.(bson.M)["id"] == int32(27) {
			movie_vec[8] = float64(1)
		} else if s.(bson.M)["id"] == int32(28) {
			movie_vec[9] = float64(1)
		} else if s.(bson.M)["id"] == int32(35) {
			movie_vec[10] = float64(1)
		} else if s.(bson.M)["id"] == int32(36) {
			movie_vec[11] = float64(1)
		} else if s.(bson.M)["id"] == int32(37) {
			movie_vec[12] = float64(1)
		} else if s.(bson.M)["id"] == int32(53) {
			movie_vec[13] = float64(1)
		} else if s.(bson.M)["id"] == int32(80) {
			movie_vec[14] = float64(1)
		} else if s.(bson.M)["id"] == int32(9648) {
			movie_vec[15] = float64(1)
			} else if s.(bson.M)["id"] == int32(10402) {
				movie_vec[15] = float64(1)
			} else if s.(bson.M)["id"] == int32(10749) {
				movie_vec[16] = float64(1)
			} else if s.(bson.M)["id"] == int32(10752) {
				movie_vec[17] = float64(1)
			} else if s.(bson.M)["id"] == int32(10770) {
				movie_vec[18] = float64(1)
			} else if s.(bson.M)["id"] == int32(878) {
				movie_vec[19] = float64(1)
			} else if s.(bson.M)["id"] == int32(10751) {
				movie_vec[20] = float64(1)
			} else if s.(bson.M)["id"] == int32(99) {
				movie_vec[21] = float64(1)
			}
	}
	for i, _ := range movie["description_vector"].(bson.A){
			movie_vec[i+22] = movie["description_vector"].(bson.A)[i].(float64)

		}


	return &movie_vec

}

func DotProduct(film1 [406]float64, film2 [406]float64) float64 {
	var s float64
	for i, _ := range film1 {
		s += film1[i] * film2[i]

	}
	return s
}

func Norm(film1 [406]float64) float64 {
	var res float64
	for i, _ := range film1 {
		res += math.Pow(film1[i],2)

	}
	return math.Pow(res,0.5)
}

func Cosine(film1 [406]float64, film2 [406]float64) float64 {
	return DotProduct(film1,film2) / (Norm(film1) * Norm(film2))
}





func FromTitle(title string) bson.M {

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

	var movie bson.M
	coll.FindOne(context.TODO(), filter).Decode(&movie)

	return movie
}

