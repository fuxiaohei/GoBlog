package handler

import (
	"github.com/fuxiaohei/GoBlog/app/handler/helper"
	"github.com/fuxiaohei/GoInk"
)

// Json creates a json context response.
func Json(context *GoInk.Context, res bool) *helper.JsonContext {
	return helper.Json(context, res)
}

// Theme creates themed context response.
func Theme(context *GoInk.Context) *helper.ThemeContext {
	return helper.Theme(context)
}
