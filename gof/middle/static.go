package middle

import (
	"github.com/Unknwon/com"
	"github.com/fuxiaohei/GoBlog/gof"
	"net/http"
	"path/filepath"
	"strings"
)

type StaticOption struct {
	Dir    string
	Suffix []string
	Index  string
}

func newDefaultStaticOption() *StaticOption {
	return &StaticOption{
		Dir:    "static",
		Suffix: []string{".html", ".css", ".js", ".json", ".png", ".jpg", ".zip", ".gif", ".pdf", ".txt", ".rar"},
		Index:  "index.html",
	}
}

func Static(opt *StaticOption) gof.RouterHandler {
	if opt == nil {
		opt = newDefaultStaticOption()
	}
	return func(ctx *gof.Context) {
		url := ctx.Request().URL.Path
		ext := filepath.Ext(url)

		// if no extension, but set index page, serve as directory
		if ext == "" && opt.Index != "" {
			url = strings.TrimSuffix(url, "/") + opt.Index
			ext = filepath.Ext(url)
		}

		// check extension valid
		isExtValid := false
		for _, ex := range opt.Suffix {
			if ex == ext {
				isExtValid = true
				break
			}
		}
		if !isExtValid {
			return
		}
		file := strings.TrimPrefix(filepath.Join(opt.Dir, url), "/")

		if com.IsFile(file) {
			http.ServeFile(ctx.Response(), ctx.Request(), file)
			ctx.Status = 200
			ctx.Stop()
		}

	}
}
