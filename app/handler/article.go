package handler

import (
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoInk/Core"
	"strconv"
)

func Article(context *Core.Context) interface{} {
	if context.Ext == ".html" && context.Param(2) != "" {
		return ArticleSingle(context)
	}
	page := 1
	p1 := context.Param(1)
	if p1 == "page" {
		page, _ = strconv.Atoi(context.Param(2))
		if page < 1 {
			page = 1
		}
	}
	article, pager := model.ArticleM.GetPaged(page, 4, true)
	context.Render("theme:default/article.html", map[string]interface{}{
		"Articles": article,
		"Pager":    pager,
		"PageUrl":  "/article/page",
		"ArticleSide":articleSide(),
	})
	return nil
}

func ArticleSingle(context *Core.Context) interface{} {
	slug := context.Param(2)
	article := model.ArticleM.GetArticleBySlug(slug)
	if article == nil {
		context.Redirect("/article/")
		return nil
	}
	id, _ := strconv.Atoi(context.Param(1))
	if article.Id != id {
		context.Redirect("/article/")
		return nil
	}
	context.Render("theme:default/article_single.html", map[string]interface{}{
		"Article": article,
	})
	return nil
}

func ArticleCategory(context *Core.Context) interface{} {
	slug := context.Param(1)
	if slug == "" {
		context.Redirect("/article/")
		return nil
	}
	category := model.CategoryM.GetCategoryBySlug(slug)
	if category == nil {
		context.Redirect("/article/")
		return nil
	}
	page := 1
	p1 := context.Param(2)
	if p1 == "page" {
		page, _ = strconv.Atoi(context.Param(3))
		if page < 1 {
			page = 1
		}
	}
	article, pager := model.ArticleM.GetCategoryPaged(category.Id, page, 4, true)
	context.Render("theme:default/article.html", map[string]interface{}{
		"Articles": article,
		"Pager":    pager,
		"PageUrl":  category.Link()+"page/",
		"Category": category,
		"ArticleSide":articleSide(),
	})
	return nil
}

func articleSide() map[string]interface {} {
	return map[string]interface {}{
		"Categories":model.CategoryM.GetCountsDesc(3600),
		"Popular":model.ArticleM.GetPopular(4),
	}
}
