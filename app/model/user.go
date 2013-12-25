package model

import (
	"github.com/fuxiaohei/GoBlog/app"
	"errors"
	"github.com/fuxiaohei/GoBlog/app/utils"
)

type User struct {
	Id         int
	Login      string
	Password   string
	Display    string
	Email      string
	Site       string
	Avatar     string
	CreateTime int64
	LoginTime  int64
	LoginIp    string
	RoleName   string
	Status     string
}

type UserModel struct {
	users map[int]*User
	loginIndex map[string]int
}

// get all users.
// its result is cached.
func (this *UserModel) GetAllUser() map[int]*User {
	if len(this.users) > 0 {
		return this.users
	}
	sql := "SELECT * FROM blog_user"
	result, _ := app.Db.Query(sql)
	users := make([]*User, 0)
	result.All(&users)
	this.users = make(map[int]*User)
	for _, u := range users {
		this.users[u.Id] = u
	}
	return this.users
}

// generate an index for login/id pair.
func (this *UserModel) generateLoginIndex() {
	if len(this.users) < 1 {
		return
	}
	this.loginIndex = make(map[string]int)
	for _, u := range this.users {
		this.loginIndex[u.Login] = u.Id
	}
}

// get one user by id.
// if no cached, query from db and cache it.
func (this *UserModel) GetUserById(id int) *User {
	if this.users[id] != nil {
		return this.users[id]
	}
	sql := "SELECT * FROM blog_user WHERE id = ?"
	res, _ := app.Db.Query(sql, id)
	u := new(User)
	res.One(u)
	// cache it
	if u.Id == id {
		this.cacheUser(u)
	}else {
		u = nil
	}
	return u
}

// get one user by login name.
// if no cached, query from db and cache it.
func (this *UserModel) GetUserByLogin(login string) *User {
	// get from login index map
	id := this.loginIndex[login]
	if id > 0 {
		return this.GetUserById(id)
	}
	// query db
	sql := "SELECT * FROM blog_user WHERE login = ?"
	res, _ := app.Db.Query(sql, login)
	u := new(User)
	res.One(u)
	// cache it
	if u.Login == login {
		this.cacheUser(u)
	}else {
		u = nil
	}
	return u
}

// cache user struct.
// put into map and login index map.
func (this *UserModel) cacheUser(u *User) {
	if u == nil {
		return
	}
	this.users[u.Id] = u
	this.loginIndex[u.Login] = u.Id
}

// no cache user struct.
// delete in map and index slice.
func (this *UserModel) nocacheUser(u *User) {
	if u == nil {
		return
	}
	delete(this.users, u.Id)
	delete(this.loginIndex, u.Login)
}

// save user profile.
// update memory cache struct immediately.
// save user data to db in routine.
// return the updated user struct
func (this *UserModel) SaveProfile(id int, login string, display string, email string, site string, avatar string) *User {
	this.nocacheUser(this.GetUserById(id))
	// save to db
	sql := "UPDATE blog_user SET login = ?,display = ?,email = ?,site = ?,avatar = ? WHERE id = ?"
	app.Db.Exec(sql, login, display, email, site, avatar, id)
	return this.GetUserById(id)
}

// save new password.
// if no user or wrong old password, return error.
func (this *UserModel) SavePassword(id int, old string, new string) error {
	user := this.GetUserById(id)
	if user == nil {
		return errors.New("无效的用户")
	}
	if user.Password != utils.Sha1(old) {
		return errors.New("旧密码错误")
	}
	user.Password = utils.Sha1(new)
	sql := "UPDATE blog_user SET password = ? WHERE id = ?"
	app.Db.Exec(sql, user.Password, id)
	return nil
}

func (this *UserModel) reset() {
	this.users = make(map[int]*User)
	this.loginIndex = make(map[string]int)
	this.GetAllUser()
	this.generateLoginIndex()
}

// create new userModel.
// load all user data for caching.
func NewUserModel() *UserModel {
	userM := new(UserModel)
	go userM.reset()
	return userM
}
