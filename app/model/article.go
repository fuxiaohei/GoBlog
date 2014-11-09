package model

type Article struct {
	Id         int64  `xorm:"pk autoincr"`
	Title      string `xorm:"not null"`
	Slug       string `xorm:"not null unqiue(article-slug)"`
	Content    string `xorm:"text"`
	Brief      string `xorm:"text"`
	CreateTime int64  `xorm:"not null"`
	EditTime   int64  `xorm:"default 0"`
	IsComment  bool
	UserId     int64  `xorm:"not null"`
	Template   string `xorm:"not null default '0-article'"`
	Status     string `xorm:"default 'private'"`
	Format     string `xorm:"default 'html'"`
	Comments   int    `xorm:"default 0"`
	Views      int    `xorm:"default 0"`
}
