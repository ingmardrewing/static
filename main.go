package main

import (
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ingmardrewing/actions"
	"github.com/ingmardrewing/staticController"
	"github.com/ingmardrewing/staticPersistence"
)

var (
	fi          = false
	fimg        = false
	fadd        = false
	fmake       = false
	fstrato     = false
	fclear      = false
	fconfigPath = ""
	conf        []staticPersistence.JsonConfig
	configFile  = "configNew.json"

	addJsonFile         = addJsonFileFn
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
	flag.BoolVar(&fadd, "add", false, "Generate json")
	flag.BoolVar(&fmake, "make", false, "Generate local site")
	flag.BoolVar(&fstrato, "strato", false, "Upload site to strato")
	flag.BoolVar(&fclear, "clear", false, "Automatically publish the image in BLOG_DEFAULT_DIR and clear the dir afterwards")
	flag.StringVar(&fconfigPath, "configPath", "./testResources/", "path to config file")
	flag.Parse()
	conf = staticPersistence.ReadConfig(fconfigPath, configFile)
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
	if fadd {
		addJsonFile()
	}
	if fmake {
		generateSiteLocally()
	}
	if fstrato {
		upload()
	}
	if fclear {
		clear()
	}
}

func generateSiteLocallyFn() {
	sc := staticController.NewSitesController(conf)
	sc.UpdateStaticSites()
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
		"json",
		"Add a json blog file",
		addJsonFile)
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

func addJsonFileFn() {
	aj := NewAddJson("AWS_BUCKET", conf[0].AddPostDir, conf[0].WritePostDir, conf[0].DefaultMeta.BlogExcerpt, "https://drewing.de/blog/")
	aj.GenerateDto()
	aj.WriteToFs()
}

func uploadFn() {
	c := newCommand("blogUpload.pl")
	c.run()
}
