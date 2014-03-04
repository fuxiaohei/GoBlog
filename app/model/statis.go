package model

import (
	"github.com/fuxiaohei/GoBlog/app/model/content"
	"github.com/fuxiaohei/GoBlog/app/model/file"
	"github.com/fuxiaohei/GoBlog/app/model/setting"
)

type Statis struct {
	Comments int
	Articles int
	Pages    int
	Files    int
	Version  int
	Readers  int
}

func NewStatis() *Statis {
	s := new(Statis)
	s.Comments = content.Len("comment")
	s.Articles = content.Len("article")
	s.Pages = content.Len("page")
	s.Files = file.Len()
	s.Version = setting.GetVersion().Version
	s.Readers = len(content.GetReaders())
	return s
}
