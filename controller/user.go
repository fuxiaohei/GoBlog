package controller

import (
	."github.com/fuxiaohei/gorink/app"
	"github.com/fuxiaohei/goink/app"
	"github.com/fuxiaohei/gorink/model"
	"github.com/fuxiaohei/gorink/lib"
	"strings"
	"strconv"
)

func init() {
	App.GET("/admin/profile", func(context *app.InkContext) interface {} {
			context.Render("user/profile.html,admin/header.html,admin/footer.html", map[string]interface {}{
					"Title":"个人资料",
					"Rel":"profile",
					"Profile":model.GetUserById(model.GetCurrentUserId(context)),
					"IsCurrent":true,
				})
			return nil
		})
	App.POST("/admin/user/edit", func(context *app.InkContext) interface {} {
			validateErrors := validateProfileData(context)
			if len(validateErrors) > 0 {
				context.Render("admin/alert.html", map[string]interface {}{
						"Errors":validateErrors,
					})
				return nil
			}
			login := context.String("login")
			nick := context.String("nick")
			email := context.String("email")
			url := context.String("url")
			bio := context.String("bio")
			id, _ := strconv.Atoi(context.String("id"))
			e := model.UpdateUserProfile(id, login, nick, email, url, bio)
			if e != nil {
				panic(e)
			}
			context.Redirect("/admin/user/edit?updated=1&id=" + context.String("id"), 302)
			return nil
		})
	//---------------------
	App.GET("/admin/user", func(context *app.InkContext) interface {} {
			context.Render("user/manage.html,admin/header.html,admin/footer.html", map[string]interface {}{
					"Title":"用户",
					"Rel":"user",
					"Users":model.GetUsers(),
				})
			return nil
		})
	App.GET("/admin/user/new", func(context *app.InkContext) interface {} {
			context.Render("user/new.html,admin/header.html,admin/footer.html", map[string]interface {}{
					"Title":"添加新用户",
					"Rel":"user",
				})
			return nil
		})
	App.GET("/admin/user/edit", func(context *app.InkContext) interface {} {
			id, _ := strconv.Atoi(context.String("id"))
			if id < 1 {
				context.Redirect("/admin/user", 302)
				return nil
			}
			if id == model.GetCurrentUserId(context) {
				context.Redirect("/admin/profile", 302)
				return nil
			}
			context.Render("user/profile.html,admin/header.html,admin/footer.html", map[string]interface {}{
					"Title":"个人资料",
					"Rel":"profile",
					"Profile":model.GetUserById(id),
				})
			return nil
		})
	//-------------------
	App.POST("/admin/password", func(context *app.InkContext) interface {} {
			new := strings.TrimSpace(context.String("new"))
			confirm := strings.TrimSpace(context.String("confirm"))
			if len(new) < 6 || len(new) > 20 || new != confirm {
				context.Render("admin/alert.html", map[string]interface {}{
						"Errors":[]string{"密码长度应在6-20位之间", "确认密码不匹配"},
					})
				return nil
			}
			old := strings.TrimSpace(context.String("old"))
			if len(old) < 6 || len(old) > 20 {
				context.Render("admin/alert.html", map[string]interface {}{
						"Errors":[]string{"旧密码错误"},
					})
				return nil
			}
			e := model.UpdatePassword(model.GetCurrentUserId(context), old, new, context.Ink.String("salt.password"))
			if e != nil {
				context.Render("admin/alert.html", map[string]interface {}{
						"Errors":[]string{e.Error()},
					})
				return nil
			}
			context.Redirect("/admin/profile/?password=1", 302)
			return nil
		})
}

func validateProfileData(context *app.InkContext) []string {
	messages := []string{}
	if !lib.IsASCII(context.String("login")) {
		messages = append(messages, "登录名不支持中文和特殊符号")
	}
	if !lib.IsEmail(context.String("email")) {
		messages = append(messages, "邮箱格式错误")
	}
	url := context.String("url")
	if len(url) > 0 {
		if !lib.IsURL(url) {
			messages = append(messages, "网址格式错误")
		}
	}
	return messages
}
