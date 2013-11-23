package controller

import (
	."github.com/fuxiaohei/gorink/app"
	"github.com/fuxiaohei/goink/app"
	"github.com/fuxiaohei/gorink/model"
	"github.com/fuxiaohei/gorink/lib"
	"strings"
	"strconv"
	"errors"
)

func init() {
	App.GET("/admin/profile", func(context *app.InkContext) interface {} {
			renderAdminPage(context, "user/profile.html", "个人资料", "profile", map[string]interface {}{
					"Profile":model.GetUserById(model.GetCurrentUserId(context)),
					"IsCurrent":true,
				})
			return nil
		})
	App.POST("/admin/user/edit", func(context *app.InkContext) interface {} {
			validateErrors := validateProfileData(context)
			if len(validateErrors) > 0 {
				renderAdminAlert(context, validateErrors)
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
				renderAdminAlert(context, []error{e})
				return nil
			}
			context.Redirect("/admin/user/edit?updated=1&id=" + context.String("id"), 302)
			return nil
		})
	//---------------------
	App.GET("/admin/user", func(context *app.InkContext) interface {} {
			renderAdminPage(context, "user/manage.html", "用户", "user", map[string]interface {}{
					"Users":model.GetUsers(),
				})
			return nil
		})
	App.GET("/admin/user/new", func(context *app.InkContext) interface {} {
			renderAdminPage(context, "user/new.html", "添加用户", "user", nil)
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
			renderAdminPage(context, "user/profile.html", "个人资料", "user", map[string]interface {}{
					"Profile":model.GetUserById(id),
				})
			return nil
		})
	//-------------------
	App.POST("/admin/password", func(context *app.InkContext) interface {} {
			new := strings.TrimSpace(context.String("new"))
			confirm := strings.TrimSpace(context.String("confirm"))
			if len(new) < 6 || len(new) > 20 || new != confirm {
				renderAdminAlert(context, []error{errors.New("密码长度应在6-20位之间"), errors.New("确认密码不匹配")})
				return nil
			}
			old := strings.TrimSpace(context.String("old"))
			if len(old) < 6 || len(old) > 20 {
				renderAdminAlert(context, []error{errors.New("旧密码")})
				return nil
			}
			e := model.UpdatePassword(model.GetCurrentUserId(context), old, new, context.Ink.String("salt.password"))
			if e != nil {
				renderAdminAlert(context, []error{e})
				return nil
			}
			context.Redirect("/admin/profile/?password=1", 302)
			return nil
		})
}

func validateProfileData(context *app.InkContext) []error {
	messages := []error{}
	if !lib.IsASCII(context.String("login")) {
		messages = append(messages, errors.New("登录名不支持中文和特殊符号"))
	}
	if !lib.IsEmail(context.String("email")) {
		messages = append(messages, errors.New("邮箱格式错误"))
	}
	url := context.String("url")
	if len(url) > 0 {
		if !lib.IsURL(url) {
			messages = append(messages, errors.New("网址格式错误"))
		}
	}
	return messages
}
