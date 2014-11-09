package model

import (
	"github.com/fuxiaohei/GoBlog/gof"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var DB *xorm.Engine

func InitDB(cfg gof.ConfigInterface) error {
	// load config items
	driver := cfg.String("database.driver")
	file := cfg.String("database.file")

	// connect database
	var err error
	DB, err = xorm.NewEngine(driver, file)
	if err != nil {
		return err
	}

	// set db options
	DB.ShowSQL = true

	// try connect really
	err = DB.Ping()
	return err
}

func CreateDB(cfg gof.ConfigInterface) error {
	// connect database
	var err error
	if err = InitDB(cfg); err != nil {
		return err
	}

	// sync table schema
	err = DB.Sync2(new(User), new(UserToken), new(Setting), new(Version), new(Message), new(Attach),
		new(Article), new(ArticleCategory), new(ArticleTag), new(Page),
		new(Comment), new(CommentArticle), new(CommentPage))
	return err
}
