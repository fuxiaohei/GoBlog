package comment

import (
	. "github.com/fuxiaohei/GoBlog/app/model/storage"
)

// Comment Reader struct.
// Saving comment reader for visiting wall usage or other statics.
type Reader struct {
	Author   string
	Email    string
	Url      string
	Active   bool
	Comments int
	Rank     int
}

// Inc increases Reader's rank.
func (r *Reader) Inc() {
	r.Rank++
	if r.Rank > 1 {
		r.Active = true
	}
}

// Dec decreases Reader's rank.
func (r *Reader) Dec() {
	r.Rank--
	if r.Rank < 1 {
		r.Active = false
	}
}

// CreateReader creates a reader from a comment.
func CreateReader(c *Comment) {
	r := new(Reader)
	r.Author = c.Author
	r.Email = c.Email
	r.Url = c.Url
	r.Active = false
	r.Comments = 1
	r.Rank = 0
	readers[r.Email] = r
	go SyncReaders()
}

// SyncReaders writes all readers data.
func SyncReaders() {
	Storage.Set("readers", readers)
}

// LoadReaders loads all readers from storage json.
func LoadReaders() {
	readers = make(map[string]*Reader)
	Storage.Get("readers", &readers)
}

// GetReaders returns slice of all readers
func GetReaders() []*Reader {
	r, i := make([]*Reader, len(readers)), 0
	for _, rd := range readers {
		r[i] = rd
		i++
	}
	return r
}

// RemoveReader removes a reader by his email.
func RemoveReader(email string) {
	delete(readers, email)
	SyncReaders()
}
