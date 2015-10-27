// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/jroimartin/gocui"

	"core/codecomplete"
	"core/project"
)

func codeComplete(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()

		x, y := ox+cx, oy+cy

		f := proj.GetFile("test.go")
		off := f.GetOffset(x, y)
		cands := codecomplete.Complete(proj, f, off)

		maxX, maxY := g.Size()
		maxY = maxY/2 - (len(cands) / 2)
		if v, err := g.SetView("msg", maxX/2-30, maxY, maxX/2+30, maxY+len(cands)+2); err != nil {
			if err != gocui.ErrorUnkView {
				return err
			}
			fmt.Fprintln(v, x, y, off)
			for _, c := range cands {
				fmt.Fprintln(v, c.Name, c.Perc)
			}

			if err := g.SetCurrentView("msg"); err != nil {
				return err
			}
		}
	}
	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "side" {
		return g.SetCurrentView("main")
	}
	return g.SetCurrentView("side")
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintln(v, l)
		if err := g.SetCurrentView("msg"); err != nil {
			return err
		}
	}

	if v, err := g.View("main"); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		f := proj.GetFile(l)
		if f != nil {
			v.Clear()
			fmt.Fprintf(v, "%s", f.GetBody())
		}
		v.Editable = true
		v.Wrap = true
	}
	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	if err := g.SetCurrentView("main"); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.Quit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("side", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delMsg); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyCtrlSpace, gocui.ModNone, codeComplete); err != nil {
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("side", -1, -1, 30, maxY); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Highlight = true
		list := proj.GetFileList()
		for _, f := range list {
			fmt.Fprintln(v, filepath.Base(f.GetPath()))
		}
	}
	if v, err := g.SetView("main", 30, -1, maxX, maxY); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		f := proj.GetFile("test.go")
		if f != nil {
			fmt.Fprintf(v, "%s", f.GetBody())
		}
		v.Editable = true
		v.Wrap = false
		if err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	return nil
}

var (
	proj *project.Project
)

func main() {
	var err error

	proj, err = project.OpenProject("testprj/test.gop")
	if err != nil {
		log.Println(err)
		return
	}
	codecomplete.Init(proj)

	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetLayout(layout)
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}
	g.SelBgColor = gocui.ColorGreen
	g.SelFgColor = gocui.ColorBlack
	g.ShowCursor = true

	err = g.MainLoop()
	if err != nil && err != gocui.Quit {
		log.Panicln(err)
	}
}
