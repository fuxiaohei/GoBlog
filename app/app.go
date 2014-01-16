package app

import (
	"fmt"
	"github.com/fuxiaohei/GoBlog/GoInk"
	"github.com/fuxiaohei/GoBlog/app/handler"
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"strings"
)

var (
	VERSION          = 20140116
	App              *GoInk.App
	staticFileSuffix = ".css,.js,.jpg,.jpeg,.png,.gif,.ico,.xml,.zip,.txt,.html,.otf,.svg,.eot,.woff,.ttf,.doc,.ppt,.xls,.docx,.pptx,.xlsx"
	uploadFileSuffix = ".jpg,.png,.gif,.zip,.txt,.doc,.docx,.xls,.xlsx,.ppt,.pptx"
)

func init() {
	// init application
	App = GoInk.New()

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
}

func Init() {

	// init some settings
	if App.Get("upload_files") == "" {
		App.Set("upload_files", uploadFileSuffix)
	}
	if App.Get("upload_size") == "" {
		App.Set("upload_size", 1024*1024 * 10)
	}
	if App.Get("upload_dir") == "" {
		App.Set("upload_dir", "static/upload")
		os.MkdirAll("static/upload", os.ModePerm)
	}

	// init temp dir
	if App.Get("log_dir") == "" {
		App.Set("log_dir", "tmp/log")
		os.MkdirAll("tmp/log", os.ModePerm)
	}

	// set static files handler
	if App.Get("static_files") != "" {
		staticFileSuffix = App.Get("static_files")
	}
	App.Static(func(context *GoInk.Context) {
		static := App.Config().StringOr("static_dir", "static")
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
		if App.View().Has("error/error.html") {
			context.Layout("")
			context.Render("error/error", map[string]interface{}{
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
		if App.View().Has("error/notfound.html") {
			context.Layout("")
			context.Render("error/notfound", map[string]interface{}{
					"context": context,
				})
		}
		context.End()
	})

	// init storage
	model.Init()

	// set version
	model.SetVersion(VERSION)
}

func registerAdminHandler() {
	// add admin handlers
	App.Get("/admin/", handler.Auth, handler.Admin)

	App.Get("/admin/profile/", handler.Auth, handler.AdminProfile)
	App.Post("/admin/profile/", handler.Auth, handler.AdminProfile)

	App.Get("/admin/password/", handler.Auth, handler.AdminPassword)
	App.Post("/admin/password/", handler.Auth, handler.AdminPassword)

	App.Get("/admin/article/write/", handler.Auth, handler.ArticleWrite)
	App.Post("/admin/article/write/", handler.Auth, handler.ArticleWrite)
	App.Get("/admin/articles/", handler.Auth, handler.AdminArticle)
	App.Get("/admin/article/:id/", handler.Auth, handler.ArticleEdit)
	App.Post("/admin/article/:id/", handler.Auth, handler.ArticleEdit)
	App.Delete("/admin/article/:id/", handler.Auth, handler.ArticleEdit)

	App.Get("/admin/page/write/", handler.Auth, handler.PageWrite)
	App.Post("/admin/page/write/", handler.Auth, handler.PageWrite)
	App.Get("/admin/pages/", handler.Auth, handler.AdminPage)
	App.Get("/admin/page/:id/", handler.Auth, handler.PageEdit)
	App.Post("/admin/page/:id/", handler.Auth, handler.PageEdit)
	App.Delete("/admin/page/:id/", handler.Auth, handler.PageEdit)

	App.Get("/admin/comments/", handler.Auth, handler.AdminComments)
	App.Delete("/admin/comments/", handler.Auth, handler.AdminComments)
	App.Put("/admin/comments/", handler.Auth, handler.AdminComments)
	App.Post("/admin/comments/", handler.Auth, handler.AdminComments)

	App.Get("/admin/setting/", handler.Auth, handler.AdminSetting)
	App.Post("/admin/setting/", handler.Auth, handler.AdminSetting)

	App.Get("/admin/setting/custom/", handler.Auth, handler.CustomSetting)
	App.Post("/admin/setting/custom/", handler.Auth, handler.CustomSetting)

	App.Get("/admin/files/", handler.Auth, handler.AdminFiles)
	App.Delete("/admin/files/", handler.Auth, handler.AdminFiles)
	App.Post("/admin/files/upload/", handler.Auth, handler.FileUpload)
}

func registerCmdHandler() {
	App.Get("/cmd/backup/", handler.Auth, handler.CmdBackup)
	App.Post("/cmd/backup/", handler.Auth, handler.CmdBackup)
	App.Delete("/cmd/backup/", handler.Auth, handler.CmdBackup)
}

func registerHomeHandler() {
	App.Get("/login/", handler.Login)
	App.Post("/login/", handler.Login)
	App.Get("/logout/", handler.Logout)

	App.Get("/article/:id/:slug", handler.Article)
	App.Get("/p/:page/", handler.Home)
	App.Post("/comment/:id/", handler.Comment)

	App.Get("/rss", handler.Feed)

	App.Get("/:slug", handler.TopPage)
	App.Get("/", handler.Home)
}

func Run() {

	App.View().FuncMap["DateInt64"] = utils.DateInt64
	App.View().FuncMap["DateString"] = utils.DateString
	App.View().FuncMap["DateTime"] = utils.DateTime
	App.View().FuncMap["Now"] = utils.Now
	App.View().FuncMap["Html2str"] = utils.Html2str
	App.View().FuncMap["FileSize"] = utils.FileSize
	App.View().FuncMap["Setting"] = model.GetSetting

	registerAdminHandler()
	registerCmdHandler()
	registerHomeHandler()

	App.Run()
}
