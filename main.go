package main

import (
	"net/http"

	"github.com/RZXBxie/web_server/framework"
	"github.com/RZXBxie/web_server/middleware"
)

func main() {
	core := framework.NewCore()
	core.Use(middleware.Recovery())
	core.Use(middleware.Cost())
	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8080",
	}

	server.ListenAndServe()

}
