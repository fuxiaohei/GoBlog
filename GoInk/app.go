package GoInk

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
)

const ()

type App struct {
	router *Router
	view   *View
	middle []Handler
	inter map[string]Handler
	config *Config
}

func New() *App {
	a := new(App)
	a.router = NewRouter()
	a.middle = make([]Handler, 0)
	a.inter = make(map[string]Handler)
	a.config, _ = NewConfig("config.json")
	a.view = NewView(a.config.StringOr("app.view_dir", "view"))

	// add empty handler
	/*a.Get("/", func(context *Context) {
			context.Body = []byte("It Works!")
		})*/
	return a
}

func (app *App) Use(h Handler) {
	app.middle = append(app.middle, h)
}

func (app *App) Config() *Config {
	return app.config
}

func (app *App) View() *View {
	return app.view
}

func (app *App) handler(res http.ResponseWriter, req *http.Request) {
	context := NewContext(app, res, req)

	defer func() {
		e := recover()
		if e == nil {
			return
		}
		context.Body = []byte(fmt.Sprint(e))
		context.Status = 503
		println(string(context.Body))
		debug.PrintStack()
		if _, ok := app.inter["recover"]; ok {
			app.inter["recover"](context)
		}
		if !context.IsEnd {
			context.End()
		}
	}()

	if _, ok := app.inter["static"]; ok {
		app.inter["static"](context)
		if context.IsEnd {
			return
		}
	}

	if len(app.middle) > 0 {
		for _, h := range app.middle {
			h(context)
			if context.IsEnd {
				break
			}
		}
	}
	params, fn := app.router.Find(req.URL.Path, req.Method)
	if params != nil && fn != nil {
		context.routeParams = params
		for _, f := range fn {
			f(context)
			if context.IsEnd {
				break
			}
		}
		if !context.IsEnd {
			context.End()
		}
	} else {
		println("router is missing at "+req.URL.Path)
		context.Status = 404
		if _, ok := app.inter["notfound"]; ok {
			app.inter["notfound"](context)
			if !context.IsEnd {
				context.End()
			}
		}else {
			context.Throw(404)
		}
	}
}

func (app *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	app.handler(res, req)
}

func (app *App) Run() {
	addr := app.config.StringOr("app.server", "localhost:9000")
	println("http server run at "+addr)
	e := http.ListenAndServe(addr, app)
	panic(e)
}

func (app *App) Set(key string, v interface {}) {
	app.config.Set("app." + key, v)
}

func (app *App) Get(key string, fn ...Handler) string {
	if len(fn) > 0 {
		app.router.Get(key, fn...)
		return ""
	}
	return app.config.String("app."+key)
}

func (app *App) Post(key string, fn ...Handler) {
	app.router.Post(key, fn...)
}

func (app *App) Put(key string, fn ...Handler) {
	app.router.Put(key, fn...)
}

func (app *App) Delete(key string, fn ...Handler) {
	app.router.Delete(key, fn...)
}

func (app *App) Route(method string, key string, fn ...Handler) {
	methods := strings.Split(method, ",")
	for _, m := range methods {
		switch m{
		case "GET":
			app.Get(key, fn...)
		case "POST":
			app.Post(key, fn...)
		case "PUT":
			app.Put(key, fn...)
		case "DELETE":
			app.Delete(key, fn...)
		default:
			println("unknow route method "+m)
		}
	}
}

func (app *App) Static(h Handler) {
	app.inter["static"] = h
}

func (app *App) Recover(h Handler) {
	app.inter["recover"] = h
}

func (app *App) NotFound(h Handler) {
	app.inter["notfound"] = h
}
