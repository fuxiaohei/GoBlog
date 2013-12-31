package handler

import (
	"github.com/fuxiaohei/GoInk/Core"
	"strconv"
)

func AdminComment(context *Core.Context) interface {} {
	page := 1
	if context.Param(2) == "page" {
		page, _ = strconv.Atoi(context.Param(3))
		if page < 1 {
			page = 1
		}
	}
	context.Render("admin:admin/comment.html", map[string]interface{}{
			"Title":     "评论",
			"IsComment": true,
		})
	return nil
}

