package handler

import (
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"github.com/fuxiaohei/GoInk"
	"strings"
	"time"
)

func SiteMap(ctx *GoInk.Context) {
	baseUrl := model.GetSetting("site_url")
	println(baseUrl)
	article, _ := model.GetPublishArticleList(1, 50)
	navigators := model.GetNavigators()
	now := time.Unix(utils.Now(), 0).Format(time.RFC3339)

	articleMap := make([]map[string]string, len(article))
	for i, a := range article {
		m := make(map[string]string)
		m["Link"] = strings.Replace(baseUrl+a.Link(), baseUrl+"/", baseUrl, -1)
		m["Created"] = time.Unix(a.CreateTime, 0).Format(time.RFC3339)
		articleMap[i] = m
	}

	navMap := make([]map[string]string, 0)
	for _, n := range navigators {
		m := make(map[string]string)
		if n.Link == "/" {
			continue
		}
		if strings.HasPrefix(n.Link, "/") {
			m["Link"] = strings.Replace(baseUrl+n.Link, baseUrl+"/", baseUrl, -1)
		} else {
			m["Link"] = n.Link
		}
		m["Created"] = now
		navMap = append(navMap, m)
	}

	ctx.ContentType("text/xml")
	bytes, e := ctx.App().View().Render("sitemap.xml", map[string]interface{}{
		"Title":      model.GetSetting("site_title"),
		"Link":       baseUrl,
		"Created":    now,
		"Articles":   articleMap,
		"Navigators": navMap,
	})
	if e != nil {
		panic(e)
	}
	ctx.Body = bytes

}

func Rss(ctx *GoInk.Context) {
	baseUrl := model.GetSetting("site_url")
	article, _ := model.GetPublishArticleList(1, 20)
	author := model.GetUsersByRole("ADMIN")[0]

	articleMap := make([]map[string]string, len(article))
	for i, a := range article {
		m := make(map[string]string)
		m["Title"] = a.Title
		m["Link"] = strings.Replace(baseUrl+a.Link(), baseUrl+"/", baseUrl, -1)
		m["Author"] = author.Nick
		str := utils.Markdown2Html(a.Content())
		str = strings.Replace(str, `src="/`, `src="`+strings.TrimSuffix(baseUrl, "/")+"/", -1)
		str = strings.Replace(str, `href="/`, `href="`+strings.TrimSuffix(baseUrl, "/")+"/", -1)
		m["Desc"] = str
		m["Created"] = time.Unix(a.CreateTime, 0).Format(time.RFC822)
		articleMap[i] = m
	}

	ctx.ContentType("application/rss+xml;charset=UTF-8")

	bytes, e := ctx.App().View().Render("rss.xml", map[string]interface{}{
		"Title":    model.GetSetting("site_title"),
		"Link":     baseUrl,
		"Desc":     model.GetSetting("site_description"),
		"Created":  time.Unix(utils.Now(), 0).Format(time.RFC822),
		"Articles": articleMap,
	})
	if e != nil {
		panic(e)
	}
	ctx.Body = bytes
}
