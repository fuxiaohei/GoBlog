package controller

import (
	"fmt"
	"github.com/fuxiaohei/goink/app"
	. "github.com/fuxiaohei/gorink/app"
	"github.com/fuxiaohei/gorink/model"
	"strings"
)

func init() {
	App.GET("/admin", func(context *app.InkContext) interface{} {
			renderAdminPage(context, "admin/index.html", "管理", "panel", map[string]interface {}{
					"Stats": map[string]interface{}{
						"Tokens": model.LoginStats(),
					},
				})
			return nil
		})
	App.On("router.run.before@authorize", func(context *app.InkContext) {
			if strings.HasPrefix(context.URI, "/admin") {
				token := context.Cookie("gorink_id")
				if len(token) < 10 {
					context.Redirect("/login", 302)
					context.Send("")
					return
				}
				r, tokenData := model.CheckLoginByToken(token)
				if !r {
					context.Redirect("/login/out?refer=" + context.Path, 302)
					context.Send("")
					return
				}
				context.Flash("CurrentUser", tokenData.Name)
				context.Flash("CurrentUserId", fmt.Sprint(tokenData.UserId))
				context.Flash("CurrentUserAvatar", tokenData.Avatar)
			}
		})
}

func renderAdminPage(context *app.InkContext, template string, title string, rel string, data map[string]interface {}) {
	if data == nil {
		data = make(map[string]interface {})
	}
	data["Title"] = title
	data["Rel"] = rel
	context.Render(template + ",admin/header.html,admin/footer.html", data)
}

func renderAdminAlert(context *app.InkContext, errors []error) {
	context.Render("admin/alert.html", map[string]interface {}{
			"Errors":errors,
		})
}
