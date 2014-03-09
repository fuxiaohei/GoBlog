package helper

import (
	"github.com/fuxiaohei/GoBlog/app/model/content"
	"github.com/fuxiaohei/GoInk"
)

func CommentHtml(context *GoInk.Context, c *content.Content) string {
	if c.Type == "page" && !c.IsComment {
		// hide comment if page's comment is closed
		return ""
	}
	thm := Theme(context)
	if !thm.HasSection("comment.html") {
		return ""
	}
	return thm.Section("comment", map[string]interface{}{
		"Content":  c,
		"Comments": c.Comments,
	})
}
