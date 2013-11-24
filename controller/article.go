package controller

import (
	."github.com/fuxiaohei/gorink/app"
	"github.com/fuxiaohei/goink/app"
	"github.com/fuxiaohei/gorink/model"
	"github.com/fuxiaohei/gorink/lib"
	"errors"
	"strconv"
	"strings"
	"time"
	"fmt"
)

func init() {
	App.GET("/article", func(context *app.InkContext) interface {} {
			context.Render("article/list.html", nil)
			return nil
		})
	App.GET("/admin/article", func(context *app.InkContext) interface {} {
			page, _ := strconv.Atoi(context.String("page"))
			if page < 1 {
				page = 1
			}
			articles, counter := model.GetArticleAllList(page, 10)
			renderAdminPage(context, "article/manage.html", "文章", "article", map[string]interface {}{
					"Articles":articles,
					"ArticlesLength":len(articles),
					"ArticleCounter":counter,
				})
			return nil
		})
	App.GET("/admin/article/new", func(context *app.InkContext) interface {} {
			renderAdminPage(context, "article/new.html", "撰写文章", "article", map[string]interface {}{
					"Categories":model.GetCategories(),
				})
			return nil
		})
	App.GET("/admin/article/edit", func(context *app.InkContext) interface {} {
			id, _ := strconv.Atoi(context.String("id"))
			if id < 1 {
				context.Redirect("/admin/article", 302)
				return nil
			}
			renderAdminPage(context, "article/edit.html", "编辑文章", "article", map[string]interface {}{
					"Categories":model.GetCategories(),
					"Article":model.GetArticleById(id),
				})
			return nil
		})
	App.POST("/admin/article/new", func(context *app.InkContext) interface {} {
			e := validateArticleData(context)
			if len(e) > 0 {
				renderAdminAlert(context, e)
				return nil
			}
			i, e2 := model.AddArticle(createNewArticle(context))
			if e2 != nil {
				renderAdminAlert(context, []error{e2})
				return nil
			}
			App.Trigger("model.article.add", i)
			context.Redirect("/admin/article/edit?saved=1&id=" + fmt.Sprint(i), 302)
			return nil
		})
	App.POST("/admin/article/edit", func(context *app.InkContext) interface {} {
			id, _ := strconv.Atoi(context.String("id"))
			if id < 1 {
				renderAdminAlert(context, []error{errors.New("参数错误")})
				return nil
			}
			e := validateArticleData(context)
			if len(e) > 0 {
				renderAdminAlert(context, e)
				return nil
			}
			article := createNewArticle(context)
			article.Id = id
			e2 := model.UpdateArticle(article)
			if e2 != nil {
				renderAdminAlert(context, []error{e2})
				return nil
			}
			App.Trigger("model.article.update", article.Id)
			context.Redirect("/admin/article/edit?updated=1&id=" + fmt.Sprint(id), 302)
			return nil
		})
}

func createNewArticle(context *app.InkContext) *model.Article {
	a := model.Article{}
	a.AuthorId = model.GetCurrentUserId(context)
	a.CategoryId, _ = strconv.Atoi(context.String("category"))
	a.Title = strings.TrimSpace(context.String("title"))
	a.Slug = strings.TrimSpace(context.String("slug"))
	a.Context = context.String("context")
	a.Excerpt = strings.Split(a.Context, "[break]")[0]
	a.CreateTime = time.Now().Unix()
	a.EditTime = a.CreateTime
	a.StatusText = context.String("status")
	a.AllowComment, _ = strconv.Atoi(context.String("comment"))
	return &a
}

func validateArticleData(context *app.InkContext) []error {
	e := make([]error, 0)
	if lib.IsEmptyString(context.String("title")) {
		e = append(e, errors.New("文章标题必填"))
	}
	if !lib.IsContain(context.String("context"), "[break]") {
		e = append(e, errors.New("全文需要分隔符[break]"))
	}
	if !lib.IsASCII(context.String("slug")) {
		e = append(e, errors.New("文章缩略名不支持中文和特殊字符"))
	}
	return e
}


