package main

import (
	"fmt"
	"strings"
)

func main() {
	l := stringsToList("hello;world")
	fmt.Println(l)
}

func stringsToList(str string) []string {
	l := strings.Split(str, ";")
	if l[len(l)-1] == "" {
		l = l[:len(l)-1]
	}
	return l
}
