package gof

import (
	"fmt"
	"net/http"
)

type Context struct {
	Server       *HttpServer
	handlers     []RouterHandler
	handlerIndex int
	params       map[string]string
	injector     *Injector

	responseWriter http.ResponseWriter
	request        *http.Request

	Status int
	Body   []byte
}

func NewContext(s *HttpServer, w http.ResponseWriter, r *http.Request, fn []RouterHandler, params map[string]string, in *Injector) *Context {
	c := new(Context)
	c.Server = s
	c.handlerIndex = 0
	c.handlers = fn
	c.params = params
	c.injector = in
	c.responseWriter = w
	c.request = r
	return c
}

func (ctx *Context) Run() {
	if ctx.handlerIndex >= len(ctx.handlers) {
		return
	}
	fn := ctx.handlers[ctx.handlerIndex]
	if fn == nil {
		return
	}
	fn(ctx)
	ctx.handlerIndex++
}

func (ctx *Context) Next() {
	ctx.handlerIndex++
	ctx.Run()
}

func (ctx *Context) Request() *http.Request {
	return ctx.request
}

func (ctx *Context) Response() http.ResponseWriter {
	return ctx.responseWriter
}

func (ctx *Context) WriteHeader(key, value string) {
	ctx.responseWriter.Header().Set(key, value)
}

func (ctx *Context) SendResponse() {
	if ctx.Status > 0 {
		return
	}
	if ctx.Status == 0 {
		ctx.Status = 200
	}
	ctx.responseWriter.WriteHeader(ctx.Status)
	ctx.responseWriter.Write(ctx.Body)
}
