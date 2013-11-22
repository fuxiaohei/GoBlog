package model

import "fmt"

type Counter struct {
	All     int
	Current int
	Pages   int
	Size    int
	Begin   int
	End     int
}

func NewCounter(all, page, size int) *Counter {
	c := &Counter{}
	c.All = all
	c.Size = size
	c.Current = page
	if all%size > 0 {
		c.Pages = all/size + 1
	}else {
		c.Pages = all/size
	}
	c.Begin = (page - 1)*size
	if c.Begin < 1 {
		c.Begin = 1
	}
	if c.Begin > all {
		c.Begin = all
	}
	c.End = page*size
	if c.End > all {
		c.End = all
	}
	return c
}

func (this *Counter) Html(pageUrl string) string {
	s := ""
	for i := 1; i <= this.Pages; i++ {
		if i == this.Current {
			s += `<a href="` + pageUrl + `?page=` + fmt.Sprint(i) + `" class="i in">` + fmt.Sprint(i) + `</a>`
			continue
		}
		s += `<a href="` + pageUrl + `?page=` + fmt.Sprint(i) + `" class="i">` + fmt.Sprint(i) + `</a>`
	}
	return s
}

