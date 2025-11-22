package controller

import (
	"github.com/RZXBxie/web_server/framework/gin"
	"github.com/RZXBxie/web_server/provider/demo"
)

func SubjectListController(c *gin.Context) {
	demoService := c.MustMake(demo.Key).(demo.Service)
	c.ISetOkStatus().IJson(demoService.GetFoo())
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
