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
			context.Render("admin/index.html,admin/header.html,admin/footer.html", map[string]interface{}{
					"Title": "管理",
					"Rel":   "panel",
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
