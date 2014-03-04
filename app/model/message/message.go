package message

import (
	"github.com/fuxiaohei/GoBlog/app/model/content"
	. "github.com/fuxiaohei/GoBlog/app/model/storage"
	"github.com/fuxiaohei/GoBlog/app/model/timer"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"strings"
)

var (
	messages         []*Message
	messageMaxId     int
	messageGenerator map[string]func(v interface{}) string
)

func init() {
	messageGenerator = make(map[string]func(v interface{}) string)
	messageGenerator["comment"] = generateCommentMessage
	messageGenerator["backup"] = generateBackupMessage
}

type Message struct {
	Id         int
	Type       string
	CreateTime int64
	Data       string
	IsRead     bool
}

func Create(tp string, data interface{}) *Message {
	m := new(Message)
	m.Type = tp
	m.Data = messageGenerator[tp](data)
	if m.Data == "" {
		println("message generator returns empty")
		return nil
	}
	m.CreateTime = utils.Now()
	m.IsRead = false
	messageMaxId += Storage.TimeInc(3)
	m.Id = messageMaxId
	messages = append([]*Message{m}, messages...)
	Sync()
	return m
}

func SetGenerator(name string, fn func(v interface{}) string) {
	messageGenerator[name] = fn
}

func ById(id int) *Message {
	for _, m := range messages {
		if m.Id == id {
			return m
		}
	}
	return nil
}

func Unread() []*Message {
	ms := make([]*Message, 0)
	for _, m := range messages {
		if m.IsRead {
			continue
		}
		ms = append(ms, m)
	}
	return ms
}

func All() []*Message {
	return messages
}

func Typed(tp string, unread bool) []*Message {
	ms := make([]*Message, 0)
	for _, m := range messages {
		if m.Type == tp {
			if unread {
				if !m.IsRead {
					ms = append(ms, m)
				}
			} else {
				ms = append(ms, m)
			}
		}
	}
	return ms
}

func SetRead(m *Message) {
	m.IsRead = true
	Sync()
}

func Sync() {
	Storage.Set("messages", messages)
}

func Load() {
	messages = make([]*Message, 0)
	messageMaxId = 0
	Storage.Get("messages", &messages)
	for _, m := range messages {
		if m.Id > messageMaxId {
			messageMaxId = m.Id
		}
	}
	startMessageTimer()
}

func Recycle() {
	for i, m := range messages {
		if m.CreateTime+3600*24*3 < utils.Now() {
			messages = messages[:i]
			return
		}
	}
}

func generateCommentMessage(co interface{}) string {
	c, ok := co.(*content.Comment)
	if !ok {
		return ""
	}
	cnt := content.ById(c.Cid)
	s := ""
	if c.Pid < 1 {
		s = "<p>" + c.Author + "同学，在文章《" + cnt.Title + "》发表评论："
		s += utils.Html2str(c.Content) + "</p>"
	} else {
		p := content.CommentById(c.Pid)
		s = "<p>" + p.Author + "同学，在文章《" + cnt.Title + "》的评论："
		s += utils.Html2str(p.Content) + "</p>"
		s += "<p>" + c.Author + "同学的回复："
		s += utils.Html2str(c.Content) + "</p>"
	}
	return s
}

func generateBackupMessage(co interface{}) string {
	str := co.(string)
	if strings.HasPrefix(str, "[0]") {
		return "备份全站失败: " + strings.TrimPrefix(str, "[0]") + "."
	}
	return "备份全站到 " + strings.TrimPrefix(str, "[1]") + " 成功."
}

func startMessageTimer() {
	timer.SetFunc("message-sync", 9, func() {
		println("write messages in 1.5 hour timer")
		Recycle()
		Sync()
	})
}
