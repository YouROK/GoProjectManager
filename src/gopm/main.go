package main

import (
	"core/build"
	"core/project"

	"log"
)

func main() {
	p, err := project.OpenProject("testprj/test.gop")
	if err != nil {
		log.Println(err)
		return
	}
	err = build.Build(p)
	if err != nil {
		log.Println(err)
	}
}
