package app

import (
	"github.com/fuxiaohei/GoInk"
	GoInkDb "github.com/fuxiaohei/GoInk/Db"
	_ "github.com/mattn/go-sqlite3"
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

}
