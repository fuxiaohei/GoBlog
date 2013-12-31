package model

import (
	"fmt"
	"github.com/fuxiaohei/GoBlog/app"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"time"
)

type Session struct {
	Token      string
	Ip         string
	UserAgent  string
	CreateTime int64
	ExpireTime int64
	UserId     int
}

type SessionModel struct {
	sessions map[string]*Session
}

func (this *Session) IsValid(uid int, ip string) bool {
	if uid != this.UserId {
		return false
	}
	if ip != this.Ip {
		return false
	}
	if this.ExpireTime < time.Now().Unix() {
		return false
	}
	return true
}

// get not expired sessions.
// they are cached in memory.
func (this *SessionModel) GetAvailableSessions() []*Session {
	sql := "SELECT * FROM blog_session WHERE expire_time > ?"
	res, _ := app.Db.Query(sql, time.Now().Unix())
	sessions := make([]*Session, 0)
	res.All(&sessions)
	return sessions
}

// create new session.
// save in memory and db together.
func (this *SessionModel) CreateSession(userId int, ip string, userAgent string, expire int64) *Session {
	s := new(Session)
	s.UserAgent = userAgent
	s.Ip = ip
	s.UserId = userId
	s.CreateTime = time.Now().Unix()
	s.ExpireTime = s.CreateTime + expire
	s.Token = utils.Sha1(fmt.Sprint(userId, ip, userAgent, s.ExpireTime))
	sql := " INSERT INTO blog_session(token,ip,user_agent,create_time,expire_time,user_id) VALUES(?,?,?,?,?,?)"
	app.Db.Exec(sql, s.Token, ip, userAgent, s.CreateTime, s.ExpireTime, userId)
	this.sessions[s.Token] = s
	return s
}

// get one session by token.
func (this *SessionModel) GetByToken(token string) *Session {
	s := this.sessions[token]
	if s != nil {
		return s
	}
	sql := "SELECT * FROM blog_session WHERE token = ?"
	res, _ := app.Db.Query(sql, token)
	s = new(Session)
	res.One(s)
	if s.Token != token {
		s = nil
		return s
	}
	this.sessions[s.Token] = s
	return s
}

// recycle sessions.
// just clean the expired session in memory.
// do not clean any data in db.
func (this *SessionModel) Recycle() {
	now := time.Now().Unix()
	for _, s := range this.sessions {
		if s.ExpireTime < now {
			delete(this.sessions, s.Token)
		}
	}
}

func (this *SessionModel) reset() {
	this.sessions = make(map[string]*Session)
	sessions := this.GetAvailableSessions()
	this.sessions = make(map[string]*Session)
	for _, session := range sessions {
		this.sessions[session.Token] = session
	}
}

// create new session model.
func NewSessionModel() *SessionModel {
	s := new(SessionModel)
	go s.reset()
	return s
}
