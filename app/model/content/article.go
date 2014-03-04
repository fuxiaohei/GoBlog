package content

import "github.com/fuxiaohei/GoBlog/app/utils"

func generatePublishArticleIndex() {
	arr := make([]int, 0)
	for _, id := range contentsIndex["article"] {
		c := ById(id)
		if c.Status == "publish" {
			arr = append(arr, id)
		}
	}
	contentsIndex["article-publish"] = arr
}

// GetPublishArticleList gets published article list and pager.
func PublishArticleList(page, size int) ([]*Content, *utils.Pager) {
	index := contentsIndex["article-publish"]
	pager := utils.NewPager(page, size, len(index))
	articles := make([]*Content, 0)
	if len(index) < 1 {
		return articles, pager
	}
	if page > pager.Pages {
		return articles, pager
	}
	for i := pager.Begin; i <= pager.End; i++ {
		articles = append(articles, ById(index[i-1]))
	}
	return articles, pager
}

// GetArticleList gets articles list and pager no matter article status.
func ArticleList(page, size int) ([]*Content, *utils.Pager) {
	index := contentsIndex["article"]
	pager := utils.NewPager(page, size, len(index))
	articles := make([]*Content, 0)
	if len(index) < 1 {
		return articles, pager
	}
	if page > pager.Pages {
		return articles, pager
	}
	for i := pager.Begin; i <= pager.End; i++ {
		articles = append(articles, ById(index[i-1]))
	}
	return articles, pager
}

// GetPopularArticleList returns popular articles list.
// Popular articles are ordered by comment number.
func PopularArticleList(size int) []*Content {
	index := contentsIndex["article-pop"]
	pager := utils.NewPager(1, size, len(index))
	articles := make([]*Content, 0)
	if len(index) < 1 {
		return articles
	}
	if 1 > pager.Pages {
		return articles
	}
	for i := pager.Begin; i <= pager.End; i++ {
		articles = append(articles, ById(index[i-1]))
	}
	return articles
}

// GetTaggedArticleList returns tagged articles list.
// These articles contains same one tag.
func TaggedArticleList(tag string, page, size int) ([]*Content, *utils.Pager) {
	index := contentsIndex["t-"+tag]
	pager := utils.NewPager(page, size, len(index))
	articles := make([]*Content, 0)
	if len(index) < 1 {
		return articles, pager
	}
	if page > pager.Pages {
		return articles, pager
	}
	for i := pager.Begin; i <= pager.End; i++ {
		articles = append(articles, ById(index[i-1]))
	}
	return articles, pager
}
