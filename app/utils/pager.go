package utils

type Pager struct {
	Current   int
	Size      int
	Total     int
	Pages     int
	PageSlice []int
	Begin     int
	End       int
	Prev      int
	Next      int
	IsPrev    bool
	IsNext    bool
}

func NewPager(page, size, total int) *Pager {
	if page < 1 {
		page = 1
	}
	p := new(Pager)
	p.Current = page
	p.Size = size
	p.Total = total
	p.Pages = total / size
	if total%size > 0 {
		p.Pages += 1
	}
	p.PageSlice = make([]int, p.Pages)
	for i := 1; i <= p.Pages; i++ {
		p.PageSlice[i-1] = i
	}
	p.Begin = (page-1)*size + 1
	if p.Begin < 1 {
		p.Begin = 1
	}
	if p.Begin > p.Total {
		p.Begin = p.Total
	}
	p.End = page * size
	if p.End > p.Total {
		p.End = p.Total
	}
	p.Prev = p.Current - 1
	p.IsPrev = true
	if p.Prev < 1 {
		p.Prev = 1
		p.IsPrev = false
	}
	p.Next = p.Current + 1
	p.IsNext = true
	if p.Next > p.Pages {
		p.Next = p.Pages
		p.IsNext = false
	}
	return p
}
