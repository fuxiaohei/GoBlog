package handler

import (
	"errors"
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"github.com/fuxiaohei/GoInk/Core"
	"strconv"
)

func Login(context *Core.Context) interface{} {
	if len(context.Cookie("admin-user-token")) > 10 {
		context.Redirect("/admin/")
		return nil
	}
	context.Render("admin/login.html", nil)
	return nil
}

func Logout(context *Core.Context) interface {} {
	context.Cookie("admin-user", "", "-3600")
	context.Cookie("admin-user-token", "", "-3600")
	context.Redirect("/login.html")
	return nil
}

func LoginPost(context *Core.Context) interface{} {
	if !context.IsAjax {
		context.Status = 400
		return nil
	}
	data, expire := context.Input(), 3600
	if data["remember"] == "on" {
		expire = 3600*24 * 3
	}
	session, err := doLogin(data["login"], data["password"], context, int64(expire))
	if err != nil {
		context.Json(map[string]interface{}{
			"res": false,
			"msg": err.Error(),
		})
		return nil
	}
	uid := strconv.Itoa(session.UserId)
	expireStr := strconv.Itoa(expire)
	context.Cookie("admin-user", uid, expireStr)
	context.Cookie("admin-user-token", session.Token, expireStr)
	context.Json(map[string]interface{}{
		"res": true,
	})
	return nil
}

func doLogin(login string, password string, context *Core.Context, expireTime int64) (*model.Session, error) {
	user := model.UserM.GetUserByLogin(login)
	if user == nil {
		return nil, errors.New("无效的用户")
	}
	if user.Password != utils.Sha1(password) {
		return nil, errors.New("错误的密码")
	}
	return model.SessionM.CreateSession(user.Id, context.Ip, context.UserAgent, expireTime), nil
}
