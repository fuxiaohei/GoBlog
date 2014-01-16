package model

import (
	"git.oschina.net/fuxiaohei/GoBlog.git/app/utils"
	"sort"
)

var (
	readers       map[string]*Reader
	comments      map[int]*Comment
	commentsIndex []int
	commentMaxId  int
)

type Reader struct {
	Author   string
	Email    string
	Url      string
	Active   bool
	Comments int
	Rank     int
}

func (r *Reader) Inc() {
	r.Rank++
	if r.Rank > 1 {
		r.Active = true
	}
}

func (r *Reader) Dec() {
	r.Rank--
	if r.Rank < 1 {
		r.Active = false
	}
}

type Comment struct {
	Id         int
	Author     string
	Email      string
	Url        string
	Avatar     string
	Content    string
	CreateTime int64
	Cid        int
	Pid        int
	Status     string
	Ip         string
	UserAgent  string
	IsAdmin    bool
}

func (c *Comment) ParentMd() string {
	if c.Pid < 1 {
		return ""
	}
	co := GetCommentById(c.Pid)
	if co == nil {
		return "> 已失效"
	}
	str := "> @"+co.Author + "\n\n"
	str += "> "+co.Content + "\n"
	return str
}

func (c *Comment) ToJson() map[string]interface{} {
	m := make(map[string]interface{})
	m["id"] = c.Id
	m["author"] = c.Author
	//m["email"] = c.Email
	m["url"] = c.Url
	m["avatar"] = c.Avatar
	m["content"] = c.Content
	m["create_time"] = c.CreateTime
	m["pid"] = c.Pid
	m["status"] = c.Status
	m["ip"] = c.Ip
	m["user_agent"] = c.UserAgent
	m["parent_md"] = c.ParentMd()
	return m
}

func (c *Comment) IsValid() bool {
	if c.Status != "approved" {
		return false
	}
	if c.Pid > 0 {
		if GetCommentById(c.Pid) == nil {
			return false
		}
	}
	return true
}

func (c *Comment) GetReader() *Reader {
	for _, r := range readers {
		if r.Email == c.Email {
			return r
		}
	}
	return nil
}

func (c *Comment) GetContent() *Content {
	return GetContentById(c.Cid)
}

func CreateReader(c *Comment) {
	r := new(Reader)
	r.Author = c.Author
	r.Email = c.Email
	r.Url = c.Url
	r.Active = false
	r.Comments = 1
	r.Rank = 0
	readers[r.Email] = r
	go SyncReaders()
}

func CreateComment(cid int, c *Comment) {
	commentMaxId += Storage.TimeInc(4)
	c.Id = commentMaxId
	c.CreateTime = utils.Now()
	c.Status = "check"
	c.Cid = cid
	if c.Url == "" {
		c.Url = "#"
	}
	if c.IsAdmin {
		c.Status = "approved"
	} else {
		r := c.GetReader()
		if r != nil {
			if r.Active {
				c.Status = "approved"
			}
		} else {
			CreateReader(c)
		}
	}
	// update comment memory data
	comments[c.Id] = c
	commentsIndex = append([]int{c.Id}, commentsIndex...)
	// append to content
	content := GetContentById(cid)
	content.Comments = append(content.Comments, c)
	go SyncContent(content)
}

func SaveComment(c *Comment) {
	cnt := GetContentById(c.Cid)
	go SyncContent(cnt)
}

func RemoveComment(cid int, id int) {
	delete(comments, id)
	for n, c := range commentsIndex {
		if c == id {
			commentsIndex = append(commentsIndex[:n], commentsIndex[n+1:]...)
			break
		}
	}
	cnt := GetContentById(cid)
	if cnt == nil {
		return
	}
	for n, c := range cnt.Comments {
		if c.Id == id {
			cnt.Comments = append(cnt.Comments[:n], cnt.Comments[n+1:]...)
		}
	}
	go SyncContent(cnt)
}

func GetCommentById(id int) *Comment {
	return comments[id]
}

func GetCommentList(page, size int) ([]*Comment, *utils.Pager) {
	index := commentsIndex
	pager := utils.NewPager(page, size, len(index))
	comments := make([]*Comment, 0)
	if page > pager.Pages {
		return comments, pager
	}
	for i := pager.Begin; i <= pager.End; i++ {
		comments = append(comments, GetCommentById(index[i-1]))
	}
	return comments, pager
}

func SyncReaders() {
	Storage.Set("readers", readers)
}

func LoadReaders() {
	readers = make(map[string]*Reader)
	Storage.Get("readers", &readers)
}

func LoadComments() {
	comments = make(map[int]*Comment)
	commentsIndex = make([]int, 0)
	commentMaxId = 0
	for _, c := range contents {
		if len(c.Comments) < 1 {
			continue
		}
		for _, cm := range c.Comments {
			comments[cm.Id] = cm
			commentsIndex = append(commentsIndex, cm.Id)
			if cm.Id > commentMaxId {
				commentMaxId = cm.Id
			}
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(commentsIndex)))
}
