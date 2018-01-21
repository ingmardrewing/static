package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ingmardrewing/aws"
	"github.com/ingmardrewing/fs"
	"github.com/ingmardrewing/img"
	"github.com/ingmardrewing/staticPersistence"

	"gopkg.in/russross/blackfriday.v2"
)

/* Date Strings */

type DateStrings struct{}

func (d *DateStrings) generateDatePath() string {
	now := time.Now()
	return fmt.Sprintf("%d/%d/%d/", now.Year(), now.Month(), now.Day())
}

func (d *DateStrings) getDate() string {
	n := time.Now()
	return fmt.Sprintf("%d-%d-%d %d:%d:%d", n.Year(), n.Month(), n.Day(), n.Hour(), n.Minute(), n.Second())
}

/* Image Manager */

func NewImageManager(awsbucket, sourceimagepath string) *ImageManager {
	im := new(ImageManager)
	im.awsbucket = awsbucket
	im.sourceimagepath = sourceimagepath
	return im
}

type ImageManager struct {
	DateStrings
	sourceimagepath   string
	uploadimgagepaths []string
	awsimageurls      []string
	imagesizes        []int
	awsbucket         string
}

func (i *ImageManager) prepareImages() {
	imgdir := fs.GetPathWithoutFilename(i.sourceimagepath)
	img := img.NewImg(i.sourceimagepath, imgdir)
	paths := img.PrepareResizeTo(i.imagesizes...)
	i.uploadimgagepaths = append(paths, i.sourceimagepath)
	img.Resize()
}

func (i *ImageManager) uploadImages() {
	for _, filepath := range i.uploadimgagepaths {
		filename := fs.GetFilenameFromPath(filepath)
		key := i.getS3Key(filename)
		url := aws.UploadFile(filepath, i.awsbucket, key)
		i.awsimageurls = append(i.awsimageurls, url)
	}
}

func (i *ImageManager) getS3Key(filename string) string {
	return "blog/" + i.generateDatePath() + filename
}

func (i *ImageManager) GetImageUrls() []string {
	return i.awsimageurls
}

func (i *ImageManager) AddImageSize(size int) {
	i.imagesizes = append(i.imagesizes, size)
}

/* page json factory */

func NewPageJsonFactory(originalmd string, awsbucket, blogUrl,
	sourceimagepath, markdownfilepath string,
	imgs ...string) *pageJsonFactory {
	if !strings.HasSuffix(blogUrl, "/") {
		blogUrl += "/"
	}
	p := new(pageJsonFactory)
	p.awsbucket = awsbucket
	p.blogUrl = blogUrl
	p.sourceimagepath = sourceimagepath
	p.markdownfilepath = markdownfilepath
	p.originalmd = originalmd
	if len(imgs) > 1 {
		p.thumburl = imgs[0]
		p.mediumurl = imgs[1]
	}
	return p
}

type pageJsonFactory struct {
	DateStrings
	awsbucket         string
	blogUrl           string
	sourceimagepath   string
	markdownfilepath  string
	uploadimgagepaths []string
	awsimageurls      []string
	imagesizes        []int
	thumburl          string
	mediumurl         string
	originalmd        string
}

func (p *pageJsonFactory) GetDto(domain, title, titlePlain, jsondir string) (staticPersistence.DTO, string) {
	url := p.generateBlogUrl(titlePlain)
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	filename, idstr := p.getPostJsonFilename(jsondir)
	id, _ := strconv.Atoi(idstr)
	md := string(fs.ReadByteArrayFromFile(p.markdownfilepath))
	content := p.generateContentFromMarkdown(md)
	excerpt := generateExzerpt(p.originalmd)
	disqusId := domain + " " + p.generateDatePath() + title
	date := p.getDate()

	d := staticPersistence.NewDto()
	d.ThumbUrl(p.thumburl)
	d.ImageUrl(p.mediumurl)
	d.Filename(filename)
	d.Id(id)
	d.CreateDate(date)
	d.Url(url)
	d.Title(title)
	d.TitlePlain(titlePlain)
	d.Description(excerpt)
	d.Content(content)
	d.DisqusId(disqusId)

	return d, filename
}

func (p *pageJsonFactory) getPostJsonFilename(dir string) (string, string) {
	dirs := fs.ReadDirEntries(dir, false)
	sort.Strings(dirs)
	lastFile := dirs[len(dirs)-1]
	rx := regexp.MustCompile("(\\d+)")
	m := rx.FindStringSubmatch(lastFile)
	i, _ := strconv.Atoi(m[1])
	i++
	return fmt.Sprintf("page%d.json", i), fmt.Sprintf("%d", 10000+i)
}

func (p *pageJsonFactory) generateBlogUrl(title string) string {
	return p.blogUrl + p.generateDatePath() + title
}

func (p *pageJsonFactory) generateContentFromMarkdown(input string) string {
	bytes := []byte(input)
	htmlBytes := blackfriday.Run(bytes, blackfriday.WithNoExtensions())
	htmlString := string(htmlBytes)
	trimmed := strings.TrimSuffix(htmlString, "\n")
	escaped := strings.Replace(trimmed, `"`, `\"`, -1)
	return strings.Replace(escaped, "\n", " ", -1)

}

func generateExzerpt(text string) string {
	if len(text) > 155 {
		return fmt.Sprintf("%.155s ...", text)
	} else if len(text) == 0 {
		return conf.Read("excerpt")
	}
	return text
}
