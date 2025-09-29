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
	// router key对应HTTP的Method，value就是一棵路由树
	router      map[string]*Trie
	middlewares []ControllerHandler
}

// NewCore 初始化框架核心结构
func NewCore() *Core {

	// 将二级map写入一级map
	router := map[string]*Trie{}
	router[HTTP_METHOD_GET] = NewTrie()
	router[HTTP_METHOD_POST] = NewTrie()
	router[HTTP_METHOD_PUT] = NewTrie()
	router[HTTP_METHOD_DELETE] = NewTrie()

	return &Core{router: router}
}

func (c *Core) Use(middlewares ...ControllerHandler) {
	c.middlewares = append(c.middlewares, middlewares...)
}

// Get 注册GET方法路由
func (c *Core) Get(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router[HTTP_METHOD_GET].AddRoute(url, allHandlers); err != nil {
		log.Fatalf("add route fail: %v", err)
	}
}

// Post 注册POST方法路由
func (c *Core) Post(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router[HTTP_METHOD_POST].AddRoute(url, allHandlers); err != nil {
		log.Fatalf("add route fail: %v", err)
	}
}

// Put 注册Put方法路由
func (c *Core) Put(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router[HTTP_METHOD_PUT].AddRoute(url, allHandlers); err != nil {
		log.Fatalf("add route fail: %v", err)
	}
}

// Delete 注册DELETE方法路由
func (c *Core) Delete(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router[HTTP_METHOD_DELETE].AddRoute(url, allHandlers); err != nil {
		log.Fatalf("add route fail: %v", err)
	}
}

// Group 实现路由组功能
func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

// FindRouterNodeByRequest 根据req查找指定handler
func (c *Core) FindRouterNodeByRequest(req *http.Request) *Node {
	// uri和method转为大写，保证大小写不敏感
	uri := strings.ToUpper(req.URL.Path)
	method := strings.ToUpper(req.Method)

	if methodHandlers, ok := c.router[method]; ok {
		return methodHandlers.root.matchNode(uri)
	}

	return nil
}

// ServeHTTP 框架核心结构实现Handler功能，该方法负责路由分发
func (c *Core) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Println("core.ServeHTTP")
	ctx := NewContext(req, rw)

	node := c.FindRouterNodeByRequest(req)
	if node == nil {
		ctx.SetStatus(http.StatusNotFound).Json("not found")
		return
	}

	ctx.SetHandlers(node.handlers)

	log.Println("core.handlers")

	// 设置路由参数
	params := node.findParamsFromEndNode(req.URL.Path)
	ctx.SetParams(params)

	if err := ctx.Next(); err != nil {
		ctx.SetStatus(http.StatusInternalServerError).Json("inner error")
		return
	}
}
