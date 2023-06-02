package main

import (
	"fmt"
	"time"
	"wassimovie-algo/internal/http"
	"wassimovie-algo/internal/indexing"
)

func main() {

	l := &indexing.ModelLoader{}
	model := l.LoadModel()

	server := http.InitServer(model.Handle)
	go server.Start(":8080")

	go func() {
		for {
			select {
			case <-time.After(7 * time.Minute):
				model.ReloadModel(l)
				fmt.Printf("Reloaded model")
			}
		}
	}()

	select {}
}
