package main

import (
	"log"
	"os/exec"
	"strings"
)

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
	command := exec.Command(c.name, c.arguments...)

	err := command.Run()
	if err != nil {
		log.Println(c.name, strings.Join(c.arguments, " "))
		log.Fatalln(err)
	}
}
