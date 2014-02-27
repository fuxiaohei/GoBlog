package model

import (
	"github.com/fuxiaohei/GoBlog/app/utils"
	"html/template"
	"sort"
	"strings"
)

var (
	readers       map[string]*Reader
	comments      map[int]*Comment
	commentsIndex []int
	commentMaxId  int
)

// Comment Reader struct.
// Saving comment reader for visiting wall usage or other statics.
type Reader struct {
	Author   string
	Email    string
	Url      string
	Active   bool
	Comments int
	Rank     int
}

// Inc increases Reader's rank.
func (r *Reader) Inc() {
	r.Rank++
	if r.Rank > 1 {
		r.Active = true
	}
}

// Dec decreases Reader's rank.
func (r *Reader) Dec() {
	r.Rank--
	if r.Rank < 1 {
		r.Active = false
	}
}

// Comment struct defines a comment item data.
type Comment struct {
	Id         int
	Author     string
	Email      string
	Url        string
	Avatar     string
	Content    string
	CreateTime int64
	// Content id
	Cid int
	// Parent Comment id
	Pid       int
	Status    string
	Ip        string
	UserAgent string
	// Is comment of admin
	IsAdmin bool
}

// ParentMd returns parent comment simple message as markdown text.
func (c *Comment) ParentMd() string {
	if c.Pid < 1 {
		return ""
	}
	co := GetCommentById(c.Pid)
	if co == nil {
		return "> 已失效"
	}
	str := "> @" + co.Author + "\n\n"
	str += "> " + co.Content + "\n"
	return str
}

// ToJson converts comment struct to public json map.
// It can hide some private fields, such as email.
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

// IsValid returns whether this comment is valid to show.
// If this comment is not approved or its parent is missing, return false.
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

// IsRemovable returns whether this comment can remove.
// If content or parent comment of this comment is removed, return true.
func (c *Comment) IsRemovable() bool {
	if GetContentById(c.Cid) == nil {
		return true
	}
	if c.Pid > 0 {
		if GetCommentById(c.Pid) == nil {
			return true
		}
	}
	return false
}

// GetReader returns the reader item of this comment.
func (c *Comment) GetReader() *Reader {
	for _, r := range readers {
		if r.Email == c.Email {
			return r
		}
	}
	return nil
}

// GetContent returns the content item of this comment.
func (c *Comment) GetContent() *Content {
	return GetContentById(c.Cid)
}

// CreateReader creates a reader from a comment.
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

// CreateComment creates a comment and links it to the cid content.
func CreateComment(cid int, c *Comment) {
	commentMaxId += Storage.TimeInc(4)
	c.Id = commentMaxId
	c.CreateTime = utils.Now()
	c.Status = "check"
	c.Cid = cid
	// escape content
	c.Content = strings.Replace(utils.Html2str(template.HTMLEscapeString(c.Content)), "\n", "<br/>", -1)
	// if empty url, use # instead.
	if c.Url == "" {
		c.Url = "#"
	}
	// if admin comment, must be approved.
	if c.IsAdmin {
		c.Status = "approved"
	} else {
		// if common comment, get reader status for checking status.
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

// SaveComment saves a comment and related updates content and reader data.
func SaveComment(c *Comment) {
	cnt := GetContentById(c.Cid)
	go SyncContent(cnt)
	go SyncReaders()
}

func removeOneComment(cid int, id int) {
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
}

// RemoveComment removes a comment by id and updates content iten by cid.
func RemoveComment(cid int, id int) {
	removeOneComment(cid, id)
	cnt := GetContentById(cid)
	if cnt == nil {
		return
	}
	go SyncContent(cnt)
}

// GetCommentById returns a comment by id.
func GetCommentById(id int) *Comment {
	return comments[id]
}

// GetCommentList returns a comments list and pager.
// This list scans all comments no matter its status.
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

// GetCommentRecentList returns a comments list of recent comments.
// Recent comments are approved and no parent and not admin comment.
// It's ordered by comment id desc.
func GetCommentRecentList(size int) []*Comment {
	comments, i := make([]*Comment, 0), 0
	for _, id := range commentsIndex {
		if i >= size {
			return comments
		}
		c := GetCommentById(id)
		if c.Pid < 1 && c.IsAdmin == false && c.Status == "approved" {
			comments = append(comments, c)
			i++
		}
	}
	return comments
}

// SyncReaders writes all readers data.
func SyncReaders() {
	Storage.Set("readers", readers)
}

// LoadReaders loads all readers from storage json.
func LoadReaders() {
	readers = make(map[string]*Reader)
	Storage.Get("readers", &readers)
}

// GetReaders returns slice of all readers
func GetReaders() []*Reader {
	r, i := make([]*Reader, len(readers)), 0
	for _, rd := range readers {
		r[i] = rd
		i++
	}
	return r
}

// RemoveReader removes a reader by his email.
func RemoveReader(email string) {
	delete(readers, email)
	SyncReaders()
}

// LoadComments loads all comments from contents.
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

// UpdateCommentAdmin updates comment author data if admin user data updated.
// It only updates admin comments.
func UpdateCommentAdmin(user *User) {
	for _, co := range comments {
		if co.IsAdmin {
			co.Author = user.Nick
			co.Email = user.Email
			co.Url = user.Url
			co.Avatar = utils.Gravatar(co.Email, "50")
		}
	}
}

// RecycleComments cleans removable comments.
func RecycleComments() {
	readerTmp := make(map[string][]int)
	for _, co := range comments {
		if co.IsRemovable() {
			removeOneComment(co.Cid, co.Id)
			continue
		}
		r := co.GetReader()
		if r == nil {
			continue
		}
		if _, ok := readerTmp[r.Email]; !ok {
			readerTmp[r.Email] = make([]int, 0)
		}
		readerTmp[r.Email] = append(readerTmp[r.Email], co.Id)
	}
	for _, r := range readers {
		if _, ok := readerTmp[r.Email]; ok {
			r.Comments = len(readerTmp[r.Email])
		}
	}
	SyncContents()
	SyncReaders()
}

func startCommentsTimer() {
	SetTimerFunc("comment-recycle", 36, func() {
		println("recycle comments in 6 hours timer")
		RecycleComments()
	})
}
