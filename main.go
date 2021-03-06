package main

import (
	"flag"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ingmardrewing/actions"
	"github.com/ingmardrewing/fs"
	"github.com/ingmardrewing/staticPersistence"
	log "github.com/sirupsen/logrus"
)

var (
	fi          = false
	fimg        = false
	fadd        = false
	fmake       = false
	fupdatejson = false
	fstrato     = false
	fclear      = false
	fconfigPath = ""
	conf        []staticPersistence.Config
	configFile  = "configNew.json"
	debug       = false

	generateSiteLocally = generateSiteLocallyFn
	upload              = uploadFn
	clear               = clearFn
	configureActions    = configureActionsFn
	checkFlags          = checkFlagsFn
	interactive         = interactiveFn
	exit                = func() { os.Exit(0) }
)

func init() {
	flag.BoolVar(&fi, "i", false, "Interactive mode")
	flag.BoolVar(&debug, "debug", false, "Run in debug mode")
	flag.BoolVar(&fmake, "make", false, "Generate local site")
	flag.BoolVar(&fupdatejson, "updatejson", false, "Updates to new json format")
	flag.BoolVar(&fstrato, "strato", false, "Upload site to strato")
	flag.BoolVar(&fclear, "clear", false, "Automatically publish the image in BLOG_DEFAULT_DIR and clear the dir afterwards")
	flag.StringVar(&fconfigPath, "configPath", os.Getenv("BLOG_CONFIG_DIR"), "path to config file")
	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
	}

	log.Debug("config dir:", fconfigPath)
	log.Debug("config file:", fconfigPath)

	exists, _ := fs.PathExists(path.Join(fconfigPath, configFile))
	if exists {
		conf = staticPersistence.ReadConfig(fconfigPath, configFile)
	} else {
		conf = staticPersistence.ReadConfig("./testResources/", configFile)
	}
}

func main() {
	if fi {
		interactive()
	} else {
		checkFlags()
	}
}

func interactiveFn() {
	a := configureActions()
	for {
		a.AskUser()
	}
}

func checkFlagsFn() {
	if fmake {
		generateSiteLocally()
	}
	if fupdatejson {
		updateJsonFiles()
	}
	if fstrato {
		upload()
	}
	if fclear {
		clear()
	}
}

func generateSiteLocallyFn() {
	log.Debug("main:generateSiteLocallyFn")
	log.Debug(conf)
	sc := NewSitesController(conf)
	sc.UpdateStaticSites()
}

func updateJsonFiles() {
	log.Debug("main:updateJsonFiles")
	sc := NewSitesController(conf)
	sc.UpdateJsonFiles()
}

func configureActionsFn() actions.Choice {
	c := actions.NewChoice()
	c.AddAction(
		"exit",
		"Exits the Application",
		exit)
	c.AddAction(
		"make",
		"Generate website locally",
		generateSiteLocally)
	c.AddAction(
		"upload",
		"Upload generated html, css and js to strato (www.drewing.de)",
		upload)
	c.AddAction(
		"clear",
		"clear auto blog dir",
		clear)
	return c
}

func inferBlogTitleFromFilename(filename string) (string, string) {
	fname := strings.TrimSuffix(filename, filepath.Ext(filename))
	return inferBlogTitle(fname), inferBlogTitlePlain(fname)
}

func inferBlogTitle(filename string) string {
	rx := regexp.MustCompile("(^[a-zäüöß]+)|([A-ZÄÜÖ][a-zäüöß,]*)|([0-9,]+)")
	parts := rx.FindAllString(filename, -1)
	spaceSeparated := strings.Join(parts, " ")
	return strings.Title(spaceSeparated)
}

func inferBlogTitlePlain(filename string) string {
	rx := regexp.MustCompile("(^[a-z]+)|([A-Z][a-z]*)|([0-9]+)")
	parts := rx.FindAllString(filename, -1)
	dashSeparated := strings.Join(parts, "-")
	return strings.ToLower(dashSeparated)
}

func clearFn() {
	c := newCommand("cleardir.pl")
	c.run()
}

func askUserForTitle() (string, string) {
	i := NewInput("Enter a title:")
	i.AskUser()
	return i.Regular(), i.Sanitized()
}

func uploadFn() {
	c := newCommand("blogUpload.pl")
	c.run()
}
