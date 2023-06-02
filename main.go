package main

import (
	"context"
	"fmt"
	"time"
	"wassimovie-algo/internal/http"
	"wassimovie-algo/internal/indexing"

	"github.com/labstack/echo/v4"
)

func main() {

	stop := make(chan bool, 1)
	stopped := make(chan bool, 1)
	stopped <- true

	go LoadServer(stop, stopped)

	go func() {
		for {
			select {
			case <-time.After(10 * time.Minute):
				LoadServer(stop, stopped)
			}
		}
	}()

	select {}
}

func Test(stop chan bool, stopped chan bool) {
	fmt.Println("test")

	stop <- true
	<-stopped

	fmt.Println("test2")

	<-stop
	stopped <- true
}

func LoadServer(stop chan bool, stopped chan bool) {

	l := &indexing.ModelLoader{}

	model := l.LoadModel()

	var server *echo.Echo

	fmt.Println("Initialized new dataset")
	stop <- true
	<-stopped

	go func() { server = http.InitServer(model.Handle) }()
	fmt.Println("HTTP server initialized")

	<-stop

	fmt.Println("Switching to new dataset")
	server.Shutdown(context.TODO())

	stopped <- true

}
