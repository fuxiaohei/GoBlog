package model

type Category struct {
	Id          int64  `xorm:"pk autoincr"`
	Name        string `xorm:"not null"`
	Slug        string `xorm:"not null"`
	Description string
	UserId      int `xorm:"not null"`
	Articles    int `xorm:"default 0"`
}

type ArticleCategory struct {
	ArticleId  int64 `xorm:"not null unique(article-category)"`
	CategoryId int64 `xorm:"not null unique(article-category)"`
}
