package main

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/ingmardrewing/fs"
	"github.com/ingmardrewing/staticPersistence"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	conf = staticPersistence.NewConfig("testResources/config.json")
}

func tearDown() {
	filepath := path.Join(getTestFileDirPath(), conf.Read("deploy", "localDir"))
	fs.RemoveDirContents(filepath)
}

func getTestFileDirPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func TestConfRead(t *testing.T) {
	expected := "styles.css"
	actual := conf.Read("deploy", "cssFileName")

	if expected != actual {
		t.Errorf("Expected %s but got %s\n", expected, actual)
	}
}

func TestGenSite(t *testing.T) {
	generateSiteLocally()

	deployDir := path.Join(getTestFileDirPath(),
		conf.Read("deploy", "localDir"))

	cssPath := path.Join(deployDir, "styles.css")
	cssFileExists, _ := fs.PathExists(cssPath)

	if !cssFileExists {
		t.Error("No css file found at:", cssPath)
	}

	indexPath := path.Join(deployDir, "blog", "index.html")
	indexExists, _ := fs.PathExists(indexPath)

	if !indexExists {
		t.Error("No css file found at:", indexPath)
	}

	index0Path := path.Join(deployDir, "blog", "index0.html")
	index0Exists, _ := fs.PathExists(index0Path)

	if !index0Exists {
		t.Error("No css file found at:", index0Path)
	}

	postPath := path.Join(deployDir, "blog", "test", "index.html")
	postExists, _ := fs.PathExists(postPath)

	if !postExists {
		t.Error("No css file found at:", postExists)
	}

	tearDown()
}

func TestGeneratePages(t *testing.T) {
	expected := "styles.css"
	actual := conf.Read("deploy", "cssFileName")

	if expected != actual {
		t.Errorf("Expected %s but got %s\n", expected, actual)
	}
}

func TestInferBlogTitleFromFilename(t *testing.T) {
	title, titlePlain := inferBlogTitleFromFilename("ATest29,This.png")

	titleExpected := "A Test 29, This"
	if title != titleExpected {
		t.Errorf("Expected %s but got %s\n", titleExpected, title)
	}

	titlePlainExpected := "a-test-29-this"
	if titlePlain != titlePlainExpected {
		t.Errorf("Expected %s but got %s\n", titlePlainExpected, titlePlain)
	}
}

func TestInferBlogTitle(t *testing.T) {
	title := inferBlogTitle("ATest29,This.png")
	expected := "A Test 29, This"

	if title != expected {
		t.Errorf("Expected %s but got %s\n", expected, title)
	}

	title = inferBlogTitle("aTest")
	expected = "A Test"

	if title != expected {
		t.Errorf("Expected %s but got %s\n", expected, title)
	}
}

func TestInferBlogTitlePlain(t *testing.T) {
	title := inferBlogTitlePlain("ATest29äüöß,This.png")
	expected := "a-test-29-this"
	if title != expected {
		t.Errorf("Expected %s but got %s\n", expected, title)
	}
}

func TestFlagsPresent(t *testing.T) {
	actual := flagsPresent()
	expected := false

	if actual != expected {
		t.Errorf("Expected %t but got %t\n", expected, actual)
	}

	fmake = true
	actual = flagsPresent()
	expected = true

	if actual != expected {
		t.Errorf("Expected %t but got %t\n", expected, actual)
	}
}

func TestCheckFlags(t *testing.T) {
	fimg, fjson, fmake, fstrato, fclear = true, true, true, true, true

	addImgCalled := false
	addJsonCalled := false
	genSiteCalled := false
	uploadCalled := false
	clearCalled := false

	addJsonFn := func() { addJsonCalled = true }
	genSiteFn := func() { genSiteCalled = true }
	uploadFn := func() { uploadCalled = true }
	clearFn := func() { clearCalled = true }

	if !(addImgCalled && addJsonCalled && genSiteCalled && uploadCalled && clearCalled) {
		t.Error("Expected no function to be executed, but one was")
	}

	checkFlags(addJsonFn, uploadFn, clearFn, genSiteFn)

	if !(addImgCalled && addJsonCalled && genSiteCalled && uploadCalled && clearCalled) {
		t.Error("Expected all functions to be executed, but they weren't")
	}
}
