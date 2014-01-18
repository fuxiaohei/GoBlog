package handler

import (
	"github.com/fuxiaohei/GoBlog/GoInk"
	"github.com/fuxiaohei/GoBlog/app/model"
	"path"
)

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

type themeContext struct {
	context *GoInk.Context
	theme   string
}

func Theme(context *GoInk.Context) *themeContext {
	t := new(themeContext)
	t.context = context
	t.theme = model.GetSetting("site_theme")
	if t.theme == "" {
		t.theme = "default"
	}
	return t
}

func (tc *themeContext) Layout(layout string) *themeContext {
	if layout == "" {
		tc.context.Layout("")
		return tc
	}
	tc.context.Layout(path.Join(tc.theme, layout))
	return tc
}

func (tc *themeContext) Render(tpl string, data map[string]interface{}) {
	tc.context.Render(path.Join(tc.theme, tpl), data)
}

func (tc *themeContext) Tpl(tpl string, data map[string]interface{}) string {
	return tc.context.Tpl(path.Join(tc.theme, tpl), data)
}

func (tc *themeContext) Has(tpl string) bool {
	file := path.Join(tc.theme, tpl)
	return tc.context.App().View().Has(file)
}
