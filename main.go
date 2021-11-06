package main

import (
	"fmt"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
)

func main() {
	app.Main(func(a app.App) {
		game := Game{}
		var sz size.Event
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glc, _ := e.DrawContext.(gl.Context)
					if glc == nil {
						continue
					}
					game.Init(glc)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					game.Stop()
				}
			case size.Event:
				sz = e
			case paint.Event:
				if game.glc == nil || e.External {
					continue
				}
				game.Draw()
				a.Publish()
				a.Send(paint.Event{})
			case touch.Event:
				game.Click(sz, e.X, e.Y)
			case key.Event:
				fmt.Printf("KEY!\n")
				if e.Code != key.CodeSpacebar {
					break
				}
				if down := e.Direction == key.DirPress; down || e.Direction == key.DirRelease {
					game.menu.KeyDown()
				}
			}
		}
	})
}
