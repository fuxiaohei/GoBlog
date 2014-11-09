package model

type Tag struct {
	Id       int64  `xorm:"pk autoincr"`
	Name     string `xorm:"not null"`
	Slug     string `xorm:"not null"`
	UserId   int    `xorm:"not null"`
	Articles int    `xorm:"default 0"`
}

type ArticleTag struct {
	ArticleId int64 `xorm:"not null unique(article-tag)"`
	TagId     int64 `xorm:"not null unique(article-tag)"`
}
