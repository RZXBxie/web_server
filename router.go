package main

import (
	"time"

	"github.com/RZXBxie/web_server/controller"
	"github.com/RZXBxie/web_server/framework/gin"
	"github.com/RZXBxie/web_server/framework/middleware"
)

func registerRouter(core *gin.Engine) {
	// 静态路由匹配
	duration := time.Second * 5
	core.GET("/user/login", middleware.Timeout(duration), controller.UserLoginController)
	// 路由组+动态路由匹配
	subjectGroup := core.Group("/subject")
	{
		subjectGroup.DELETE("/:id", controller.SubjectDelController)
		subjectGroup.GET("/:id", controller.SubjectGetController)
		subjectGroup.PUT("/:id", controller.SubjectUpdateController)
		subjectGroup.GET("/list/all", controller.SubjectListController)
		subjectInnerGroup := subjectGroup.Group("/info")
		{
			subjectInnerGroup.Use(controller.UserLoginController)
			subjectInnerGroup.GET("/name", controller.SubjectDelController)
		}
	}
}
