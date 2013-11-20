package model

import (
	. "github.com/fuxiaohei/gorink/app"
	"github.com/fuxiaohei/goink/db"
	"github.com/fuxiaohei/goink/app"
	"errors"
	"fmt"
)

type User struct {
	Id           int `col:"id" tbl:"gorink_user"`
	LoginName    string `col:"login"`
	Password     string `col:"password"`
	NickName     string `col:"nick"`
	Email        string `col:"email"`
	Url          string `col:"url"`
	RegisterTime int64 `col:"register_time"`
	LoginTime    int64 `col:"login_time"`
	Status       string `col:"status"`
	Role         string `col:"role_name"`
	PasswordKey  string `col:"password_key"`
	Bio          string `col:"biography"`
}

func (this *User) DisplayName() string {
	if len(this.NickName) < 1 {
		return this.LoginName
	}
	return this.NickName
}

func (this *User) RoleName() string {
	if this.Role == "super" {
		return "超级管理员"
	}
	if this.Role == "admin" {
		return "管理员"
	}
	if this.Role == "editor" {
		return "编辑"
	}
	if this.Role == "reader" {
		return "读者"
	}
	return "未知"
}

func (this *User) Avatar() string {
	return GetGravatar(this.Email, "40")
}

func init() {
	db.Define(User{})
}

func GetUserByLogin(login string) *User {
	data, e := Orm.FindOne("model.User", db.NewSql("").Where("login = ?"), login)
	if e != nil {
		App.LogErr(e)
		return nil
	}
	return data.(*User)
}

func GetUserById(id int) *User {
	data, e := Orm.FindOne("model.User", db.NewSql("").Where("id = ?"), id)
	if e != nil {
		App.LogErr(e)
		return nil
	}
	return data.(*User)
}

func GetCurrentUserId(context *app.InkContext) int {
	token := context.Cookie("gorink_id")
	tokenData := tokens[token]
	if tokenData == nil {
		return 0
	}
	return tokenData.UserId
}

func UpdateUserProfile(userId int, login string, nick string, email string, url string, bio string) error {
	sql := db.NewSql("gorink_user", "id").Where("email = ?").Select()
	result, e := Db.Query(sql, email)
	data := result.Map()
	if data != nil {
		if len(data["id"]) > 0 && data["id"] != fmt.Sprint(userId) {
			return errors.New("邮箱和别的用户重复")
		}
	}
	sql = db.NewSql("gorink_user", "login", "nick", "email", "url", "biography").Where("id = ?").Update()
	_, e = Db.Exec(sql, login, nick, email, url, bio, userId)
	return e
}

func UpdatePassword(userId int, old string, new string, salt string) error {
	user := GetUserById(userId)
	if user.Id < 1 {
		return errors.New("用户无效")
	}
	if user.Password != encryptPassword(old, salt) {
		return errors.New("旧密码错误")
	}
	new = encryptPassword(new, salt)
	sql := db.NewSql("gorink_user", "password").Where("id = ?").Update()
	Db.Exec(sql, new, userId)
	return nil
}

func GetUsers() []*User {
	data, e := Orm.Find("model.User", nil)
	if e != nil {
		App.LogErr(e)
		return make([]*User, 0)
	}
	res := make([]*User, len(data))
	for i, v := range data {
		res[i] = v.(*User)
	}
	return res
}
