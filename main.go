package main

import (
	"wassimovie-algo/internal/http"
	"wassimovie-algo/internal/indexing"
)

func main() {
	go http.InitServer()
	loader := &indexing.ModelLoader{}
	loader.UserIndexGeneration()

}
