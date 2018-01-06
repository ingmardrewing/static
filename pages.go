package static

import (
	"strconv"

	"github.com/ingmardrewing/htmlDoc"
)

func NewLocation(url, prodDomain, title, thumbnailUrl, fsPath, fsFilename string) *Loc {
	return &Loc{url, prodDomain, title, thumbnailUrl, fsPath, fsFilename}
}

type Loc struct {
	url          string
	prodDomain   string
	title        string
	thumbnailUrl string
	fsPath       string
	fsFilename   string
}

func (l *Loc) GetDomain() string {
	return l.prodDomain
}

func (l *Loc) GetFsPath() string {
	return l.fsPath
}

func (l *Loc) GetFsFilename() string {
	return l.fsFilename
}

func (l *Loc) GetPath() string {
	return l.url
}

func (l *Loc) GetTitle() string {
	return l.title
}

func (l *Loc) GetThumbnailUrl() string {
	return l.thumbnailUrl
}

/* Page */

type Page struct {
	Loc
	doc           *htmlDoc.HtmlDoc
	id            int
	Content       string
	Description   string
	ImageUrl      string
	PublishedTime string
	DisqusId      string
}

func NewPage(
	id int,
	title, description, content,
	thumbUrl, imageUrl, prodDomain,
	path, filename, publishedTime,
	disqusId string) *Page {
	return &Page{
		Loc: Loc{
			title:        title,
			url:          path + filename,
			prodDomain:   prodDomain,
			thumbnailUrl: thumbUrl,
			fsPath:       path,
			fsFilename:   filename},
		id:            id,
		Description:   description,
		Content:       content,
		ImageUrl:      imageUrl,
		PublishedTime: publishedTime,
		DisqusId:      disqusId,
		doc:           htmlDoc.NewHtmlDoc()}
}

func (p *Page) Render() string {
	return p.doc.Render()
}

func (p *Page) GetId() int {
	return p.id
}

func (p *Page) GetDisqusId() string {
	return p.DisqusId
}

func (p *Page) GetContent() string {
	return p.Content
}

func (p *Page) GetDescription() string {
	if p.Description != "" {
		return p.Description
	}
	return " "
}

func (p *Page) GetImageUrl() string {
	return p.ImageUrl
}

func (p *Page) GetPublishedTime() string {
	return p.PublishedTime
}

func (p *Page) AcceptVisitor(v Component) {
	v.VisitPage(p)
}

func (p *Page) AddHeaderNodes(nodes []*htmlDoc.Node) {
	for _, n := range nodes {
		p.doc.AddHeadNode(n)
	}
}

func (p *Page) AddBodyNodes(nodes []*htmlDoc.Node) {
	for _, n := range nodes {
		p.doc.AddBodyNode(n)
	}
}

/* PageManager */

func NewPageManager() *PageManager {
	return new(PageManager)
}

// TODO: move all configuration to main function
// - AddMarginal, AddPost, AddPage - each with an own splice
// - AddMarginalContext, AddPostContext, AddPageContext, AddPostsNaviContext
// - put pages into their contexts, here
// - Handle file creation here
type PageManager struct {
	marginal      []Element
	posts         []Element
	postNaviPages []Element
	pages         []Element
}

// TODO Factor out Generation into another struct / func
func (p *PageManager) GeneratePostNaviPages(atPath string, posts []Element, title, description, domain, classDiv, classA, classSpan string) []Element {
	pnps := []Element{}
	elementss := generateElementBundles(posts)
	last := len(elementss) - 1
	for i, elems := range elementss {
		naviPageContent := p.generateNaviPageContent(elems, classDiv, classA, classSpan)
		filename := "index" + strconv.Itoa(i) + ".html"
		if i == last {
			filename = "index.html"
		}
		pnp := NewPage(i, title, description,
			naviPageContent, "", "", domain,
			atPath, filename, "", "")
		pnps = append(pnps, pnp)
	}

	return pnps
}

func (p *PageManager) generateNaviPageContent(elems []Element, classDiv, classA, classSpan string) string {
	n := htmlDoc.NewNode("div", "", "class", classDiv)
	for _, e := range elems {
		ta := e.GetThumbnailUrl()
		if ta == "" {
			ta = e.GetImageUrl()
		}
		path := e.GetPath()
		a := htmlDoc.NewNode("a", " ",
			"href", path,
			"class", classA)
		span := htmlDoc.NewNode("span", " ",
			"style", "background-image: url("+e.GetThumbnailUrl()+")",
			"class", classSpan)
		h2 := htmlDoc.NewNode("h2", e.GetTitle())
		a.AddChild(span)
		a.AddChild(h2)
		n.AddChild(a)
	}
	n.AddChild(htmlDoc.NewNode("div", "", "style", "clear: both"))
	return n.Render()
}

// util

func ElementsToLocations(elements []Element) []Location {
	locs := []Location{}
	for _, p := range elements {
		locs = append(locs, p)
	}
	return locs
}

func generateElementBundles(pages []Element) [][]Element {
	length := len(pages)
	reversed := []Element{}
	for i := length - 1; i >= 0; i-- {
		reversed = append(reversed, pages[i])
	}

	b := newElementBundle()
	bundles := []*elementBundle{}
	for _, p := range reversed {
		b.addElement(p)
		if b.full() {
			bundles = append(bundles, b)
			b = newElementBundle()
		}
	}
	if !b.full() {
		bundles = append(bundles, b)
	}

	length = len(bundles)
	revbundles := []*elementBundle{}
	elementss := [][]Element{}
	for i := length - 1; i >= 0; i-- {
		revbundles = append(revbundles, bundles[i])
		elementss = append(elementss, bundles[i].getElements())
	}

	return elementss
}

func newElementBundle() *elementBundle {
	return new(elementBundle)
}

type elementBundle struct {
	elements []Element
}

func (l *elementBundle) addElement(e Element) {
	l.elements = append(l.elements, e)
}

func (l *elementBundle) full() bool {
	return len(l.elements) >= 10
}

func (l *elementBundle) getElements() []Element {
	return l.elements
}
