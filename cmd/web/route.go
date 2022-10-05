package main

import (
	"github.com/bmizerany/pat"
	"github.com/rezaDastrs/simple-websocket/internal/handler"
	"net/http"
)

func routes() http.Handler  {
	//using pat package for routing
	mux := pat.New()

	mux.Get("/",http.HandlerFunc(handler.Home))
	mux.Get("/ws",http.HandlerFunc(handler.WsEndpint))

	//reconnection websocket with javascript package
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Get("/static/",http.StripPrefix("/static",fileServer))


	return mux
}
