package main

import (
	"strconv"

	"github.com/ingmardrewing/staticIntf"
	"github.com/ingmardrewing/staticModel"
)

func NewBlogNaviPageGenerator(
	site staticIntf.Site,
	path string,
	container staticIntf.PagesContainer) *blogNaviPageGenerator {

	b := new(blogNaviPageGenerator)
	b.site = site
	b.path = path
	b.container = container
	return b
}

type blogNaviPageGenerator struct {
	pages     []staticIntf.Page
	site      staticIntf.Site
	path      string
	container staticIntf.PagesContainer
}

func (n *blogNaviPageGenerator) Createpages() []staticIntf.Page {
	bundles := n.generateBundles()
	last := len(bundles) - 1
	naviPages := make([]staticIntf.Page, 0)

	if n.site == nil {
		return naviPages
	}
	for i, bundle := range bundles {
		filename := "index" + strconv.Itoa(i) + ".html"
		if i == last {
			filename = "index.html"
		}

		pm := staticModel.NewPageMaker()
		pm.Title(n.site.Domain() + " Overview")
		pm.Category("blog post navi")
		pm.PathFromDocRoot(n.path)
		pm.FileName(filename)
		pm.Site(n.site)
		pm.NavigatedPages(bundle...)

		naviPages = append(naviPages, pm.Make())
	}

	return naviPages
}

func (n *blogNaviPageGenerator) getReversedPages() []staticIntf.Page {
	pages := n.container.Pages()
	length := len(pages)
	reversed := make([]staticIntf.Page, 0)
	for i := length - 1; i >= 0; i-- {
		reversed = append(reversed, pages[i])
	}
	return reversed
}

func (n *blogNaviPageGenerator) generateBundles() [][]staticIntf.Page {
	reversedBundles := n.generateReversedBundles()
	length := len(reversedBundles)
	pageBundles := [][]staticIntf.Page{}
	for i := length - 1; i >= 0; i-- {
		pageBundles = append(pageBundles, reversedBundles[i].getElements())
	}
	return pageBundles
}

func (n *blogNaviPageGenerator) generateReversedBundles() []*elementBundle {
	reversed := n.getReversedPages()
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
	return bundles
}
