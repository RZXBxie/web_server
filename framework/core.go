package framework

import (
	"log"
	"net/http"
	"strings"
)

const (
	HTTP_METHOD_GET    = "GET"
	HTTP_METHOD_POST   = "POST"
	HTTP_METHOD_PUT    = "PUT"
	HTTP_METHOD_DELETE = "DELETE"
)

// Core 框架核心结构
type Core struct {
	// router 两级map方便路由寻址
	// 一级为http method，二级为路由以及对应的handler
	router map[string]*Trie
}

// NewCore 初始化框架核心结构
func NewCore() *Core {
	
	// 将二级map写入一级map
	router := map[string]*Trie{}
	router[HTTP_METHOD_GET] = NewTrie()
	router[HTTP_METHOD_POST] = NewTrie()
	router[HTTP_METHOD_PUT] = NewTrie()
	router[HTTP_METHOD_DELETE] = NewTrie()
	
	return &Core{router}
}

// Get 注册GET方法路由
func (c *Core) Get(url string, handler ControllerHandler) {
	if err := c.router[HTTP_METHOD_GET].AddRoute(url, handler); err != nil {
		log.Fatalf("add route fail: %v", err)
	}
}

// Post 注册POST方法路由
func (c *Core) Post(url string, handler ControllerHandler) {
	if err := c.router[HTTP_METHOD_POST].AddRoute(url, handler); err != nil {
		log.Fatalf("add route fail: %v", err)
	}
}

// Put 注册Put方法路由
func (c *Core) Put(url string, handler ControllerHandler) {
	if err := c.router[HTTP_METHOD_PUT].AddRoute(url, handler); err != nil {
		log.Fatalf("add route fail: %v", err)
	}
}

// Delete 注册DELETE方法路由
func (c *Core) Delete(url string, handler ControllerHandler) {
	if err := c.router[HTTP_METHOD_DELETE].AddRoute(url, handler); err != nil {
		log.Fatalf("add route fail: %v", err)
	}
}

// Group 实现路由组功能
func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

// FindRouterByRequest 根据req查找指定handler
func (c *Core) FindRouterByRequest(req *http.Request) ControllerHandler {
	// uri和method转为大写，保证大小写不敏感
	uri := strings.ToUpper(req.URL.Path)
	method := strings.ToUpper(req.Method)
	
	if methodHandlers, ok := c.router[method]; ok {
		return methodHandlers.FindHandler(uri)
	}
	
	return nil
}

// ServeHTTP 框架核心结构实现Handler功能，该方法负责路由分发
func (c *Core) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Println("core.ServeHTTP")
	ctx := NewContext(req, rw)
	
	router := c.FindRouterByRequest(req)
	if router == nil {
		ctx.Json(http.StatusNotFound, "404 page not found")
		return
	}
	
	log.Println("core.router")
	
	if err := router(ctx); err != nil {
		ctx.Json(http.StatusInternalServerError, "500 Internal server error")
		return
	}
}
