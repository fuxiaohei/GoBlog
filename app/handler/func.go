package handler

import "github.com/fuxiaohei/GoBlog/GoInk"

type jsonContext struct {
	context *GoInk.Context
	data    map[string]interface{}
}

func Json(context *GoInk.Context, res bool) *jsonContext {
	c := new(jsonContext)
	c.context = context
	c.data = make(map[string]interface{})
	c.data["res"] = res
	return c
}

func (jc *jsonContext) Set(key string, v interface{}) *jsonContext {
	jc.data[key] = v
	return jc
}

func (jc *jsonContext) End() {
	jc.context.Json(jc.data)
}
