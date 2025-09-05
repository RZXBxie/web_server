package main

import (
	"net/http"
	
	"github.com/RZXBxie/web_server/framework"
)

func main() {
	server := &http.Server{
		Handler: framework.NewCore(),
		Addr:    ":8080",
	}
	
	server.ListenAndServe()
	
}
