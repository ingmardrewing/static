package main

import (
	"fmt"
	"path"
	"reflect"
	"strings"

	"github.com/ingmardrewing/fs"
	"github.com/ingmardrewing/staticIntf"
	"github.com/ingmardrewing/staticModel"
	"github.com/ingmardrewing/staticPersistence"
	log "github.com/sirupsen/logrus"
)

// Creates a new siteCreator from the part of the
// JsonConfig specific to one site. The complete
// config can define several sites.
func NewSiteCreator(config staticPersistence.JsonConfig) *siteCreator {
	siteCreator := new(siteCreator)
	siteCreator.config = config
	return siteCreator
}

// Intended for migrating old Json data, which still
// has the form of wordpress exported json data, to
// the newer streamlined version.
func RewriteJson(config staticPersistence.JsonConfig) {
	siteCreator := new(siteCreator)
	siteCreator.config = config
	siteCreator.rewritePages()
}

// The siteCreator handles the creation of one
// web site, located under one domain.
type siteCreator struct {
	site           staticIntf.Site
	config         staticPersistence.JsonConfig
	sources        []source
	contexts       []staticIntf.Context
	fileContainers []fs.FileContainer
}

// Creates and adds a siteDto with the data
// read from the corresponding part of the config
func (s *siteCreator) addSite() {
	log.Debug("siteCreator.addSite()")
	s.site = staticModel.NewSiteDto(
		s.config.Context.TwitterHandle,
		s.config.Context.Topic,
		s.config.Context.Tags,
		s.config.Domain,
		s.config.BasePath,
		s.config.Context.CardType,
		s.config.Context.Section,
		s.config.Context.FbPage,
		s.config.Context.TwitterPage,
		s.config.Deploy.RssPath,
		s.config.Deploy.RssFilename,
		s.config.Deploy.CssFileName,
		s.config.Context.DisqusShortname,
		s.config.Deploy.TargetDir,
		s.config.DefaultMeta.BlogExcerpt,
		s.config.DefaultMeta.KeyWords,
		s.config.DefaultMeta.Subject,
		s.config.DefaultMeta.Author,
		s.config.HomeText,
		s.config.HomeHeadline,
		s.config.SvgLogo)
}

// Adds a single context to the slice of contexts
// within the siteCreator
func (s *siteCreator) addContext(cg staticIntf.Context) {
	if !s.contextExists(cg) {
		s.contexts = append(s.contexts, cg)
	}
}

// Reads, creates and adds the locations from the
// given config part
func (s *siteCreator) addLocations() {

	// add configured main navigation
	if s.site == nil {
		return
	}
	log.Debug("siteCreator.addLocations()")
	for _, fl := range s.config.Context.MainLinks {

		pth := path.Join(fl.Path, fl.FileName)
		if !strings.HasPrefix(pth, "/") {
			pth = "/" + pth
		}

		url := "https://" + s.config.Domain + pth

		l := staticModel.NewLocation(
			fl.ExternalLink,
			s.config.Domain,
			fl.Label,
			"",
			fl.Path,
			fl.FileName,
			"",
			pth,
			url)
		s.site.AddMain(l)
	}

	// add configured marginal navigation
	for _, fl := range s.config.Context.MarginalLinks {

		pth := path.Join(fl.Path, fl.FileName)
		if !strings.HasPrefix(pth, "/") {
			pth = "/" + pth
		}
		url := "https://" + s.config.Domain + pth
		l := staticModel.NewLocation(
			fl.ExternalLink,
			s.config.Domain,
			fl.Label,
			"",
			fl.Path,
			fl.FileName,
			"",
			pth,
			url)
		s.site.AddMarginal(l)
	}
}

// Intended for migrational purposes
func (s *siteCreator) rewritePages() {
	for _, src := range s.config.Src {
		dtos := staticPersistence.ReadPagesFromDir(src.Dir)
		dirname := src.Dir + "/" + "migrated"
		staticPersistence.WritePagesToDir(dtos, dirname)
	}
}

// Reads the list of sources from the config and creates
// source structs from them.
func (s *siteCreator) addSources() {

	log.Debugf("siteCreator.addSources(), amount: %d\n", len(s.config.Src))
	for _, srcCfg := range s.config.Src {
		src := NewSource(
			srcCfg.Type,
			srcCfg.Dir,
			srcCfg.SubDir,
			srcCfg.Headline,
			s.site,
			s.config)
		s.sources = append(s.sources, src)
	}
}

// Generates and stores the containers generated
// from the sources
func (s *siteCreator) addContainers() {
	if s.site != nil {
		for _, src := range s.sources {
			src.generate()
			s.site.AddContainer(src.Container())
		}
		log.Debugf("siteCreator.addContainers(), nr of added containers: %d\n", len(s.site.Containers()))
	}
}

// Generates various render contexts from and for the sources
func (s *siteCreator) addContexts() {
	log.Debug("siteCreator.addContexts()")
	for _, src := range s.sources {
		s.addContext(src.CreateContext())
	}
}

// Checks if a context already exists, to
// avoid redundancy and double output
func (s *siteCreator) contextExists(cg staticIntf.Context) bool {
	contextName := reflect.TypeOf(cg).Elem().Name()
	for _, ctx := range s.contexts {
		if contextName == reflect.TypeOf(ctx).Elem().Name() {
			return true
		}
	}
	return false
}

// Fills the file containers with the data to
// be written
func (s *siteCreator) fillFileContainers(config staticPersistence.JsonConfig) {
	collector := NewComponentCollector()
	for _, ctx := range s.contexts {
		cmps := ctx.GetComponents()
		collector.AddComponents(cmps)
		fcs := ctx.RenderPages()
		s.fileContainers = append(s.fileContainers, fcs...)
	}

	css := ""
	for _, cmp := range collector.GetComponents() {
		css += cmp.GetCss()
	}

	cssFc := fs.NewFileContainer()
	cssFc.SetDataAsString(css)
	cssFc.SetPath(config.Deploy.TargetDir)
	cssFc.SetFilename(config.Deploy.CssFileName)
	s.fileContainers = append(s.fileContainers, cssFc)
}

// Actually writes the files of the website to
// the local file system
func (s *siteCreator) writeFiles() {
	msg := fmt.Sprintf("Number of files to write: %d", len(s.fileContainers))
	log.Debug(msg)
	for _, f := range s.fileContainers {
		log.Debug("Writing file: " + f.GetPath() + "/" + f.GetFilename())
		//log.Debug(f.GetDataAsString())
		f.Write()
	}
}
