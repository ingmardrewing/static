package main

import (
	"fmt"

	"github.com/ingmardrewing/staticIntf"
	"github.com/ingmardrewing/staticModel"
	"github.com/ingmardrewing/staticPersistence"
	"github.com/ingmardrewing/staticPresentation"

	log "github.com/sirupsen/logrus"
)

//
type source interface {
	generate()
	Container() staticIntf.PagesContainer
	CreateContext() staticIntf.Context
	SetData(variant, headline, dir, subDir string, site staticIntf.Site, config staticPersistence.JsonConfig)
}

func NewSource(
	variant, dir, subDir, headline string,
	site staticIntf.Site,
	config staticPersistence.JsonConfig) source {

	log.Debugf("NewSource() called for variant %s\n", variant)
	var s source
	switch variant {
	case staticIntf.HOME:
		s = new(homeSource)
	case staticIntf.BLOG:
		s = new(blogSource)
	case staticIntf.PORTFOLIO:
		s = new(portfolioSource)
	case staticIntf.MARGINALS:
		s = new(marginalSource)
	case staticIntf.NARRATIVES:
		s = new(narrativeSource)
	case staticIntf.NARRATIVEMARGINALS:
		s = new(narrativeMarginalSource)
	default:
		log.Fatal("unaccounted variant:", variant)
	}

	s.SetData(variant, headline, dir, subDir, site, config)

	return s
}

type blogSource struct {
	defaultSource
}

func (bs *blogSource) generate() {
	bs.generateContainer()

	bnpg := NewBlogNaviPageGenerator(
		bs.site,
		bs.subDir,
		bs.container)
	naviPages := bnpg.Createpages()
	for _, p := range naviPages {
		bs.container.AddNaviPage(p)
		p.Container(bs.container)
	}
	pages := bs.container.Pages()
	nrOfRepPages := 4
	if len(pages) > nrOfRepPages {
		for _, pg := range pages[len(pages)-nrOfRepPages:] {
			bs.container.AddRepresentational(pg)
		}
	}
}

func (bs *blogSource) CreateContext() staticIntf.Context {
	return staticPresentation.NewBlogContext(bs.site)
}

//
type portfolioSource struct {
	defaultSource
}

func (ps *portfolioSource) generate() {
	ps.generateContainer()

	pages := ps.container.Pages()
	log.Debugf("portfolioSource.generate() with %d pages\n", len(pages))
	for _, pg := range pages {
		ps.container.AddRepresentational(pg)
	}
}

func (ps *portfolioSource) CreateContext() staticIntf.Context {
	return staticPresentation.NewPortfolioContext(ps.site)
}

//
type homeSource struct {
	defaultSource
}

func (hs *homeSource) generate() { hs.generateContainer() }

func (hs *homeSource) CreateContext() staticIntf.Context {
	return staticPresentation.NewHomeContextGroup(hs.site)
}

//
type narrativeMarginalSource struct {
	defaultSource
}

func (nms *narrativeMarginalSource) generate() {
	nms.generateContainer()
}

func (nms *narrativeMarginalSource) CreateContext() staticIntf.Context {
	return staticPresentation.NewMarginalContextGroup(nms.site)
}

//
type marginalSource struct {
	defaultSource
}

func (mrs *marginalSource) generate() {
	mrs.generateContainer()
	locs := ElementsToLocations(mrs.container.Pages())
	for _, l := range locs {
		mrs.site.AddMarginal(l)
	}

}

func (mrs *marginalSource) CreateContext() staticIntf.Context {
	return staticPresentation.NewMarginalContextGroup(mrs.site)
}

//
type narrativeSource struct {
	defaultSource
}

func (ns *narrativeSource) generate() {
	ns.generateContainer()

	pages := ns.container.Pages()
	nrOfRepPages := 4
	if len(pages) > nrOfRepPages {
		for _, pg := range pages[len(pages)-nrOfRepPages:] {
			ns.container.AddRepresentational(pg)
		}
	}
}

func (ns *narrativeSource) CreateContext() staticIntf.Context {
	return staticPresentation.NewNarrativeContextGroup(ns.site)
}

//

type defaultSource struct {
	variant   string
	headline  string
	dir       string
	subDir    string
	site      staticIntf.Site
	config    staticPersistence.JsonConfig
	container staticIntf.PagesContainer
}

func (a *defaultSource) CreateContext() staticIntf.Context {
	return nil
}

func (a *defaultSource) Container() staticIntf.PagesContainer {
	return a.container
}

func (a *defaultSource) generate() {}

func (a *defaultSource) SetData(variant, headline, dir, subDir string, site staticIntf.Site, config staticPersistence.JsonConfig) {
	a.variant = variant
	a.headline = headline
	a.dir = dir
	a.subDir = subDir
	a.site = site
	a.config = config
}

func (a *defaultSource) generateContainer() {
	log.Debug(fmt.Sprintf("-- new container, type %s, headline %s", a.variant, a.headline))
	a.container = staticModel.NewPagesContainer(a.variant, a.headline)
	pageDtos := staticPersistence.ReadPagesFromDir(a.dir)
	log.Debugf("defaultSource.generateContainer() with %d pageDtos", len(pageDtos))
	for _, dto := range pageDtos {
		a.createPage(dto)
	}
}

func (a *defaultSource) createPage(dto staticIntf.PageDto) {
	p := staticModel.NewPage(dto, a.config.Domain, a.site)
	if p == nil {
		log.Error("defaultSource.createPage() - newly created page is nil")
	}
	log.Debugf("createPage(), %s", p.Url())
	a.container.AddPage(p)
	p.Container(a.container)
}
