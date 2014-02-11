package model

import (
	"errors"
	"github.com/fuxiaohei/GoBlog/app/utils"
)

var (
	users     []*User
	userMaxId int
)

type User struct {
	Id            int
	Name          string
	Password      string
	Nick          string
	Email         string
	Avatar        string
	Url           string
	Bio           string
	CreateTime    int64
	LastLoginTime int64
	Role          string
}

// check user password.
func (u *User) CheckPassword(pwd string) bool {
	return utils.Sha1(pwd+"xxxxx") == u.Password
}

// change user email.
// check unique.
func (u *User) ChangeEmail(email string) bool {
	u2 := GetUserByEmail(u.Email)
	if u2.Id != u.Id {
		return false
	}
	u.Email = email
	return true
}

// change user password.
func (u *User) ChangePassword(pwd string) {
	u.Password = utils.Sha1(pwd + "xxxxx")
}

// get a user by given id.
func GetUserById(id int) *User {
	for _, u := range users {
		if u.Id == id {
			return u
		}
	}
	return nil
}

// get a user by given name.
func GetUserByName(name string) *User {
	for _, u := range users {
		if u.Name == name {
			return u
		}
	}
	return nil
}

// get a user by given email.
func GetUserByEmail(email string) *User {
	for _, u := range users {
		if u.Email == email {
			return u
		}
	}
	return nil
}

// get users of given role.
func GetUsersByRole(role string) []*User {
	us := make([]*User, 0)
	for _, u := range users {
		if u.Role == role {
			us = append(us, u)
		}
	}
	return us
}

// create new user.
func CreateUser(u *User) error {
	if GetUserByName(u.Email) != nil {
		return errors.New("email-repeat")
	}
	userMaxId += Storage.TimeInc(5)
	u.Id = userMaxId
	u.CreateTime = utils.Now()
	u.LastLoginTime = u.CreateTime
	users = append(users, u)
	go SyncUsers()
	return nil
}

// remove a user.
func RemoveUser(u *User) {
	for i, u2 := range users {
		if u2.Id == u.Id {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}
	go SyncUsers()
}

// write users to json.
func SyncUsers() {
	Storage.Set("users", users)
}

func LoadUsers() {
	users = make([]*User, 0)
	userMaxId = 0
	Storage.Get("users", &users)
	for _, u := range users {
		if u.Id > userMaxId {
			userMaxId = u.Id
		}
	}
}
