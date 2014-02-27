package model

import (
	"errors"
	"fmt"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	contents      map[int]*Content
	contentMaxId  int
	contentsIndex map[string][]int
	tags          []*Tag
)

// Content instance, defines content data items.
type Content struct {
	Id    int
	Title string
	Slug  string
	Text  string
	// Rendered string
	textRendered string
	//Category   string
	Tags []string

	CreateTime int64
	EditTime   int64
	UpdateTime int64

	// IsComment opens or closes comment
	IsComment bool
	// IsLinked makes pager link as top level link /link.html
	IsLinked bool
	//IsTop     bool

	AuthorId int

	// Template makes pager use own template file
	Template string

	Type   string
	Status string

	// Format defines the content text format type. Now only support markdown.
	Format string

	Comments []*Comment
	Hits     int
}

// TagString returns content tags in a string that's joined by ",".
func (cnt *Content) TagString() string {
	return strings.Join(cnt.Tags, ",")
}

// GetTags returns tags struct in this content.
func (cnt *Content) GetTags() []*Tag {
	tgs := make([]*Tag, len(cnt.Tags))
	for i, t := range cnt.Tags {
		tgs[i] = &Tag{Name: t}
	}
	return tgs
}

// Link returns content link as {type}/{id}/{slug}.html.
// If content isn't published, return "#".
// If content is page and top linked, return {slug}.html as top level link.
func (cnt *Content) Link() string {
	if cnt.Status != "publish" {
		return "#"
	}
	if cnt.IsLinked {
		return "/" + cnt.Slug + ".html"
	}
	return fmt.Sprintf("/%s/%d/%s.html", cnt.Type, cnt.Id, cnt.Slug)
}

// Content returns whole content text.
// If enable go-markdown, return markdown-rendered content.
func (cnt *Content) Content() string {
	txt := strings.Replace(cnt.Text, "<!--more-->", "", -1)
	if GetSetting("enable_go_markdown") == "true" {
		if cnt.textRendered == "" {
			cnt.textRendered = utils.Markdown2Html(txt)
		}
		return cnt.textRendered
	}
	return txt
}

// Summary returns content summary.
// Summary text means the part before page-break <!--more-->.
// It can be go-markdown rendered.
func (cnt *Content) Summary() string {
	text := strings.Split(cnt.Text, "<!--more-->")[0]
	if GetSetting("enable_go_markdown") == "true" {
		return utils.Markdown2Html(text)
	}
	return text
}

// CommentNum returns content comments number.
// If comment are checking or, its parent are lost, do not count it.
func (cnt *Content) CommentNum() int {
	max := len(cnt.Comments)
	for _, c := range cnt.Comments {
		if !c.IsValid() {
			max--
		}
	}
	return max
}

// ChangeSlug changes content's slug.
// It checks whether this slug is unique.
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

// User returns content author user instance.
func (cnt *Content) User() *User {
	return GetUserById(cnt.AuthorId)
}

// Content Tag struct. It convert tag string to proper struct or link.
type Tag struct {
	Name string
	Cid  []int
}

// Link returns tag name url-encoded link.
func (t *Tag) Link() string {
	return "/tag/" + url.QueryEscape(strings.Replace(t.Name, ".", "-", -1))
}

// GetContentById gets a content by given id.
func GetContentById(id int) *Content {
	return contents[id]
}

// GetContentBySlug gets a content by given slug.
func GetContentBySlug(slug string) *Content {
	for _, c := range contents {
		if c.Slug == slug {
			return c
		}
	}
	return nil
}

// CreateContent creates new content.
// t means content type, article or page.
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
	c.textRendered = ""
	contents[c.Id] = c
	contentsIndex[c.Type] = append([]int{c.Id}, contentsIndex[c.Type]...)
	generatePublishArticleIndex()
	go SyncContent(c)
	return c, nil
}

// SaveContent saves changed content.
// It will re-generate related indexes.
func SaveContent(c *Content) {
	c.EditTime = utils.Now()
	// clean rendered cache text
	c.textRendered = ""
	generatePublishArticleIndex()
	go SyncContent(c)
}

// RemoveContent removes a content.
// Not delete file really, just change status to DELETE.
// This content can't be loaded in memory from storage json.
func RemoveContent(c *Content) {
	delete(contents, c.Id)
	for i, id := range contentsIndex[c.Type] {
		if c.Id == id {
			contentsIndex[c.Type] = append(contentsIndex[c.Type][:i], contentsIndex[c.Type][i+1:]...)
			break
		}
	}
	c.Status = "DELETE"
	generatePublishArticleIndex()
	go SyncContent(c)
}

// SyncContent writes a content to storage json.
func SyncContent(c *Content) {
	key := fmt.Sprintf("content/%s-%d", c.Type, c.Id)
	Storage.Set(key, c)
}

// SyncContents writes all contents to storage json.
func SyncContents() {
	for _, c := range contents {
		SyncContent(c)
	}
}

// LoadContents loads all contents, then generates indexes.
func LoadContents() {
	contents = make(map[int]*Content)
	contentsIndex = make(map[string][]int)
	contentMaxId = 0
	articleIndex := make([]int, 0)
	pageIndex := make([]int, 0)
	// walk files in directory
	filepath.Walk(filepath.Join(Storage.dir, "content"), func(_ string, info os.FileInfo, err error) error {
		if err == nil {
			// ignore dir and sub-dir
			if info.IsDir() {
				return nil
			}
			c := new(Content)
			file := strings.Replace("content/"+info.Name(), filepath.Ext(info.Name()), "", -1)
			Storage.Get(file, c)
			if c.Id > 0 {
				// ignore DELETE status
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
	// sort indexes as desc order
	sort.Sort(sort.Reverse(sort.IntSlice(articleIndex)))
	sort.Sort(sort.Reverse(sort.IntSlice(pageIndex)))
	contentsIndex["article"] = articleIndex
	contentsIndex["page"] = pageIndex
}

func startContentSyncTimer() {
	SetTimerFunc("content-sync", 6, func() {
		println("write contents in 1 hour timer")
		SyncContents()
	})
}

func generateContentTmpIndexes() {
	var (
		popTmp     = make([][2]int, 0)
		popIndex   = make([]int, 0)
		tagIndexes = make(map[string][]int)
		data       = make(map[string][]int)
	)
	for _, c := range contents {
		if c.Status == "publish" && c.Type == "article" {
			// pop temp index
			popTmp = append(popTmp, [2]int{c.Id, c.CommentNum()})
			if len(c.Tags) > 0 {
				for _, t := range c.Tags {
					if tagIndexes[t] == nil {
						tagIndexes[t] = make([]int, 0)
					}
					tagIndexes[t] = append(tagIndexes[t], c.Id)
				}
			}
		}
	}

	// sort popular list
	utils.SortInt(popTmp)
	for _, p := range popTmp {
		popIndex = append(popIndex, p[0])
	}
	contentsIndex["article-pop"] = popIndex

	// assemble indexes map
	data["pop-index"] = popIndex
	tags = make([]*Tag, 0)
	for tag, index := range tagIndexes {
		if strings.TrimSpace(tag) == "" {
			continue
		}
		sort.Sort(sort.Reverse(sort.IntSlice(index)))
		data["t-"+tag] = index
		contentsIndex["t-"+tag] = index
		t := new(Tag)
		t.Name = tag
		t.Cid = index
		tags = append(tags, t)
	}

	// write to tmp data
	TmpStorage.Set("contents", data)
}

// GetContentTags returns all tags.
func GetContentTags() []*Tag {
	return tags
}

func startContentTmpIndexesTimer() {
	SetTimerFunc("content-indexes", 36, func() {
		println("write content indexes in 6 hours timer")
		generateContentTmpIndexes()
	})
}
