package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"wassimovie-algo/internal/database"
	"wassimovie-algo/internal/http"

	"go.mongodb.org/mongo-driver/bson"
)

type MovieDescription struct {
	Budget            int64   `bson:"budget"`
	Original_language string  `bson:"original_language"`
	Description       string  `bson:"overview"`
	Popularity        float32 `bson:"popularity"`
	Release_date      string  `bson:"release_date"`
	Revenue           int32   `bson:"revenue"`
	Runtime           int32   `bson:"runtime"`
	Title             string  `bson:"title"`
	Vote_average      float32 `bson:"vote_average"`
	Vote_count        int32   `bson:"vote_count"`
}

type MovieVector = [406]float32

const MAX_CONCURRENT_JOBS = 500

type User struct {
	UserId string `bson:"userId"`
}

func main() {
	go http.InitServer()
	fmt.Println("frost")
	//UserVectorGeneration()
	//fmt.Println(FromTitle("Nemo"))
	fmt.Println(ComputeUserVector("2",RetrieveRatingsDatabase(),RetrieveMoviesDatabase()))
}

func RetrieveMoviesDatabase() map[string]bson.M {
	client, err := database.MongoConnect("wassidb")

	if err != nil {
		panic(err)
	}

	coll := client.Database("wassidb").Collection("movies")
	filter := bson.M{}

	cursor, err := coll.Find(context.TODO(), filter)
	db_movies := make(map[string]bson.M)

	for cursor.Next(context.TODO()) {
		var temp_movie bson.M
		cursor.Decode(&temp_movie)
		db_movies[fmt.Sprintf("%v", temp_movie["imdb_id"])] = temp_movie
	}

	return db_movies

}

func RetrieveRatingsDatabase() map[string][]bson.M {
	client, err := database.MongoConnect("wassidb")

	if err != nil {
		panic(err)
	}

	coll := client.Database("wassidb").Collection("ratings")
	filter := bson.M{}

	cursor, err := coll.Find(context.TODO(), filter)
	db_ratings := make(map[string][]bson.M)

	for cursor.Next(context.TODO()) {
		var temp_rating bson.M
		cursor.Decode(&temp_rating)
		db_ratings[fmt.Sprintf("%v", temp_rating["userId"])] = append(db_ratings[fmt.Sprintf("%v", temp_rating["userId"])], temp_rating)
	}

	return db_ratings

}

func UserVectorGeneration() *map[string][406]float32 {

	var results map[string][406]float32
	results = make(map[string][406]float32)

	client, err := database.MongoConnect("wassidb")

	waitChan := make(chan struct{}, MAX_CONCURRENT_JOBS)

	if err != nil {
		panic(err)
	}

	coll := client.Database("wassidb").Collection("users")
	filter := bson.M{}

	cursor, err := coll.Find(context.TODO(), filter)

	for cursor.Next(context.TODO()) {
		waitChan <- struct{}{}
		var user bson.M
		if err := cursor.Decode(&user); err != nil {
			log.Fatal(err)
		}
		if isOutdated, ok := user["outdated"]; ok && !isOutdated.(bool) {
			// just store user vector in memory
		} else {
			// goroutine for vector calculation
			// store in memory, update in db through another goroutine
			go func() {
				// results[user["userId"].(string)] = ComputeUserVector(user["userId"].(string))
			//	ComputeUserVector(user["userId"].(string))
				<-waitChan
			}()
		}
	}

	return &results

}

func BuildMovieVector(movie bson.M) *MovieVector {
	var movie_vec [406]float32
	movie_vec[0] = movie["popularity"].(float32) / 100
	movie_vec[1] = float32(movie["runtime"].(int32)) / 95
	movie_vec[2] = movie["vote_average"].(float32) / 10
	for _, s := range movie["genres"].(bson.A) {
		if s.(bson.M)["id"] == int32(12) {
			movie_vec[4] = float32(1)
		} else if s.(bson.M)["id"] == int32(14) {
			movie_vec[5] = float32(1)
		} else if s.(bson.M)["id"] == int32(16) {
			movie_vec[6] = float32(1)
		} else if s.(bson.M)["id"] == int32(18) {
			movie_vec[7] = float32(1)
		} else if s.(bson.M)["id"] == int32(27) {
			movie_vec[8] = float32(1)
		} else if s.(bson.M)["id"] == int32(28) {
			movie_vec[9] = float32(1)
		} else if s.(bson.M)["id"] == int32(35) {
			movie_vec[10] = float32(1)
		} else if s.(bson.M)["id"] == int32(36) {
			movie_vec[11] = float32(1)
		} else if s.(bson.M)["id"] == int32(37) {
			movie_vec[12] = float32(1)
		} else if s.(bson.M)["id"] == int32(53) {
			movie_vec[13] = float32(1)
		} else if s.(bson.M)["id"] == int32(80) {
			movie_vec[14] = float32(1)
		} else if s.(bson.M)["id"] == int32(9648) {
			movie_vec[15] = float32(1)
		} else if s.(bson.M)["id"] == int32(10402) {
			movie_vec[15] = float32(1)
		} else if s.(bson.M)["id"] == int32(10749) {
			movie_vec[16] = float32(1)
		} else if s.(bson.M)["id"] == int32(10752) {
			movie_vec[17] = float32(1)
		} else if s.(bson.M)["id"] == int32(10770) {
			movie_vec[18] = float32(1)
		} else if s.(bson.M)["id"] == int32(878) {
			movie_vec[19] = float32(1)
		} else if s.(bson.M)["id"] == int32(10751) {
			movie_vec[20] = float32(1)
		} else if s.(bson.M)["id"] == int32(99) {
			movie_vec[21] = float32(1)
		}
	}
	for i, _ := range movie["description_vector"].(bson.A) {
		movie_vec[i+22] = movie["description_vector"].(bson.A)[i].(float32)

	}

	return &movie_vec

}

func DotProduct(film1 [406]float32, film2 [406]float32) float32 {
	var s float32
	for i, _ := range film1 {
		s += film1[i] * film2[i]

	}
	return s
}

func Norm(film1 [406]float32) float32 {
	var res float64
	for i, _ := range film1 {
		res += math.Pow(float64(film1[i]), 2)

	}
	return float32(math.Pow(float64(res), 0.5))
}

func Cosine(film1 [406]float32, film2 [406]float32) float32 {
	return DotProduct(film1, film2) / (Norm(film1) * Norm(film2))
}

func ComputeUserVector(id string, db_ratings map[string][]bson.M, db_movies map[string]bson.M) [406]float64 {

	var count float64
	var user_vector [406]float64

	for _, s := range db_ratings[id] {
		coeff := float64(s["rating"].(int32))
		count += math.Abs(coeff)
		fmt.Println(db_movies[s["movieId"].(string)])
		film_id, ok := s["movieId"].(string)
		if (!ok) {
			continue
		}
		temp_movie_vec := *BuildMovieVector(db_movies[film_id])
		for i, _ := range temp_movie_vec {
			user_vector[i] += temp_movie_vec[i]*coeff

		}

	}
	for i, _ := range user_vector {
		user_vector[i] /= count
	}
	return user_vector
}

func FromTitle(title string) bson.M {

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
	filter := bson.M{"title": title}

	var movie bson.M
	coll.FindOne(context.TODO(), filter).Decode(&movie)

	return movie
}