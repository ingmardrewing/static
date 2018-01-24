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
	dirpath := os.Getenv("BLOG_DEFAULT_DIR")
	pa := staticPersistence.NewPostAdder(dirpath)
	pa.Read()
	checkFlags(pa, addJsonFile, strato, clear, generateSiteLocally)
	enterInteractiveMode(pa)

}

func flagsPresent() bool {
	return fimg || fjson || fmake || fstrato || fclear
}

func checkFlags(pa staticPersistence.PostAdder, addJson, upload, clr func(adder staticPersistence.PostAdder), genSite func()) {
	if flagsPresent() {
		if fjson {
			addJson(pa)
		}
		if fmake {
			genSite()
		}
		if fstrato {
			upload(pa)
		}
		if fclear {
			clr(pa)
		}
		os.Exit(0)
	}
}

func generateSiteLocally() {
	sc := staticController.NewSiteController(conf)
	sc.UpdateStaticSite()
}

func enterInteractiveMode(adder staticPersistence.PostAdder) {
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
			addJsonFile(adder)
		})
	c.AddAction(
		"strato",
		"Upload generated html, css and js to strato (www.drewing.de)",
		func() {
			strato(adder)
		})
	c.AddAction(
		"auto",
		"generate images, json and generate local site",
		func() {
			auto(adder)
		})
	c.AddAction(
		"clear",
		"clear auto blog dir",
		func() {
			clear(adder)
		})

	for {
		c.AskUser()
	}
}

func strato(adder staticPersistence.PostAdder) {
	fmt.Println("Uploading content to strato .. may take a while")
	c := newCommand("blogUpload.pl")
	c.run()
}

func auto(adder staticPersistence.PostAdder) {
	if adder.GetImgFileName() == "" {
		log.Println("No image file in default dir. Nothing to do.")
		return
	}
	generateSiteLocally()
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

func clear(adder staticPersistence.PostAdder) {
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

func addJsonFile(adder staticPersistence.PostAdder) {
	bucket := os.Getenv("AWS_BUCKET")
	bdg := staticController.NewBlogDataGenerator(bucket, conf.Read("src", "addDir"), conf.Read("src", "postsDir"), conf.Read("defaultContent", "blogExcerpt"))
	bdg.Generate()
}

/*
func addimage(adder staticPersistence.PostAdder) {
	imgfile := adder.GetImgFileName()
	tmpl := `{"postImg":"%s","thumbImg":"%s","fullImg":"%s"}`
	path := adder.GetImgFilePath()
	bucket := os.Getenv("AWS_BUCKET")

	im := staticController.NewImageManager(bucket, path)
	transformImages(im)
	im.AddImageSize(800)
	im.AddImageSize(390)
	im.PrepareImages()
	im.UploadImages()

	urls := im.GetImageUrls()
	json := fmt.Sprintf(tmpl, urls[0], urls[1], urls[2])

	fc := fs.NewFileContainer()
	fc.SetPath(adder.GetDir())
	fc.SetFilename(imgfile + ".json")
	fc.SetDataAsString(json)
	fc.Write()
}
*/

func transformImages(im *staticController.ImageManager) {
	im.AddImageSize(800)
	im.AddImageSize(390)
	im.PrepareImages()
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
