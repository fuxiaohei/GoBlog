package main

import (
	"github.com/fuxiaohei/GoBlog/app"
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoBlog/app/handler"
)

func main() {
	// init application
	app.Init()

	// init model
	model.Init()

	// init handler
	handler.Init()

	// run *GoInk.Simple application
	app.Ink.Run()
}

