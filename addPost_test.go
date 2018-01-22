package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestGenerateExzerpt_truncates_empty_texts(t *testing.T) {
	emptyTxt := ""
	actual := generateExzerpt(emptyTxt)
	expected := "A blog containing texts, drawings, graphic narratives/novels and (rarely) code snippets by Ingmar Drewing."
	if actual != expected {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestGenerateExzerpt_truncates_long_texts(t *testing.T) {
	longTxt := "Duis venenatis massa non ex aliquam, sed tempus mi scelerisque. Sed ultricies metus purus, at accumsan lacus venenatis in. Ut a scelerisque justo. Praesent quis erat euismod, dapibus magna non, tristique velit. Maecenas eu ex tristique eros eleifend auctor a non justo. Nulla pulvinar porta ipsum id molestie. Integer a suscipit velit, ac sollicitudin tortor. Aliquam erat volutpat. Nunc elementum ipsum efficitur, egestas augue varius, dignissim dui. Donec tempus eros eget congue vehicula. Vestibulum lobortis elementum magna, non semper felis rutrum at."
	expected := "Duis venenatis massa non ex aliquam, sed tempus mi scelerisque. Sed ultricies metus purus, at accumsan lacus venenatis in. Ut a scelerisque justo. Praesent ..."
	actual := generateExzerpt(longTxt)
	if expected != actual {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestGenerateExzerpt_truncates_short_texts(t *testing.T) {
	txt := "Hello World"
	actual := generateExzerpt(txt)
	expected := txt
	if expected != actual {
		t.Error("Expected", expected, "but got", actual)
	}
}

/*
func TestGetPostJsonFilename(t *testing.T) {
	b := NewPageJsonFactory(
		"", "", "", "")
	actual, _ := b.getPostJsonFilename("/Users/drewing/Desktop/drewing2018/posts/")
	expected := "page338.json"
	if actual != expected {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestAddJson(t *testing.T) {
	b := NewPageJsonFactory(
		"", "https://drewing.de/blog/", "",
		"/Users/drewing/Desktop/drewing2018/add/test.md")
	actual, _:= b.GetJson("drewing.de", "test", "/Users/drewing/Desktop/drewing2018/posts/")
	expected := "wurst"
	if actual != expected {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestUploadImages(t *testing.T) {
	b := NewImageManager(os.Getenv("AWS_BUCKET"),
		"/Users/drewing/Desktop/drewing2018/add/atthezoo.png")
	b.AddImageSize(800)
	b.AddImageSize(390)
	b.treatImages()
	b.uploadImages()
}
*/

func TestS3KeyGenerationFromDate(t *testing.T) {
	b := NewImageManager(os.Getenv("AWS_BUCKET"),
		"/Users/drewing/Desktop/drewing2018/add/atthezoo.png")
	actual := b.generateDatePath()
	now := time.Now()
	expected := fmt.Sprintf("%d/%d/%d/", now.Year(), now.Month(), now.Day())

	if actual != expected {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestGenerateContentFromMarkdown(t *testing.T) {
	input := `test`
	p := NewPageJsonFactory("", "", "https://drewing.de/", "", "")
	actual := p.generateContentFromMarkdown(input)
	expected := "<p>test</p>"

	if actual != expected {
		t.Error("Expected", expected, "but got", ">"+actual+"<")
	}
}

func TestGenerateBlogUrl(t *testing.T) {
	now := time.Now()
	d := "https://drewing.de/blog"
	k := fmt.Sprintf("%d/%d/%d/", now.Year(), now.Month(), now.Day())
	title := "just-a-test/"

	p := NewPageJsonFactory("", "", d, "", "")
	actual := p.generateBlogUrl(title)
	expected := d + "/" + k + title

	if actual != expected {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestStripLinksAndImages(t *testing.T) {

	text := "[weafasdfasdfali](asdfasdfasdf)wurst"
	actual := stripLinksAndImages(text)
	expected := "wurst"

	if actual != expected {
		t.Error("Expected", expected, "but got", actual)
	}
}
