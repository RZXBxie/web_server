package framework

import (
	"log"
	"net/http"
)

// Core 框架核心结构
type Core struct {
	router map[string]ControllerHandler
}

// NewCore 初始化框架核心结构
func NewCore() *Core {
	return &Core{router: map[string]ControllerHandler{}}
}

func (c *Core) Get(url string, handler ControllerHandler) {
	c.router[url] = handler
}

// ServeHTTP 框架核心结构实现Handler功能
func (c *Core) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Println("core.ServeHTTP")
	ctx := NewContext(req, rw)
	
	// 先测试 "foo" 这个路由
	router := c.router["foo"]
	if router == nil {
		return
	}
	
	log.Println("core.router")
	router(ctx)
}
