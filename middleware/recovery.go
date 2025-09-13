package middleware

import (
	"net/http"
	
	"github.com/RZXBxie/web_server/framework"
)

func Recovery() framework.ControllerHandler {
	return func(c *framework.Context) error {
		defer func() {
			if err := recover(); err != nil {
				c.Json(http.StatusInternalServerError, "internal server error")
				
			}
		}()
		c.Next()
		
		return nil
	}
}
