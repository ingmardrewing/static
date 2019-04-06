package main

import (
	"github.com/ingmardrewing/staticPersistence"
	log "github.com/sirupsen/logrus"
)

// Creates a new sitesController, which creates
// multiple sites based on the given json config
func NewSitesController(configs []staticPersistence.Config) *sitesController {
	c := new(sitesController)
	c.configs = configs
	return c
}

// the sitesController struct
type sitesController struct {
	configs []staticPersistence.Config
}

// Intended for migrational purposes
func (s *sitesController) UpdateJsonFiles() {
	for _, config := range s.configs {
		RewriteJson(config)
	}
}

// Renders the sites defined by the Json config
func (s *sitesController) UpdateStaticSites() {
	for _, config := range s.configs {
		log.Debug("sites.Controller.UpdateStaticSites - Creating Site:" + config.Domain)
		siteCreator := NewSiteCreator(config)
		siteCreator.addSite()
		siteCreator.addSources()
		siteCreator.addContainers()
		siteCreator.addLocations()
		siteCreator.addContexts()
		siteCreator.fillFileContainers(config)
		siteCreator.writeFiles()
	}
}
