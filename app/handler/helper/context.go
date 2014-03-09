package helper

import (
	"github.com/fuxiaohei/GoBlog/app/model/setting"
	"github.com/fuxiaohei/GoInk"
	"path"
)

type JsonContext struct {
	context *GoInk.Context
	data    map[string]interface{}
}

// Json creates a json context response.
func Json(context *GoInk.Context, res bool) *JsonContext {
	c := new(JsonContext)
	c.context = context
	c.data = make(map[string]interface{})
	c.data["res"] = res
	return c
}

func (jc *JsonContext) Set(key string, v interface{}) *JsonContext {
	jc.data[key] = v
	return jc
}

func (jc *JsonContext) End() {
	jc.context.Json(jc.data)
}

type ThemeContext struct {
	context *GoInk.Context
	theme   string
}

// Theme creates themed context response.
func Theme(context *GoInk.Context) *ThemeContext {
	t := new(ThemeContext)
	t.context = context
	t.theme = setting.Get("site_theme")
	if t.theme == "" {
		t.theme = "default"
	}
	return t
}

func (tc *ThemeContext) Layout(layout string) *ThemeContext {
	if layout == "" {
		tc.context.Layout("")
		return tc
	}
	tc.context.Layout(path.Join(tc.theme, layout))
	return tc
}

func (tc *ThemeContext) Render(tpl string, data map[string]interface{}) {
	tc.context.Render(path.Join(tc.theme, tpl), data)
}

func (tc *ThemeContext) Tpl(tpl string, data map[string]interface{}) string {
	return tc.context.Tpl(path.Join(tc.theme, tpl), data)
}

func (tc *ThemeContext) Section(tpl string, data map[string]interface{}) string {
	return tc.Tpl("section_"+tpl, data)
}

func (tc *ThemeContext) Has(tpl string) bool {
	file := path.Join(tc.theme, tpl)
	return tc.context.App().View().Has(file)
}

func (tc *ThemeContext) HasSection(tpl string) bool {
	return tc.Has("section_" + tpl)
}
