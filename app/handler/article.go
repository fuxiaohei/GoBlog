package handler

import (
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"github.com/fuxiaohei/GoInk/Core"
	"strconv"
	"strings"
)

func Article(context *Core.Context) interface{} {
	if context.Ext == ".html" && context.Param(2) != "" {
		if context.IsAjax {
			return ArticleComments(context)
		}
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
		"Articles":    article,
		"Pager":       pager,
		"PageUrl":     "/article/page",
		"ArticleSide": articleSide(),
	})
	return nil
}

func checkArticleFromContext(context *Core.Context) (*model.Article, bool) {
	slug := context.Param(2)
	article := model.ArticleM.GetArticleBySlug(slug)
	if article == nil {
		//context.Redirect("/article/")
		return nil, false
	}
	id, _ := strconv.Atoi(context.Param(1))
	if article.Id != id {
		//context.Redirect("/article/")
		return nil, false
	}
	return article, true
}

func ArticleSingle(context *Core.Context) interface{} {
	article, ok := checkArticleFromContext(context)
	if !ok {
		context.Redirect("/article/")
		return nil
	}
	context.Render("theme:default/article_single.html", map[string]interface{}{
		"Article": article,
	})
	model.ArticleM.IncreaseView(article.Id)
	return nil
}

func ArticleComments(context *Core.Context) interface{} {
	article, ok := checkArticleFromContext(context)
	if !ok {
		context.Status = 400
		return nil
	}
	comments := model.CommentM.GetAllOfContent(article.Id, true)
	context.Json(map[string]interface{}{
		"res":      len(comments) > 0,
		"comments": comments,
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
		"Articles":    article,
		"Pager":       pager,
		"PageUrl":     category.Link() + "page/",
		"Category":    category,
		"ArticleSide": articleSide(),
	})
	return nil
}

func articleSide() map[string]interface{} {
	return map[string]interface{}{
		"Categories": model.CategoryM.GetCountsDesc(3600),
		"Popular":    model.ArticleM.GetPopular(4),
	}
}

func ArticleCommentPost(context *Core.Context) interface{} {
	if !context.IsAjax {
		context.Status = 400
		return nil
	}
	data := context.Input()
	msg, ok := validateCommentData(data)
	if !ok {
		context.Json(map[string]interface{}{
			"res": true,
			"msg": msg,
		})
		return nil
	}
	c := new(model.Comment)
	c.Author = data["author"]
	c.Email = data["email"]
	c.Site = data["site"]
	c.Content = strings.Replace(data["content"], "\n", "<br/>", -1)
	c.ContentId, _ = strconv.Atoi(data["article"])
	c.Pid, _ = strconv.Atoi(data["pid"])
	c.Avatar = utils.Gravatar(c.Email, "50")
	c.UserId = 0
	c.IsAdmin = false
	c = model.CommentM.CreateComment(c)
	if c != nil {
		context.Json(map[string]interface{}{
			"res":     true,
			"comment": c,
		})
		model.ArticleM.CountComments()
		return nil
	}
	context.Json(map[string]interface{}{
		"res": true,
		"msg": "评论失败",
	})
	return nil
}

func validateCommentData(data map[string]string) (string, bool) {
	if !utils.IsEmail(data["email"]) {
		return "邮箱错误", false
	}
	if len(data["site"]) > 0 {
		if !utils.IsURL(data["site"]) {
			return "网址错误", false
		}
	}
	if len(data["content"]) < 3 {
		return "内容至少三个字", false
	}
	return "", true
}
