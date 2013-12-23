package handler

import (
	"github.com/fuxiaohei/GoInk/Core"
	"strconv"
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"github.com/fuxiaohei/GoBlog/app"
)

func AdminProfile(context *Core.Context) interface {} {
	uid, _ := strconv.Atoi(context.Cookie("admin-user"))
	context.Render("admin:admin/profile.html", map[string]interface {}{
			"IsProfile":true,
			"User":model.UserM.GetUserById(uid),
			"Title":"个人资料",
			"Update":context.String("update"),
		})
	return nil
}

func AdminProfilePost(context *Core.Context) interface {} {
	data := context.Input()
	uid, _ := strconv.Atoi(context.Cookie("admin-user"))
	model.UserM.SaveProfile(uid, data["login"], data["display"], data["email"], data["site"], utils.Gravatar(data["email"], "180"))
	context.Redirect("/admin/profile?update=1")
	app.Ink.Listener.EmitAll("model.user.profile.update", uid)
	return nil
}

func AdminPassword(context *Core.Context) interface {} {
	context.Render("admin:admin/password.html", map[string]interface {}{
			"IsProfile":true,
			"IsPasswprd":true,
			"Title":"修改密码",
			"Update":context.String("update"),
			"Error":context.String("error"),
		})
	return nil
}

func AdminPasswordPost(context *Core.Context) interface {} {
	data := context.Input()
	if data["new"] != data["confirm"] {
		context.Redirect("/admin/password?error=1")
		return nil
	}
	uid, _ := strconv.Atoi(context.Cookie("admin-user"))
	e := model.UserM.SavePassword(uid, data["old"], data["new"])
	if e != nil {
		context.Redirect("/admin/password?error=1")
		return nil
	}
	context.Redirect("/admin/password?update=1")
	app.Ink.Listener.EmitAll("model.user.password.update", uid)
	return nil
}
