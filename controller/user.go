package controller

import (
	"time"

	"github.com/RZXBxie/web_server/framework/gin"
)

func UserLoginController(c *gin.Context) {
	foo, _ := c.DefaultQueryString("foo", "def")
	time.Sleep(10 * time.Second)
	c.ISetOkStatus().IJson("ok, UserLoginController" + foo)
}
