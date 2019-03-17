package main

import "github.com/ingmardrewing/staticIntf"

// element bundle constructor
func newElementBundle() *elementBundle {
	return new(elementBundle)
}

type elementBundle struct {
	elements []staticIntf.Page
}

func (l *elementBundle) addElement(e staticIntf.Page) {
	l.elements = append(l.elements, e)
}

func (l *elementBundle) full() bool {
	return len(l.elements) >= 10
}

func (l *elementBundle) getElements() []staticIntf.Page {
	return l.elements
}

func ElementsToLocations(elements []staticIntf.Page) []staticIntf.Location {
	locs := []staticIntf.Location{}
	for _, p := range elements {
		locs = append(locs, p)
	}
	return locs
}
