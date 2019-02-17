package model

import (
	"fmt"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"github.com/fuxiaohei/GoInk"
)

var tokens map[string]*Token

type Token struct {
	Value      string
	UserId     int
	CreateTime int64
	ExpireTime int64
}

// IsValid checks token is valid or expired.
func (t *Token) IsValid() bool {
	if GetUserById(t.UserId) == nil {
		return false
	}
	return t.ExpireTime > utils.Now()
}

// CreateToken: create new token from user and context.
func CreateToken(u *User, context *GoInk.Context, expire int64) *Token {
	t := new(Token)
	t.UserId = u.Id
	t.CreateTime = utils.Now()
	t.ExpireTime = t.CreateTime + expire
	t.Value = utils.Sha1(fmt.Sprintf("%s-%s-%d-%d", context.Ip, context.UserAgent, t.CreateTime, t.UserId))
	tokens[t.Value] = t
	go SyncTokens()
	return t
}

// GetTokenByValue gets token by token value.
func GetTokenByValue(v string) *Token {
	return tokens[v]
}

// GetTokensByUser gets tokens of given user.
func GetTokensByUser(u *User) []*Token {
	ts := make([]*Token, 0)
	for _, t := range tokens {
		if t.UserId == u.Id {
			ts = append(ts, t)
		}
	}
	return ts
}

// RemoveToken removes a token by token value.
func RemoveToken(v string) {
	delete(tokens, v)
	go SyncTokens()
}

// CleanTokens: clean all expired tokens in memory.
// do not write to json.
func CleanTokens() {
	for k, t := range tokens {
		if !t.IsValid() {
			delete(tokens, k)
		}
	}
}

// SyncTokens writes tokens to json.
// it calls CleanTokens before writing.
func SyncTokens() {
	CleanTokens()
	Storage.Set("tokens", tokens)
}

// LoadTokens: load all tokens from json.
func LoadTokens() {
	tokens = make(map[string]*Token)
	Storage.Get("tokens", &tokens)
}
