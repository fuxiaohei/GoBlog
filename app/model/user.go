package model

type User struct {
	Id            int64  `xorm:"pk autoincr`
	Name          string `xorm:"not null"`
	Password      string `xorm:"not null"`
	Email         string `xorm:"unique not null"`
	Avatar        string
	Url           string `xorm:"default '#'"`
	Bio           string `xorm:"text"`
	CreateTime    int64
	LastLoginTime int64
	Role          string `xorm:"varchar(20) not null default 'reader'"`
	Comments      int    `xorm:"default 0"`
	Rank          int    `xorm:"default 0"`
	Active        int    `xorm:"default 11"` // 11 - active, 9 - forbidden, 7 - spam
}

type UserToken struct {
	Id         int64 `xorm:"pk autoincr`
	UserId     int64
	CreateTime int64
	ExpireTime int64
	Value      string `xorm:"unique not null`
	From       string
}
