package model

import (
	"fmt"
	"github.com/fuxiaohei/GoBlog/app"
	"github.com/fuxiaohei/GoBlog/app/utils"
	//"github.com/fuxiaohei/GoInk/Db"
	"strconv"
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

type CommentContent struct{
	Title string
	Link  string
}

func (this *Comment) ContentNode() *CommentContent {
	article := ArticleM.GetArticleById(this.ContentId)
	if article != nil {
		c := new(CommentContent)
		c.Title = article.Title
		if article.Status == "publish" {
			c.Link = article.Link()
		}
		return c
	}
	return nil
}

func (this *Comment) Parent() *Comment {
	if this.Pid < 1 {
		return nil
	}
	return CommentM.GetCommentById(this.Pid)
}

type CommentModel struct {
	commentsCache   map[int]*Comment
	pagedCache map[string][]*Comment
	pagerCache map[string]int
	emailCountCache map[string]int
}

func (this *CommentModel) GetCommentById(id int) *Comment {
	if this.commentsCache[id] == nil {
		sql := "SELECT * FROM blog_comment WHERE type = ? AND id = ?"
		res, _ := app.Db.Query(sql, "comment", id)
		c := new(Comment)
		res.One(c)
		if c.Id < 1 {
			return nil
		}
		this.cacheComment(c)
	}
	return this.commentsCache[id]
}

func (this *CommentModel) cacheComment(comments... *Comment) {
	for _, c := range comments {
		this.commentsCache[c.Id] = c
	}
}

// get paged comments.
func (this *CommentModel) GetPaged(page, size int, onlySpam bool) ([]*Comment, *Pager) {
	key := fmt.Sprintf("%d-%d-spam-%t", page, size, onlySpam)
	if this.pagedCache[key] == nil {
		sql := "SELECT * FROM blog_comment WHERE type = ?"
		args := []interface{}{"comment"}
		limit := (page-1) * size
		if onlySpam {
			sql += " AND status = ?"
			args = append(args, "spam")
		}
		sql += " ORDER BY id DESC LIMIT " + fmt.Sprintf("%d,%d", limit, size)
		res, e := app.Db.Query(sql, args...)
		if len(res.Data) > 0 && e == nil {
			comments := make([]*Comment, 0)
			res.All(&comments)
			this.pagedCache[key] = comments
			this.cacheComment(comments...)
		}
	}
	pagerKey := fmt.Sprintf("counter-spam-%t", onlySpam)
	if this.pagerCache[pagerKey] == 0 {
		sql := "SELECT count(*) AS c FROM blog_comment WHERE type = ?"
		args := []interface{}{"comment"}
		if onlySpam {
			sql += " AND status = ?"
			args = append(args, "spam")
		}
		res, e := app.Db.Query(sql, args...)
		if e != nil {
			return nil, nil
		}
		all, _ := strconv.Atoi(res.Data[0]["c"])
		this.pagerCache[pagerKey] = all
	}
	return this.pagedCache[key], newPager(page, size, this.pagerCache[pagerKey])
}

func (this *CommentModel) GetAllOfContent(contentId int, noDraft bool) []*Comment {
	key := fmt.Sprintf("content-%d-spam-%t", contentId, noDraft)
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
			this.cacheComment(comments...)
		}
	}
	return this.pagedCache[key]
}

func (this *CommentModel) CreateComment(c *Comment) *Comment {
	sql := "INSERT INTO blog_comment(author,email,site,avatar,create_time,content,content_id,user_id,pid,is_admin,type,status) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)"
	var status string
	if c.IsAdmin {
		status = "approved"
	}else {
		status = this.getEmailStatus(c.Email)
	}
	res, _ := app.Db.Exec(sql, c.Author, c.Email, c.Site, c.Avatar, utils.Now(), c.Content, c.ContentId, c.UserId, c.Pid, false, "comment", status)
	this.Reset()
	if res.LastInsertId < 1 {
		return nil
	}
	return this.GetCommentById(res.LastInsertId)
}

func (this *CommentModel) getEmailStatus(email string) string {
	_, ok := this.emailCountCache[email]
	if ok {
		if this.emailCountCache[email] > 0 {
			return "approved"
		}
		return "spam"
	}
	sql := "SELECT count(*) AS c FROM blog_comment WHERE email = ? AND type = ? AND status = ?"
	res, _ := app.Db.Query(sql, email, "comment", "approved")
	if len(res.Data) > 0 {
		all, _ := strconv.Atoi(res.Data[0]["c"])
		this.emailCountCache[email] = all
		if all > 0 {
			return "approved"
		}
		return "spam"
	}
	return "spam"
}

func (this *CommentModel) DeleteComment(id int) {
	sql := " DELETE FROM blog_comment WHERE id = ?"
	app.Db.Exec(sql, id)
	this.Reset()
}

func (this *CommentModel) DeleteCommentsInContent(contentId int) {
	sql := "DELETE FROM blog_comment WHERE content_id = ?"
	app.Db.Exec(sql, contentId)
	this.Reset()
}

func (this *CommentModel) ChangeCommentStatus(id int, status string) {
	sql := "UPDATE blog_comment SET status = ? WHERE id = ?"
	app.Db.Exec(sql, status, id)
	this.Reset()
}

func (this *CommentModel) Reset() {
	this.pagedCache = make(map[string][]*Comment)
	this.pagerCache = make(map[string]int)
	this.commentsCache = make(map[int]*Comment)
	this.emailCountCache = make(map[string]int)
}

func NewCommentModel() *CommentModel {
	c := new(CommentModel)
	c.Reset()
	return c
}
