package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func NewInput(prompt string) *input {
	i := new(input)
	i.prompt = prompt
	i.whitespaceRegex = regexp.MustCompile("\\s+|[,.!?:;_]+")
	i.fsRegex = regexp.MustCompile("[^-a-zA-Z0-9]+")
	return i
}

type input struct {
	prompt          string
	userInput       string
	whitespaceRegex *regexp.Regexp
	fsRegex         *regexp.Regexp
}

func (i *input) AskUser() {
	fmt.Println(i.prompt)
	reader := bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')
	i.userInput = strings.TrimSpace(userInput)
}

func (i *input) Regular() string {
	return i.userInput
}

func (i *input) Sanitized() string {
	lcInp := strings.ToLower(i.userInput)
	lcInp = i.whitespaceRegex.ReplaceAllString(lcInp, " ")
	lcInp = strings.TrimSpace(lcInp)
	lcInp = i.whitespaceRegex.ReplaceAllString(lcInp, "-")
	return i.fsRegex.ReplaceAllString(lcInp, "")
}
