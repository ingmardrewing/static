package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ingmardrewing/actions"
	"github.com/ingmardrewing/staticBlogAdd"
	"github.com/ingmardrewing/staticController"
	"github.com/ingmardrewing/staticPersistence"
)

var (
	fimg        = false
	fadd        = false
	fmake       = false
	fstrato     = false
	fclear      = false
	fconfigPath = ""
	conf        []staticPersistence.JsonConfig
	configFile  = "configNew.json"
)

func init() {
	flag.BoolVar(&fadd, "add", false, "Generate json")
	flag.BoolVar(&fmake, "make", false, "Generate local site")
	flag.BoolVar(&fstrato, "strato", false, "Upload site to strato")
	flag.BoolVar(&fclear, "clear", false, "Automatically publish the image in BLOG_DEFAULT_DIR and clear the dir afterwards")
	flag.StringVar(&fconfigPath, "configPath", "./testResources/", "path to config file")
	flag.Parse()
	conf = readConfig()
}

func readConfig() []staticPersistence.JsonConfig {
	return staticPersistence.ReadConfig(fconfigPath, configFile)
}

func main() {
	checkFlags(addJsonFile, strato, clear, generateSiteLocally)
	enterInteractiveMode()
}

func flagsPresent() bool {
	return fadd || fmake || fstrato || fclear
}

func checkFlags(addJson, upload, clr, genSite func()) {
	if flagsPresent() {
		if fadd {
			addJson()
		}
		if fmake {
			genSite()
		}
		if fstrato {
			upload()
		}
		if fclear {
			clr()
		}
		os.Exit(0)
	}
}

func generateSiteLocally() {
	sc := staticController.NewSitesController(conf)
	sc.UpdateStaticSites()
}

func enterInteractiveMode() {
	c := actions.NewChoice()
	c.AddAction(
		"exit",
		"Exits the Application",
		func() { os.Exit(0) })
	c.AddAction(
		"make",
		"Generate website locally",
		func() {
			generateSiteLocally()
		})
	c.AddAction(
		"json",
		"Add a json blog file",
		func() {
			addJsonFile()
		})
	c.AddAction(
		"strato",
		"Upload generated html, css and js to strato (www.drewing.de)",
		func() {
			strato()
		})
	c.AddAction(
		"clear",
		"clear auto blog dir",
		func() {
			clear()
		})

	for {
		c.AskUser()
	}
}

func strato() {
	fmt.Println("Uploading content to strato .. may take a while")
	c := newCommand("blogUpload.pl")
	c.run()
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

func clear() {
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

func addJsonFile() {
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

func newCommand(name string, args ...string) *command {
	c := new(command)
	c.name = name
	c.setArgs(args...)
	return c
}

type command struct {
	name      string
	arguments []string
}

func (c *command) setArgs(args ...string) {
	for _, a := range args {
		c.arguments = append(c.arguments, a)
	}
}

func (c *command) run() {
	err := exec.Command(c.name, c.arguments...).Run()
	if err != nil {
		log.Println(c.name, strings.Join(c.arguments, " "))
		log.Fatalln(err)
	}
}
