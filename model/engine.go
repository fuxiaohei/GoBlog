package model

import (
	"github.com/fuxiaohei/GoBlog/gof"
	"github.com/go-xorm/xorm"
)

var DB *xorm.Engine

func InitDB(cfg gof.ConfigInterface) error {
	return nil
}
