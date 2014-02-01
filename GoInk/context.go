package GoInk

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	CONTEXT_RENDERED = "context_rendered"
	CONTEXT_END      = "context_end"
	CONTEXT_SEND     = "context_send"
)

type Context struct {
	Request    *http.Request
	Base       string
	Url        string
	RequestUrl string
	Method     string
	Ip         string
	UserAgent  string
	Referer    string
	Host       string
	Ext        string
	IsSSL      bool
	IsAjax     bool

	Response http.ResponseWriter
	Status   int
	Header   map[string]string
	Body     []byte

	routeParams map[string]string
	flashData   map[string]interface{}

	eventsFunc map[string]reflect.Value

	IsSend bool
	IsEnd  bool

	app    *App
	layout string
}

func NewContext(app *App, res http.ResponseWriter, req *http.Request) *Context {
	context := new(Context)
	context.flashData = make(map[string]interface{})
	context.eventsFunc = make(map[string]reflect.Value)
	context.IsSend = false
	context.IsEnd = false

	context.Request = req
	context.Url = req.URL.Path
	context.RequestUrl = req.RequestURI
	context.Method = req.Method
	context.Ext = path.Ext(req.URL.Path)
	context.Host = req.Host
	context.Ip = strings.Split(req.RemoteAddr, ":")[0]
	context.IsAjax = req.Header.Get("X-Requested-With") == "XMLHttpRequest"
	context.IsSSL = req.TLS != nil
	context.Referer = req.Referer()
	context.UserAgent = req.UserAgent()
	context.Base = "://" + context.Host + "/"
	if context.IsSSL {
		context.Base = "https" + context.Base
	} else {
		context.Base = "http" + context.Base
	}

	context.Response = res
	context.Status = 200
	context.Header = make(map[string]string)
	context.Header["Content-Type"] = "text/html;charset=UTF-8"

	context.app = app

	req.ParseForm()

	return context
}

func (ctx *Context) Param(key string) string {
	return ctx.routeParams[key]
}

func (ctx *Context) Flash(key string, v ...interface{}) interface{} {
	if len(v) > 0 {
		return ctx.flashData[key]
	}
	ctx.flashData[key] = v[0]
	return nil
}

func (ctx *Context) On(e string, fn interface{}) {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		println("only support function type for Context.On method")
		return
	}
	ctx.eventsFunc[e] = reflect.ValueOf(fn)
}

func (ctx *Context) Do(e string, args ...interface{}) []interface{} {
	_, ok := ctx.eventsFunc[e]
	if !ok {
		return nil
	}
	if !ctx.eventsFunc[e].IsValid() {
		println("invalid function call for Context.Do(" + e + ")")
		return nil
	}
	fn := ctx.eventsFunc[e]
	numIn := fn.Type().NumIn()
	if numIn > len(args) {
		println("not enough parameters for Context.Do(" + e + ")")
		return nil
	}
	rArgs := make([]reflect.Value, numIn)
	for i := 0; i < numIn; i++ {
		rArgs[i] = reflect.ValueOf(args[i])
	}
	resValue := fn.Call(rArgs)
	if len(resValue) < 1 {
		return nil
	}
	res := make([]interface{}, len(resValue))
	for i, v := range resValue {
		res[i] = v.Interface()
	}
	return res
}

func (ctx *Context) Input() map[string]string {
	data := make(map[string]string)
	for key, v := range ctx.Request.Form {
		data[key] = v[0]
	}
	return data
}

func (ctx *Context) Strings(key string) []string {
	return ctx.Request.Form[key]
}

func (ctx *Context) String(key string) string {
	return ctx.Request.FormValue(key)
}

func (ctx *Context) StringOr(key string, def string) string {
	value := ctx.String(key)
	if value == "" {
		return def
	}
	return value
}

func (ctx *Context) Int(key string) int {
	str := ctx.String(key)
	i, _ := strconv.Atoi(str)
	return i
}

func (ctx *Context) IntOr(key string, def int) int {
	i := ctx.Int(key)
	if i == 0 {
		return def
	}
	return i
}

func (ctx *Context) Float(key string) float64 {
	str := ctx.String(key)
	f, _ := strconv.ParseFloat(str, 64)
	return f
}

func (ctx *Context) FloatOr(key string, def float64) float64 {
	f := ctx.Float(key)
	if f == 0.0 {
		return def
	}
	return f
}

func (ctx *Context) Bool(key string) bool {
	str := ctx.String(key)
	b, _ := strconv.ParseBool(str)
	return b
}

func (ctx *Context) Cookie(key string, value ...string) string {
	if len(value) < 1 {
		c, e := ctx.Request.Cookie(key)
		if e != nil {
			return ""
		}
		return c.Value
	}
	if len(value) == 2 {
		t := time.Now()
		expire, _ := strconv.Atoi(value[1])
		t = t.Add(time.Duration(expire) * time.Second)
		cookie := &http.Cookie{
			Name:    key,
			Value:   value[0],
			Path:    "/",
			MaxAge:  expire,
			Expires: t,
		}
		http.SetCookie(ctx.Response, cookie)
		return ""
	}
	return ""
}

func (ctx *Context) GetHeader(key string) string {
	return ctx.Request.Header.Get(key)
}

func (ctx *Context) Redirect(url string, status ...int) {
	ctx.Header["Location"] = url
	if len(status) > 0 {
		ctx.Status = status[0]
		return
	}
	ctx.Status = 302
}

func (ctx *Context) ContentType(contentType string) {
	ctx.Header["Content-Type"] = contentType
}

func (ctx *Context) Json(data interface{}) {
	bytes, e := json.MarshalIndent(data, "", "    ")
	if e != nil {
		panic(e)
	}
	ctx.ContentType("application/json;charset=UTF-8")
	ctx.Body = bytes
}

func (ctx *Context) Send() {
	if ctx.IsSend {
		return
	}
	for name, value := range ctx.Header {
		ctx.Response.Header().Set(name, value)
	}
	ctx.Response.WriteHeader(ctx.Status)
	ctx.Response.Write(ctx.Body)
	ctx.IsSend = true
	ctx.Do(CONTEXT_SEND)
}

func (ctx *Context) End() {
	if ctx.IsEnd {
		return
	}
	if !ctx.IsSend {
		ctx.Send()
	}
	ctx.IsEnd = true
	ctx.Do(CONTEXT_END)
}

func (ctx *Context) Throw(status int, message ...interface{}) {
	e := strconv.Itoa(status)
	ctx.Status = status
	ctx.Do(e, message...)
	ctx.End()
}

func (ctx *Context) Layout(str string) {
	ctx.layout = str
}

func (ctx *Context) Tpl(tpl string, data map[string]interface{}) string {
	b, e := ctx.app.view.Render(tpl+".html", data)
	if e != nil {
		panic(e)
	}
	return string(b)
}

func (ctx *Context) Render(tpl string, data map[string]interface{}) {
	b, e := ctx.app.view.Render(tpl+".html", data)
	if e != nil {
		panic(e)
	}
	if ctx.layout != "" {
		l, e := ctx.app.view.Render(ctx.layout+".layout", data)
		if e != nil {
			panic(e)
		}
		b = bytes.Replace(l, []byte("{@Content}"), b, -1)
	}
	ctx.Body = b
	ctx.Do(CONTEXT_RENDERED)
}

func (ctx *Context) Func(name string, fn interface{}) {
	ctx.app.view.FuncMap[name] = fn
}

func (ctx *Context) App() *App {
	return ctx.app
}

func (ctx *Context) Download(file string) {
	f, e := os.Stat(file)
	if e != nil {
		ctx.Status = 404
		return
	}
	if f.IsDir() {
		ctx.Status = 403
		return
	}
	output := ctx.Response.Header()
	output.Set("Content-Type", "application/octet-stream")
	output.Set("Content-Disposition", "attachment; filename="+path.Base(file))
	output.Set("Content-Transfer-Encoding", "binary")
	output.Set("Expires", "0")
	output.Set("Cache-Control", "must-revalidate")
	output.Set("Pragma", "public")
	http.ServeFile(ctx.Response, ctx.Request, file)
	ctx.IsSend = true
}
