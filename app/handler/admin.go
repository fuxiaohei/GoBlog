package handler

import (
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoBlog/app/plugin"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"github.com/fuxiaohei/GoInk"
	"strconv"
	"strings"
)

func Admin(context *GoInk.Context) {
	uid, _ := strconv.Atoi(context.Cookie("token-user"))
	user := model.GetUserById(uid)
	context.Layout("admin/admin")
	context.Render("admin/home", map[string]interface{}{
		"Title":    "控制台",
		"Statis":   model.NewStatis(),
		"User":     user,
		"Messages": model.GetUnreadMessages(),
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
		go model.UpdateCommentAdmin(user)
		context.Do("profile_update", user)
		return
	}
	context.Layout("admin/admin")
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
		context.Do("password_update", user)
		return
	}
	context.Layout("admin/admin")
	context.Render("admin/password", map[string]interface{}{
		"Title": "修改密码",
		//"User":user,
	})
}

func AdminArticle(context *GoInk.Context) {
	articles, pager := model.GetArticleList(context.Int("page"), 10)
	context.Layout("admin/admin")
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
		context.Do("article_created", c)
		//c.Type = "article"
		return
	}
	context.Layout("admin/admin")
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
		context.Do("article_modified", c)
		//c.Type = "article"
		return
	}
	context.Layout("admin/admin")
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
		context.Do("page_created", c)
		return
	}
	context.Layout("admin/admin")
	context.Render("admin/write_page", map[string]interface{}{
		"Title": "撰写页面",
	})
}

func AdminPage(context *GoInk.Context) {
	pages, pager := model.GetPageList(context.Int("page"), 10)
	context.Layout("admin/admin")
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
		context.Do("page_modified", c)
		//c.Type = "article"
		return
	}
	context.Layout("admin/admin")
	context.Render("admin/edit_page", map[string]interface{}{
		"Title": "编辑文章",
		"Page":  c,
	})
}

func AdminSetting(context *GoInk.Context) {
	if context.Method == "POST" {
		data := context.Input()
		for k, v := range data {
			if v == "" {
				if data[k+"_def"] != "" {
					v = data[k+"_def"]
				}
			}
			model.SetSetting(k, v)
		}
		model.SyncSettings()
		Json(context, true).End()
		context.Do("setting_saved")
		return
	}
	context.Layout("admin/admin")
	context.Render("admin/setting", map[string]interface{}{
		"Title":      "配置",
		"Custom":     model.GetCustomSettings(),
		"Navigators": model.GetNavigators(),
	})
}

func CustomSetting(context *GoInk.Context) {
	keys := context.Strings("key")
	values := context.Strings("value")
	for i, k := range keys {
		if len(k) < 1 {
			continue
		}
		model.SetSetting("c_"+k, values[i])
	}
	model.SyncSettings()
	Json(context, true).End()
	context.Do("setting_saved")
	return
}

func NavigatorSetting(context *GoInk.Context) {
	order := context.Strings("order")
	text := context.Strings("text")
	title := context.Strings("title")
	link := context.Strings("link")
	model.SetNavigators(order, text, title, link)
	Json(context, true).End()
	context.Do("setting_saved")
	return
}

func AdminComments(context *GoInk.Context) {
	if context.Method == "DELETE" {
		id := context.Int("id")
		cmt := model.GetCommentById(id)
		model.RemoveComment(cmt.Cid, id)
		Json(context, true).End()
		context.Do("comment_delete", id)
		return
	}
	if context.Method == "PUT" {
		id := context.Int("id")
		cmt2 := model.GetCommentById(id)
		cmt2.Status = "approved"
		cmt2.GetReader().Active = true
		model.SaveComment(cmt2)
		Json(context, true).End()
		context.Do("comment_change_status", cmt2)
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
		co.Content = context.String("content")
		co.Avatar = utils.Gravatar(co.Email, "50")
		co.Pid = pid
		co.Ip = context.Ip
		co.UserAgent = context.UserAgent
		co.IsAdmin = true
		model.CreateComment(cid, co)
		Json(context, true).Set("comment", co.ToJson()).End()
		model.CreateMessage("comment", co)
		context.Do("comment_reply", co)
		return
	}
	page := context.IntOr("page", 1)
	comments, pager := model.GetCommentList(page, 10)
	context.Layout("admin/admin")
	context.Render("admin/comments", map[string]interface{}{
		"Title":    "评论",
		"Comments": comments,
		"Pager":    pager,
	})
}

func AdminPlugin(context *GoInk.Context) {
	if context.Method == "POST" {
		action := context.String("action")
		if action == "" {
			Json(context, false).End()
			return
		}
		pln := context.String("plugin")
		if action == "activate" {
			plugin.Activate(pln)
			plugin.Update(context.App())
			Json(context, true).End()
			context.Do("plugin_activated", pln)
			return
		}
		if action == "deactivate" {
			plugin.Deactivate(pln)
			Json(context, true).End()
			context.Do("plugin_deactivated", pln)
			return
		}
		context.Status = 405
		Json(context, false).End()
		return
	}
	context.Layout("admin/admin")
	context.Render("admin/plugin", map[string]interface{}{
		"Title":   "插件",
		"Plugins": plugin.GetPlugins(),
	})
}

func PluginSetting(context *GoInk.Context) {
	key := context.Param("plugin_key")
	if key == "" {
		context.Redirect("/admin/plugins/")
		return
	}
	p := plugin.GetPluginByKey(key)
	if p == nil {
		context.Redirect("/admin/plugins/")
		return
	}
	if context.Method == "POST" {
		p.SetSetting(context.Input())
		Json(context, true).End()
		context.Do("plugin_setting_saved", p)
		return
	}
	context.Layout("admin/admin")
	context.Render("admin/plugin_setting", map[string]interface{}{
		"Title": "插件 - " + p.Name(),
		"Form":  p.Form(),
	})
}

func AdminMessageRead(context *GoInk.Context) {
	id := context.Int("id")
	if id < 0 {
		Json(context, false).End()
		return
	}
	m := model.GetMessage(id)
	if m == nil {
		Json(context, false).End()
		return
	}
	model.SaveMessageRead(m)
	Json(context, true).End()
}
