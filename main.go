package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ingmardrewing/actions"
	"github.com/ingmardrewing/staticBlogAdd"
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
	fmt.Println("Enter a title:")
	reader := bufio.NewReader(os.Stdin)

	title, _ := reader.ReadString('\n')
	title = strings.TrimSuffix(title, "\n")

	whitespace := regexp.MustCompile("\\s+")
	preptitle := whitespace.ReplaceAllString(strings.ToLower(title), "-")
	r := regexp.MustCompile("[^-a-zA-Z0-9]+")
	title_plain := r.ReplaceAllString(preptitle, "")
	return title, title_plain
}

func addJsonFileFn() {
	fmt.Println("addJsonFile")
	bucket := os.Getenv("AWS_BUCKET")
	addDir := conf[0].AddPostDir
	postsDir := conf[0].Src.PostsDir
	defaultExcerpt := conf[0].DefaultMeta.BlogExcerpt

	bda := staticBlogAdd.NewBlogDataAbstractor(bucket, addDir, postsDir, defaultExcerpt, "https://drewing.de/blog/")
	dto := bda.GeneratePostDto()

	filename := fmt.Sprintf("page%d.json", dto.Id())

	fmt.Println("Writing ...", dto, postsDir, filename)
	staticPersistence.WritePageDtoToJson(dto, postsDir, filename)
}

func exit() { os.Exit(0) }

func uploadFn() {
	fmt.Println("Uploading content to strato .. may take a while")
	c := newCommand("blogUpload.pl")
	c.run()
}
