package app

import "github.com/fuxiaohei/GoBlog/app/handler"

func registerHomeHandler() {
	App.Route("GET,POST", "/login/", handler.Login)
	App.Get("/logout/", handler.Logout)

	App.Get("/article/:id/:slug", handler.Article)
	App.Get("/page/:id/:slug", handler.Page)
	App.Get("/p/:page/", handler.Articles)
	App.Post("/comment/:id/", handler.Comment)
	App.Get("/tag/:tag/", handler.TagArticles)
	App.Get("/tag/:tag/p/:page/", handler.TagArticles)
	App.Get("/robots", handler.Robots)

	App.Get("/redirect/", handler.Redirect)
	App.Get("/feed/", handler.Rss)
	App.Get("/sitemap", handler.SiteMap)
	App.Get("/upload/:id/:name", handler.Upload)

	App.Get("/:slug", handler.TopPage)
	App.Get("/", handler.Articles)
}

// Run begins Fxh.Go http server.
func Run() {

	registerAdminHandler()
	registerCmdHandler()
	registerHomeHandler()

	App.Run()
}
