package controller

import (
	"net/http"

	"github.com/RZXBxie/web_server/framework"
)

func UserLoginController(core *framework.Context) error {
	core.Json(http.StatusOK, "ok, UserLoginController")
	return nil
}
