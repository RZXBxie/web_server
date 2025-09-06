package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RZXBxie/web_server/framework"
)

// FooControllerHandler 业务函数，测试 context.WithTimeout
// 在浏览器中访问127.0.0.1:8080/foo，过一秒就会得到"time out"的结果
func FooControllerHandler(c *framework.Context) error {
	finishChan := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)

	durationCtx, cancel := context.WithTimeout(c.BaseContext(), 1*time.Second)
	defer cancel()

	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		// Do real action
		time.Sleep(10 * time.Second)
		c.Json(http.StatusOK, "ok")
		finishChan <- struct{}{}
	}()

	select {
	case p := <-panicChan:
		c.WriteMux().Lock()
		defer c.WriteMux().Unlock()
		log.Println(p)
		c.Json(http.StatusInternalServerError, "panic")
	case <-finishChan:
		fmt.Println("finish")
	case <-durationCtx.Done():
		c.WriteMux().Lock()
		defer c.WriteMux().Unlock()
		c.Json(http.StatusInternalServerError, "time out")
		c.SetIsTimeOut()
	}

	return nil
}
