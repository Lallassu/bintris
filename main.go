package main

import (
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
)

func main() {
	game := Game{}

	app.Main(func(a app.App) {
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glc, _ := e.DrawContext.(gl.Context)
					game.Init(glc)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					game.Stop()
				}
			case size.Event:
				game.Resize(e)
			case paint.Event:
				if game.glc == nil || e.External {
					continue
				}
				game.Draw()
				a.Publish()
				a.Send(paint.Event{})
			case touch.Event:
				game.Click(e.X, e.Y)
			}
		}
	})
}
