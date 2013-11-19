package app

import (
	"github.com/fuxiaohei/goink"
	"github.com/fuxiaohei/goink/app"
	"github.com/fuxiaohei/goink/db"
	_ "github.com/mattn/go-sqlite3"
)

const (
	VERSION = "0.5"
)

var (
	App *app.InkApp
	Db  *db.InkDatabase
	Orm *db.InkOrm
)

func init() {
	App = goink.NewApp("config.json")
	App.Static(App.Config().StringOr("view.static", "static"))
	App.On("router.run.null@default", func(context *app.InkContext) {
			context.Render(App.Config().StringOr("view.404", "404.html"), nil)
			context.Send("", 404)
		})
	initDatabase()
}

func initDatabase() {
	options := &db.InkDatabaseOption{}
	options.Driver = App.String("database.driver")
	options.Dsn = App.String("database.dsn")
	options.MaxConnection = App.Config().IntOr("database.max_conn", 10)
	options.IdleConnection = App.Config().IntOr("database.idle_conn", 5)
	options.Mode = App.Config().StringOr("database.mode", app.MODE_DEBUG)
	var err error
	Db, err = db.NewDatabase(options, App)
	if err != nil {
		App.Crash(err)
	}
	Orm = db.NewOrm(Db)
}
