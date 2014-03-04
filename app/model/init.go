package model

import (
	"github.com/fuxiaohei/GoBlog/app/model/content"
	"github.com/fuxiaohei/GoBlog/app/model/file"
	"github.com/fuxiaohei/GoBlog/app/model/message"
	"github.com/fuxiaohei/GoBlog/app/model/setting"
	"github.com/fuxiaohei/GoBlog/app/model/storage"
	"github.com/fuxiaohei/GoBlog/app/model/timer"
	"github.com/fuxiaohei/GoBlog/app/model/user"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"os"
	"path"
	"strconv"
)

func loadAllData() {
	setting.LoadVersion()
	setting.Load()
	setting.LoadNavigators()
	user.Load()
	user.LoadTokens()
	content.Load()
	message.Load()
	content.LoadReaders()
	content.LoadComments()
	file.Load()
}

func writeDefaultData() {
	// write user
	u := new(user.User)
	u.Id = Storage.TimeInc(10)
	u.Name = "admin"
	u.Password = utils.Sha1("adminxxxxx")
	u.Nick = "管理员"
	u.Email = "admin@example.com"
	u.Url = "http://example.com/"
	u.CreateTime = utils.Now()
	u.Bio = "这是站点的管理员，你可以添加一些个人介绍，支持换行不支持markdown"
	u.LastLoginTime = u.CreateTime
	u.Role = "ADMIN"
	Storage.Set("users", []*user.User{u})

	// write token
	Storage.Set("tokens", map[string]*user.Token{})

	// write contents
	a := new(content.Content)
	a.Id = Storage.TimeInc(9)
	a.Title = "欢迎使用 Fxh.Go"
	a.Slug = "welcome-fxh-go"
	a.Text = "如果您看到这篇文章,表示您的 blog 已经安装成功."
	a.Tags = []string{"Fxh.Go"}
	a.CreateTime = utils.Now()
	a.EditTime = a.CreateTime
	a.UpdateTime = a.CreateTime
	a.IsComment = true
	a.IsLinked = false
	a.AuthorId = u.Id
	a.Type = "article"
	a.Status = "publish"
	a.Format = "markdown"
	a.Template = "blog.html"
	a.Hits = 1
	// write comments
	co := new(content.Comment)
	co.Author = u.Nick
	co.Email = u.Email
	co.Url = u.Url
	co.Content = "欢迎加入使用 Fxh.Go"
	co.Avatar = utils.Gravatar(co.Email, "50")
	co.Pid = 0
	co.Ip = "127.0.0.1"
	co.UserAgent = "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.17 (KHTML, like Gecko) Chrome/24.0.1312.57 Safari/537.17"
	co.IsAdmin = true
	co.Id = Storage.TimeInc(7)
	co.CreateTime = utils.Now()
	co.Status = "approved"
	co.Cid = a.Id
	a.Comments = []*content.Comment{co}
	Storage.Set("content/article-"+strconv.Itoa(a.Id), a)

	// write pages
	p := new(Content)
	p.Id = a.Id + Storage.TimeInc(6)
	p.Title = "关于"
	p.Slug = "about-me"
	p.Text = "本页面由 Fxh.Go 创建, 这只是个测试页面."
	p.Tags = []string{}
	p.CreateTime = utils.Now()
	p.EditTime = p.CreateTime
	p.UpdateTime = p.UpdateTime
	p.IsComment = true
	p.IsLinked = true
	p.AuthorId = u.Id
	p.Type = "page"
	p.Status = "publish"
	p.Format = "markdown"
	p.Comments = make([]*content.Comment, 0)
	p.Template = "page.html"
	p.Hits = 1
	Storage.Set("content/page-"+strconv.Itoa(p.Id), p)
	p2 := new(Content)
	p2.Id = p.Id + Storage.TimeInc(6)
	p2.Title = "好友"
	p2.Slug = "friends"
	p2.Text = "本页面由 Fxh.Go 创建, 这只是个测试页面."
	p2.Tags = []string{}
	p2.CreateTime = utils.Now()
	p2.EditTime = p2.CreateTime
	p2.UpdateTime = p2.UpdateTime
	p2.IsComment = true
	p2.IsLinked = true
	p2.AuthorId = u.Id
	p2.Type = "page"
	p2.Status = "publish"
	p2.Format = "markdown"
	p2.Comments = make([]*content.Comment, 0)
	p2.Template = "page.html"
	p2.Hits = 1
	Storage.Set("content/page-"+strconv.Itoa(p2.Id), p2)

	// write new reader
	Storage.Set("readers", map[string]*content.Reader{})

	// write version
	v := new(setting.Version)
	v.Name = "Fxh.Go"
	v.BuildTime = utils.Now()
	v.Version = AppVersion
	Storage.Set("version", v)

	// write settings
	s := map[string]string{
		"site_title":         "Fxh.Go",
		"site_sub_title":     "Go开发的简单博客",
		"site_keywords":      "Fxh.Go,Golang,Blog",
		"site_description":   "Go语言开发的简单博客程序",
		"site_url":           "http://localhost/",
		"article_size":       "4",
		"site_theme":         "default",
		"enable_go_markdown": "false",
		"c_footer_weibo":     "#",
		"c_footer_github":    "#",
		"c_footer_email":     "#",
		"c_home_avatar":      "/static/img/site.png",
		"c_footer_ga":        "<!-- google analytics or other -->",
	}
	Storage.Set("settings", s)

	// write files
	Storage.Set("files", []*File{})

	// write message
	Storage.Set("messages", []*message.Message{})

	// write navigators
	n := new(setting.NavItem)
	n.Order = 1
	n.Text = "文章"
	n.Title = "文章"
	n.Link = "/"
	n2 := new(setting.NavItem)
	n2.Order = 2
	n2.Text = "关于"
	n2.Title = "关于"
	n2.Link = "/about-me.html"
	n3 := new(setting.NavItem)
	n3.Order = 3
	n3.Text = "好友"
	n3.Title = "好友"
	n3.Link = "/friends.html"
	Storage.Set("navigators", []*setting.NavItem{n, n2, n3})

	// write default tmp data
	writeDefaultTmpData()
}

func writeDefaultTmpData() {
	TmpStorage.Set("contents", make(map[string][]int))
}

// Init does model initialization.
// If first run, write default data.
// v means app.Version number. It's needed for version data.
func Init(v int) {
	storage.AppVersion = v
	storage.Storage = new(storage.JsonStorage)
	Storage.Init("data")
	storage.TmpStorage = new(storage.JsonStorage)
	storage.TmpStorage.dir = "tmp/data"
	if !Storage.Has("version") {
		os.Mkdir(Storage.dir, os.ModePerm)
		os.Mkdir(path.Join(Storage.dir, "content"), os.ModePerm)
		os.Mkdir(path.Join(Storage.dir, "plugin"), os.ModePerm)
		writeDefaultData()
	}
}

// All loads all data from storage to memory.
// Start timers for content, comment and message.
func All() {
	loadAllData()
	// generate indexes
	SyncIndexes()
	// start model timer, do all timer stuffs
	timer.Start()
}

func SyncIndexes() {
	// generate indexes
	content.UpdatePublishIndex()
	content.UpdateTmpIndex()
}

// SyncAll writes all current memory data to storage files.
func SyncAll() {
	content.Sync()
	message.Sync()
	file.Sync()
	content.SyncReaders()
	user.Sync()
	user.SyncTokens()
	setting.Sync()
	setting.SyncNavigators()
	setting.SyncVersion()
}
