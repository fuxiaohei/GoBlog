package model

import (
	"fmt"
	"github.com/fuxiaohei/GoBlog/app"
	"github.com/fuxiaohei/GoBlog/app/utils"
	//"github.com/fuxiaohei/GoInk/Db"
)

type Comment struct {
	Id         int
	Author     string
	Email      string `json:"-"`
	Site       string
	Avatar     string
	CreateTime int64
	Content    string
	ContentId  int
	UserId     int
	Pid        int
	IsAdmin    bool
	Status     string
}

type CommentModel struct {
	comments   map[int]*Comment
	pagedCache map[string][]*Comment
	pagerCache map[string]int
	emailCount map[string]int
}

func (this *CommentModel) GetCommentById(id int) *Comment {
	if this.comments[id] == nil {
		sql := "SELECT * FROM blog_comment WHERE type = ? AND id = ?"
		res, _ := app.Db.Query(sql, "comment", id)
		c := new(Comment)
		res.One(c)
		if c.Id < 1 {
			return nil
		}
		this.cacheComment(c)
	}
	return this.comments[id]
}

func (this *CommentModel) cacheComment(c *Comment) {
	this.comments[c.Id] = c
}

func (this *CommentModel) nocacheComment(c *Comment) {
	delete(this.comments, c.Id)
}

func (this *CommentModel) nocachePaged() {
	this.pagedCache = make(map[string][]*Comment)
	this.pagerCache = make(map[string]int)
}

func (this *CommentModel) GetAllOfContent(contentId int, noDraft bool) []*Comment {
	key := fmt.Sprintf("content-%d-draft-%t", contentId, noDraft)
	if len(this.pagedCache[key]) < 1 {
		sql := "SELECT * FROM blog_comment WHERE type = ? AND content_id = ?"
		args := []interface{}{"comment", contentId}
		if noDraft {
			sql += " AND status = ?"
			args = append(args, "approved")
		}
		sql += " ORDER BY id ASC"
		res, e := app.Db.Query(sql, args...)
		if len(res.Data) > 0 && e == nil {
			comments := make([]*Comment, 0)
			res.All(&comments)
			this.pagedCache[key] = comments
		}
	}
	return this.pagedCache[key]
}

func (this *CommentModel) CreateComment(c *Comment) *Comment {
	sql := "INSERT INTO blog_comment(author,email,site,avatar,create_time,content,content_id,user_id,pid,is_admin,type,status) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)"
	c.CreateTime = utils.Now()
	res, _ := app.Db.Exec(sql, c.Author, c.Email, c.Site, c.Avatar, c.CreateTime, c.Content, c.ContentId, c.UserId, c.Pid, false, "comment", "spam")
	if res.LastInsertId < 1 {
		return nil
	}
	c.Id = res.LastInsertId
	this.nocachePaged()
	return c
}

func NewCommentModel() *CommentModel {
	c := new(CommentModel)
	c.comments = make(map[int]*Comment)
	c.nocachePaged()
	return c
}
