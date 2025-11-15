package controller

import (
	"github.com/RZXBxie/web_server/framework/gin"
)

func SubjectListController(c *gin.Context) {
	c.ISetOkStatus().IJson("ok, SubjectListController")
}

func SubjectDelController(c *gin.Context) {
	c.ISetOkStatus().IJson("ok, SubjectDelController")

}

func SubjectUpdateController(c *gin.Context) {
	c.ISetOkStatus().IJson("ok, SubjectUpdateController")

}

func SubjectGetController(c *gin.Context) {
	c.ISetOkStatus().IJson("ok, SubjectGetController")

}

func SubjectNameController(c *gin.Context) {
	c.ISetOkStatus().IJson("ok, SubjectNameController")

}
