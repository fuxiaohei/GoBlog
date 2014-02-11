package handler

import (
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoInk"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

func AdminFiles(context *GoInk.Context) {
	if context.Method == "DELETE" {
		id := context.Int("id")
		model.RemoveFile(id)
		Json(context, true).End()
		context.Do("attach_delete", id)
		return
	}
	files, pager := model.GetFileList(context.Int("page"), 10)
	context.Layout("admin/admin")
	context.Render("admin/files", map[string]interface{}{
		"Title": "媒体文件",
		"Files": files,
		"Pager": pager,
	})
}

func FileUpload(context *GoInk.Context) {
	var req *http.Request
	req = context.Request
	req.ParseMultipartForm(32 << 20)
	f, h, e := req.FormFile("file")
	if e != nil {
		Json(context, false).Set("msg", e.Error()).End()
		return
	}
	data, _ := ioutil.ReadAll(f)
	maxSize := context.App().Config().Int("app.upload_size")
	defer func() {
		f.Close()
		data = nil
		h = nil
	}()
	if len(data) >= maxSize {
		Json(context, false).Set("msg", "文件应小于10M").End()
		return
	}
	if !strings.Contains(context.App().Config().String("app.upload_files"), path.Ext(h.Filename)) {
		Json(context, false).Set("msg", "文件只支持Office文件，图片和zip存档").End()
		return
	}
	ff := new(model.File)
	ff.Name = h.Filename
	ff.Type = context.StringOr("type", "image")
	ff.Size = int64(len(data))
	ff.ContentType = h.Header["Content-Type"][0]
	ff.Author, _ = strconv.Atoi(context.Cookie("token-user"))
	ff.Url = model.CreateFilePath(context.App().Get("upload_dir"), ff)
	e = ioutil.WriteFile(ff.Url, data, os.ModePerm)
	if e != nil {
		Json(context, false).Set("msg", e.Error()).End()
		return
	}
	model.CreateFile(ff)
	Json(context, true).Set("file", ff).End()
	context.Do("attach_created", ff)
}
