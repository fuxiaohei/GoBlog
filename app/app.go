package app

import (
	"fmt"
	"github.com/fuxiaohei/GoBlog/app/handler"
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoInk"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime/debug"
	"strings"
	"syscall"
)

var (
	// APP VERSION, as date version
	VERSION = 20140228
	// Global GoInk application
	App              *GoInk.App
	staticFileSuffix = ".css,.js,.jpg,.jpeg,.png,.gif,.ico,.xml,.zip,.txt,.html,.otf,.svg,.eot,.woff,.ttf,.doc,.ppt,.xls,.docx,.pptx,.xlsx,.xsl"
	uploadFileSuffix = ".jpg,.png,.gif,.zip,.txt,.doc,.docx,.xls,.xlsx,.ppt,.pptx"
)

func init() {
	// init application
	App = GoInk.New()

	// init some settings
	App.Config().StringOr("app.static_dir", "static")
	App.Config().StringOr("app.log_dir", "tmp/log")
	os.MkdirAll(App.Get("log_dir"), os.ModePerm)
	os.MkdirAll("tmp/data", os.ModePerm)

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
