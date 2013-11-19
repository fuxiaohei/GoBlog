
package controller

import (
	."github.com/fuxiaohei/gorink/app"
	"github.com/fuxiaohei/goink/app"
)

func init() {
	App.GET("/admin/category", func(context *app.InkContext) interface {} {
			context.Render("category/manage.html,admin/header.html,admin/footer.html", map[string]interface {}{
					"title":"分类",
					"rel":"category",
				})
			return nil
		})
	App.GET("/admin/category/new", func(context *app.InkContext) interface {} {
			context.Render("category/new.html,admin/header.html,admin/footer.html", map[string]interface {}{
					"Title":"新建分类",
					"Rel":"category",
				})
			return nil
		})
	App.GET("/admin/category/edit", func(context *app.InkContext) interface {} {
			context.Render("category/edit.html,admin/header.html,admin/footer.html", map[string]interface {}{
					"Title":"编辑分类",
					"Rel":"category",
				})
			return nil
		})
}

