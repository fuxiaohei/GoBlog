package model

type Message struct {
	Id         int64 `xorm:"pk autoincr"`
	UserId     int64
	Content    string
	Type       string `xorm:"varchar(20)"`
	CreateTime int64
	ReadTime   int64
}
