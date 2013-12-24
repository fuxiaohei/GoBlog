package model

import (
	"github.com/fuxiaohei/GoBlog/app"
)

type Tag struct {
	Id     int
	Name   string
	Slug   string
	Counts int
}

type TagModel struct {
	tags map[string]*Tag
	idIndex map[int]string
}

// get all tags. no ordered.
func (this *TagModel) GetAll() []*Tag {
	if this.tags != nil {
		tags := make([]*Tag, len(this.tags))
		i := 0
		for _, tag := range this.tags {
			tags[i] = tag
			i++
		}
		return tags
	}
	sql := "SELECT * FROM blog_meta WHERE type = ?"
	res, _ := app.Db.Query(sql, "tag")
	tags := make([]*Tag, 0)
	res.All(&tags)
	return tags
}

// get one tag by slug.
func (this *TagModel) GetTagBySlug(slug string) *Tag {
	tag := this.tags[slug]
	if tag != nil {
		return tag
	}
	sql := "SELECT * FROM blog_meta WHERE type = ? AND slug = ?"
	res, _ := app.Db.Query(sql, "tag", slug)
	if len(res.Data) < 1 {
		return nil
	}
	t := new(Tag)
	res.One(t)
	this.cacheTag(t)
	return t
}

// get one tag by id.
func (this *TagModel) GetTagById(id int) *Tag {
	slug := this.idIndex[id]
	if slug != "" {
		return this.GetTagBySlug(slug)
	}
	sql := "SELECT * FROM blog_meta WHERE type = ? AND id = ?"
	res, _ := app.Db.Query(sql, "tag", id)
	if len(res.Data) < 1 {
		return nil
	}
	t := new(Tag)
	res.One(t)
	this.cacheTag(t)
	return t
}

// cache tag struct to map and index.
func (this *TagModel) cacheTag(t *Tag) {
	if t == nil {
		return
	}
	this.tags[t.Slug] = t
	this.idIndex[t.Id] = t.Slug
}

// delete cached tag in map and index.
func (this *TagModel) nocacheTag(t *Tag) {
	if t == nil {
		return
	}
	delete(this.tags, t.Slug)
	delete(this.idIndex, t.Id)
}

// create new tag.
// return the new tag.
func (this *TagModel) CreateTag(name string, slug string) *Tag {
	sql := "INSERT INTO blog_meta(name,slug,type) VALUES(?,?,?,?)"
	res, _ := app.Db.Exec(sql, name, slug, "tag")
	return this.GetTagById((res.LastInsertId))
}

// save tag data.
// return the updated tag.
func (this *TagModel) SaveTag(id int, name string, slug string) *Tag {
	this.nocacheTag(this.GetTagById(id))
	// update and re-cache it
	sql := "UPDATE blog_meta SET name = ?,slug = ? WHERE id = ?"
	app.Db.Exec(sql, name, slug, id)
	return this.GetTagById(id)
}

// delete tag.
// delete relationship together.
func (this *TagModel) DeleteTag(id int) {
	sql := "DELETE FROM blog_relationship WHERE meta_id = ?"
	app.Db.Exec(sql, id)
	this.nocacheTag(this.GetTagById(id))
	sql = "DELETE FROM blog_meta WHERE id = ?"
	app.Db.Exec(sql, id)
}

func (this *TagModel) reset() {
	this.tags = make(map[string]*Tag)
	this.idIndex = make(map[int]string)
	tags := this.GetAll()
	for _, tag := range tags {
		this.tags[tag.Slug] = tag
		this.idIndex[tag.Id] = tag.Slug
	}
}

// create new tag model.
func NewTagModel() *TagModel {
	t := new(TagModel)
 	go t.reset()
	return t
}
