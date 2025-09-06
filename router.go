package main

import "github.com/RZXBxie/web_server/framework"

func registerRouter(core *framework.Core) {
	core.Get("foo", FooControllerHandler)
}
