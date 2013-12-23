package model

import (
	"github.com/fuxiaohei/GoBlog/app"
	"time"
)

type Category struct {
	Id          int
	Name        string
	Slug        string
	Description string
	Counts      int
}

type CategoryModel struct {
	categories       map[string]*Category
	idIndex          map[int]string
	countsDescResult []*Category
	countsDescExpire int64
}

// get all categories. no ordered.
func (this *CategoryModel) GetAll() []*Category {
	// get from memory cache.
	if len(this.categories) > 0 {
		c := make([]*Category, len(this.categories))
		i := 0
		for _, ca := range this.categories {
			c[i] = ca
			i++
		}
		return c
	}
	sql := "SELECT * FROM blog_meta WHERE type = ?"
	res, _ := app.Db.Query(sql, "category")
	c := make([]*Category, 0)
	res.All(&c)
	return c
}

// get categories ordered by counts.
// the expire time means how long to cache the ordered result.
// if the last create time + expire time is less than now, create new result.
// so that it is afford to be affected when counts changed.
func (this *CategoryModel) GetCountsDesc(expire int64) []*Category {
	now := time.Now().Unix()
	if this.countsDescExpire+expire > now {
		return this.countsDescResult
	}
	sql := "SELECT * FROM blog_meta WHERE type = ? ORDER BY counts DESC"
	res, _ := app.Db.Query(sql, "category")
	c := make([]*Category, 0)
	res.All(&c)
	this.countsDescExpire = now
	this.countsDescResult = c
	return c
}

// get one category by slug string.
func (this *CategoryModel) GetCategoryBySlug(slug string) *Category {
	c := this.categories[slug]
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
	slug := this.idIndex[id]
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
func (this *CategoryModel) cacheCategory(c *Category) {
	if c == nil {
		return
	}
	this.categories[c.Slug] = c
	this.idIndex[c.Id] = c.Slug
}

// remove one cached category in memory map or index slice.
func (this *CategoryModel) nocacheCategory(c *Category) {
	if c == nil {
		return
	}
	delete(this.categories, c.Slug)
	delete(this.idIndex, c.Id)
}

// create a new category.
// cache it in memory automatically.
func (this *CategoryModel) CreateCategory(name string, slug string, desc string) *Category {
	sql := "INSERT INTO blog_meta(name,slug,description,type) VALUES(?,?,?,?)"
	res, _ := app.Db.Exec(sql, name, slug, desc, "category")
	this.countsDescExpire = 0
	return this.GetCategoryById(res.LastInsertId)
}

// save a category with new data.
// cache it right now.
func (this *CategoryModel) SaveCategory(id int, name string, slug string, desc string) *Category {
	this.nocacheCategory(this.GetCategoryById(id))
	// update and re-cache it
	sql := "UPDATE blog_meta SET name = ?,slug = ?,description = ? WHERE id = ?"
	app.Db.Exec(sql, name, slug, desc, id)
	this.countsDescExpire = 0
	return this.GetCategoryById(id)
}

// delete one category.
// move relationship to new category id.
func (this *CategoryModel) DeleteCategory(id int, move int) {
	sql := "UPDATE blog_relationship SET meta_id = ? WHERE meta_id = ?"
	app.Db.Exec(sql, move, id)
	this.nocacheCategory(this.GetCategoryById(id))
	sql = "DELETE FROM blog_meta WHERE id = ?"
	this.countsDescExpire = 0
	app.Db.Exec(sql, id)
}

// create new category model.
func NewCategoryModel() *CategoryModel {
	c := new(CategoryModel)
	c.categories = make(map[string]*Category)
	c.idIndex = make(map[int]string)
	go func() {
		categories := c.GetAll()
		for _, ca := range categories {
			c.categories[ca.Slug] = ca
			c.idIndex[ca.Id] = ca.Slug
		}
	}()
	return c
}
