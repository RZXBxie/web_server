// Copyright 2021 jianfengye.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package gin

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

// IResponse IResponse代表返回方法
type IResponse interface {
	IJson(obj interface{}) IResponse

	IJsonp(obj interface{}) IResponse

	IXml(obj interface{}) IResponse

	IHtml(template string, obj interface{}) IResponse

	IText(format string, values ...interface{}) IResponse

	IRedirect(path string) IResponse

	ISetHeader(key string, val string) IResponse

	ISetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse

	ISetStatus(code int) IResponse

	ISetOkStatus() IResponse
}

// IJsonp Jsonp输出
func (c *Context) IJsonp(obj interface{}) IResponse {
	// 获取请求参数callback
	callbackFunc := c.Query("callback")
	c.ISetHeader("Content-Type", "application/javascript")
	// 输出到前端页面的时候需要注意下进行字符过滤，否则有可能造成xss攻击
	callback := template.JSEscapeString(callbackFunc)

	// 输出函数名
	_, err := c.Writer.Write([]byte(callback))
	if err != nil {
		return c
	}
	// 输出左括号
	_, err = c.Writer.Write([]byte("("))
	if err != nil {
		return c
	}
	// 数据函数参数
	ret, err := json.Marshal(obj)
	if err != nil {
		return c
	}
	_, err = c.Writer.Write(ret)
	if err != nil {
		return c
	}
	// 输出右括号
	_, err = c.Writer.Write([]byte(")"))
	if err != nil {
		return c
	}
	return c
}

// IXml xml输出
func (c *Context) IXml(obj interface{}) IResponse {
	byt, err := xml.Marshal(obj)
	if err != nil {
		return c.ISetStatus(http.StatusInternalServerError)
	}
	c.ISetHeader("Content-Type", "application/html")
	c.Writer.Write(byt)
	return c
}

// IHtml html输出
func (c *Context) IHtml(file string, obj interface{}) IResponse {
	// 读取模版文件，创建template实例
	t, err := template.New("output").ParseFiles(file)
	if err != nil {
		return c
	}
	// 执行Execute方法将obj和模版进行结合
	if err := t.Execute(c.Writer, obj); err != nil {
		return c
	}

	c.ISetHeader("Content-Type", "application/html")
	return c
}

// IText string
func (c *Context) IText(format string, values ...interface{}) IResponse {
	out := fmt.Sprintf(format, values...)
	c.ISetHeader("Content-Type", "application/text")
	c.Writer.Write([]byte(out))
	return c
}

// IRedirect 重定向
func (c *Context) IRedirect(path string) IResponse {
	http.Redirect(c.Writer, c.Request, path, http.StatusMovedPermanently)
	return c
}

// ISetHeader header
func (c *Context) ISetHeader(key string, val string) IResponse {
	c.Writer.Header().Add(key, val)
	return c
}

// ISetCookie Cookie
func (c *Context) ISetCookie(key string, val string, maxAge int, path string, domain string, secure bool, httpOnly bool) IResponse {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     key,
		Value:    url.QueryEscape(val),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: 1,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
	return c
}

// ISetStatus 设置状态码
func (c *Context) ISetStatus(code int) IResponse {
	c.Writer.WriteHeader(code)
	return c
}

// ISetOkStatus 设置200状态
func (c *Context) ISetOkStatus() IResponse {
	c.Writer.WriteHeader(http.StatusOK)
	return c
}

func (c *Context) IJson(obj interface{}) IResponse {
	byt, err := json.Marshal(obj)
	if err != nil {
		return c.ISetStatus(http.StatusInternalServerError)
	}
	c.ISetHeader("Content-Type", "application/json")
	c.Writer.Write(byt)
	return c
}
