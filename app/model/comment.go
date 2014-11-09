package model

type Comment struct {
	Id         int64  `xorm:"pk autoincr"`
	Author     string `xorm:"not null"`
	Email      string `xorm:"not null"`
	Url        string
	Avatar     string
	Content    string `xorn:"text"`
	CreateTime int64
	ParentId   int64  `xorm:"not null default 0"`
	Status     string `xorm:"not null default 'wait'`
	Ip         string
	UserAgent  string
	IsAdmin    bool
}

type CommentArticle struct {
	ArticleId int64 `xorm:"not null unique(comment-article)"`
	CommentId int64 `xorm:"not null unique(comment-article)"`
}

type CommentPage struct {
	PageId    int64 `xorm:"not null unique(comment-page)"`
	CommentId int64 `xorm:"not null unique(comment-page)"`
}
