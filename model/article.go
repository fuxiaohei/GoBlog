package model

import (
	. "github.com/fuxiaohei/gorink/app"
	"github.com/fuxiaohei/goink/db"
	"errors"
	"fmt"
	"strconv"
)

type Article struct {
	Id               int `col:"id" tbl:"gorink_article"`
	Title            string `col:"title"`
	Slug             string `col:"slug"`
	Excerpt          string `col:"excerpt"`
	Context          string `col:"context"`
	AuthorId         int `col:"author_id"`
	Author 			 *User
	CategoryId       int `col:"category_id"`
	Category 		 *Category
	CreateTime       int64 `col:"create_time"`
	EditTime         int64 `col:"edit_time"`
	StatusText       string `col:"status"`
	AllowComment     int `col:"allow_comment"`
	Comments         int `col:"comment_counts"`
	Views            int `col:"view_counts"`
}

func (this *Article) Link() string {
	return "/article/" + this.Slug + ".html"
}

func (this *Article) IsAllowComment() bool {
	return this.AllowComment == 1
}

func (this *Article) Status() string {
	if this.StatusText == "publish" {
		return "已发布"
	}
	if this.StatusText == "draft" {
		return "草稿"
	}
	if this.StatusText == "archived" {
		return "归档"
	}
	return ""
}

func init() {
	db.Define(Article{})
}

func AddArticle(a *Article) (int, error) {
	sql := db.NewSql("gorink_article", "id").Where("slug = ?").Select()
	result, e := Db.Query(sql, a.Slug)
	if e != nil {
		return -1, e
	}
	data := result.Map()
	if data != nil {
		return -1, errors.New("文章缩略名重复")
	}
	i, e := Orm.Insert(a)
	return i, e
}

func GetArticleById(id int) *Article {
	data, e := Orm.FindOne("model.Article", db.NewSql("").Where("id = ?"), id)
	if e != nil {
		App.LogErr(e)
		return nil
	}
	article := data.(*Article)
	article.Author = GetUserById(article.AuthorId)
	article.Category = GetCategoryById(article.CategoryId)
	return article
}

func UpdateArticle(a *Article) error {
	sql := db.NewSql("gorink_article", "id").Where("slug = ?").Select()
	result, e := Db.Query(sql, a.Slug)
	if e != nil {
		return e
	}
	data := result.Map()
	if data != nil {
		if len(data["id"]) > 0 && data["id"] != fmt.Sprint(a.Id) {
			return errors.New("文章缩略名重复")
		}
	}
	_, e = Orm.Update(a, "id", "title", "slug", "excerpt", "context", "category_id", "edit_time", "status", "allow_comment")
	return e
}

func GetArticleAllList(page int, size int) ([]*Article, *Counter) {
	data, e := Orm.Find("model.Article", db.NewSql("").Where("format_type != ?").Page(page, size).Order("id DESC"), "trash")
	if e != nil {
		App.LogErr(e)
		return make([]*Article, 0), nil
	}
	res := make([]*Article, len(data))
	tmpCategory := make(map[int]*Category)
	tmpUser := make(map[int]*User)
	for i, v := range data {
		a := v.(*Article)
		if tmpCategory[a.CategoryId] == nil {
			tmpCategory[a.CategoryId] = GetCategoryById(a.CategoryId)
		}
		a.Category = tmpCategory[a.CategoryId]
		if tmpUser[a.AuthorId] == nil {
			tmpUser[a.AuthorId] = GetUserById(a.AuthorId)
		}
		a.Author = tmpUser[a.AuthorId]
		res[i] = a
	}
	countRes, _ := Db.Query(db.NewSql("gorink_article").Where("format_type != ?").Count(), "trash")
	count, _ := strconv.Atoi(countRes.Map()["countNum"])
	return res, NewCounter(count, page, size)
}
