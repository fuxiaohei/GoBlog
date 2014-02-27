package model

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
	s.Comments = len(commentsIndex)
	s.Articles = len(contentsIndex["article"])
	s.Pages = len(contentsIndex["page"])
	s.Files = len(files)
	s.Version = GetVersion().Version
	s.Readers = len(GetReaders())
	return s
}
