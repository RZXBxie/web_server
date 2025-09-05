package framework

import "net/http"

// Core 框架核心结构
type Core struct {
}

// NewCore 初始化框架核心结构
func NewCore() *Core {
	return &Core{}
}

// ServeHTTP 框架核心结构实现Handler功能
func (c *Core) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// TODO
}
