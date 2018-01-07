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
	Id(...int) int
	Title(...string) string
	TitlePlain(...string) string
	ThumbUrl(...string) string
	ImageUrl(...string) string
	Description(...string) string
	DisqusId(...string) string
	CreateDate(...string) string
	Content(...string) string
	Url(...string) string
	PathFromDocRoot(...string) string
	HtmlFilename(...string) string
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

func (p *docDAO) Id(id ...int) int {
	if len(id) > 0 {
		p.id = id[0]
	}
	return p.id
}

func (p *docDAO) Title(title ...string) string {
	if len(title) > 0 {
		p.title = title[0]
	}
	return p.title
}

func (p *docDAO) TitlePlain(titlePlain ...string) string {
	if len(titlePlain) > 0 {
		p.titlePlain = titlePlain[0]
	}
	return p.titlePlain
}

func (p *docDAO) ThumbUrl(thumbUrl ...string) string {
	if len(thumbUrl) > 0 {
		p.thumbUrl = thumbUrl[0]
	}
	return p.thumbUrl
}

func (p *docDAO) ImageUrl(imageUrl ...string) string {
	if len(imageUrl) > 0 {
		p.imageUrl = imageUrl[0]
	}
	return p.imageUrl
}

func (p *docDAO) Description(desc ...string) string {
	if len(desc) > 0 {
		p.description = desc[0]
	}
	return p.description
}

func (p *docDAO) DisqusId(disqusId ...string) string {
	if len(disqusId) > 0 {
		p.disqusId = disqusId[0]
	}
	return p.disqusId
}

func (p *docDAO) CreateDate(date ...string) string {
	if len(date) > 0 {
		p.createDate = date[0]
	}
	return p.createDate
}

func (p *docDAO) Content(content ...string) string {
	if len(content) > 0 {
		p.content = content[0]
	}
	return p.content
}

func (p *docDAO) Url(url ...string) string {
	if len(url) > 0 {
		p.url = url[0]
	}
	return p.url
}

func (p *docDAO) PathFromDocRoot(path ...string) string {
	if len(path) > 0 {
		p.path = path[0]
	}
	return p.path
}

func (p *docDAO) HtmlFilename(htmlFilename ...string) string {
	if len(htmlFilename) > 0 {
		p.filename = htmlFilename[0]
	}
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
