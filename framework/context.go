package framework

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// Context 自定义Context接口
// 通过阅读源码发现，http.ListenAndServe 方法在Serve函数中产生了baseCtx，并通过c.serve方法往传递，
// 最终会将ctx传递给req.ctx。因此我们自定义的Context结构就继承req.ctx
type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	ctx            context.Context

	// 当前请求的handler链条
	handlers []ControllerHandler
	// 当前请求调用到调用链的哪个节点
	index int

	// 超时标记
	IsTimeOut bool

	// 写锁
	writeMutex *sync.Mutex

	params map[string]string
}

func NewContext(r *http.Request, rw http.ResponseWriter) *Context {
	return &Context{
		request:        r,
		responseWriter: rw,
		// 由于每次调用Next时index都会自增1，所以需要从-1开始加，这样才能正确使用handler
		index:      -1,
		ctx:        r.Context(),
		writeMutex: &sync.Mutex{},
	}
}

// #region base func

func (ctx *Context) WriteMux() *sync.Mutex {
	return ctx.writeMutex
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.responseWriter
}

func (ctx *Context) SetIsTimeOut() {
	ctx.IsTimeOut = true
}

func (ctx *Context) IsTimeout() bool {
	return ctx.IsTimeOut
}

// SetHandlers 将路由中handlers注册到ctx中，后续通过 Next 函数依次调用
func (ctx *Context) SetHandlers(handlers []ControllerHandler) {
	ctx.handlers = handlers
}

func (ctx *Context) SetParams(params map[string]string) {
	ctx.params = params
}

// Next 核心函数，调用context的下一个函数
// 调用时机：1、第一次启动服务的时候
// 2、每个中间件都会调用 Next ，因此ctx.index是自增的
func (ctx *Context) Next() error {
	ctx.index++
	if ctx.index < len(ctx.handlers) {
		if err := ctx.handlers[ctx.index](ctx); err != nil {
			return err
		}
	}

	return nil
}

// #endrigion

// #region 自定义Context实现context.Context接口

// BaseContext server启动后，所有的context最终都保存在req中，所以baseContext就是req中的context
func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.BaseContext().Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

func (ctx *Context) Err() error {
	return ctx.BaseContext().Err()
}

func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.BaseContext().Value(key)
}

// #endreigon
