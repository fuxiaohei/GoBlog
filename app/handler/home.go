package handler

import (
	"github.com/fuxiaohei/GoBlog/app/model/content"
	"github.com/fuxiaohei/GoBlog/app/model/message"
	"github.com/fuxiaohei/GoBlog/app/model/setting"
	mUser "github.com/fuxiaohei/GoBlog/app/model/user"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"github.com/fuxiaohei/GoInk"
	"net/url"
	"strconv"
	"strings"
)

// Login is login page handler, pattern /login/.
func Login(context *GoInk.Context) {
	if context.Method == "POST" {
		data := context.Input()
		user := mUser.ByName(data["user"])
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
		s := mUser.CreateToken(user, context, int64(exp))
		context.Cookie("token-user", strconv.Itoa(s.UserId), expStr)
		context.Cookie("token-value", s.Value, expStr)
		Json(context, true).End()
		return
	}
	if context.Cookie("token-value") != "" {
		context.Redirect("/admin/")
		return
	}
	context.Render("admin/login", nil)
}

// Auth is authorization checking handler, use for middleware.
func Auth(context *GoInk.Context) {
	tokenValue := context.Cookie("token-value")
	token := mUser.TokenByValue(tokenValue)
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

// Logout is safely log out page, pattern /logout/.
func Logout(context *GoInk.Context) {
	context.Cookie("token-user", "", "-3600")
	context.Cookie("token-value", "", "-3600")
	context.Redirect("/login/")
}

// TagArticles is tag article list page, pattern /tag/:tag_name/.
func TagArticles(ctx *GoInk.Context) {
	ctx.Layout("home")
	page, _ := strconv.Atoi(ctx.Param("page"))
	tag, _ := url.QueryUnescape(ctx.Param("tag"))
	size := getArticleListSize()
	articles, pager := content.TaggedArticleList(tag, page, getArticleListSize())
	// fix dotted tag
	if len(articles) < 1 && strings.Contains(tag, "-") {
		articles, pager = content.TaggedArticleList(strings.Replace(tag, "-", ".", -1), page, size)
	}
	Theme(ctx).Layout("home").Render("index", map[string]interface{}{
		"Articles":    articles,
		"Pager":       pager,
		"SidebarHtml": SidebarHtml(ctx),
		"Tag":         tag,
		"Title":       tag,
	})
}

// Home is home page handler, pattern /.
func Home(context *GoInk.Context) {
	context.Layout("home")
	page, _ := strconv.Atoi(context.Param("page"))
	articles, pager := content.PublishArticleList(page, getArticleListSize())
	data := map[string]interface{}{
		"Articles":    articles,
		"Pager":       pager,
		"SidebarHtml": SidebarHtml(context),
	}
	if page > 1 {
		data["Title"] = "第 " + strconv.Itoa(page) + " 页"
	}
	Theme(context).Layout("home").Render("index", data)
}

// Article is single article page, pattern /article/:article_id/:article_slug.
func Article(context *GoInk.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	slug := context.Param("slug")
	article := content.ById(id)
	if article == nil {
		context.Redirect("/")
		return
	}
	if article.Slug != slug || article.Type != "article" {
		context.Redirect("/")
		return
	}
	article.Hits++
	Theme(context).Layout("home").Render("article", map[string]interface{}{
		"Title":       article.Title,
		"Article":     article,
		"CommentHtml": CommentHtml(context, article),
	})
}

// Page is single page showing page, pattern /page/:page_id/:page_slug.
func Page(context *GoInk.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	slug := context.Param("slug")
	article := content.ById(id)
	if article == nil || article.Status != "publish" {
		context.Redirect("/")
		return
	}
	if article.Slug != slug || article.Type != "page" {
		context.Redirect("/")
		return
	}
	article.Hits++
	Theme(context).Layout("home").Render("page", map[string]interface{}{
		"Title": article.Title,
		"Page":  article,
		//"CommentHtml": Comments(context, article),
	})
}

// TopPage is top level page handler, pattern /:page_slug.
func TopPage(context *GoInk.Context) {
	slug := context.Param("slug")
	page := content.BySlug(slug)
	if page == nil || page.Status != "publish" {
		context.Redirect("/")
		return
	}
	if page.IsLinked && page.Type == "page" {
		Theme(context).Layout("home").Render("page", map[string]interface{}{
			"Title": page.Title,
			"Page":  page,
		})
		page.Hits++
		return
	}
	context.Redirect("/")
}

// Comment is ajax comment post handler, pattern /comment/:content_id/.
func Comment(context *GoInk.Context) {
	cid, _ := strconv.Atoi(context.Param("id"))
	if cid < 1 {
		Json(context, false).End()
		return
	}
	if content.ById(cid) == nil {
		Json(context, false).End()
		return
	}
	data := context.Input()
	msg := validateComment(data)
	if msg != "" {
		Json(context, false).Set("msg", msg).End()
		return
	}
	co := new(content.Comment)
	co.Author = data["user"]
	co.Email = data["email"]
	co.Url = data["url"]
	co.Content = data["content"]
	co.Avatar = utils.Gravatar(co.Email, "50")
	co.Pid, _ = strconv.Atoi(data["pid"])
	co.Ip = context.Ip
	co.UserAgent = context.UserAgent
	co.IsAdmin = false
	content.CreateComment(cid, co)
	Json(context, true).Set("comment", co.ToJson()).End()
	message.Create("comment", co)
	context.Do("comment_created", co)
}

func Redirect(ctx *GoInk.Context) {
	to := ctx.StringOr("to", "/")
	ctx.Redirect(to, 302)
}

func validateComment(data map[string]string) string {
	if utils.IsEmptyString(data["user"]) || utils.IsEmptyString(data["content"]) {
		return "称呼，邮箱，内容必填"
	}
	if !utils.IsEmail(data["email"]) {
		return "邮箱格式错误"
	}
	if !utils.IsEmptyString(data["url"]) && !utils.IsURL(data["url"]) {
		return "网址格式错误"
	}
	return ""
}

func getArticleListSize() int {
	size, _ := strconv.Atoi(setting.Get("article_size"))
	if size < 1 {
		size = 5
	}
	return size
}
