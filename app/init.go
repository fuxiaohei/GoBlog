package app

import (
	"github.com/fuxiaohei/GoInk"
	GoInkDb "github.com/fuxiaohei/GoInk/Db"
	_ "github.com/mattn/go-sqlite3"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"os"
)

var (
	Ink *GoInk.Simple
	Db  *GoInkDb.Engine
)

func Init() {
	var err error

	// create *GoInk.Simple application
	Ink, err = GoInk.NewSimple("config.json")
	if err != nil {
		panic(err)
		os.Exit(1)
	}

	// create *Db.Engine object
	Db, err = GoInkDb.NewEngine("sqlite3", "sqlite.db", 20, 30)
	if err != nil {
		panic(err)
		os.Exit(1)
	}

	// init something
	Ink.View.NewFunc("DateInt64", utils.DateInt64)
	Ink.View.NewFunc("DateString", utils.DateString)
	Ink.View.NewFunc("DateTime", utils.DateTime)
	Ink.View.NewFunc("Now", utils.Now)
}
