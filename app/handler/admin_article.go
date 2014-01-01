package handler

import (
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoInk/Core"
	"strconv"
	"strings"
	"time"
)

func AdminArticle(context *Core.Context) interface{} {
	page := 1
	if context.Param(2) == "page" {
		page, _ = strconv.Atoi(context.Param(3))
		if page < 1 {
			page = 1
		}
	}
	articles, pager := model.ArticleM.GetPaged(page, 10, false)
	context.Render("admin:admin/article.html", map[string]interface{}{
			"Title":     "文章",
			"IsArticle": true,
			"Articles":  articles,
			"Pager":     pager,
		})
	return nil
}

func AdminArticleNew(context *Core.Context) interface{} {
	context.Render("admin:admin/article_new.html", map[string]interface{}{
			"Title":      "写文章",
			"Categories": model.CategoryM.GetAll(),
		})
	return nil
}

func AdminArticleNewPost(context *Core.Context) interface{} {
	if !context.IsAjax {
		context.Status = 400
		return nil
	}
	data := context.Input()
	a := model.ArticleM.GetArticleBySlug(data["slug"])
	if a != nil {
		context.Json(map[string]interface{}{
			"res": false,
			"msg": "链接重复",
		})
		return nil
	}
	article := new(model.Article)
	article.Title = data["title"]
	article.Slug = data["slug"]
	article.Summary = strings.Split(data["content"], "[break]")[0]
	article.Content = data["content"]
	article.CreateTime = time.Now().Unix()
	article.EditTime = article.CreateTime
	article.CategoryId, _ = strconv.Atoi(data["category"])
	article.AuthorId, _ = strconv.Atoi(context.Cookie("admin-user"))
	article.Status = data["status"]
	article.IsComment, _ = strconv.Atoi(data["comment"])
	article.IsFeed, _ = strconv.Atoi(data["feed"])
	article = model.ArticleM.CreateArticle(article)
	if article != nil {
		context.Json(map[string]interface{}{
			"res": true,
			"id":  article.Id,
		})
		go model.CategoryM.CountArticle()
		return nil
	}
	context.Json(map[string]interface{}{
		"res": false,
		"msg": "保存失败",
	})
	return nil
}

func AdminArticleEdit(context *Core.Context) interface{} {
	id, _ := strconv.Atoi(context.Param(3))
	article := model.ArticleM.GetArticleById(id)
	if article == nil {
		context.Redirect("/admin/article/")
		return nil
	}
	context.Render("admin:admin/article_edit.html", map[string]interface{}{
			"Title":      "修改文章",
			"IsArticle":  true,
			"Article":    article,
			"Categories": model.CategoryM.GetAll(),
		})
	return nil
}

func AdminArticleEditPost(context *Core.Context) interface{} {
	if !context.IsAjax {
		context.Status = 400
		return nil
	}
	data := context.Input()
	id, _ := strconv.Atoi(context.Param(3))
	a := model.ArticleM.GetArticleBySlug(data["slug"])
	if a != nil && a.Id != id {
		context.Json(map[string]interface{}{
			"res": false,
			"msg": "链接重复",
		})
		return nil
	}
	article := new(model.Article)
	article.Id = id
	article.Title = data["title"]
	article.Slug = data["slug"]
	article.Summary = strings.Split(data["content"], "[break]")[0]
	article.Content = data["content"]
	article.EditTime = time.Now().Unix()
	article.CategoryId, _ = strconv.Atoi(data["category"])
	article.Status = data["status"]
	article.IsComment, _ = strconv.Atoi(data["comment"])
	article.IsFeed, _ = strconv.Atoi(data["feed"])
	article = model.ArticleM.SaveArticle(article)
	if article != nil {
		context.Json(map[string]interface{}{
			"res": true,
			"id":  article.Id,
		})
		go model.CategoryM.CountArticle()
		return nil
	}
	context.Json(map[string]interface{}{
		"res": false,
		"msg": "保存失败",
	})
	return nil
}

func AdminArticleDelete(context *Core.Context) interface {} {
	id, _ := strconv.Atoi(context.Param(3))
	article := model.ArticleM.GetArticleById(id)
	if article == nil {
		context.Redirect("/admin/article/")
		return nil
	}
	// delete comments
	model.CommentM.DeleteCommentsInContent(article.Id)
	// delete article
	model.ArticleM.DeleteArticle(article.Id)
	// update category counts
	model.CategoryM.CountArticle()
	context.Redirect(context.Referer)
	return nil
}
