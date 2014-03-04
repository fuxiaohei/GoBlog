package handler

import (
	"github.com/fuxiaohei/GoBlog/app/model/content"
	"github.com/fuxiaohei/GoBlog/app/model/setting"
	"github.com/fuxiaohei/GoBlog/app/model/user"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"github.com/fuxiaohei/GoInk"
	"strings"
	"time"
)

// SiteMap is sitemap xml handler, pattern /sitemap.xml.
func SiteMap(ctx *GoInk.Context) {
	baseUrl := setting.Get("site_url")
	article, _ := content.PublishArticleList(1, 50)
	navigators := setting.GetNavigators()
	now := time.Unix(utils.Now(), 0).Format(time.RFC3339)

	// article links
	articleMap := make([]map[string]string, len(article))
	for i, a := range article {
		m := make(map[string]string)
		m["Link"] = strings.Replace(baseUrl+a.Link(), baseUrl+"/", baseUrl, -1)
		m["Created"] = time.Unix(a.CreateTime, 0).Format(time.RFC3339)
		articleMap[i] = m
	}

	// nav links
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
		"Title":      setting.Get("site_title"),
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

// RSS is feed generator handler, pattern /feed/.
func Rss(ctx *GoInk.Context) {
	baseUrl := setting.Get("site_url")
	article, _ := content.PublishArticleList(1, 20)
	author := user.ListByRole("ADMIN")[0]

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
		"Title":    setting.Get("site_title"),
		"Link":     baseUrl,
		"Desc":     setting.Get("site_description"),
		"Created":  time.Unix(utils.Now(), 0).Format(time.RFC822),
		"Articles": articleMap,
	})
	if e != nil {
		panic(e)
	}
	ctx.Body = bytes
}
