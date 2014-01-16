package handler

import (
	"git.oschina.net/fuxiaohei/GoBlog.git/GoInk"
	"git.oschina.net/fuxiaohei/GoBlog.git/app/model"
	"git.oschina.net/fuxiaohei/GoBlog.git/app/utils"
	"strconv"
	"strings"
)

func Login(context *GoInk.Context) {
	if context.Method == "POST" {
		data := context.Input()
		user := model.GetUserByName(data["user"])
		if user == nil {
			Json(context, false).End()
			return
		}
		if !user.CheckPassword(data["password"]) {
			Json(context, false).End()
			return
		}
		exp := 3600 * 24 * 3
		expStr := strconv.Itoa(exp)
		s := model.CreateToken(user, context, int64(exp))
		context.Cookie("token-user", strconv.Itoa(s.UserId), expStr)
		context.Cookie("token-value", s.Value, expStr)
		Json(context, true).End()
		return
	}
	if context.Cookie("token-value") != "" {
		context.Redirect("/admin/")
		return
	}
	context.Render("home/login", nil)
}

func Auth(context *GoInk.Context) {
	tokenValue := context.Cookie("token-value")
	token := model.GetTokenByValue(tokenValue)
	if token == nil {
		context.Redirect("/logout/")
		context.End()
		return
	}
	if !token.IsValid() {
		context.Redirect("/logout/")
		context.End()
		return
	}
}

func Logout(context *GoInk.Context) {
	context.Cookie("token-user", "", "-3600")
	context.Cookie("token-value", "", "-3600")
	context.Redirect("/login/")
}

func Home(context *GoInk.Context) {
	context.Layout("home")
	page, _ := strconv.Atoi(context.Param("page"))
	size, _ := strconv.Atoi(model.GetSetting("article_size"))
	articles, pager := model.GetArticleList(page, size)
	context.Render("home/home", map[string]interface{}{
		"Articles": articles,
		"Pager":    pager,
	})
}

func Article(context *GoInk.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	slug := context.Param("slug")
	article := model.GetContentById(id)
	if article == nil {
		context.Redirect("/")
		return
	}
	if article.Slug != slug || article.Type != "article" {
		context.Redirect("/")
		return
	}
	article.Hits++
	context.Layout("home")
	context.Render("home/article", map[string]interface{}{
		"Title":       article.Title,
		"Article":     article,
		"CommentHtml": Comments(context, article),
	})
}

func TopPage(context *GoInk.Context) {
	slug := context.Param("slug")
	page := model.GetContentBySlug(slug)
	if page == nil {
		context.Redirect("/")
		return
	}
	if page.IsLinked && page.Type == "page" {
		context.Layout("home")
		context.Render("home/page", map[string]interface{}{
			"Title": page.Title,
			"Page":  page,
		})
		page.Hits++
		return
	}
	context.Redirect("/")
}

func Comments(context *GoInk.Context, c *model.Content) string {
	return context.Tpl("home/comment", map[string]interface{}{
		"Content":  c,
		"Comments": c.Comments,
	})
}

func Comment(context *GoInk.Context) {
	cid, _ := strconv.Atoi(context.Param("id"))
	if cid < 1 {
		Json(context, false).End()
		return
	}
	if model.GetContentById(cid) == nil {
		Json(context, false).End()
		return
	}
	data := context.Input()
	co := new(model.Comment)
	co.Author = data["user"]
	co.Email = data["email"]
	co.Url = data["url"]
	co.Content = strings.Replace(utils.Html2str(data["content"]), "\n", "<br/>", -1)
	co.Avatar = utils.Gravatar(co.Email, "50")
	co.Pid, _ = strconv.Atoi(data["pid"])
	co.Ip = context.Ip
	co.UserAgent = context.UserAgent
	co.IsAdmin = false
	model.CreateComment(cid, co)
	Json(context, true).Set("comment", co.ToJson()).End()
}
