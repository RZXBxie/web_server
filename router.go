package main

import (
	"github.com/RZXBxie/web_server/controller"
	"github.com/RZXBxie/web_server/framework"
)

func registerRouter(core *framework.Core) {
	// 静态路由匹配
	core.Get("/user/login", controller.UserLoginController)
	// 路由组+动态路由匹配
	subjectGroup := core.Group("/subject")
	{
		subjectGroup.Delete("/:id", controller.SubjectDelController)
		subjectGroup.Get("/:id", controller.SubjectGetController)
		subjectGroup.Put("/:id", controller.SubjectUpdateController)
		subjectGroup.Get("/list/all", controller.SubjectListController)
	}
}
