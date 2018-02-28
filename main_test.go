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

func TestReadNewConfig(t *testing.T) {
	configs := readNewConfig()

	actual := configs[0].Domain
	expected := "drewing.de"

	if expected != actual {
		t.Errorf("Expected %s but got %s\n", expected, actual)
	}

	actual = configs[1].Domain
	expected = "devabo.de"

	if expected != actual {
		t.Errorf("Expected %s but got %s\n", expected, actual)
	}
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
	fadd, fmake, fstrato, fclear = true, true, true, true

	addJsonCalled := false
	genSiteCalled := false
	uploadCalled := false
	clearCalled := false

	addJsonFn := func() { addJsonCalled = true }
	genSiteFn := func() { genSiteCalled = true }
	uploadFn := func() { uploadCalled = true }
	clearFn := func() { clearCalled = true }

	if !(addJsonCalled && genSiteCalled && uploadCalled && clearCalled) {
		t.Error("Expected no function to be executed, but one was")
	}

	checkFlags(addJsonFn, uploadFn, clearFn, genSiteFn)

	if !(addJsonCalled && genSiteCalled && uploadCalled && clearCalled) {
		t.Error("Expected all functions to be executed, but they weren't")
	}
}

func getExpectedPageHtml() string {
	return `<!doctype html><html itemscope lang="en"><head><meta content="width=device-width, initial-scale=1.0" name="viewport"/><meta content="index,follow" name="robots"/><meta content="Ingmar Drewing" name="author"/><meta content="Ingmar Drewing" name="publisher"/><meta content="storytelling, illustration, drawing, web comic, comic, cartoon, caricatures" name="keywords"/><meta content="storytelling, illustration, drawing, web comic, comic, cartoon, caricatures" name="DC.subject"/><meta content="art" name="page-topic"/><meta charset="UTF-8"/><link href="/icons/favicon-16x16.png" rel="icon" sizes="16x16" type="image/png"/><link href="/icons/favicon-32x32.png" rel="icon" sizes="32x32" type="image/png"/><link href="/icons/android-192x192.png" rel="icon" sizes="192x192" type="image/png"/><link href="/icons/apple-touch-icon-180x180.png" rel="apple-touch-icon" sizes="180x180" type="image/png"/><meta content="/icons/browserconfig.xml" name="msapplication-config"/><meta content="A Little Test" itemprop="name"/><meta content="A blog containing texts, drawings, graphic novels and code snippets by Ingmar Drewing." itemprop="description"/><meta content="https://drewingde.s3.us-west-1.amazonaws.com/blog/2018/1/13/ALittleTest-w800.png" itemprop="image"/><meta content="summary_large_image" name="t:card"/><meta content="@ingmardrewing" name="t:site"/><meta content="A Little Test" name="t:title"/><meta content="A blog containing texts, drawings, graphic novels and code snippets by Ingmar Drewing." name="t:text:description"/><meta content="@ingmardrewing" name="t:creator"/><meta content="https://drewingde.s3.us-west-1.amazonaws.com/blog/2018/1/13/ALittleTest-w800.png" name="t:image"/><meta content="A Little Test" property="og:title"/><meta content="2018/1/13/a-little-test/index.html" property="og:url"/><meta content="https://drewingde.s3.us-west-1.amazonaws.com/blog/2018/1/13/ALittleTest-w800.png" property="og:image"/><meta content="A blog containing texts, drawings, graphic novels and code snippets by Ingmar Drewing." property="og:description"/><meta content="" property="og:site_name"/><meta content="article" property="og:type"/><meta content="2018-1-7 1:50:42" property="article:published_time"/><meta content="2018-1-7 1:50:42" property="article:modified_time"/><meta content="Illustration" property="article:section"/><meta content="comic, graphic novel, webcomic, science-fiction, sci-fi" property="article:tag"/><title>A Little Test</title></head><body><div class="wrapperOuter "><div class="wrapperInner"><main class="maincontent"><h1 class="maincontent__h1">A Little Test</h1><h2 class="maincontent__h2">2018-1-7 1:50:42</h2><p><a href="https://drewingde.s3.us-west-1.amazonaws.com/blog/2018/1/13/ALittleTest.png"><img src="https://drewingde.s3.us-west-1.amazonaws.com/blog/2018/1/13/ALittleTest-w800.png" alt=""/></a></p></main></div></div><div class="wrapperOuter "><div class="wrapperInner"><div class="disqus" id="disqus_thread"> </div></div></div><script>var disqus_config = function () { this.page.title= "A Little Test"; this.page.url = 'https://drewing.de2018/1/13/a-little-test/index.html'; this.page.identifier =  'drewing.de 2018/1/13/A Little Test'; }; (function() { var d = document, s = d.createElement('script'); s.src = 'https://drewing.disqus.com/embed.js'; s.setAttribute('data-timestamp', +new Date()); (d.head || d.body).appendChild(s); })();</script><div class="wrapperOuter headerbar__wrapper"><div class="wrapperInner"><header class="headerbar"><div class="headerbar__logocontainer"><a class="headerbar__logo" href="https://drewing.de"><!-- logo --></a></div></header></div></div><div class="wrapperOuter mainnavi__wrapper"><div class="wrapperInner"><div class=""><nav class="mainnavi"><a class="mainnavi__navelement" href="testResources/deploy/blog/index.html">Blog</a><a class="mainnavi__navelement" href="https://www.facebook.com/drewing.de">Facebook</a><a class="mainnavi__navelement" href="https://twitter.com/ingmardrewing">Twitter</a></nav></div></div></div><div class="wrapperOuter "><div class="wrapperInner"><div class="copyright"><a rel="license" class="copyright__cc" href="https://creativecommons.org/licenses/by-nc-nd/3.0/"></a><p class="copyright__license">© 2017 by Ingmar Drewing </p><p class="copyright__license">Except where otherwise noted, content on this site is licensed under a <a rel="license" href="https://creativecommons.org/licenses/by-nc-nd/3.0/">Creative Commons Attribution-NonCommercial-NoDerivs 3.0 Unported (CC BY-NC-ND 3.0) license</a>.</p><p class="copyright__license">Soweit nicht anders explizit ausgewiesen, stehen die Inhalte auf dieser Website unter der <a rel="license" href="https://creativecommons.org/licenses/by-nc-nd/3.0/">Creative Commons Namensnennung-NichtKommerziell-KeineBearbeitung (CC BY-NC-ND 3.0)</a> Lizenz. Unless otherwise noted the author of the content on this page is <a href="https://plus.google.com/113943655600557711368?rel=author">Ingmar Drewing</a></p></div></div></div><div class="wrapperOuter footernavi__wrapper"><div class="wrapperInner"><div class=""><nav class="footernavi"><a class="footernavi__navelement" href="https://www.facebook.com/sharer.php?u=https%3A%2F%2Fwww.drewing.de%2Fblog&amp;t=drewing.de">Share on Facebook</a><a class="footernavi__navelement" href="mailto:blank?subject=I%20found%20this%20on%20drewing.de&amp;body=I%20thought%20you%20might%20be%20interested%20in%20this%20blogpost:%20https://www.drewing.de/blog/ingmars-booklist/">Tell a friend</a><a class="footernavi__navelement" href="/blog/about.html">Me</a><a class="footernavi__navelement" href="/blog/booklist.html">booklist</a><a class="footernavi__navelement" href="/blog/imprint.html">imprint</a></nav></div></div></div></body></html>`
}
