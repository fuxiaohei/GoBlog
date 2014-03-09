package helper

import (
	"github.com/fuxiaohei/GoBlog/app/model/content"
	"github.com/fuxiaohei/GoBlog/app/model/setting"
	"github.com/fuxiaohei/GoInk"
)

// SidebarHtml returns rendered sidebar template html.
func SidebarHtml(context *GoInk.Context) string {
	thm := Theme(context)
	if !thm.HasSection("sidebar.html") {
		return ""
	}
	popSize := setting.Int("popular_size")
	if popSize < 1 {
		popSize = 4
	}
	cmtSize := setting.Int("recent_comment_size")
	if cmtSize < 1 {
		cmtSize = 3
	}
	return thm.Section("sidebar", map[string]interface{}{
		"Popular":       content.PopularArticleList(popSize),
		"RecentComment": content.RecentCommentList(cmtSize),
		"Tags":          content.GetContentTags(),
	})
}
