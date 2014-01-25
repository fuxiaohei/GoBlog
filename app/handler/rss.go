package handler

import (
	"github.com/fuxiaohei/GoBlog/GoInk"
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"github.com/gorilla/feeds"
	"path"
	"time"
)

func Feed(context *GoInk.Context) {
	baseUrl := model.GetSetting("site_url")
	article, _ := model.GetArticleList(1, 20)
	feed := new(feeds.Feed)
	feed.Title = model.GetSetting("site_title")
	feed.Link = &feeds.Link{Href: baseUrl}
	feed.Description = model.GetSetting("site_description")
	author := model.GetUsersByRole("ADMIN")[0]
	feed.Author = &feeds.Author{author.Nick, author.Email}
	feed.Items = make([]*feeds.Item, 0)
	var create int64
	if len(article) > 0 {
		create = article[0].EditTime
	} else {
		create = utils.Now()
	}
	feed.Created = time.Unix(create, 0)
	for _, a := range article {
		item := new(feeds.Item)
		item.Title = a.Title
		item.Link = &feeds.Link{Href: path.Join(baseUrl, a.Link())}
		item.Author = feed.Author
		item.Created = time.Unix(a.CreateTime, 0)
		item.Description = utils.Markdown2Html(a.Summary())
		feed.Items = append(feed.Items, item)
	}
	str, e := feed.ToRss()
	if e != nil {
		panic(e)
	}
	context.ContentType("application/rss+xml;charset=UTF-8")
	context.Body = []byte(str)
}
