package indexing

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"sync"
	"wassimovie-algo/internal/database"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

type ModelLoader struct{}
type Model struct {
	movies     *MovieIndex
	user_index map[string][406]float32
}

const MAX_CONCURRENT_JOBS = 2000

func (m *Model) Handle(ctx echo.Context) error {

	username := ctx.Param("username")

	user_vec, ok := m.user_index[username]
	if !ok {
		return ctx.String(http.StatusNotFound, "User does not exist")
	}

	unn, _ := m.movies.VecANN(user_vec[:], 100)
	var id_list []string
	id_list = make([]string, 100)

	for i, v := range unn {
		id_list[i] = m.movies.Tmdb_map[v]
	}

	return ctx.JSON(http.StatusOK, id_list)

}

func (m *Model) ReloadModel(l *ModelLoader) {
	user_idx := l.UserIndexGeneration()
	movie_db := l.RetrieveMoviesDatabase()
	film_idx := CreateIndex(movie_db)

	m.movies = film_idx
	m.user_index = user_idx
}

func (l *ModelLoader) LoadModel() *Model {
	usr_idx := l.UserIndexGeneration()
	movie_db := l.RetrieveMoviesDatabase()
	film_idx := CreateIndex(movie_db)

	model := &Model{
		movies:     film_idx,
		user_index: usr_idx,
	}

	return model
}

func (l *ModelLoader) RetrieveMoviesDatabase() map[string]bson.M {
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

func (l *ModelLoader) RetrieveRatingsDatabase() map[string][]bson.M {
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

func (l *ModelLoader) UserIndexGeneration() map[string][406]float32 {

	var results map[string][406]float32
	results = make(map[string][406]float32)
	var mutex = &sync.Mutex{}

	client, err := database.MongoConnect("wassidb")

	ratings := l.RetrieveRatingsDatabase()
	movies := l.RetrieveMoviesDatabase()

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
				temp_value := l.ComputeUserVector(user["userId"].(string), ratings, movies)
				mutex.Lock()
				results[user["userId"].(string)] = temp_value
				mutex.Unlock()
				<-waitChan
			}()
		}
	}

	return results

}

func (l *ModelLoader) ComputeUserVector(id string, db_ratings map[string][]bson.M, db_movies map[string]bson.M) [406]float32 {

	var count float32
	var user_vector [406]float32

	for _, s := range db_ratings[id] {
		coeff := float32(s["rating"].(int32))
		count += float32(math.Abs(float64(coeff)))
		film_id, ok := s["movieId"].(string)
		if !ok {
			continue
		}
		temp_movie_vec := *BuildMovieVector(db_movies[film_id])
		for i, _ := range temp_movie_vec {
			user_vector[i] += float32(temp_movie_vec[i]) * float32(coeff)

		}

	}
	for i, _ := range user_vector {
		user_vector[i] /= count
	}
	return user_vector
}
