package controller

import (
	."github.com/fuxiaohei/gorink/app"
	"github.com/fuxiaohei/goink/app"
)

func renderHomePage(context *app.InkContext, template string, title string, rel string, data map[string]interface {}) {
	if data == nil {
		data = make(map[string]interface {})
	}
	data["Title"] = title
	data["Rel"] = rel
	context.Render(template + ",home/header.html,home/footer.html", data)
}

func renderNotFoundPage(context *app.InkContext) {
	context.Render(App.Config().StringOr("view.404", "404.html"), nil)
	context.Send("", 404)
}
