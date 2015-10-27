package main

import (
	//	"core/build"
	"core/codecomplete"
	"core/project"

	"log"
)

func main1() {
	p, err := project.OpenProject("testprj/test.gop")
	if err != nil {
		log.Println(err)
		return
	}
	codecomplete.Init(p)

	file := p.GetFile("test.go")
	if file != nil {
		x := 4
		y := 33
		off := file.GetOffset(x, y)
		log.Println(x, y, off)
		log.Println(codecomplete.Complete(p, file, off))
	}

	/*err = build.GetDepends(p, false)
	if err != nil {
		log.Println(err)
	}
	err = build.Build(p)
	if err != nil {
		log.Println(err)
	}*/
}
