package utils

import "math"

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
