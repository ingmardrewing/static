package static

import (
	"fmt"

	"github.com/buger/jsonparser"
)

// Json
type Json struct{}

func (j *Json) ReadString(value []byte, keys ...string) string {
	v, err := jsonparser.GetString(value, keys...)
	if err != nil {
		return ""
	}
	return v
}

func (j *Json) ReadInt(value []byte, keys ...string) int {
	v, err := jsonparser.GetInt(value, keys...)
	if err != nil {
		return 0
	}
	return int(v)
}

type DAO interface {
	ExtractFromJson()
	FillJson() []byte
	Id() int
	Title() string
	TitlePlain() string
	ThumbUrl() string
	ImageUrl() string
	Description() string
	DisqusId() string
	CreateDate() string
	Content() string
	Url() string
	PathFromDocRoot() string
	HtmlFilename() string
}

// docDAO
type docDAO struct {
	data []byte
	id   int
	title, titlePlain, thumbUrl,
	imageUrl, description, disqusId,
	createDate, content, url,
	path, filename,
	fspath, fsfilename string
}

func (p *docDAO) Id() int {
	return p.id
}

func (p *docDAO) Title() string {
	return p.title
}

func (p *docDAO) TitlePlain() string {
	return p.titlePlain
}

func (p *docDAO) ThumbUrl() string {
	return p.thumbUrl
}

func (p *docDAO) ImageUrl() string {
	return p.imageUrl
}

func (p *docDAO) Description() string {
	return p.description
}

func (p *docDAO) DisqusId() string {
	return p.disqusId
}

func (p *docDAO) CreateDate() string {
	return p.createDate
}

func (p *docDAO) Content() string {
	return p.content
}

func (p *docDAO) Url() string {
	return p.url
}

func (p *docDAO) PathFromDocRoot() string {
	return p.path
}

func (p *docDAO) HtmlFilename() string {
	return p.filename
}

func (p *docDAO) Template() string {
	return `{
	"thumbImg":"%s",
	"postImg":"%s",
	"filename":"%s",
	"post":{
		"post_id":"%s",
		"date":"%s",
		"url":"%s",
		"title":"%s",
		"title_plain":"%s",
		"excerpt":"%s",
		"content":"%s",
		"custom_fields":{
			"dsq_thread_id":["%s"]
		}
	}
}`
}

// Post Dawo
func NewPostDAO(json []byte, path, filename string) DAO {
	p := new(postDAO)
	p.data = json
	p.path = path
	p.filename = filename
	return p
}

type postDAO struct {
	Json
	docDAO
}

func (p *postDAO) ExtractFromJson() {
	p.id = p.ReadInt(p.data, "post", "post_id")
	p.title = p.ReadString(p.data, "post", "title")
	p.thumbUrl = p.ReadString(p.data, "thumbImg")
	p.imageUrl = p.ReadString(p.data, "postImg")
	p.description = p.ReadString(p.data, "post", "excerpt")
	p.disqusId = p.ReadString(p.data, "post", "custom_fields", "dsq_thread_id", "[0]")
	p.createDate = p.ReadString(p.data, "post", "date")
	p.content = p.ReadString(p.data, "post", "content")
	p.url = p.ReadString(p.data, "post", "url")
	p.path = p.ReadString(p.data, "path")
	p.filename = p.ReadString(p.data, "filename")
}

func (p *postDAO) FillJson() []byte {
	json := fmt.Sprintf(p.Template(),
		p.thumbUrl, p.imageUrl, p.filename,
		p.id, p.createDate, p.url,
		p.title, p.titlePlain, p.description,
		p.content, p.disqusId)
	return []byte(json)
}

// marginalDAO
func NewMarginalDAO(json []byte, path, filename string) DAO {
	p := new(marginalDAO)
	p.data = json
	p.fspath = path
	p.fsfilename = filename
	return p
}

type marginalDAO struct {
	Json
	docDAO
}

func (p *marginalDAO) ExtractFromJson() {
	p.id = p.ReadInt(p.data, "page", "post_id")
	p.title = p.ReadString(p.data, "title")
	p.filename = p.ReadString(p.data, "filename")
	p.thumbUrl = p.ReadString(p.data, "thumbImg")
	p.imageUrl = p.ReadString(p.data, "postImg")
	p.description = p.ReadString(p.data, "page", "excerpt")
	p.disqusId = p.ReadString(p.data, "page", "custom_fields", "dsq_thread_id", "[0]")
	p.createDate = p.ReadString(p.data, "page", "date")
	p.content = p.ReadString(p.data, "content")
	p.path = p.ReadString(p.data, "path")
}

func (p *marginalDAO) FillJson() []byte {
	json := fmt.Sprintf(p.Template(),
		p.thumbUrl, p.imageUrl, p.filename,
		p.id, p.createDate, p.url,
		p.title, p.titlePlain, p.description,
		p.content, p.disqusId)
	return []byte(json)
}
