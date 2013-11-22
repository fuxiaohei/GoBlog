package model

import (
	. "github.com/fuxiaohei/gorink/app"
	"github.com/fuxiaohei/goink/db"
	"fmt"
	"errors"
)

type Category struct {
	Id       int `col:"id" tbl:"gorink_category"`
	Name     string `col:"name"`
	Slug     string `col:"slug"`
	Desc     string `col:"desc"`
	Articles int `col:"articles"`
	Parent   int `col:"parent"`
}

func (this *Category) Link() string {
	return "/category/" + this.Slug + "/"
}

func init() {
	db.Define(Category{})
	App.On("model.article.add@category_update", func() {
			go UpdateCategoryArticles()
		});
	App.On("model.article.update@category_update", func() {
			go UpdateCategoryArticles()
		});
}

func GetCategories() []*Category {
	data, e := Orm.Find("model.Category", db.NewSql("").Where("id > 10").Where("type = ?"), "category")
	if e != nil {
		App.LogErr(e)
		return make([]*Category, 0)
	}
	res := make([]*Category, len(data))
	for i, v := range data {
		res[i] = v.(*Category)
	}
	return res
}

func GetCategoryById(id int) *Category {
	data, e := Orm.FindOne("model.Category", db.NewSql("").Where("id = ?"), id)
	if e != nil {
		App.LogErr(e)
		return nil
	}
	return data.(*Category)
}

func UpdateCategory(c *Category) error {
	sql := db.NewSql("gorink_category", "id").Where("slug = ?").Select()
	result, e := Db.Query(sql, c.Slug)
	if e != nil {
		return e
	}
	data := result.Map()
	if data != nil {
		if len(data["id"]) > 0 && data["id"] != fmt.Sprint(c.Id) {
			return errors.New("缩略名重复")
		}
	}
	_, e = Orm.Update(c, "id", "name", "slug", "desc")
	return e
}

func AddCategory(c *Category) (int, error) {
	sql := db.NewSql("gorink_category", "id").Where("slug = ?").Select()
	result, e := Db.Query(sql, c.Slug)
	if e != nil {
		return -1, e
	}
	data := result.Map()
	if data != nil {
		return -1, errors.New("缩略名重复")
	}
	i, e := Orm.Insert(c)
	return i, e
}

func DeleteCategoryById(id int) error {
	sql := db.NewSql("gorink_category").Where("id = ?").Delete()
	_ , e := Db.Exec(sql, id)
	return e
}

func UpdateCategoryArticles() {
	sql := "UPDATE gorink_category SET articles = (SELECT count(*) FROM gorink_article WHERE gorink_article.category_id = gorink_category.id) AND gorink_article.format_type != 'trash';"
	Db.Exec(sql)
}
