package handler

import (
	"github.com/fuxiaohei/GoBlog/app/model/file"
	"github.com/fuxiaohei/GoInk"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

// AdminFiles is uploaded file list and operation page, pattern /admin/files/.
func AdminFiles(context *GoInk.Context) {
	// delete file
	if context.Method == "DELETE" {
		id := context.Int("id")
		file.Remove(id)
		Json(context, true).End()
		context.Do("attach_delete", id)
		return
	}
	files, pager := file.List(context.Int("page"), 10)
	context.Layout("admin/admin")
	context.Render("admin/files", map[string]interface{}{
		"Title": "媒体文件",
		"Files": files,
		"Pager": pager,
	})
}

// FileUpload is file upload post handler, pattern /admin/files/upload/.
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
	ff := new(file.File)
	ff.Name = h.Filename
	ff.Type = context.StringOr("type", "image")
	ff.Size = int64(len(data))
	ff.ContentType = h.Header["Content-Type"][0]
	ff.Author, _ = strconv.Atoi(context.Cookie("token-user"))
	ff.Url = file.CreatePath(context.App().Get("upload_dir"), ff)
	e = ioutil.WriteFile(ff.Url, data, os.ModePerm)
	if e != nil {
		Json(context, false).Set("msg", e.Error()).End()
		return
	}
	file.Create(ff)
	Json(context, true).Set("file", ff).Set("link",ff.Link()).End()
	context.Do("attach_created", ff)
}

func Upload(ctx *GoInk.Context) {
	idStr, name := ctx.Param("id"), ctx.Param("name")
	id, _ := strconv.Atoi(idStr)
	f := file.ById(id)
	if f == nil {
		ctx.Status = 404
		return
	}
	if !strings.Contains(f.Url, name) {
		ctx.Status = 404
		return
	}
	http.ServeFile(ctx.Response, ctx.Request, f.Url)
	f.Hits++
	ctx.IsSend = true
}
