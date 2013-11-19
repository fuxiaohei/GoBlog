package model

import (
	"crypto/sha1"
	"fmt"
	"io"
	"time"
)

var tokens map[string]*LoginResult

func init() {
	tokens = make(map[string]*LoginResult)
	loginGc()
}

type LoginResult struct {
	UserId int
	Name   string
	Token  string
	expire int64
	Avatar string
}

func Login(user string, password string, salt string) (*LoginResult, string) {
	userData := GetUserByLogin(user)
	if userData.Id < 1 {
		return nil, "用户无效"
	}
	if userData.Role != "super" || userData.Status != "activate" {
		return nil, "用户失效"
	}
	if userData.Password != encryptPassword(password, salt) {
		return nil, "密码错误"
	}
	loginResult := &LoginResult{}
	loginResult.UserId = userData.Id
	loginResult.Name = userData.DisplayName()
	loginResult.expire = time.Now().Unix() + 3600*24*7
	loginResult.Token = encryptPassword(userData.Password, fmt.Sprint(loginResult.expire))
	loginResult.Avatar = GetGravatar(userData.Email, "40")
	tokens[loginResult.Token] = loginResult
	return loginResult, ""
}

func CheckLoginByToken(token string) (bool, *LoginResult) {
	tokenData := tokens[token]
	if tokenData == nil {
		return false, nil
	}
	if tokenData.expire < 1 {
		return false, nil
	}
	if tokenData.expire < time.Now().Unix() {
		return false, nil
	}
	return true, tokenData
}

func LoginStats() int {
	return len(tokens)
}

func encryptPassword(raw string, salt string) string {
	m := sha1.New()
	io.WriteString(m, raw)
	io.WriteString(m, salt)
	return fmt.Sprintf("%x", m.Sum(nil))
}

func loginGc() {
	if len(tokens) > 0 {
		for k, token := range tokens {
			if token.expire < time.Now().Unix() {
				delete(tokens, k)
			}
		}
	}
	time.AfterFunc(time.Duration(3600)*time.Second, func() {
			loginGc()
		})
}
