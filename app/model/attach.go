package model

type Attach struct {
	Id          int64  `xorm:"pk autoincr"`
	UserId      int64  `xorm:"not null"`
	Name        string `xorm:"not null"`
	UploadTime  int64
	Description string
	LocalUrl    string `xorm:"not null"`
	OtherUrls   string
	ContentType string `xorm:"not null"`
	FileType    string `xorm:"not null"`
	Size        int64
	Downloads   int64 `xorm:"default 0"`
}
