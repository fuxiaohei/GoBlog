package handler

import (
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoInk/Core"
	"strconv"
	"strings"
)

func Admin(context *Core.Context) interface{} {
	context.Render("admin:admin/dashboard.html", map[string]interface{}{
		"IsDashboard": true,
		"Title":       "控制台",
	})
	return nil
}

func authAdmin(context *Core.Context) {
	if strings.HasPrefix(context.Url, "/admin") {
		token := context.Cookie("admin-user-token")
		session := model.SessionM.GetByToken(token)
		if session != nil {
			uid, _ := strconv.Atoi(context.Cookie("admin-user"))
			if session.IsValid(uid, context.Ip) {
				return
			}
		}
		context.Redirect("/logout.html")
		context.Send()
	}
}
