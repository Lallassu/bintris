package main

import (
	"fmt"
	"os"
)

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

	offsetX := float32(0.02)
	m.logo = m.gh.tex.AddText("bintris", 0.14, 0.75, 0.6, 0.1, 0.19, EffectMetaballsBlue)
	m.start = m.gh.tex.AddText("start", 0.34+offsetX, 0.65, 0.6, 0.05, 0.05, EffectMetaballsBlue)
	m.how = m.gh.tex.AddText("How to play", 0.16+offsetX, 0.58, 0.6, 0.05, 0.05, EffectMetaballsBlue)
	m.scoreBoard = m.gh.tex.AddText("Scoreboard", 0.20+offsetX, 0.51, 0.6, 0.05, 0.05, EffectMetaballsBlue)
	m.about = m.gh.tex.AddText("about", 0.34+offsetX, 0.44, 0.6, 0.05, 0.05, EffectMetaballsBlue)
	m.quit = m.gh.tex.AddText("quit", 0.37+offsetX, 0.37, 0.6, 0.05, 0.05, EffectMetaballsBlue)

	for i := range m.logo {
		m.logo[i].ChangeEffect(EffectMetaballs)
	}
}
func (m *Menu) Show() {
	if !m.hidden {
		return
	}
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
	if m.hidden {
		return
	}
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

func (m *Menu) KeyDown(x, y float32) {
	if x > 0.36 && x < 0.36+(float32(len(m.start))*0.05) && y > 0.65 && y < 0.65+0.05 {
		m.Hide()
		m.gh.Show()
		m.gh.mode.Start(GameModeNormal, m.gh)
	} else if x > 0.18 && x < 0.18+(float32(len(m.how))*0.05) && y > 0.58 && y < 0.58+0.05 {
		fmt.Printf("HOW\n")
	} else if x > 0.22 && x < 0.22+(float32(len(m.scoreBoard))*0.05) && y > 0.51 && y < 0.51+0.05 {
		fmt.Printf("Score\n")
	} else if x > 0.36 && x < 0.36+(float32(len(m.about))*0.05) && y > 0.44 && y < 0.44+0.05 {
		fmt.Printf("Menu: %v\n", m.about[0].Texture)
		m.about[0].ChangeTexture("b")
	} else if x > 0.39 && x < 0.39+(float32(len(m.quit))*0.05) && y > 0.37 && y < 0.37+0.05 {
		os.Exit(0)
	}
}