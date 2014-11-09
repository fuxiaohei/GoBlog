package model

type Setting struct {
	Id         int64  `xorm:"pk autoincr"`
	Key        string `xorm:"unique(setting-item)`
	Value      string
	UserId     int64 `xorm:"unique(setting-item) default 0"`
	AddTime    int64
	UpdateTime int64
	IsAutoLoad bool
}
