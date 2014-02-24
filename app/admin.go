package app

import "github.com/fuxiaohei/GoBlog/app/handler"

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
	App.Post("/admin/setting/nav/", handler.Auth, handler.NavigatorSetting)

	App.Route("GET,DELETE", "/admin/files/", handler.Auth, handler.AdminFiles)
	App.Post("/admin/files/upload/", handler.Auth, handler.FileUpload)

	App.Route("GET,POST", "/admin/plugins/", handler.Auth, handler.AdminPlugin)
	App.Route("GET,POST", "/admin/plugins/:plugin_key/", handler.Auth, handler.PluginSetting)

	App.Post("/admin/message/read/", handler.Auth, handler.AdminMessageRead)
}
