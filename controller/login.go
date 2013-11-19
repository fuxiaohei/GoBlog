package controller

import (
	"fmt"
	"github.com/fuxiaohei/goink/app"
	. "github.com/fuxiaohei/gorink/app"
	"github.com/fuxiaohei/gorink/model"
)

func init() {
	App.GET("/login", func(context *app.InkContext) interface{} {
			if len(context.Cookie("gorink_id")) > 10 {
				context.Redirect("/admin/", 302)
				return nil
			}
			context.Render("login.html", nil)
			return nil
		})
	App.GET("/login/out", func(context *app.InkContext) interface{} {
			context.Cookie("gorink_id", "", fmt.Sprint(-3600))
			url := "/login/"
			if len(context.String("refer")) > 0 {
				url += "?r=" + context.String("refer")
			}
			context.Redirect(url, 302)
			return nil
		})
	App.POST("/login", func(context *app.InkContext) interface{} {
			user := context.String("user")
			loginInfo, message := model.Login(user, context.String("password"), context.Ink.String("salt.password"))
			if message == "" {
				context.Cookie("gorink_id", loginInfo.Token, fmt.Sprint(3600*24*7))
				redirect := context.String("r")
				if len(redirect) > 0 {
					context.Redirect(redirect, 302)
					return nil
				}
				context.Redirect("/admin/", 302)
				return nil
			}
			context.Render("login.html", map[string]interface{}{"error": message,
					"user": user})
			return nil
		})
}
