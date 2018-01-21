package main

import (
	"os"
	"testing"

	"github.com/ingmardrewing/staticPersistence"
)

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

	addImgFn := func(a staticPersistence.PostAdder) { addImgCalled = true }
	addJsonFn := func(a staticPersistence.PostAdder) { addJsonCalled = true }
	genSiteFn := func() { genSiteCalled = true }
	uploadFn := func(a staticPersistence.PostAdder) { uploadCalled = true }
	clearFn := func(a staticPersistence.PostAdder) { clearCalled = true }

	if !(addImgCalled && addJsonCalled && genSiteCalled && uploadCalled && clearCalled) {
		t.Error("Expected no function to be executed, but one was")
	}

	dirpath := os.Getenv("BLOG_DEFAULT_DIR")
	pa := staticPersistence.NewPostAdder(dirpath)
	checkFlags(pa, addImgFn, addJsonFn, uploadFn, clearFn, genSiteFn)

	if !(addImgCalled && addJsonCalled && genSiteCalled && uploadCalled && clearCalled) {
		t.Error("Expected all functions to be executed, but they weren't")
	}
}
