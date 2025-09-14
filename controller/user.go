package controller

import (
	"github.com/RZXBxie/web_server/framework"
)

func UserLoginController(c *framework.Context) error {
	c.SetOkStatus().Json("ok, UserLoginController")
	return nil
}
