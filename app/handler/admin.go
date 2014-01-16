package handler

import (
	"github.com/fuxiaohei/GoBlog/GoInk"
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"strconv"
	"strings"
)

func Admin(context *GoInk.Context) {
	context.Layout("admin")
	context.Render("admin/home", map[string]interface{}{
			"Title": "控制台",
		})
}

func AdminProfile(context *GoInk.Context) {
	uid, _ := strconv.Atoi(context.Cookie("token-user"))
	user := model.GetUserById(uid)
	if context.Method == "POST" {
		data := context.Input()
		if !user.ChangeEmail(data["email"]) {
			Json(context, false).Set("msg", "邮箱与别的用户重复").End()
			return
		}
		user.Name = data["user"]
		user.Email = data["email"]
		user.Avatar = utils.Gravatar(user.Email, "180")
		user.Url = data["url"]
		user.Nick = data["nick"]
		user.Bio = data["bio"]
		Json(context, true).End()
		go model.SyncUsers()
		return
	}
	context.Layout("admin")
	context.Render("admin/profile", map[string]interface{}{
			"Title": "个性资料",
			"User":  user,
		})
}

func AdminPassword(context *GoInk.Context) {
	if context.Method == "POST" {
		uid, _ := strconv.Atoi(context.Cookie("token-user"))
		user := model.GetUserById(uid)
		if !user.CheckPassword(context.String("old")) {
			Json(context, false).Set("msg", "旧密码错误").End()
			return
		}
		user.ChangePassword(context.String("new"))
		go model.SyncUsers()
		Json(context, true).End()
		return
	}
	context.Layout("admin")
	context.Render("admin/password", map[string]interface{}{
			"Title": "修改密码",
			//"User":user,
		})
}

func AdminArticle(context *GoInk.Context) {
	articles, pager := model.GetArticleList(context.Int("page"), 10)
	context.Layout("admin")
	context.Render("admin/articles", map[string]interface{}{
			"Title":    "文章",
			"Articles": articles,
			"Pager":    pager,
		})
}

func ArticleWrite(context *GoInk.Context) {
	if context.Method == "POST" {
		c := new(model.Content)
		c.Id = 0
		data := context.Input()
		if !c.ChangeSlug(data["slug"]) {
			Json(context, false).Set("msg", "固定链接重复").End()
			return
		}
		c.Title = data["title"]
		c.Text = data["content"]
		c.Tags = strings.Split(strings.Replace(data["tag"], "，", ",", -1), ",")
		c.IsComment = data["comment"] == "1"
		c.IsLinked = false
		c.AuthorId, _ = strconv.Atoi(context.Cookie("token-user"))
		c.Template = "blog.html"
		c.Status = data["status"]
		c.Format = "markdown"
		c.Hits = 1
		var e error
		c, e = model.CreateContent(c, "article")
		if e != nil {
			Json(context, false).Set("msg", e.Error()).End()
			return
		}
		Json(context, true).Set("content", c).End()
		//c.Type = "article"
		return
	}
	context.Layout("admin")
	context.Render("admin/write_article", map[string]interface{}{
			"Title": "撰写文章",
		})
}

func ArticleEdit(context *GoInk.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	c := model.GetContentById(id)
	if c == nil {
		context.Redirect("/admin/articles/")
		return
	}
	if context.Method == "DELETE" {
		model.RemoveContent(c)
		Json(context, true).End()
		return
	}
	if context.Method == "POST" {
		data := context.Input()
		if !c.ChangeSlug(data["slug"]) {
			Json(context, false).Set("msg", "固定链接重复").End()
			return
		}
		c.Title = data["title"]
		c.Text = data["content"]
		c.Tags = strings.Split(strings.Replace(data["tag"], "，", ",", -1), ",")
		c.IsComment = data["comment"] == "1"
		//c.IsLinked = false
		//c.AuthorId, _ = strconv.Atoi(context.Cookie("token-user"))
		//c.Template = "blog.html"
		c.Status = data["status"]
		//c.Format = "markdown"
		model.SaveContent(c)
		Json(context, true).Set("content", c).End()
		//c.Type = "article"
		return
	}
	context.Layout("admin")
	context.Render("admin/edit_article", map[string]interface{}{
			"Title":   "编辑文章",
			"Article": c,
		})
}

func PageWrite(context *GoInk.Context) {
	if context.Method == "POST" {
		c := new(model.Content)
		c.Id = 0
		data := context.Input()
		if !c.ChangeSlug(data["slug"]) {
			Json(context, false).Set("msg", "固定链接重复").End()
			return
		}
		c.Title = data["title"]
		c.Text = data["content"]
		c.Tags = make([]string, 0)
		c.IsComment = data["comment"] == "1"
		c.IsLinked = data["link"] == "1"
		c.AuthorId, _ = strconv.Atoi(context.Cookie("token-user"))
		c.Template = "page.html"
		c.Status = data["status"]
		c.Format = "markdown"
		c.Hits = 1
		var e error
		c, e = model.CreateContent(c, "page")
		if e != nil {
			Json(context, false).Set("msg", e.Error()).End()
			return
		}
		Json(context, true).Set("content", c).End()
		//c.Type = "article"
		return
	}
	context.Layout("admin")
	context.Render("admin/write_page", map[string]interface{}{
			"Title": "撰写页面",
		})
}

func AdminPage(context *GoInk.Context) {
	pages, pager := model.GetPageList(context.Int("page"), 10)
	context.Layout("admin")
	context.Render("admin/pages", map[string]interface{}{
			"Title": "页面",
			"Pages": pages,
			"Pager": pager,
		})
}

func PageEdit(context *GoInk.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	c := model.GetContentById(id)
	if c == nil {
		context.Redirect("/admin/pages/")
		return
	}
	if context.Method == "DELETE" {
		model.RemoveContent(c)
		Json(context, true).End()
		return
	}
	if context.Method == "POST" {
		data := context.Input()
		if !c.ChangeSlug(data["slug"]) {
			Json(context, false).Set("msg", "固定链接重复").End()
			return
		}
		c.Title = data["title"]
		c.Text = data["content"]
		//c.Tags = strings.Split(strings.Replace(data["tag"], "，", ",", -1), ",")
		c.IsComment = data["comment"] == "1"
		c.IsLinked = data["link"] == "1"
		//c.AuthorId, _ = strconv.Atoi(context.Cookie("token-user"))
		//c.Template = "blog.html"
		c.Status = data["status"]
		//c.Format = "markdown"
		model.SaveContent(c)
		Json(context, true).Set("content", c).End()
		//c.Type = "article"
		return
	}
	context.Layout("admin")
	context.Render("admin/edit_page", map[string]interface{}{
			"Title": "编辑文章",
			"Page":  c,
		})
}

func AdminSetting(context *GoInk.Context) {
	if context.Method == "POST" {
		data := context.Input()
		for k, v := range data {
			model.SetSetting(k, v)
		}
		model.SyncSettings()
		Json(context, true).End()
		return
	}
	context.Layout("admin")
	context.Render("admin/setting", map[string]interface{}{
			"Title": "配置",
		})
}

func CustomSetting(context *GoInk.Context) {
	if context.Method == "POST" {
		keys := context.Strings("key")
		values := context.Strings("value")
		for i, k := range keys {
			model.SetSetting("c_" + k, values[i])
		}
		model.SyncSettings()
		Json(context, true).End()
		return
	}
	context.Layout("admin")
	context.Render("admin/custom_setting", map[string]interface{}{
			"Title":    "自定义配置",
			"Settings": model.GetCustomSettings(),
		})
}

func AdminComments(context *GoInk.Context) {
	if context.Method == "DELETE" {
		id := context.Int("id")
		cmt := model.GetCommentById(id)
		model.RemoveComment(cmt.Cid, id)
		Json(context, true).End()
		return
	}
	if context.Method == "PUT" {
		id := context.Int("id")
		cmt2 := model.GetCommentById(id)
		cmt2.Status = "approved"
		cmt2.GetReader().Active = true
		model.SaveComment(cmt2)
		Json(context, true).End()
		return
	}
	if context.Method == "POST" {
		// get required data
		pid := context.Int("pid")
		cid := model.GetCommentById(pid).Cid
		uid, _ := strconv.Atoi(context.Cookie("token-user"))
		user := model.GetUserById(uid)

		co := new(model.Comment)
		co.Author = user.Nick
		co.Email = user.Email
		co.Url = user.Url
		co.Content = strings.Replace(utils.Html2str(context.String("content")), "\n", "<br/>", -1)
		co.Avatar = utils.Gravatar(co.Email, "50")
		co.Pid = pid
		co.Ip = context.Ip
		co.UserAgent = context.UserAgent
		co.IsAdmin = true
		model.CreateComment(cid, co)
		Json(context, true).Set("comment", co.ToJson()).End()
		return
	}
	page := context.IntOr("page", 1)
	comments, pager := model.GetCommentList(page, 6)
	context.Layout("admin")
	context.Render("admin/comments", map[string]interface{}{
			"Title":    "评论",
			"Comments": comments,
			"Pager":    pager,
		})
}
