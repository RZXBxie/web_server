package framework

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
)

type IResponse interface {
	// Json 将所有的响应方法的返回值设置成为 IResponse 类型，便于链式调用
	Json(obj interface{}) IResponse

	Jsonp(obj interface{}) IResponse

	Xml(obj interface{}) IResponse

	Html(file string, obj interface{}) IResponse

	Text(format string, values ...interface{}) IResponse

	Redirect(path string) IResponse

	SetHeader(key string, val string) IResponse

	SetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse

	SetStatus(code int) IResponse

	SetOkStatus() IResponse
}

func (ctx *Context) SetHeader(key string, val string) IResponse {
	ctx.responseWriter.Header().Add(key, val)
	return ctx
}

func (ctx *Context) SetCookie(key, val string, maxAge int, path string, domain string, secure, httpOnly bool) IResponse {
	if path == "" {
		path = "/"
	}
	http.SetCookie(ctx.responseWriter, &http.Cookie{
		Name:     key,
		Value:    val,
		Path:     path,
		Domain:   domain,
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
	return ctx
}

func (ctx *Context) SetStatus(code int) IResponse {
	ctx.responseWriter.WriteHeader(code)
	return ctx
}

func (ctx *Context) SetOkStatus() IResponse {
	ctx.responseWriter.WriteHeader(http.StatusOK)
	return ctx
}

func (ctx *Context) Redirect(path string) IResponse {
	http.Redirect(ctx.responseWriter, ctx.request, path, http.StatusMovedPermanently)
	return ctx
}

func (ctx *Context) Json(obj interface{}) IResponse {
	bte, err := json.Marshal(obj)
	if err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}
	ctx.SetHeader("Content-Type", "application/json")
	ctx.responseWriter.Write(bte)

	return ctx
}

func (ctx *Context) Jsonp(obj interface{}) IResponse {
	// 获取请求参数callback
	callbackFunc, _ := ctx.QueryString("callback", "callback_func")
	ctx.SetHeader("Content-Type", "application/javascript")

	// 输出到前端页面的时候必须对字符进行过滤，否则有可能造成xss攻击
	callback := template.JSEscapeString(callbackFunc)

	// 输出函数名
	_, err := ctx.responseWriter.Write([]byte(callback))
	if err != nil {
		return ctx
	}

	// 输出左括号
	_, err = ctx.responseWriter.Write([]byte("("))
	if err != nil {
		return ctx
	}

	ret, err := json.Marshal(obj)
	if err != nil {
		return ctx
	}

	// 输出实际数据
	_, err = ctx.responseWriter.Write(ret)
	if err != nil {
		return ctx
	}

	// 输出右括号
	_, err = ctx.responseWriter.Write([]byte(")"))
	if err != nil {
		return ctx
	}

	return ctx
}

func (ctx *Context) Xml(obj interface{}) IResponse {
	bye, err := xml.Marshal(obj)
	if err != nil {
		return ctx
	}
	ctx.SetHeader("Content-Type", "application/xml")
	ctx.responseWriter.Write(bye)

	return ctx
}

func (ctx *Context) Html(file string, obj interface{}) IResponse {
	// 读取模板文件，创建template示例
	t, err := template.New("output").ParseFiles(file)
	if err != nil {
		return ctx
	}

	// 执行Execute方法将obj和模板进行结合
	if err := t.Execute(ctx.responseWriter, obj); err != nil {
		return ctx
	}
	ctx.SetHeader("Content-Type", "application/html")

	return ctx
}

func (ctx *Context) Text(format string, values ...interface{}) IResponse {
	out := fmt.Sprintf(format, values...)
	ctx.SetHeader("Content-Type", "application/text")
	ctx.responseWriter.Write([]byte(out))

	return ctx
}
