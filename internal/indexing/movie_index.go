package indexing

import (
	"annoyindex"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovieVector = [406]float32

type MovieIndex struct {
	Movie_index   annoyindex.AnnoyIndexAngular
	Tmdb_map      map[int]string
	Movie_vectors map[int]MovieVector
}

func CreateIndex(movies map[string]bson.M) *MovieIndex {

	var movie_vecs map[int]MovieVector
	var tmdb_map map[int]string
	tmdb_map = make(map[int]string)
	movie_vecs = make(map[int]MovieVector)
	i := 0

	for _, film := range movies {
		imdb, ok := film["imdb_id"].(string)
		if !ok {
			imdb = ""
		}
		tmdb_map[i] = imdb
		movie_vecs[i] = *BuildMovieVector(film)
	}

	idx := &MovieIndex{
		Movie_index:   nil,
		Tmdb_map:      tmdb_map,
		Movie_vectors: movie_vecs,
	}

	idx.rebuildIndex()
	return idx

}

func (index *MovieIndex) rebuildIndex() error {

	f := 406
	index.Movie_index = annoyindex.NewAnnoyIndexAngular(f)
	for i, v := range index.Movie_vectors {
		index.Movie_index.AddItem(i, v[:])
	}

	index.Movie_index.Build(30)
	index.Movie_index.Save("user_index.bin")

	return nil
}

func (index *MovieIndex) UserANN(user_id int, count int) ([]int, error) {

	var result []int
	index.Movie_index.GetNnsByItem(user_id, count, -1, &result)

	return result, nil

}

func (index *MovieIndex) VecANN(vec []float32, count int) ([]int, error) {
	var result []int
	index.Movie_index.GetNnsByVector(vec, count, -1, &result)
	index.Movie_index.GetNItems()

	return result, nil
}

func BuildMovieVector(movie bson.M) *MovieVector {
	var movie_vec [406]float32
	popularity, ok := movie["popularity"].(float64)
	if !ok {
		popularity = float64(0)
	}
	runtime, ok := movie["runtime"].(int32)
	if !ok {
		runtime = int32(0)
	}
	vote_average, ok := movie["vote_average"].(float64)
	if !ok {
		vote_average = float64(0)
	}
	movie_vec[0] = float32(popularity / 100)
	movie_vec[1] = float32(runtime) / 95
	movie_vec[2] = float32(vote_average / 10)
	genres, ok := movie["genres"].(bson.A)
	if !ok {
		genres = primitive.A{}
	}
	for _, s := range genres {
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
	desc_vector, ok := movie["description_vector"].(bson.A)
	if !ok {
		desc_vector = bson.A{}
	}
	for i, _ := range desc_vector {
		elem, ok := desc_vector[i].(float64)
		if !ok {
			elem = float64(0)
		}
		movie_vec[i+22] = float32(elem)

	}

	return &movie_vec

}
