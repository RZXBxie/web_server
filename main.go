package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	
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
	
	go func() {
		server.ListenAndServe()
	}()
	
	// 新建一个Signal类型的channel
	quit := make(chan os.Signal)
	// 订阅这三种类型的关闭信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("shutdown server error: %v", err)
	}
	
}
