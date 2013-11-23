package controller

import (
	."github.com/fuxiaohei/gorink/app"
	"github.com/fuxiaohei/goink/app"
	"github.com/fuxiaohei/gorink/model"
	"strconv"
	"github.com/fuxiaohei/gorink/lib"
	"fmt"
	"errors"
)

func init() {
	App.GET("/admin/category", func(context *app.InkContext) interface {} {
			renderAdminPage(context, "category/manage.html", "分类", "category", map[string]interface {}{
					"Categories":model.GetCategories(),
				})
			return nil
		})
	App.GET("/admin/category/new", func(context *app.InkContext) interface {} {
			renderAdminPage(context, "category/new.html", "新建分类", "category", nil)
			return nil
		})
	App.GET("/admin/category/edit", func(context *app.InkContext) interface {} {
			id, _ := strconv.Atoi(context.String("id"))
			if id < 1 {
				context.Redirect("/admin/category", 302)
				return nil
			}
			renderAdminPage(context, "category/edit.html", "编辑分类", "category", map[string]interface {}{
					"Category":model.GetCategoryById(id),
				})
			return nil
		})
	App.POST("/admin/category/edit", func(context *app.InkContext) interface {} {
			id, _ := strconv.Atoi(context.String("id"))
			if id < 1 {
				context.Redirect(context.Refer, 302)
				return nil
			}
			msg := validateCategoryData(context)
			if len(msg) > 1 {
				renderAdminAlert(context, []error{errors.New(msg)})
				return nil
			}
			c := &model.Category{
				Id:id,
				Name:context.String("name"),
				Slug:context.String("slug"),
				Desc:context.String("desc"),
			}
			e := model.UpdateCategory(c)
			if e != nil {
				renderAdminAlert(context, []error{e})
				return nil
			}
			context.Redirect("/admin/category/edit?updated=1&id=" + context.String("id"), 302)
			return nil
		})
	App.POST("/admin/category/new", func(context *app.InkContext) interface {} {
			msg := validateCategoryData(context)
			if len(msg) > 1 {
				renderAdminAlert(context, []error{errors.New(msg)})
				return nil
			}
			c := &model.Category{}
			c.Name = context.String("name")
			c.Slug = context.String("slug")
			c.Desc = context.String("desc")
			i, e := model.AddCategory(c)
			if e != nil {
				renderAdminAlert(context, []error{e})
				return nil
			}
			context.Redirect("/admin/category/edit?id=" + fmt.Sprint(i), 302)
			return nil
		})
	App.GET("/admin/category/delete", func(context *app.InkContext) interface {} {
			id, _ := strconv.Atoi(context.String("id"))
			if id < 10 {
				renderAdminAlert(context, []error{errors.New("参数错误")})
				return nil
			}
			e := model.DeleteCategoryById(id)
			if e != nil {
				renderAdminAlert(context, []error{e})
				return nil
			}
			context.Redirect("/admin/category?deleted=1", 302)
			return nil
		})
}

func validateCategoryData(context *app.InkContext) string {
	if lib.IsEmptyString(context.String("name")) {
		return "分类名称必填"
	}
	if !lib.IsASCII(context.String("slug")) {
		return "分类缩略名不支持中文和特殊字符"
	}
	return ""
}
