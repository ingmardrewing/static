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
	"github.com/ingmardrewing/staticController"
	"github.com/ingmardrewing/staticPersistence"
)

// TODO one-stop post adding:
// only image given -> auto complete text fields
// only md given -> error
// image and md given -> autocompose content

var (
	fimg       = false
	fjson      = false
	fmake      = false
	fstrato    = false
	fclear     = false
	conf       *staticPersistence.Config
	configPath = "/Users/drewing/Desktop/drewing2018/config.json"
)

func init() {
	flag.BoolVar(&fimg, "img", false, "Generate json from image")
	flag.BoolVar(&fjson, "json", false, "Generate json")
	flag.BoolVar(&fmake, "make", false, "Generate local site")
	flag.BoolVar(&fstrato, "strato", false, "Upload site to strato")
	flag.BoolVar(&fclear, "clear", false, "Automatically publish the image in BLOG_DEFAULT_DIR and clear the dir afterwards")
	flag.Parse()
	conf = staticPersistence.NewConfig(configPath)
}

func main() {
	checkFlags(addJsonFile, strato, clear, generateSiteLocally)
	enterInteractiveMode()
}

func flagsPresent() bool {
	return fimg || fjson || fmake || fstrato || fclear
}

func checkFlags(addJson, upload, clr, genSite func()) {
	if flagsPresent() {
		if fjson {
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
	sc := staticController.NewSiteController(conf)
	sc.UpdateStaticSite()
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
	bucket := os.Getenv("AWS_BUCKET")
	bdg := staticController.NewBlogDataGenerator(bucket, conf.Read("src", "addDir"), conf.Read("src", "postsDir"), conf.Read("defaultContent", "blogExcerpt"))
	bdg.Generate()
}

/* util */

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
