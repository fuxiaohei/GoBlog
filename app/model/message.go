package model

import (
	"github.com/fuxiaohei/GoBlog/app/utils"
	"time"
)

var (
	messages         []*Message
	messageMaxId     int
	messageGenerator map[string]func(v interface{}) map[string]interface{}
)

func init() {
	messageGenerator = make(map[string]func(v interface{}) map[string]interface{})
	messageGenerator["comment"] = generateCommentMessage
}

type Message struct {
	Id         int
	Type       string
	CreateTime int64
	Data       map[string]interface{}
	IsRead     bool
}

func CreateMessage(tp string, data interface{}) *Message {
	m := new(Message)
	m.Type = tp
	m.Data = messageGenerator[tp](data)
	m.CreateTime = utils.Now()
	m.IsRead = false
	messageMaxId += Storage.TimeInc(3)
	messages = append([]*Message{m}, messages...)
	SyncMessages()
	return m
}

func GetMessage(id int) *Message {
	for _, m := range messages {
		if m.Id == id {
			return m
		}
	}
	return nil
}

func SaveMessageRead(m *Message) {
	m.IsRead = true
	SyncMessages()
}

func SyncMessages() {
	Storage.Set("messages", messages)
}

func LoadMessages() {
	messages = make([]*Message, 0)
	messageMaxId = 0
	Storage.Get("messages", &messages)
	for _, m := range messages {
		if m.Id > messageMaxId {
			messageMaxId = m.Id
		}
	}
}

func RecycleMessages() {
	for i, m := range messages {
		if m.CreateTime+3600*24*3 < utils.Now() {
			messages = messages[:i]
			return
		}
	}
}

func generateCommentMessage(co interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	return m
}

func StartMessageTimer() {
	time.AfterFunc(time.Duration(1)*time.Hour, func() {
		println("write messages in 1 hours timer")
		RecycleMessages()
		SyncMessages()
	})
}
