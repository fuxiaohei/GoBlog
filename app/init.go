package app

import (
	"github.com/fuxiaohei/GoBlog/app/handler/helper"
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoBlog/app/model/setting"
	"github.com/fuxiaohei/GoBlog/app/plugin"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"strconv"
)

// Init starts Fxh.Go application preparation.
// Load models and plugins, update views.
func Init() {

	// init storage
	model.Init(VERSION)

	// load all data
	model.All()

	// init plugin
	plugin.Init()

	// update plugin handlers
	plugin.Update(App)

	// utils view func
	App.View().FuncMap["DateInt64"] = utils.DateInt64
	App.View().FuncMap["DateString"] = utils.DateString
	App.View().FuncMap["DateTime"] = utils.DateTime
	App.View().FuncMap["Now"] = utils.Now
	App.View().FuncMap["Html2str"] = utils.Html2str
	App.View().FuncMap["FileSize"] = utils.FileSize
	App.View().FuncMap["Md2html"] = utils.Markdown2HtmlTemplate
	App.View().FuncMap["Redirect"] = func(url string) string {
		return "/redirect?to=" + url
	}
	/*App.View().FuncMap["Render"] = func(str string) string {
		bytes, e := App.View().RenderString(str, nil)
		if e != nil {
			println(e.Error())
			return str
		}
		return string(bytes)
	}*/

	// model view func
	App.View().FuncMap["Setting"] = setting.Get
	App.View().FuncMap["Navigator"] = setting.GetNavigators

	// helper view func
	App.View().FuncMap["CommentHtml"] = helper.CommentHtml
	App.View().FuncMap["SidebarHtml"] = helper.SidebarHtml
	App.View().FuncMap["ContentHtml"] = helper.ContentHtml
	App.View().FuncMap["ArticleListHtml"] = helper.ArticleListHtml

	App.View().IsCache = (setting.Get("theme_cache") == "true")

	println("app version @ " + strconv.Itoa(setting.GetVersion().Version))
}
