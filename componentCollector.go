package main

import (
	"reflect"

	"github.com/ingmardrewing/staticIntf"
)

func NewComponentCollector() *componentCollector {
	return new(componentCollector)
}

type componentCollector struct {
	components []staticIntf.Component
}

func (c *componentCollector) GetComponents() []staticIntf.Component {
	return c.components
}

func (c *componentCollector) AddComponents(comps []staticIntf.Component) {
	for _, comp := range comps {
		if !c.componentExists(comp) {
			c.components = append(c.components, comp)
		}
	}
}

func (c *componentCollector) componentExists(givenComp staticIntf.Component) bool {
	for _, comp := range c.components {
		if reflect.TypeOf(comp) == reflect.TypeOf(givenComp) {
			return true
		}
	}
	return false
}
