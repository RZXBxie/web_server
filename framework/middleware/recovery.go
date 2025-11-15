package middleware

import (
	"net/http"

	"github.com/RZXBxie/web_server/framework/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.ISetStatus(http.StatusInternalServerError).IJson("internal server error")

			}
		}()
		c.Next()
	}
}
