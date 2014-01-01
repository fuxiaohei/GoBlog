package model

import (
	"github.com/fuxiaohei/GoBlog/app"
)

type Category struct {
	Id          int
	Name        string
	Slug        string
	Description string
	Counts      int
}

func (this *Category) Link() string {
	return "/category/"+this.Slug + "/"
}

type CategoryModel struct {
	categoriesCache       map[string]*Category
	idIndexCache          map[int]string
	countsDescResult []*Category
}

// get all categories. no ordered.
func (this *CategoryModel) GetAll() []*Category {
	// get from memory cache.
	if len(this.categoriesCache) > 0 {
		c := make([]*Category, len(this.categoriesCache))
		i := 0
		for _, ca := range this.categoriesCache {
			c[i] = ca
			i++
		}
		return c
	}
	sql := "SELECT * FROM blog_meta WHERE type = ?"
	res, _ := app.Db.Query(sql, "category")
	c := make([]*Category, 0)
	res.All(&c)
	this.cacheCategory(c...)
	return c
}

// get categories ordered by counts.
// the expire time means how long to cache the ordered result.
// if the last create time + expire time is less than now, create new result.
// so that it is afford to be affected when counts changed.
func (this *CategoryModel) GetCountsDesc(expire int64) []*Category {
	if len(this.countsDescResult) > 0 {
		return this.countsDescResult
	}
	sql := "SELECT * FROM blog_meta WHERE type = ? ORDER BY counts DESC"
	res, _ := app.Db.Query(sql, "category")
	c := make([]*Category, 0)
	res.All(&c)
	this.countsDescResult = c
	return c
}

// get one category by slug string.
func (this *CategoryModel) GetCategoryBySlug(slug string) *Category {
	c := this.categoriesCache[slug]
	if c != nil {
		return c
	}
	sql := "SELECT * FROM blog_meta WHERE type = ? AND slug = ?"
	res, _ := app.Db.Query(sql, "category", slug)
	c = new(Category)
	res.One(c)
	if c.Slug != slug {
		return nil
	}
	this.cacheCategory(c)
	return c
}

// get one category by id int.
func (this *CategoryModel) GetCategoryById(id int) *Category {
	slug := this.idIndexCache[id]
	if slug != "" {
		return this.GetCategoryBySlug(slug)
	}
	sql := "SELECT * FROM blog_meta WHERE type = ? AND id = ?"
	res, _ := app.Db.Query(sql, "category", id)
	c := new(Category)
	res.One(c)
	if c.Id != id {
		return nil
	}
	this.cacheCategory(c)
	return c
}

// cache category if it's not in memory.
func (this *CategoryModel) cacheCategory(categories... *Category) {
	for _, c := range categories {
		if c == nil {
			return
		}
		this.categoriesCache[c.Slug] = c
		this.idIndexCache[c.Id] = c.Slug
	}
}

// create a new category.
// cache it in memory automatically.
func (this *CategoryModel) CreateCategory(name string, slug string, desc string) *Category {
	sql := "INSERT INTO blog_meta(name,slug,description,type) VALUES(?,?,?,?)"
	res, _ := app.Db.Exec(sql, name, slug, desc, "category")
	this.Reset()
	return this.GetCategoryById(res.LastInsertId)
}

// save a category with new data.
// cache it right now.
func (this *CategoryModel) SaveCategory(id int, name string, slug string, desc string) *Category {
	// update and re-cache it
	sql := "UPDATE blog_meta SET name = ?,slug = ?,description = ? WHERE id = ?"
	app.Db.Exec(sql, name, slug, desc, id)
	this.Reset()
	return this.GetCategoryById(id)
}

// delete one category.
// move relationship to new category id.
func (this *CategoryModel) DeleteCategory(id int, move int) {
	sql := "UPDATE blog_content SET category_id = ? WHERE category_id = ?"
	app.Db.Exec(sql, move, id)
	sql = "DELETE FROM blog_meta WHERE id = ?"
	app.Db.Exec(sql, id)
	this.Reset()
	this.CountArticle()
}

func (this *CategoryModel) CountArticle() {
	sql := "UPDATE blog_meta SET counts = (SELECT count(*) FROM blog_content WHERE blog_content.category_id = blog_meta.id AND blog_content.status = 'publish' AND blog_content.type = 'article')"
	app.Db.Exec(sql)
	this.Reset()
}

func (this *CategoryModel) Reset() {
	this.categoriesCache = make(map[string]*Category)
	this.idIndexCache = make(map[int]string)
	categories := this.GetAll()
	for _, ca := range categories {
		this.categoriesCache[ca.Slug] = ca
		this.idIndexCache[ca.Id] = ca.Slug
	}
	this.countsDescResult = make([]*Category, 0)
}

// create new category model.
func NewCategoryModel() *CategoryModel {
	c := new(CategoryModel)
	go c.Reset()
	return c
}
