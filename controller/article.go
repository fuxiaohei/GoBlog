
package controller

import (
	."github.com/fuxiaohei/gorink/app"
	"github.com/fuxiaohei/goink/app"
)

func init() {
	App.GET("/admin/article", func(context *app.InkContext) interface {} {
			context.Render("article/manage.html,admin/header.html,admin/footer.html", map[string]interface {}{
					"Title":"文章",
					"Rel":"article",
				})
			return nil
		})
	App.GET("/admin/article/new", func(context *app.InkContext) interface {} {
			context.Render("article/new.html,admin/header.html,admin/footer.html", map[string]interface {}{
					"Title":"撰写文章",
					"Rel":"article",
				})
			return nil
		})
}

