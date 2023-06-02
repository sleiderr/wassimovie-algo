package main

import (
	"wassimovie-algo/internal/http"
	"wassimovie-algo/internal/indexing"
)

func main() {
	l := &indexing.ModelLoader{}

	model := l.LoadModel()

	http.InitServer(model.Handle)

}
