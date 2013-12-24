package handler

import (
	"github.com/fuxiaohei/GoInk/Core"
	"github.com/fuxiaohei/GoBlog/app/model"
	"strings"
	"time"
	"strconv"
)

func AdminArticle(context *Core.Context) interface {} {
	articles, pager := model.ArticleM.GetPaged(1, 8, false)
	context.Render("admin:admin/article.html", map[string]interface {}{
			"Title":"文章",
			"IsArticle":true,
			"Articles":articles,
			"Pager":pager,
		})
	return nil
}

func AdminArticleNew(context *Core.Context) interface {} {
	context.Render("admin:admin/article_new.html", map[string]interface {}{
			"Title":"写文章",
			"Categories":model.CategoryM.GetAll(),
		})
	return nil
}

func AdminArticleNewPost(context *Core.Context) interface {} {
	if !context.IsAjax {
		context.Status = 400
		return nil
	}
	data := context.Input()
	a := model.ArticleM.GetArticleBySlug(data["slug"])
	if a != nil {
		context.Json(map[string]interface {}{
			"res":false,
			"msg":"链接重复",
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
	article = model.ArticleM.SaveArticle(article)
	if article != nil {
		context.Json(map[string]interface {}{
			"res":true,
			"id":article.Id,
		})
		return nil
	}
	context.Json(map[string]interface {}{
		"res":false,
		"msg":"保存失败",
	})
	return nil
}
