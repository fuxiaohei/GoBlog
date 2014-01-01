package handler

import (
	"github.com/fuxiaohei/GoInk/Core"
	"strconv"
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"strings"
)

func AdminComment(context *Core.Context) interface{} {
	page := 1
	if context.Param(2) == "page" {
		page, _ = strconv.Atoi(context.Param(3))
		if page < 1 {
			page = 1
		}
	}
	comments, pager := model.CommentM.GetPaged(page, 8, false)
	context.Render("admin:admin/comment.html", map[string]interface{}{
			"Title":     "评论",
			"IsComment": true,
			"Comments":comments,
			"Pager":pager,
		})
	return nil
}

func AdminCommentStatusPost(context *Core.Context) interface {} {
	status := context.String("status")
	id := context.Int("id")
	if len(status) < 1 || id < 1 {
		context.Json(map[string]interface {}{
			"res":false,
			"msg":"审核失败",
		})
		return nil
	}
	model.CommentM.ChangeCommentStatus(id, status)
	model.ArticleM.CountComments()
	context.Json(map[string]interface {}{
		"res":true,
	})
	return nil
}

func AdminCommentReplyPost(context *Core.Context) interface {} {
	pid := context.Int("pid")
	pComment := model.CommentM.GetCommentById(pid)
	if pComment == nil {
		context.Json(map[string]interface {}{
			"res":false,
			"msg":"回复失败",
		})
		return nil
	}
	uid, _ := strconv.Atoi(context.Cookie("admin-user"))
	user := model.UserM.GetUserById(uid)
	data := context.Input()
	c := new(model.Comment)
	c.Author = user.Display
	c.Email = user.Email
	c.Site = user.Site
	c.Content = strings.Replace(data["content"], "\n", "<br/>", -1)
	c.ContentId = pComment.ContentId
	c.Pid = pid
	c.Avatar = utils.Gravatar(c.Email, "50")
	c.UserId = user.Id
	c.IsAdmin = true
	c = model.CommentM.CreateComment(c)
	context.Json(map[string]interface {}{
		"res":true,
		"comment":c,
	})
	return nil
}

func AdminCommentDeletePost(context *Core.Context) interface {} {
	id := context.Int("id")
	if id < 1 {
		context.Json(map[string]interface {}{
			"res":false,
			"msg":"删除",
		})
		return nil
	}
	model.CommentM.DeleteComment(id)
	model.ArticleM.CountComments()
	context.Json(map[string]interface {}{
		"res":true,
	})
	return nil
}
