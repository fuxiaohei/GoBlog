package handler

import "github.com/fuxiaohei/GoBlog/app"

func Init() {
	// define layouts
	e := app.Ink.View.NewLayout("admin", "admin/admin.layout")
	if e != nil {
		app.Ink.Crash(e)
	}
	e = app.Ink.View.NewLayout("theme", "default/index.layout")
	if e != nil {
		app.Ink.Crash(e)
	}

	// login handlers
	app.Ink.Router.Get("/login.html", Login)
	app.Ink.Router.Get("/logout.html", Logout)
	app.Ink.Router.Post("/login.html", LoginPost)
	app.Ink.Listener.AddListener("server.dynamic.before", "auth", authAdmin)

	// admin handlers
	app.Ink.Router.Get("/admin", Admin)

	// admin profile handlers
	app.Ink.Router.Get("/admin/profile", AdminProfile)
	app.Ink.Router.Post("/admin/profile", AdminProfilePost)
	app.Ink.Router.Get("/admin/password", AdminPassword)
	app.Ink.Router.Post("/admin/password", AdminPasswordPost)

	// admin category handlers
	app.Ink.Router.Get("/admin/category", AdminCategory)
	app.Ink.Router.Get("/admin/category/new", AdminCategoryNew)
	app.Ink.Router.Post("/admin/category/new", AdminCategoryNewPost)
	app.Ink.Router.Get("/admin/category/edit/", AdminCategoryEdit)
	app.Ink.Router.Post("/admin/category/edit/", AdminCategoryEditPost)
	app.Ink.Router.Post("/admin/category", AdminCategoryDelete)

	// admin article handler
	app.Ink.Router.Get("/admin/article", AdminArticle)
	app.Ink.Router.Get("/admin/article/new", AdminArticleNew)
	app.Ink.Router.Post("/admin/article/new", AdminArticleNewPost)
	app.Ink.Router.Get("/admin/article/edit", AdminArticleEdit)
	app.Ink.Router.Post("/admin/article/edit", AdminArticleEditPost)

	// admin comment handler
	app.Ink.Router.Get("/admin/comment", AdminComment)

	// article handler
	app.Ink.Router.Get("/", Article)
	app.Ink.Router.Get("/article", Article)
	app.Ink.Router.Post("/article", ArticleCommentPost)
	app.Ink.Router.Get("/category", ArticleCategory)
}
