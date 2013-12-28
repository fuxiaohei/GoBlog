package model

import (
	"fmt"
	"github.com/fuxiaohei/GoBlog/app"
	"strconv"
)

// article object
type Article struct {
	Id         int
	Title      string
	Slug       string
	Summary    string
	Content    string
	CreateTime int64
	EditTime   int64
	CategoryId int
	AuthorId   int
	Format     string
	Status     string
	IsComment  int
	IsFeed     int
	Comments   int
	Views      int
}

// return article slug link.
func (this *Article) Link() string {
	return "/article/" + fmt.Sprintf("%d/%s.html", this.Id, this.Slug)
}

// return article author *User
func (this *Article) Author() *User {
	return UserM.GetUserById(this.AuthorId)
}

// return article *Category
func (this *Article) Category() *Category {
	return CategoryM.GetCategoryById(this.CategoryId)
}

// article model
type ArticleModel struct {
	article    map[string]*Article
	idIndex    map[int]string
	pagedCache map[string][]*Article
	pagerCache map[string]int
}

// cache one article in model.
func (this *ArticleModel) cacheArticle(a *Article) {
	if a == nil {
		return
	}
	this.article[a.Slug] = a
	this.idIndex[a.Id] = a.Slug
}

// remove cached one article in model.
func (this *ArticleModel) nocacheArticle(a *Article) {
	if a == nil {
		return
	}
	delete(this.article, a.Slug)
	delete(this.idIndex, a.Id)
}

// get one article by slug string.
func (this *ArticleModel) GetArticleBySlug(slug string) *Article {
	a := this.article["slug"]
	if a == nil {
		sql := "SELECT * FROM blog_content WHERE type = ? AND slug = ?"
		res, _ := app.Db.Query(sql, "article", slug)
		a = new(Article)
		res.One(a)
		if a.Slug != slug {
			return nil
		}
		this.cacheArticle(a)
	}
	return a
}

// get one article by given id.
func (this *ArticleModel) GetArticleById(id int) *Article {
	slug := this.idIndex[id]
	if slug != "" {
		return this.GetArticleBySlug(slug)
	}
	sql := "SELECT * FROM blog_content WHERE type = ? AND id = ?"
	res, _ := app.Db.Query(sql, "article", id)
	a := new(Article)
	res.One(a)
	if a.Id != id {
		return nil
	}
	this.cacheArticle(a)
	return a
}

// remove all paged article slice cache.
func (this *ArticleModel) nocachePaged() {
	this.pagedCache = make(map[string][]*Article)
	this.pagerCache = make(map[string]int)
}

// get paged article.
func (this *ArticleModel) GetPaged(page, size int, noDraft bool) ([]*Article, *Pager) {
	key := fmt.Sprintf("%d-%d-draft-%t", page, size, noDraft)
	if this.pagedCache[key] == nil {
		sql := "SELECT * FROM blog_content WHERE type = ?"
		args := []interface{}{"article"}
		limit := (page - 1) * size
		if noDraft {
			sql += " AND status != ?"
			args = append(args, "draft")
		}
		sql += " ORDER BY id DESC LIMIT "+fmt.Sprintf("%d,%d", limit, size)
		res, e := app.Db.Query(sql, args...)
		if len(res.Data) > 0 && e == nil {
			articles := make([]*Article, 0)
			res.All(&articles)
			this.pagedCache[key] = articles
		}
	}
	pagerKey := fmt.Sprintf("counter-draft-%t", noDraft)
	if this.pagerCache[pagerKey] == 0 {
		sql := "SELECT count(*) AS c FROM blog_content WHERE type = ?"
		args := []interface{}{"article"}
		if noDraft {
			sql += " AND status != ?"
			args = append(args, "draft")
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

func (this *ArticleModel) GetCategoryPaged(categoryId, page, size int, noDraft bool) ([]*Article, *Pager) {
	key := fmt.Sprintf("%d-%d-draft-%t-category-%d", page, size, noDraft, categoryId)
	if this.pagedCache[key] == nil {
		sql := "SELECT * FROM blog_content WHERE type = ? AND category_id = ?"
		args := []interface{}{"article", categoryId}
		limit := (page - 1) * size
		if noDraft {
			sql += " AND status != ?"
			args = append(args, "draft")
		}
		sql += " ORDER BY id DESC LIMIT "+fmt.Sprintf("%d,%d", limit, size)
		res, e := app.Db.Query(sql, args...)
		if len(res.Data) > 0 && e == nil {
			articles := make([]*Article, 0)
			res.All(&articles)
			this.pagedCache[key] = articles
		}
	}
	pagerKey := fmt.Sprintf("counter-draft-%t-category-%d", noDraft, categoryId)
	if this.pagerCache[pagerKey] == 0 {
		sql := "SELECT count(*) AS c FROM blog_content WHERE type = ? AND category_id = ?"
		args := []interface{}{"article", categoryId}
		if noDraft {
			sql += " AND status != ?"
			args = append(args, "draft")
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

func (this *ArticleModel) GetPopular(size int) []*Article {
	key := fmt.Sprintf("popular-%d", size)
	if this.pagedCache[key] == nil {
		sql := "SELECT * FROM blog_content WHERE type = ? AND status = ? ORDER BY comments DESC LIMIT ?"
		res, e := app.Db.Query(sql, "article", "publish", size)
		if e != nil {
			return nil
		}
		articles := make([]*Article, 0)
		res.All(&articles)
		this.pagedCache[key] = articles
	}
	return this.pagedCache[key]
}

func (this *ArticleModel) SaveArticle(article *Article) *Article {
	sql := "UPDATE blog_content SET title = ?,slug = ?,summary = ?,content = ?,edit_time = ?,category_id = ?,status = ?,is_comment = ?,is_feed = ? WHERE id = ?"
	res, e := app.Db.Exec(sql, article.Title, article.Slug, article.Summary, article.Content, article.EditTime, article.CategoryId, article.Status, article.IsComment, article.IsFeed, article.Id)
	if e != nil {
		return nil
	}
	if res.AffectedRows > 0 {
		// clean all cache.
		this.nocachePaged()
		this.nocacheArticle(article)
		return this.GetArticleById(article.Id)
	}
	return nil
}

func (this *ArticleModel) CreateArticle(article *Article) *Article {
	sql := "INSERT INTO blog_content(title,slug,summary,content,create_time,edit_time,category_id,author_id,type,format,status,is_comment,is_feed) "
	sql += "VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)"
	res, e := app.Db.Exec(sql, article.Title, article.Slug, article.Summary, article.Content, article.CreateTime, article.EditTime, article.CategoryId, article.AuthorId, "article", "md", article.Status, article.IsComment, article.IsFeed)
	if e != nil {
		return nil
	}
	if res.LastInsertId > 0 {
		// clean all cache.
		this.nocachePaged()
		return this.GetArticleById(res.LastInsertId)
	}
	return nil
}

func NewArticleModel() *ArticleModel {
	articleM := new(ArticleModel)
	articleM.article = make(map[string]*Article)
	articleM.idIndex = make(map[int]string)
	articleM.nocachePaged()
	return articleM
}
