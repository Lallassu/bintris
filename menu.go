package main

import "fmt"

type Menu struct {
	hidden     bool
	gh         *Game
	logo       []*Sprite
	start      []*Sprite
	scoreBoard []*Sprite
	about      []*Sprite
	quit       []*Sprite
	how        []*Sprite
}

func (m *Menu) Init(g *Game) {
	m.gh = g
	// Create menu sprites and position them.

	// Handle keyboard input

	offsetX := float32(0.02)
	m.logo = m.gh.tex.AddText("bintris", 0.14, 0.75, 0.6, 0.1, 0.19, EffectMetaballsBlue)
	m.start = m.gh.tex.AddText("start", 0.34+offsetX, 0.65, 0.6, 0.05, 0.05, EffectMetaballsBlue)
	m.how = m.gh.tex.AddText("How to play", 0.16+offsetX, 0.58, 0.6, 0.05, 0.05, EffectMetaballsBlue)
	m.scoreBoard = m.gh.tex.AddText("Scoreboard", 0.20+offsetX, 0.51, 0.6, 0.05, 0.05, EffectMetaballsBlue)
	m.about = m.gh.tex.AddText("about", 0.34+offsetX, 0.44, 0.6, 0.05, 0.05, EffectMetaballsBlue)
	m.quit = m.gh.tex.AddText("quit", 0.37+offsetX, 0.37, 0.6, 0.05, 0.05, EffectMetaballsBlue)

}
func (m *Menu) Show() {
	m.hidden = false

	for i := range m.logo {
		m.logo[i].Show()
	}

	for i := range m.how {
		m.how[i].Show()
	}

	for i := range m.start {
		m.start[i].Show()
	}

	for i := range m.scoreBoard {
		m.scoreBoard[i].Show()
	}

	for i := range m.about {
		m.about[i].Show()
	}

	for i := range m.quit {
		m.quit[i].Show()
	}
}

func (m *Menu) Hide() {
	m.hidden = true

	for i := range m.logo {
		m.logo[i].Hide()
	}

	for i := range m.how {
		m.how[i].Hide()
	}

	for i := range m.start {
		m.start[i].Hide()
	}

	for i := range m.scoreBoard {
		m.scoreBoard[i].Hide()
	}

	for i := range m.about {
		m.about[i].Hide()
	}

	for i := range m.quit {
		m.quit[i].Hide()
	}
}

func (m *Menu) KeyDown() {
	fmt.Printf("KEYDOWN!\n")
	m.Hide()
	m.gh.Show()
	m.gh.mode.Start(GameModeNormal, m.gh)
}
