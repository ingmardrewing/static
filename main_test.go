package main

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/ingmardrewing/actions"
	"github.com/ingmardrewing/fs"
	"github.com/ingmardrewing/staticPersistence"
	log "github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestMainInternals(t *testing.T) {
	fi = false
	checkFlagsCalled := false
	checkFlags = func() { checkFlagsCalled = true }

	main()

	if !checkFlagsCalled {
		t.Error("expected checkFlags to be called, but it wasn't.")
	}

	fi = true
	startedInteractive := false
	interactive = func() { startedInteractive = true }

	main()

	if !startedInteractive {
		t.Error("Interactive session not started")
	}
}

func setup() {
	conf = staticPersistence.ReadConfig("testResources/", "configNew.json")
	log.SetLevel(log.DebugLevel)
}

func tearDown() {
	filepath := path.Join(getTestFileDirPath(), conf[0].Deploy.TargetDir)
	fs.RemoveDirContents(filepath)
}

func getTestFileDirPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func TestConfigureActions(t *testing.T) {
	a := configureActions()
	as := a.Actions()

	actionFunction := findActionByName("exit", as).GetFunction()
	fn := generateSiteLocally
	if &actionFunction == &fn {
		t.Error("Expected action using exit, but it doesn't.")
	}

	actionFunction = findActionByName("make", as).GetFunction()
	fn = generateSiteLocally
	if &actionFunction == &fn {
		t.Error("Expected action using generateSiteLocally, but it doesn't.")
	}

	actionFunction = findActionByName("upload", as).GetFunction()
	fn = upload
	if &actionFunction == &fn {
		t.Error("Expected action using upload, but it doesn't.")
	}

	actionFunction = findActionByName("clear", as).GetFunction()
	fn = clear
	if &actionFunction == &fn {
		t.Error("Expected action using clear, but it doesn't.")
	}
}

func findActionByName(aname string,
	actions []actions.Action) actions.Action {
	for _, a := range actions {
		if a.GetName() == aname {
			return a
		}
	}
	return nil
}

func TestConfRead(t *testing.T) {
	expected := "styles.css"
	actual := conf[0].Deploy.CssFileName

	if expected != actual {
		t.Errorf("Expected %s but got %s\n", expected, actual)
	}

}

func TestGenSite(t *testing.T) {
	generateSiteLocally()

	deployDir := path.Join(getTestFileDirPath(),
		conf[0].Deploy.TargetDir)

	cssPath := path.Join(deployDir, "styles.css")
	cssFileExists, _ := fs.PathExists(cssPath)

	if !cssFileExists {
		t.Error("No css file found at:", cssPath)
	}

	indexPath := path.Join(deployDir, "blog", "index.html")
	indexExists, _ := fs.PathExists(indexPath)

	if !indexExists {
		t.Error("No index.html file found at:", indexPath)
	}

	index0Path := path.Join(deployDir, "blog", "index0.html")
	index0Exists, _ := fs.PathExists(index0Path)

	if !index0Exists {
		t.Error("No index0.html file found at:", index0Path)
	}

	tearDown()
}

func TestGeneratePages(t *testing.T) {
	expected := "styles.css"
	actual := conf[0].Deploy.CssFileName

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

func TestCheckFlags(t *testing.T) {
	fmake, fstrato, fclear = true, true, true
	b, c, d := false, false, false

	generateSiteLocally = func() { b = true }
	upload = func() { c = true }
	clear = func() { d = true }

	checkFlagsFn()

	if !(b && c && d) {
		t.Error("checkFlags did not trigger all expected functions.")
	}
}
