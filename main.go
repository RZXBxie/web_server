package main

import (
	"net/http"

	"github.com/RZXBxie/web_server/framework"
)

func main() {
	core := framework.NewCore()
	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8080",
	}

	server.ListenAndServe()

}
