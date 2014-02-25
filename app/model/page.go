package model

import "github.com/fuxiaohei/GoBlog/app/utils"

// GetPageList gets pages list and pager no matter page status.
// In common cases, no need to get a list or pagers for public page.
func GetPageList(page, size int) ([]*Content, *utils.Pager) {
	index := contentsIndex["page"]
	pager := utils.NewPager(page, size, len(index))
	pages := make([]*Content, 0)
	if len(index) < 1 {
		return pages, pager
	}
	if page > pager.Pages {
		return pages, pager
	}
	for i := pager.Begin; i <= pager.End; i++ {
		pages = append(pages, GetContentById(index[i-1]))
	}
	return pages, pager
}
