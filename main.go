package main

import (
	. "github.com/fuxiaohei/gorink/app"
	_ "github.com/fuxiaohei/gorink/controller"
)

func main() {
	App.Listen()
}
