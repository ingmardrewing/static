package main

import (
	"os"
	"testing"

	"github.com/ingmardrewing/fs"
	"github.com/ingmardrewing/staticPersistence"
)

func TestNewAddJson(t *testing.T) {
	envName := "TEST_AWS_BUCKET"
	bucketName := "testBucketName"
	srcDir := "testResources/add"
	destDir := "testResources/deploy"
	excerpt := "Test 1, 2"
	url := "https://drewing.de/blog"
	os.Setenv(envName, bucketName)

	aj := NewAddJson(envName, srcDir, destDir, excerpt, url)
	if bucketName != aj.awsBucket {
		t.Error("Expected", aj.awsBucket, "to be", bucketName)
	}
}

func TestWriteToFs(t *testing.T) {
	dto := staticPersistence.NewFilledDto(42,
		"titleValue",
		"titlePlainValue",
		"thumbUrlValue",
		"imageUrlValue",
		"descriptionValue",
		"disqusIdValue",
		"createDateValue",
		"contentValue",
		"urlValue",
		"domainValue",
		"pathValue",
		"fspathValue",
		"htmlfilenameValue",
		"thumbBase64Value",
		"categoryValue")

	envName := "TEST_AWS_BUCKET"
	bucketName := "testBucketName"
	srcDir := "testResources/add"
	destDir := "testResources/deploy"
	excerpt := "Test 1, 2"
	url := "https://drewing.de/blog"
	os.Setenv(envName, bucketName)

	aj := NewAddJson(envName, srcDir, destDir, excerpt, url)
	aj.dto = dto
	aj.WriteToFs()

	ba := fs.ReadByteArrayFromFile("testResources/deploy/page42.json")

	actual := len(ba)
	expected := 345

	if actual != expected {
		t.Error("Expected byte array to be of length", expected, "but it was", actual)
	}

}
