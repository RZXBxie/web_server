package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/RZXBxie/web_server/framework/gin"
)

func Timeout(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		finishChan := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		durationCtx, cancel := context.WithTimeout(c.Copy(), d)
		defer cancel() // 确保释放资源

		go func() {
			defer func() {
				if err := recover(); err != nil {
					panicChan <- err
				}
			}()

			c.Next()
			finishChan <- struct{}{}
		}()

		select {
		case <-finishChan:
			fmt.Println("finish: request completed successfully")
		case p := <-panicChan:
			c.Abort()
			c.ISetStatus(http.StatusInternalServerError).IJson(gin.H{"error": "internal server error", "detail": fmt.Sprint(p)})
		case <-durationCtx.Done():
			c.Abort()
			c.ISetStatus(http.StatusGatewayTimeout).IJson(gin.H{"error": "request timeout"})
		}
	}
}
