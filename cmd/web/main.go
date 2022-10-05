package main

import (
	"fmt"
	"github.com/rezaDastrs/simple-websocket/internal/handler"
	"net/http"
)

func main()  {

	//get routes
	mux := routes()
	fmt.Println("Starting channel listener")
	go handler.ListenToWsChannel()

	fmt.Println("Starting web server on port :8080")

	_ = http.ListenAndServe(":8080",mux)
}
