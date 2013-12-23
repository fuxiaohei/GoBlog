
package handler

import (
	"github.com/fuxiaohei/GoInk/Core"
	"github.com/fuxiaohei/GoBlog/app/model"
)

func AdminArticleNew(context *Core.Context)interface {}{
	context.Render("admin:admin/article_new.html", map[string]interface {}{
			"Title":"写文章",
			"Categories":model.CategoryM.GetAll(),
		})
	return nil
}
