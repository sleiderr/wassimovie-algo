package indexing

import "annoyindex"

type MovieVector = [406]float32

type UserIndex struct {
	user_index   annoyindex.AnnoyIndexAngular
	user_vectors *map[int]MovieVector
}

func (index *UserIndex) rebuildIndex() error {

	f := 406
	index.user_index = annoyindex.NewAnnoyIndexAngular(f)
	for i, v := range *index.user_vectors {
		index.user_index.AddItem(i, v[:])
	}

	index.user_index.Build(30)
	index.user_index.Save("user_index.bin")

	return nil
}

func (index *UserIndex) UserANN(user_id int, count int) ([]int, error) {

	var result []int
	index.user_index.GetNnsByItem(user_id, count, -1, &result)

	return result, nil

}

func (index *UserIndex) VecANN(vec MovieVector, count int) ([]int, error) {
	var result []int
	index.user_index.GetNnsByVector(vec, count, -1, &result)
	index.user_index.GetNItems()

	return result, nil
}
