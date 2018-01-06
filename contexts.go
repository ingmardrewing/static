package static

import (
	"fmt"

	"github.com/ingmardrewing/fs"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
)

/* Global Context */

type ContextImpl struct {
	twitterHandle             string
	contentSection            string
	tags                      string
	siteName                  string
	twitterCardType           string
	ogType                    string
	fbPageUrl                 string
	twitterPageUrl            string
	cssUrl                    string
	rssUrl                    string
	homeUrl                   string
	disqusShortname           string
	mainNavigationLocations   []Location
	footerNavigationLocations []Location
	pages                     []Element
	components                []Component
}

func (c *ContextImpl) SetGlobalFields(
	twitterHandle,
	topic,
	tags,
	site,
	cardType,
	section,
	fbPage,
	twitterPage,
	cssUrl,
	rssUrl,
	home,
	disqusShortname string) {
	c.twitterHandle = twitterHandle
	c.contentSection = topic
	c.tags = tags
	c.siteName = site
	c.twitterCardType = cardType
	c.ogType = section
	c.fbPageUrl = fbPage
	c.twitterPageUrl = twitterPage
	c.cssUrl = cssUrl
	c.rssUrl = rssUrl
	c.homeUrl = home
	c.disqusShortname = disqusShortname
}

func (c *ContextImpl) RenderPages(targetDir string) []fs.FileContainer {
	fcs := []fs.FileContainer{}
	for _, p := range c.pages {
		for _, comp := range c.components {
			p.AcceptVisitor(comp)
		}
		fc := fs.NewFileContainer()
		fc.SetPath(targetDir + p.GetFsPath())
		fc.SetFilename(p.GetFsFilename())
		fc.SetDataAsString(p.Render() + c.GetJs())
		fcs = append(fcs, fc)
	}
	return fcs
}

func (c *ContextImpl) SetElements(pages []Element) {
	c.pages = pages
}

func (c *ContextImpl) GetComponents() []Component {
	return c.components
}

func (c *ContextImpl) GetElements() []Element {
	return c.pages
}

func (c *ContextImpl) AddPage(p Element) {
	c.pages = append(c.pages, p)
}

func (c *ContextImpl) AddComponent(comp Component) {
	c.components = append(c.components, comp)
}

func (c *ContextImpl) GetHomeUrl() string {
	return c.homeUrl
}

func (c *ContextImpl) GetRssUrl() string {
	return c.rssUrl
}

func (c *ContextImpl) GetDisqusShortname() string {
	return c.disqusShortname
}

func (c *ContextImpl) GetMainNavigationLocations() []Location {
	return c.mainNavigationLocations
}

func (c *ContextImpl) GetFooterNavigationLocations() []Location {
	return c.footerNavigationLocations
}

func (c *ContextImpl) GetCssUrl() string {
	return c.cssUrl
}

func (c *ContextImpl) GetTwitterPage() string {
	return c.twitterPageUrl
}

func (c *ContextImpl) GetFBPageUrl() string {
	return c.fbPageUrl
}

func (c *ContextImpl) GetOGType() string {
	return c.ogType
}

func (c *ContextImpl) GetTwitterCardType() string {
	return c.twitterCardType
}

func (c *ContextImpl) GetTwitterHandle() string {
	return c.twitterHandle
}

func (c *ContextImpl) GetContentSection() string {
	return c.contentSection
}

func (c *ContextImpl) GetContentTags() string {
	return c.tags
}

func (c *ContextImpl) GetSiteName() string {
	return c.siteName
}

func (c *ContextImpl) GetCss() string {
	css := ""
	for _, c := range c.components {
		css += c.GetCss()
	}
	return c.minifyCss(css)
}

func (c *ContextImpl) GetJs() string {
	jsCode := ""
	for _, c := range c.components {
		jsCode += c.GetJs()
	}

	m := minify.New()
	m.AddFunc("text/javascript", js.Minify)
	s, err := m.String("text/javascript", jsCode)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(`<script>%s</script>`, s)
}

func (c *ContextImpl) minifyCss(txt string) string {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	s, err := m.String("text/css", txt)
	if err != nil {
		panic(err)
	}
	return s
}

func (c *ContextImpl) GetReadNavigationLocations() []Location {
	return nil
}

func fillContextWithComponents(context Context, components ...Component) {
	for _, compo := range components {
		compo.SetContext(context)
		context.AddComponent(compo)
	}
}

func newContext(mainnavi, footernavi []Location, contentComponents []Component) Context {
	c := new(ContextImpl)
	c.mainNavigationLocations = mainnavi
	c.footerNavigationLocations = footernavi

	fillContextWithComponents(c,
		NewGeneralMetaComponent(),
		NewFaviconComponent(),
		NewGlobalCssComponent(),
		NewGoogleComponent(),
		NewTwitterComponent(),
		NewFBComponent(),
		NewCssLinkComponent(),
		NewTitleComponent())

	fillContextWithComponents(c, contentComponents...)

	fillContextWithComponents(c,
		NewMainHeaderComponent(),
		NewMainNaviComponent(),
		NewCopyRightComponent(),
		NewFooterNaviComponent())

	return c
}

/* Blog Context */

func NewBlogContext(mainnavi, footernavi []Location) Context {
	contentComponents := []Component{
		NewContentComponent(),
		NewDisqusComponent()}
	c := newContext(mainnavi, footernavi, contentComponents)
	return c
}

/* Footer Context */

func NewFooterContext(mainnavi, footernavi []Location) Context {
	contentComponents := []Component{
		NewContentComponent()}
	c := newContext(mainnavi, footernavi, contentComponents)
	return c
}

/* Blog Navi Context */

func NewBlogNaviContext(mainnavi, footernavi []Location) Context {
	contentComponents := []Component{
		NewBlogNaviContextComponent()}
	c := newContext(mainnavi, footernavi, contentComponents)
	return c
}
