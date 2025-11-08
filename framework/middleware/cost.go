package middleware

import (
	"log"
	"time"

	"github.com/RZXBxie/web_server/framework"
)

func Cost() framework.ControllerHandler {
	return func(c *framework.Context) error {
		start := time.Now()

		c.Next()

		end := time.Now()
		cost := end.Sub(start)
		log.Printf("api uri: %v, cost: %v", c.GetRequest().RequestURI, cost)

		return nil
	}
}
