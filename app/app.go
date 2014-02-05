package app

import (
	"fmt"
	"github.com/fuxiaohei/GoBlog/GoInk"
	"github.com/fuxiaohei/GoBlog/app/handler"
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoBlog/app/plugin"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
)

var (
	VERSION          = 20140131
	App              *GoInk.App
	staticFileSuffix = ".css,.js,.jpg,.jpeg,.png,.gif,.ico,.xml,.zip,.txt,.html,.otf,.svg,.eot,.woff,.ttf,.doc,.ppt,.xls,.docx,.pptx,.xlsx"
	uploadFileSuffix = ".jpg,.png,.gif,.zip,.txt,.doc,.docx,.xls,.xlsx,.ppt,.pptx"
)

func init() {
	// init application
	App = GoInk.New()

	// init some settings
	App.Config().StringOr("app.static_dir", "static")
	App.Config().StringOr("app.log_dir", "tmp/log")
	os.MkdirAll(App.Get("log_dir"), os.ModePerm)

	App.Config().IntOr("app.upload_size", 1024*1024*10)
	App.Config().StringOr("app.upload_files", uploadFileSuffix)
	App.Config().StringOr("app.upload_dir", path.Join(App.Get("static_dir"), "upload"))
	os.MkdirAll(App.Get("upload_dir"), os.ModePerm)

	if App.Get("static_files") != "" {
		staticFileSuffix = App.Get("static_files")
	}

	App.Static(func(context *GoInk.Context) {
		static := App.Config().String("app.static_dir")
		url := strings.TrimPrefix(context.Url, "/")
		if url == "favicon.ico" {
			url = path.Join(static, url)
		}
		if !strings.HasPrefix(url, static) {
			return
		}
		if !strings.Contains(staticFileSuffix, context.Ext) {
			context.Status = 403
			context.End()
			return
		}
		f, e := os.Stat(url)
		if e == nil {
			if f.IsDir() {
				context.Status = 403
				context.End()
				return
			}
		}
		/*_, e := os.Stat(url)
		if e != nil {
			context.Throw(404)
			return
		}*/
		http.ServeFile(context.Response, context.Request, url)
		context.IsEnd = true
	})

	// set recover handler
	App.Recover(func(context *GoInk.Context) {
		go LogError(append(append(context.Body, []byte("\n")...), debug.Stack()...))
		theme := handler.Theme(context)
		if theme.Has("error/error.html") {
			theme.Layout("").Render("error/error", map[string]interface{}{
				"error":   string(context.Body),
				"stack":   string(debug.Stack()),
				"context": context,
			})
		} else {
			context.Body = append([]byte("<pre>"), context.Body...)
			context.Body = append(context.Body, []byte("\n")...)
			context.Body = append(context.Body, debug.Stack()...)
			context.Body = append(context.Body, []byte("</pre>")...)
		}
		context.End()
	})

	// set not found handler
	App.NotFound(func(context *GoInk.Context) {
		theme := handler.Theme(context)
		if theme.Has("error/notfound.html") {
			theme.Layout("").Render("error/notfound", map[string]interface{}{
				"context": context,
			})
		}
		context.End()
	})

	// add recover defer
	defer func() {
		e := recover()
		if e != nil {
			bytes := append([]byte(fmt.Sprint(e)+"\n"), debug.Stack()...)
			LogError(bytes)
			println("panic error, crash down")
			os.Exit(1)
		}
	}()

	// catch exit command
	go catchExit()
}

// code from https://github.com/Unknwon/gowalker/blob/master/gowalker.go
func catchExit() {
	sigTerm := syscall.Signal(15)
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, sigTerm)

	for {
		switch <-sig {
		case os.Interrupt, sigTerm:
			println("before exit, saving data")
			model.SyncAll()
			println("ready to exit")
			os.Exit(0)
		}
	}
}

func Init() {

	// init storage
	model.Init(VERSION)

	// load all data
	model.All()

	// init plugin
	plugin.Init()

	// update plugin handlers
	plugin.Update(App)

	App.View().FuncMap["DateInt64"] = utils.DateInt64
	App.View().FuncMap["DateString"] = utils.DateString
	App.View().FuncMap["DateTime"] = utils.DateTime
	App.View().FuncMap["Now"] = utils.Now
	App.View().FuncMap["Html2str"] = utils.Html2str
	App.View().FuncMap["FileSize"] = utils.FileSize
	App.View().FuncMap["Setting"] = model.GetSetting
	App.View().FuncMap["Md2html"] = utils.Markdown2HtmlTemplate

	println("app version @ " + strconv.Itoa(model.GetVersion().Version))
}

func registerAdminHandler() {
	// add admin handlers
	App.Get("/admin/", handler.Auth, handler.Admin)

	App.Route("GET,POST", "/admin/profile/", handler.Auth, handler.AdminProfile)

	App.Route("GET,POST", "/admin/password/", handler.Auth, handler.AdminPassword)

	App.Route("GET,POST", "/admin/article/write/", handler.Auth, handler.ArticleWrite)
	App.Get("/admin/articles/", handler.Auth, handler.AdminArticle)
	App.Route("GET,POST,DELETE", "/admin/article/:id/", handler.Auth, handler.ArticleEdit)

	App.Route("GET,POST", "/admin/page/write/", handler.Auth, handler.PageWrite)
	App.Get("/admin/pages/", handler.Auth, handler.AdminPage)
	App.Route("GET,POST,DELETE", "/admin/page/:id/", handler.Auth, handler.PageEdit)

	App.Route("GET,POST,PUT,DELETE", "/admin/comments/", handler.Auth, handler.AdminComments)

	App.Route("GET,POST", "/admin/setting/", handler.Auth, handler.AdminSetting)
	App.Post("/admin/setting/custom/", handler.Auth, handler.CustomSetting)

	App.Route("GET,DELETE", "/admin/files/", handler.Auth, handler.AdminFiles)
	App.Post("/admin/files/upload/", handler.Auth, handler.FileUpload)

	App.Route("GET,POST", "/admin/plugins/", handler.Auth, handler.AdminPlugin)
	App.Route("GET,POST", "/admin/plugins/:plugin_key/", handler.Auth, handler.PluginSetting)

	App.Post("/admin/message/read/", handler.Auth, handler.AdminMessageRead)
}

func registerCmdHandler() {
	App.Route("GET,POST,DELETE", "/cmd/backup/", handler.Auth, handler.CmdBackup)
	App.Get("/cmd/backup/file", handler.Auth, handler.CmdBackupFile)

	App.Route("GET,POST,DELETE", "/cmd/message/", handler.Auth, handler.CmdMessage)
	App.Route("GET,DELETE", "/cmd/logs/", handler.Auth, handler.CmdLogs)
}

func registerHomeHandler() {
	App.Route("GET,POST", "/login/", handler.Login)
	App.Get("/logout/", handler.Logout)

	App.Get("/article/:id/:slug", handler.Article)
	App.Get("/page/:id/:slug", handler.Page)
	App.Get("/p/:page/", handler.Home)
	App.Post("/comment/:id/", handler.Comment)

	App.Get("/rss/", handler.Feed)
	App.Get("/feed/", handler.Feed)

	App.Get("/:slug", handler.TopPage)
	App.Get("/", handler.Home)
}

func Run() {

	registerAdminHandler()
	registerCmdHandler()
	registerHomeHandler()

	App.Run()
}
