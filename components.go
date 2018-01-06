package static

import (
	"fmt"
	"strings"

	"github.com/ingmardrewing/htmlDoc"
)

/* component */

type abstractComponent struct {
	context Context
}

func (ac *abstractComponent) SetContext(context Context) {
	ac.context = context
}

func (ac *abstractComponent) GetCss() string {
	return ""
}

func (ac *abstractComponent) GetJs() string {
	return ""
}

/* wrapper */
type wrapper struct{}

func (cw *wrapper) wrap(n *htmlDoc.Node, addedclasses ...string) *htmlDoc.Node {
	inner := htmlDoc.NewNode("div", "", "class", "wrapperInner")
	inner.AddChild(n)
	classes := "wrapperOuter " + strings.Join(addedclasses, " ")
	wrapperNode := htmlDoc.NewNode("div", "", "class", classes)
	wrapperNode.AddChild(inner)
	return wrapperNode
}

/* global css component */

func NewGlobalCssComponent() *GlobalCssComponent {
	return new(GlobalCssComponent)
}

type GlobalCssComponent struct {
	abstractComponent
}

func (gcc *GlobalCssComponent) VisitPage(p Element) {}

func (gcc *GlobalCssComponent) GetCss() string {
	return `
body, p, span {
	margin: 0;
	padding: 0;
	font-family: Arial, Helvetica, sans-serif;
}
a {
	color: grey;
	text-decoration: none;
}
a:hover {
	text-decoration: underline;
}
.wrapperOuter {
	text-align: center;
}

.wrapperInner {
	margin: 0 auto;
	width: 800px;
}
p + p {
	margin-top: 10px;
}
`
}

/* GeneralMetaComponent */

func NewGeneralMetaComponent() *GeneralMetaComponent {
	return new(GeneralMetaComponent)
}

type GeneralMetaComponent struct {
	abstractComponent
}

func (g *GeneralMetaComponent) VisitPage(p Element) {
	m := []*htmlDoc.Node{
		htmlDoc.NewNode("meta", "", "viewport", "width=device-width, initial-scale=1.0"),
		htmlDoc.NewNode("meta", "", "http-equiv", "imagetoolbar", "content", "no"),
		htmlDoc.NewNode("meta", "", "robots", "index,follow"),
		htmlDoc.NewNode("meta", "", "author", "Ingmar Drewing"),
		htmlDoc.NewNode("meta", "", "publisher", "Ingmar Drewing"),
		htmlDoc.NewNode("meta", "", "keywords", "storytelling, illustration, drawing, web comic, comic, cartoon, caricatures"),
		htmlDoc.NewNode("meta", "", "DC.subject", "storytelling, illustration, drawing, web comic, comic, cartoon, caricatures"),
		htmlDoc.NewNode("meta", "", "page-topic", "art"),
		htmlDoc.NewNode("meta", "", "http-equiv", "Language", "content", "en"),
		htmlDoc.NewNode("meta", "", "http-equiv", "Content-Language", "content", "en"),
		htmlDoc.NewNode("meta", "", "http-equiv", "cache-control", "content", "Private"),
		htmlDoc.NewNode("meta", "", "http-equiv", "pragma", "content", "Private"),
		htmlDoc.NewNode("meta", "", "http-equiv", "expires", "content", "0"),
		htmlDoc.NewNode("meta", "", "charset", "UTF-8"),
		htmlDoc.NewNode("meta", "", "http-equiv", "content-type", "content", "text/html;charset=UTF-8")}
	p.AddHeaderNodes(m)
}

/* favicon component */

func NewFaviconComponent() *FaviconComponent {
	return new(FaviconComponent)
}

type FaviconComponent struct {
	abstractComponent
}

func (f *FaviconComponent) VisitPage(p Element) {
	m := []*htmlDoc.Node{
		htmlDoc.NewNode("link", "", "rel", "icon", "type", "image/png", "sizes", "16x16", "href", "/icons/favicon-16x16.png"),
		htmlDoc.NewNode("link", "", "rel", "icon", "type", "image/png", "sizes", "32x32", "href", "/icons/favicon-32x32.png"),
		htmlDoc.NewNode("link", "", "rel", "icon", "type", "image/png", "sizes", "192x192", "href", "/icons/android-192x192.png"),
		htmlDoc.NewNode("link", "", "rel", "apple-touch-icon", "type", "image/png", "sizes", "180x180", "href", "/icons/apple-touch-icon-180x180.png"),
		htmlDoc.NewNode("meta", "", "name", "msapplication-config", "content", "/icons/browserconfig.xml")}
	p.AddHeaderNodes(m)
}

/* fb component */
type FBComponent struct {
	abstractComponent
}

func NewFBComponent() *FBComponent {
	fb := new(FBComponent)
	return fb
}

func (fbc *FBComponent) VisitPage(p Element) {
	m := []*htmlDoc.Node{
		htmlDoc.NewNode("meta", "", "property", "og:title", "content", p.GetTitle()),
		htmlDoc.NewNode("meta", "", "property", "og:url", "content", p.GetPath()),
		htmlDoc.NewNode("meta", "", "property", "og:image", "content", p.GetImageUrl()),
		htmlDoc.NewNode("meta", "", "property", "og:description", "content", p.GetDescription()),
		htmlDoc.NewNode("meta", "", "property", "og:site_name", "content", fbc.abstractComponent.context.GetSiteName()),
		htmlDoc.NewNode("meta", "", "property", "og:type", "content", fbc.abstractComponent.context.GetOGType()),
		htmlDoc.NewNode("meta", "", "property", "article:published_time", "content", p.GetPublishedTime()),
		htmlDoc.NewNode("meta", "", "property", "article:modified_time", "content", p.GetPublishedTime()),
		htmlDoc.NewNode("meta", "", "property", "article:section", "content", fbc.abstractComponent.context.GetContentSection()),
		htmlDoc.NewNode("meta", "", "property", "article:tag", "content", fbc.abstractComponent.context.GetContentTags())}

	p.AddHeaderNodes(m)
}

/* google component */

type GoogleComponent struct {
	abstractComponent
}

func NewGoogleComponent() *GoogleComponent {
	gc := new(GoogleComponent)
	return gc
}

func (goo *GoogleComponent) VisitPage(p Element) {
	m := []*htmlDoc.Node{
		htmlDoc.NewNode("meta", "", "itemprop", "name", "content", p.GetTitle()),
		htmlDoc.NewNode("meta", "", "itemprop", "description", "content", p.GetDescription()),
		htmlDoc.NewNode("meta", "", "itemprop", "image", "content", p.GetImageUrl())}
	p.AddHeaderNodes(m)
}

/* twitter component */

type TwitterComponent struct {
	abstractComponent
}

func NewTwitterComponent() *TwitterComponent {
	t := new(TwitterComponent)
	return t
}

func (tw *TwitterComponent) VisitPage(p Element) {
	m := []*htmlDoc.Node{
		htmlDoc.NewNode("meta", "",
			"name", "t:card",
			"content", tw.abstractComponent.context.GetTwitterCardType()),
		htmlDoc.NewNode("meta", "",
			"name", "t:site",
			"content", tw.abstractComponent.context.GetTwitterHandle()),
		htmlDoc.NewNode("meta", "",
			"name", "t:title",
			"content", p.GetTitle()),
		htmlDoc.NewNode("meta", "",
			"name", "t:text:description",
			"content", p.GetDescription()),
		htmlDoc.NewNode("meta", "",
			"name", "t:creator",
			"content", tw.abstractComponent.context.GetTwitterHandle()),
		htmlDoc.NewNode("meta", "",
			"name", "t:image",
			"content", p.GetImageUrl())}
	p.AddHeaderNodes(m)
}

/* title component */
type TitleComponent struct {
	abstractComponent
}

func NewTitleComponent() *TitleComponent {
	return new(TitleComponent)
}

func (tc *TitleComponent) VisitPage(p Element) {
	title := htmlDoc.NewNode("title", p.GetTitle())
	p.AddHeaderNodes([]*htmlDoc.Node{title})
}

/* css link component */

type CssLinkComponent struct {
	abstractComponent
}

func NewCssLinkComponent() *CssLinkComponent {
	clc := new(CssLinkComponent)
	return clc
}

func (clc *CssLinkComponent) VisitPage(p Element) {
	link := htmlDoc.NewNode("link", "", "href", clc.abstractComponent.context.GetCssUrl(), "rel", "stylesheet", "type", "text/css")
	p.AddHeaderNodes([]*htmlDoc.Node{link})
}

/**/

type BlogNaviComponent struct {
	wrapper
	abstractComponent
}

func NewBlogNaviContextComponent() *BlogNaviComponent {
	bnc := new(BlogNaviComponent)
	return bnc
}

func (b *BlogNaviComponent) addPrevious(p Element, n *htmlDoc.Node) {
	inx := b.getIndexOfPage(p)
	if inx == 0 {
		span := htmlDoc.NewNode("span", "< previous posts", "class", "blognavicomponent__previous")
		n.AddChild(span)
	} else {
		elems := b.abstractComponent.context.GetElements()
		pv := elems[inx-1]
		a := htmlDoc.NewNode("a", "< previous posts", "href", pv.GetPath(), "rel", "prev", "class", "blognavicomponent__previous")
		n.AddChild(a)
	}
}

func (b *BlogNaviComponent) addNext(p Element, n *htmlDoc.Node) {
	inx := b.getIndexOfPage(p)
	if inx == len(b.abstractComponent.context.GetElements())-1 {
		span := htmlDoc.NewNode("span", "next posts >", "class", "blognavicomponent__next")
		n.AddChild(span)
	} else {
		elems := b.abstractComponent.context.GetElements()
		nx := elems[inx+1]
		a := htmlDoc.NewNode("a", "next posts >", "href", nx.GetPath(), "rel", "next", "class", "blognavicomponent__next")
		n.AddChild(a)
	}
}

func (b *BlogNaviComponent) addBodyNodes(p Element) {
	nav := htmlDoc.NewNode("nav", "", "class", "blognavicomponent__nav")
	b.addPrevious(p, nav)
	b.addNext(p, nav)
	d := htmlDoc.NewNode("div", "", "class", "blognavicomponent")
	d.AddChild(htmlDoc.NewNode("div", p.GetContent()))
	d.AddChild(nav)
	wn := b.wrap(d)
	p.AddBodyNodes([]*htmlDoc.Node{wn})
}

func (b *BlogNaviComponent) VisitPage(p Element) {
	if len(b.abstractComponent.context.GetElements()) < 3 {
		return
	}
	b.addBodyNodes(p)
}

func (b *BlogNaviComponent) getIndexOfPage(p Element) int {
	for i, l := range b.abstractComponent.context.GetElements() {
		if l.GetPath() == p.GetPath() {
			return i
		}
	}
	return -1
}

func (b *BlogNaviComponent) GetCss() string {
	return `
.blognavicomponent {
	text-align: left;
	padding-top: 123px;
	padding-bottom: 20px;
}
.blognavicomponent__nav {
	text-align: center;
	color: lightgrey;
	margin-top: 40px;
	margin-bottom: 20px;
}
.blognavicomponent__nav span {
	font-family: Arial Black, Arial, Helvetica, sans-serif;
	color: lightgrey;
	font-weight: 900;
}
.blognavicomponent__next {
	margin-left: 10px;
}
.blognavicomponent__previous,
.blognavicomponent__next {
	font-family: Arial Black, Arial, Helvetica, sans-serif;
	color: grey;
	text-transform: uppercase;
	font-weight: 900;
	font-size: 16px;
}

a.blognavientry__tile {
	display: block;
	position: relative;
	width: 390px;
	height: 470px;
	margin-bottom: 20px;
	float: left;
	text-decoration: none;
}

.blognavientry__tile:nth-child(odd) {
	margin-right: 20px;
}

.blognavientry__image {
	display: block;
	width: 390px;
	height: 390px;
	background-size: cover;
}
.blognavientry__tile h2 {
	font-family: Arial Black, Arial, Helvetica, sans-serif;
	text-transform: uppercase;
	color: black;
	margin-top: 4px;
	line-height: 24px;
}
.blognavientry__tile:hover h2 {
	color: grey;
}
`
}

/* MainNaviComponent */
func NewMainNaviComponent() *MainNaviComponent {
	nc := new(MainNaviComponent)
	return nc
}

type MainNaviComponent struct {
	abstractComponent
	wrapper
	cssClass string
}

func (nv *MainNaviComponent) VisitPage(p Element) {
	nav := htmlDoc.NewNode("nav", "",
		"class", "mainnavi")
	url := p.GetPath()
	for _, l := range nv.abstractComponent.context.GetMainNavigationLocations() {
		if url == l.GetPath() {
			span := htmlDoc.NewNode("span", l.GetTitle(),
				"class", "mainnavi__navelement--current")
			nav.AddChild(span)
		} else {
			a := htmlDoc.NewNode("a", l.GetTitle(),
				"href", l.GetPath(),
				"class", "mainnavi__navelement")
			nav.AddChild(a)
		}
	}
	node := htmlDoc.NewNode("div", "", "class", nv.cssClass)
	node.AddChild(nav)
	wn := nv.wrap(node, "mainnavi__wrapper")
	p.AddBodyNodes([]*htmlDoc.Node{wn})
}

func (mhc *MainNaviComponent) GetJs() string {
	return ""
}

func (mhc *MainNaviComponent) GetCss() string {
	return `
.mainnavi {
	border-top: 1px solid black;
	border-bottom: 2px solid black;
}
.mainnavi__wrapper {
	position: fixed;
	width: 100%;
	top: 80px;
	background-color: white;
}
.mainnavi__navelement--current,
a.mainnavi__navelement {
	display: inline-block;
	font-family: Arial Black, Arial, Helvetica, sans-serif;
	font-weight: 900;
	font-size: 18px;
	line-height: 20px;
	text-transform: uppercase;
	color: black;
	padding: 10px 20px;
}
.mainnavi__navelement--current,
a.mainnavi__navelement:hover {
	text-decoration: none;
	color: gray;
}
.mainnavi__nav {
	border-bottom: 2px solid black;
}
`
}

/* FooterNaviComponent */

func NewFooterNaviComponent() *FooterNaviComponent {
	nc := new(FooterNaviComponent)
	return nc
}

type FooterNaviComponent struct {
	abstractComponent
	wrapper
	cssClass string
}

func (nv *FooterNaviComponent) VisitPage(p Element) {
	nav := htmlDoc.NewNode("nav", "",
		"class", "footernavi")
	url := p.GetPath()
	for _, l := range nv.abstractComponent.context.GetFooterNavigationLocations() {
		if url == l.GetPath() {
			span := htmlDoc.NewNode("span", l.GetTitle(),
				"class", "footernavi__navelement--current")
			nav.AddChild(span)
		} else {
			a := htmlDoc.NewNode("a", l.GetTitle(),
				"href", l.GetPath(),
				"class", "footernavi__navelement")
			nav.AddChild(a)
		}
	}
	node := htmlDoc.NewNode("div", "", "class", nv.cssClass)
	node.AddChild(nav)
	wn := nv.wrap(node, "footernavi__wrapper")
	p.AddBodyNodes([]*htmlDoc.Node{wn})
}

func (mhc *FooterNaviComponent) GetJs() string { return "" }

func (mhc *FooterNaviComponent) GetCss() string {
	return `
.footernavi {
	border-top: 1px solid black;
}
.footernavi__wrapper {
	position: fixed;
	width: 100%;
	bottom: 0;
	background-color: white;
}
.footernavi__navelement--current ,
a.footernavi__navelement {
	display: inline-block;
	font-family: Arial Black, Arial, Helvetica, sans-serif;
	font-weight: 900;
	font-size: 16px;
	line-height: 20px;
	text-transform: uppercase;
	text-decoration: none;
	color: black;
	padding: 10px 15px;
}
a.footernavi__navelement:hover,
.footernavi__navelement--current {
	color: gray;
}
`
}

/* ReadNaviComponent */

type ReadNaviComponent struct {
	abstractComponent
	wrapper
}

func NewReadNaviComponent() *ReadNaviComponent {
	rnc := new(ReadNaviComponent)
	return rnc
}

func (rnv *ReadNaviComponent) addFirst(p Element, n *htmlDoc.Node) {
	inx := rnv.getIndexOfPage(p)
	locations := rnv.abstractComponent.context.GetReadNavigationLocations()
	if inx == 0 {
		span := htmlDoc.NewNode("span", "<< first")
		n.AddChild(span)
	} else {
		f := locations[0]
		a := htmlDoc.NewNode("a", "<< first", "href", f.GetPath(), "rel", "first")
		n.AddChild(a)
	}
}

func (rnv *ReadNaviComponent) addPrevious(p Element, n *htmlDoc.Node) {
	inx := rnv.getIndexOfPage(p)
	locations := rnv.abstractComponent.context.GetReadNavigationLocations()
	if inx == 0 {
		span := htmlDoc.NewNode("span", "< previous")
		n.AddChild(span)
	} else {
		p := locations[inx-1]
		a := htmlDoc.NewNode("a", "< previous", "href", p.GetPath(), "rel", "prev")
		n.AddChild(a)
	}
}

func (rnv *ReadNaviComponent) addNext(p Element, n *htmlDoc.Node) {
	inx := rnv.getIndexOfPage(p)
	locations := rnv.abstractComponent.context.GetReadNavigationLocations()
	if inx == len(locations)-1 {
		span := htmlDoc.NewNode("span", "next >")
		n.AddChild(span)
	} else {
		nx := locations[inx+1]
		a := htmlDoc.NewNode("a", "next >", "href", nx.GetPath(), "rel", "next")
		n.AddChild(a)
	}
}

func (rnv *ReadNaviComponent) addLast(p Element, n *htmlDoc.Node) {
	inx := rnv.getIndexOfPage(p)
	locations := rnv.abstractComponent.context.GetReadNavigationLocations()
	if inx == len(locations)-1 {
		span := htmlDoc.NewNode("span", "newest >>")
		n.AddChild(span)
	} else {
		nw := locations[len(locations)-1]
		a := htmlDoc.NewNode("a", "neweset >>", "href", nw.GetPath(), "rel", "last")
		n.AddChild(a)
	}
}

func (rnv *ReadNaviComponent) addHeaderNodes(p Element) {
	inx := rnv.getIndexOfPage(p)
	locations := rnv.abstractComponent.context.GetReadNavigationLocations()
	n := []*htmlDoc.Node{}
	firstUrl := locations[0].GetPath()
	n = append(n, htmlDoc.NewNode("link", "", "rel", "first", "href", firstUrl))
	if inx > 0 {
		prevUrl := locations[inx-1].GetPath()
		n = append(n, htmlDoc.NewNode("link", "", "rel", "prev", "href", prevUrl))
	}
	if inx < len(locations)-1 {
		nextUrl := locations[inx+1].GetPath()
		n = append(n, htmlDoc.NewNode("link", "", "rel", "next", "href", nextUrl))
	}
	lastUrl := locations[len(locations)-1].GetPath()
	n = append(n, htmlDoc.NewNode("link", "", "rel", "last", "href", lastUrl))
	p.AddHeaderNodes(n)
}

func (rnv *ReadNaviComponent) addBodyNodes(p Element) {
	bodyNav := htmlDoc.NewNode("nav", "")
	rnv.addFirst(p, bodyNav)
	rnv.addPrevious(p, bodyNav)
	rnv.addNext(p, bodyNav)
	rnv.addLast(p, bodyNav)
	wn := rnv.wrap(bodyNav)
	p.AddBodyNodes([]*htmlDoc.Node{wn})
}

func (rnv *ReadNaviComponent) VisitPage(p Element) {
	if len(rnv.abstractComponent.context.GetReadNavigationLocations()) < 3 {
		return
	}
	rnv.addHeaderNodes(p)
	rnv.addBodyNodes(p)

}

func (rnv *ReadNaviComponent) getIndexOfPage(p Element) int {
	locations := rnv.abstractComponent.context.GetReadNavigationLocations()
	for i, l := range locations {
		if l.GetPath() == p.GetPath() {
			return i
		}
	}
	return -1
}

/* disqus component */

type DisqusComponent struct {
	abstractComponent
	wrapper
	configuredJs string
}

func NewDisqusComponent() *DisqusComponent {
	d := new(DisqusComponent)
	return d
}

func (dc *DisqusComponent) GetCss() string {
	return `
.diqus,
.diqus p {
	font-family: Arial, Helvetica, sans-serif;
}
`
}

func (dc *DisqusComponent) GetJs() string {
	return dc.configuredJs
}

func (dc *DisqusComponent) VisitPage(p Element) {
	dc.configuredJs = fmt.Sprintf(
		`
var disqus_config = function () {
	this.page.title= "%s";
	this.page.url = '%s';
	this.page.identifier =  '%s';
};
(function() {
	var d = document, s = d.createElement('script');
	s.src = 'https://%s.disqus.com/embed.js';
	s.setAttribute('data-timestamp', +new Date());
	(d.head || d.body).appendChild(s);
})();
`, p.GetTitle(), p.GetDomain()+p.GetPath(), p.GetDisqusId(), dc.abstractComponent.context.GetDisqusShortname())
	n := htmlDoc.NewNode("div", " ", "id", "disqus_thread", "class", "disqus")
	wn := dc.wrap(n)
	p.AddBodyNodes([]*htmlDoc.Node{wn})
}

/* main  header component */

type MainHeaderComponent struct {
	abstractComponent
	wrapper
}

func NewMainHeaderComponent() *MainHeaderComponent {
	mhc := new(MainHeaderComponent)
	return mhc
}

func (mhc *MainHeaderComponent) VisitPage(p Element) {
	logo := htmlDoc.NewNode("a", "<!-- logo -->",
		"href", mhc.abstractComponent.context.GetHomeUrl(),
		"class", "headerbar__logo")
	logocontainer := htmlDoc.NewNode("div", "",
		"class", "headerbar__logocontainer")
	logocontainer.AddChild(logo)

	header := htmlDoc.NewNode("header", "", "class", "headerbar")
	header.AddChild(logocontainer)

	wn := mhc.wrap(header, "headerbar__wrapper")
	p.AddBodyNodes([]*htmlDoc.Node{wn})
}

func (mhc *MainHeaderComponent) GetJs() string {
	return ""
}

func (mhc *MainHeaderComponent) GetCss() string {
	return `
.headerbar__wrapper {
	position: fixed;
	width: 100%;
	top: 0;
	background-color: white;
}
.headerbar__logo {
	background-image: url(https://s3.amazonaws.com/drewingdeblog/drewing_de_logo.png);
	background-repeat: no-repeat;
	background-position: center center;
	display: block;
	width: 100%;
	height: 80px;
}
.headerbar__navelement {
	display: inline-block;
	font-family: Arial Black, Arial, Helvetica, sans-serif;
	font-weight: 900;
	font-size: 18px;
	line-height: 20px;
	text-transform: uppercase;
	text-decoration: none;
	color: black;
	padding: 10px 20px;
}
`
}

/* content component */
type ContentComponent struct {
	abstractComponent
	wrapper
}

func NewContentComponent() *ContentComponent {
	return new(ContentComponent)
}

func (cc *ContentComponent) VisitPage(p Element) {
	h1 := htmlDoc.NewNode("h1", p.GetTitle(),
		"class", "maincontent__h1")
	h2 := htmlDoc.NewNode("h2", p.GetPublishedTime(),
		"class", "maincontent__h2")
	n := htmlDoc.NewNode("main", p.GetContent(),
		"class", "maincontent")
	n.AddChild(h1)
	n.AddChild(h2)
	wn := cc.wrap(n)
	p.AddBodyNodes([]*htmlDoc.Node{wn})
}

func (cc *ContentComponent) GetJs() string { return "" }

func (cc *ContentComponent) GetCss() string {
	return `
.maincontent{
	padding-top: 126px;
	text-align: left;
}
.maincontent__h1,
.maincontent__h2 {
	display: inline-block;
	font-family: Arial Black, Arial, Helvetica, sans-serif;
	text-transform: uppercase;
}
.maincontent__h1 ,
.maincontent__h2 {
	font-size: 18px;
	line-height: 20px;
}
.maincontent__h2 {
	color: grey;
	margin-left: 10px;
}
`
}

/* gallery component */

type GalleryComponent struct {
	wrapper
}

func NewGalleryComponent() *GalleryComponent {
	gc := new(GalleryComponent)
	return gc
}

func (gal *GalleryComponent) VisitPage(p Element) {
	inner := htmlDoc.NewNode("div", "", "class", "maincontent__inner")
	for i := 0; i < 5; i++ {
		title := htmlDoc.NewNode("span", "At The Zoo", "class", "portfoliothumb__title")
		subtitle := htmlDoc.NewNode("span", "Digital drawing", "class", "portfoliothumb__details")

		label := htmlDoc.NewNode("div", "", "class", "portfoliothumb__label")
		label.AddChild(title)
		label.AddChild(subtitle)

		img := htmlDoc.NewNode("img", "", "class", "portfoliothumb__image", "src", "https://s3.amazonaws.com/drewingdeblog/blog/wp-content/uploads/2017/12/02152842/atthezoo-400x400.png")

		div := htmlDoc.NewNode("a", "", "class", "portfoliothumb", "href", "https://drewing.de")
		div.AddChild(img)
		div.AddChild(label)

		inner.AddChild(div)
	}

	m := htmlDoc.NewNode("main", "", "class", "maincontent")
	m.AddChild(inner)
	wn := gal.wrap(m)
	p.AddBodyNodes([]*htmlDoc.Node{wn})
}

func (gal *GalleryComponent) getCss() string { return `` }

/* copyright component */
type CopyRightComponent struct {
	abstractComponent
	wrapper
}

func NewCopyRightComponent() *CopyRightComponent {
	c := new(CopyRightComponent)
	return c
}

func (crc *CopyRightComponent) VisitPage(p Element) {
	n := htmlDoc.NewNode("div", `
	<a rel="license" class="copyright__cc" href="https://creativecommons.org/licenses/by-nc-nd/3.0/"></a>
	<p class="copyright__license">&copy; 2017 by Ingmar Drewing </p><p class="copyright__license">Except where otherwise noted, content on this site is licensed under a <a rel="license" href="https://creativecommons.org/licenses/by-nc-nd/3.0/">Creative Commons Attribution-NonCommercial-NoDerivs 3.0 Unported (CC BY-NC-ND 3.0) license</a>.</p><p class="copyright__license">Soweit nicht anders explizit ausgewiesen, stehen die Inhalte auf dieser Website unter der <a rel="license" href="https://creativecommons.org/licenses/by-nc-nd/3.0/">Creative Commons Namensnennung-NichtKommerziell-KeineBearbeitung (CC BY-NC-ND 3.0)</a> Lizenz. Unless otherwise noted the author of the content on this page is <a href="https://plus.google.com/113943655600557711368?rel=author">Ingmar Drewing</a></p>
	`, "class", "copyright")
	wn := crc.wrap(n)
	p.AddBodyNodes([]*htmlDoc.Node{wn})
}

func (crc *CopyRightComponent) GetCss() string {
	return `
.copyright {
	background-color: black;
	color: white;
	text-align: left;
	font-family: Arial, Helvetica, sans-serif;
	font-size: 14px;
	padding: 20px 20px 50px;
	margin-top: 20px;
}
.copyright__license {
	margin-top: 20px;
	margin-bottom: 20px;
}
.copyright__cc {
    display: block;
    border-width: 0;
    background-image: url(https://i.creativecommons.org/l/by-nc-nd/3.0/88x31.png);
    width: 88px;
    height: 31px;
    margin-right: 15px;
    margin-bottom: 5px;
}
`
}

func (crc *CopyRightComponent) GetJs() string { return `` }

/* cookie notifier component */

type CookieNotifierComponent struct {
}

func (cnc *CookieNotifierComponent) VisitPage(p Element) {}

func (cnc *CookieNotifierComponent) getCss() string { return `` }

func (cnc *CookieNotifierComponent) getJs() string {
	return `
function cli_show_cookiebar(p) {
	var Cookie = {
		set: function(name,value,days) {
			if (days) {
				var date = new Date();
				date.setTime(date.getTime()+(days*24*60*60*1000));
				var expires = "; expires="+date.toGMTString();
			}
			else var expires = "";
			document.cookie = name+"="+value+expires+"; path=/";
		},
		read: function(name) {
			var nameEQ = name + "=";
			var ca = document.cookie.split(';');
			for(var i=0;i < ca.length;i++) {
				var c = ca[i];
				while (c.charAt(0)==' ') {
					c = c.substring(1,c.length);
				}
				if (c.indexOf(nameEQ) === 0) {
					return c.substring(nameEQ.length,c.length);
				}
			}
			return null;
		},
		erase: function(name) {
			this.set(name,"",-1);
		},
		exists: function(name) {
			return (this.read(name) !== null);
		}
	};

	var ACCEPT_COOKIE_NAME = 'viewed_cookie_policy',
		ACCEPT_COOKIE_EXPIRE = 365,
		json_payload = p.settings;

	if (typeof JSON.parse !== "function") {
		console.log("CookieLawInfo requires JSON.parse but your browser doesn't support it");
		return;
	}
	var settings = JSON.parse(json_payload);

	var cached_header = jQuery(settings.notify_div_id),
		cached_showagain_tab = jQuery(settings.showagain_div_id),
		btn_accept = jQuery('#cookie_hdr_accept'),
		btn_decline = jQuery('#cookie_hdr_decline'),
		btn_moreinfo = jQuery('#cookie_hdr_moreinfo'),
		btn_settings = jQuery('#cookie_hdr_settings');

	cached_header.hide();
	if ( !settings.showagain_tab ) {
		cached_showagain_tab.hide();
	}

	var hdr_args = { };

	var showagain_args = { };
	cached_header.css( hdr_args );
	cached_showagain_tab.css( showagain_args );

	if (!Cookie.exists(ACCEPT_COOKIE_NAME)) {
		displayHeader();
	}
	else {
		cached_header.hide();
	}

	if ( settings.show_once_yn ) {
		setTimeout(close_header, settings.show_once);
	}
	function close_header() {
		Cookie.set(ACCEPT_COOKIE_NAME, 'yes', ACCEPT_COOKIE_EXPIRE);
		hideHeader();
	}

	var main_button = jQuery('.cli-plugin-main-button');
	main_button.css( 'color', settings.button_1_link_colour );

	if ( settings.button_1_as_button ) {
		main_button.css('background-color', settings.button_1_button_colour);

		main_button.hover(function() {
			jQuery(this).css('background-color', settings.button_1_button_hover);
		},
		function() {
			jQuery(this).css('background-color', settings.button_1_button_colour);
		});
	}
	var main_link = jQuery('.cli-plugin-main-link');
	main_link.css( 'color', settings.button_2_link_colour );

	if ( settings.button_2_as_button ) {
		main_link.css('background-color', settings.button_2_button_colour);

		main_link.hover(function() {
			jQuery(this).css('background-color', settings.button_2_button_hover);
		},
		function() {
			jQuery(this).css('background-color', settings.button_2_button_colour);
		});
	}

	cached_showagain_tab.click(function(e) {
		e.preventDefault();
		cached_showagain_tab.slideUp(settings.animate_speed_hide, function slideShow() {
			cached_header.slideDown(settings.animate_speed_show);
		});
	});

	jQuery("#cookielawinfo-cookie-delete").click(function() {
		Cookie.erase(ACCEPT_COOKIE_NAME);
		return false;
	});
	jQuery("#cookie_action_close_header").click(function(e) {
		e.preventDefault();
		accept_close();
	});

	function accept_close() {
		Cookie.set(ACCEPT_COOKIE_NAME, 'yes', ACCEPT_COOKIE_EXPIRE);

		if (settings.notify_animate_hide) {
			cached_header.slideUp(settings.animate_speed_hide);
		}
		else {
			cached_header.hide();
		}
		cached_showagain_tab.slideDown(settings.animate_speed_show);
		return false;
	}

	function closeOnScroll() {
		if (window.pageYOffset > 100 && !Cookie.read(ACCEPT_COOKIE_NAME)) {
			accept_close();
			if (settings.scroll_close_reload === true) {
				location.reload();
			}
			window.removeEventListener("scroll", closeOnScroll, false);
		}
	}
	if (settings.scroll_close === true) {
		window.addEventListener("scroll", closeOnScroll, false);
	}

	function displayHeader() {
		if (settings.notify_animate_show) {
			cached_header.slideDown(settings.animate_speed_show);
		}
		else {
			cached_header.show();
		}
		cached_showagain_tab.hide();
	}
	function hideHeader() {
		if (settings.notify_animate_show) {
			cached_showagain_tab.slideDown(settings.animate_speed_show);
		}
		else {
			cached_showagain_tab.show();
		}
		cached_header.slideUp(settings.animate_speed_show);
	}
};

function l1hs(str){if(str.charAt(0)=="#"){str=str.substring(1,str.length);}else{return "#"+str;}return l1hs(str);}

cli_show_cookiebar({
					settings: '{"animate_speed_hide":"500","animate_speed_show":"500","background":"#fff","border":"#444","border_on":true,"button_1_button_colour":"#000","button_1_button_hover":"#000000","button_1_link_colour":"#fff","button_1_as_button":true,"button_2_button_colour":"#333","button_2_button_hover":"#292929","button_2_link_colour":"#444","button_2_as_button":false,"font_family":"inherit","header_fix":false,"notify_animate_hide":true,"notify_animate_show":false,"notify_div_id":"#cookie-law-info-bar","notify_position_horizontal":"right","notify_position_vertical":"bottom","scroll_close":false,"scroll_close_reload":false,"showagain_tab":false,"showagain_background":"#fff","showagain_border":"#000","showagain_div_id":"#cookie-law-info-again","showagain_x_position":"100px","text":"#000","show_once_yn":false,"show_once":"10000"}'
});

`
}
