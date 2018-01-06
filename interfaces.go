package static

import (
	"github.com/ingmardrewing/fs"
	"github.com/ingmardrewing/htmlDoc"
)

type Context interface {
	GetTwitterHandle() string
	GetContentSection() string
	GetContentTags() string
	GetSiteName() string
	GetTwitterCardType() string
	GetOGType() string
	GetFBPageUrl() string
	GetTwitterPage() string
	GetCssUrl() string
	GetCss() string
	GetRssUrl() string
	GetHomeUrl() string
	GetDisqusShortname() string
	GetMainNavigationLocations() []Location
	GetReadNavigationLocations() []Location
	GetFooterNavigationLocations() []Location
	GetElements() []Element
	SetElements([]Element)
	AddComponent(c Component)
	GetComponents() []Component
	RenderPages(targetDir string) []fs.FileContainer
	AddPage(p Element)
	SetGlobalFields(twitterHandle, topic, tags, site, cardType, section, fbPage, twitterPage, cssUrl, rssUrl, home, disqusShortname string)
}

type Component interface {
	VisitPage(p Element)
	GetCss() string
	GetJs() string
	SetContext(context Context)
}

type Location interface {
	GetPath() string
	GetDomain() string
	GetTitle() string
	GetThumbnailUrl() string
	GetFsPath() string
}

type Element interface {
	Location
	GetId() int
	AcceptVisitor(v Component)
	AddBodyNodes([]*htmlDoc.Node)
	AddHeaderNodes([]*htmlDoc.Node)
	GetPublishedTime() string
	GetDescription() string
	GetContent() string
	GetImageUrl() string
	GetDisqusId() string
	Render() string
	GetFsFilename() string
}
