package helper

import (
	"github.com/fuxiaohei/GoBlog/app/model/content"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"github.com/fuxiaohei/GoInk"
)

func ContentHtml(ctx *GoInk.Context, cnt *content.Content) string {
	if cnt.Status != "publish" {
		return ""
	}
	thm := Theme(ctx)
	if !thm.HasSection("article.html") {
		return ""
	}
	return thm.Section("article", map[string]interface{}{
		"Article": cnt,
	})
}

func ArticleListHtml(ctx *GoInk.Context, cnt []*content.Content, pager *utils.Pager) string {
	thm := Theme(ctx)
	if !thm.HasSection("articles.html") {
		return ""
	}
	return thm.Section("articles", map[string]interface{}{
		"Articles": cnt,
		"Pager":    pager,
	})
}
