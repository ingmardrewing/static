package main

import "testing"

func TestInputRegular(t *testing.T) {
	i := new(input)
	i.userInput = "Hello World"

	expected := "Hello World"
	actual := i.Regular()

	if actual != expected {
		t.Error("Expected", expected, ", but got", actual)
	}
}

func TestInputSanitized(t *testing.T) {
	i := NewInput("")
	i.userInput = "Hello World,42!"

	expected := "hello-world-42"
	actual := i.Sanitized()

	if actual != expected {
		t.Error("Expected", expected, ", but got", actual)
	}
}
