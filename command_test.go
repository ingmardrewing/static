package main

import "testing"

func TestNewCommand(t *testing.T) {
	c := newCommand("testCommand", "arg1", "arg2")

	if c.name != "testCommand" {
		t.Error("Expected", c.name, "to be testCommand")
	}

	if c.arguments[0] != "arg1" || c.arguments[1] != "arg2" {
		t.Error("Expected", c.arguments, "to be arg1 and arg2")
	}
}
