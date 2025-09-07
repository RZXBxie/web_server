package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Context 自定义Context接口
type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	ctx            context.Context
	handler        ControllerHandler
	
	// 超时标记
	IsTimeOut bool
	
	// 写锁
	writeMutex *sync.Mutex
}

func NewContext(r *http.Request, rw http.ResponseWriter) *Context {
	return &Context{
		request:        r,
		responseWriter: rw,
		ctx:            r.Context(),
		writeMutex:     &sync.Mutex{},
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

// #region 实现request相关功能

func (ctx *Context) QueryInt(key string, defaultValue int) int {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		length := len(vals)
		if length > 0 {
			intValue, err := strconv.Atoi(vals[length-1])
			if err != nil {
				return defaultValue
			}
			return intValue
		}
	}
	return defaultValue
}

func (ctx *Context) QueryString(key, defaultValue string) string {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		length := len(vals)
		if length > 0 {
			return vals[length-1]
		}
	}
	return defaultValue
}

func (ctx *Context) QueryArray(key string, defaultValue []string) []string {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return defaultValue
}

func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request != nil {
		return ctx.request.URL.Query()
	}
	return map[string][]string{}
}

func (ctx *Context) FormInt(key string, defaultValue int) int {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		length := len(vals)
		if length > 0 {
			intValue, err := strconv.Atoi(vals[length-1])
			if err != nil {
				return defaultValue
			}
			return intValue
		}
	}
	return defaultValue
}

func (ctx *Context) FormString(key, defaultValue string) string {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		length := len(vals)
		if length > 0 {
			return vals[length-1]
		}
	}
	return defaultValue
}

func (ctx *Context) FormArray(key string, defaultValue []string) []string {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return defaultValue
}

func (ctx *Context) FormAll() map[string][]string {
	if ctx.request != nil {
		return ctx.request.PostForm
	}
	return map[string][]string{}
}

// #endrigion

// #region application/json post

func (ctx *Context) BindJson(obj interface{}) error {
	if ctx.request != nil {
		body, err := io.ReadAll(ctx.request.Body)
		if err != nil {
			return err
		}
		
		// 重新设置Body，方便以后可以重复读取
		ctx.request.Body = io.NopCloser(bytes.NewBuffer(body))
		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ctx.request is empty")
	}
	
	return nil
}

func (ctx *Context) Json(statusCode int, obj interface{}) error {
	if ctx.IsTimeout() {
		return nil
	}
	ctx.responseWriter.Header().Set("Content-Type", "application/json")
	ctx.responseWriter.WriteHeader(statusCode)
	byt, err := json.Marshal(obj)
	if err != nil {
		ctx.responseWriter.WriteHeader(http.StatusInternalServerError)
		return err
	}
	ctx.responseWriter.Write(byt)
	return nil
}

func (ctx *Context) HTML(statusCode int, obj interface{}, html string) error {
	return nil
}

func (ctx *Context) Text(statusCode int, obj string) error {
	return nil
}

// #endregion
