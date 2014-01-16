package model

import (
	"errors"
	"fmt"
	"git.oschina.net/fuxiaohei/GoBlog.git/app/utils"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	contents      map[int]*Content
	contentMaxId  int
	contentsIndex map[string][]int
)

type Content struct {
	Id    int
	Title string
	Slug  string
	Text  string
	//Category   string
	Tags []string

	CreateTime int64
	EditTime   int64
	UpdateTime int64

	IsComment bool
	IsLinked  bool
	//IsTop     bool

	AuthorId int

	Template string

	Type   string
	Status string
	Format string

	Comments []*Comment
	Hits     int
}

// get content tags in a string.
func (cnt *Content) TagString() string {
	return strings.Join(cnt.Tags, ",")
}

// get content link.
func (cnt *Content) Link() string {
	if cnt.IsLinked {
		return "/"+cnt.Slug + ".html"
	}
	return fmt.Sprintf("/%s/%d/%s.html", cnt.Type, cnt.Id, cnt.Slug)
}

// get content text.
func (cnt *Content) Content() string {
	return cnt.Text
}

// get content summary.
func (cnt *Content) Summary() string {
	return strings.Split(cnt.Text, "<!--more-->")[0]
}

// get content comments number.
func (cnt *Content) CommentNum() int {
	max := len(cnt.Comments)
	for _, c := range cnt.Comments {
		if !c.IsValid() {
			max--
		}
	}
	return max
}

// change slug and check unique.
func (cnt *Content) ChangeSlug(slug string) bool {
	c2 := GetContentBySlug(slug)
	if c2 == nil {
		cnt.Slug = slug
		return true
	}
	if c2.Id != cnt.Id {
		return false
	}
	cnt.Slug = slug
	return true
}

// get content author user.
func (cnt *Content) User() *User {
	return GetUserById(cnt.AuthorId)
}

// get a content by given id.
func GetContentById(id int) *Content {
	return contents[id]
}

// get a content by given slug.
func GetContentBySlug(slug string) *Content {
	for _, c := range contents {
		if c.Slug == slug {
			return c
		}
	}
	return nil
}

// get articles list.
func GetArticleList(page, size int) ([]*Content, *utils.Pager) {
	index := contentsIndex["article"]
	pager := utils.NewPager(page, size, len(index))
	articles := make([]*Content, 0)
	if page > pager.Pages {
		return articles, pager
	}
	for i := pager.Begin; i <= pager.End; i++ {
		articles = append(articles, GetContentById(index[i-1]))
	}
	return articles, pager
}

// get pages list.
func GetPageList(page, size int) ([]*Content, *utils.Pager) {
	index := contentsIndex["page"]
	pager := utils.NewPager(page, size, len(index))
	pages := make([]*Content, 0)
	if page > pager.Pages {
		return pages, pager
	}
	for i := pager.Begin; i <= pager.End; i++ {
		pages = append(pages, GetContentById(index[i-1]))
	}
	return pages, pager
}

// create new content.
func CreateContent(c *Content, t string) (*Content, error) {
	c2 := GetContentBySlug(c.Slug)
	if c2 != nil {
		return nil, errors.New("slug-repeat")
	}
	contentMaxId += Storage.TimeInc(3)
	c.Id = contentMaxId
	c.CreateTime = utils.Now()
	c.EditTime = c.CreateTime
	c.UpdateTime = c.CreateTime
	c.Comments = make([]*Comment, 0)
	c.Type = t
	c.Hits = 1
	contents[c.Id] = c
	contentsIndex[c.Type] = append([]int{c.Id}, contentsIndex[c.Type]...)
	go SyncContent(c)
	return c, nil
}

// save changed content.
func SaveContent(c *Content) {
	c.EditTime = utils.Now()
	go SyncContent(c)
}

// remove a content.
// not delete file, just change status to DELETE.
// it can't be loaded in memory from json.
func RemoveContent(c *Content) {
	delete(contents, c.Id)
	for i, id := range contentsIndex[c.Type] {
		if c.Id == id {
			contentsIndex[c.Type] = append(contentsIndex[c.Type][:i], contentsIndex[c.Type][i+1:]...)
			break
		}
	}
	c.Status = "DELETE"
	go SyncContent(c)
}

// write a content to json.
func SyncContent(c *Content) {
	key := fmt.Sprintf("content/%s-%d", c.Type, c.Id)
	Storage.Set(key, c)
}

// write all contents to json.
func SyncContents() {
	for _, c := range contents {
		SyncContent(c)
	}
}

// load all contents.
// generate indexes.
func LoadContents() {
	contents = make(map[int]*Content)
	contentsIndex = make(map[string][]int)
	contentMaxId = 0
	articleIndex := make([]int, 0)
	pageIndex := make([]int, 0)
	filepath.Walk(filepath.Join(Storage.dir, "content"), func(_ string, info os.FileInfo, err error) error {
			if err == nil {
				if info.IsDir() {
					return nil
				}
				c := new(Content)
				file := strings.Replace("content/" + info.Name(), filepath.Ext(info.Name()), "", -1)
				Storage.Get(file, c)
				if c.Id > 0 {
					if c.Status != "DELETE" {
						contents[c.Id] = c
						if c.Type == "article" {
							articleIndex = append(articleIndex, c.Id)
						}
						if c.Type == "page" {
							pageIndex = append(pageIndex, c.Id)
						}
					}
					if c.Id > contentMaxId {
						contentMaxId = c.Id
					}
				}
			}
			return nil
		})
	sort.Sort(sort.Reverse(sort.IntSlice(articleIndex)))
	sort.Sort(sort.Reverse(sort.IntSlice(pageIndex)))
	contentsIndex["article"] = articleIndex
	contentsIndex["page"] = pageIndex
}

func StartContentsTimer() {
	time.AfterFunc(time.Duration(10) * time.Minute, func() {
			println("write contents in timer")
			SyncContents()
			StartContentsTimer()
		})
}
