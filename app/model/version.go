package model

type Version struct {
	Version     string `xorm:"unique not null"`
	BuildTime   int64
	InstallTime int64
	ChangeLog   string
}
