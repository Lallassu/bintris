package main

import (
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
	bg         *Sprite
}

func (m *Menu) Init(g *Game) {
	m.gh = g

	offsetX := float32(0.02)
	m.logo = m.gh.tex.AddText("bintris", 0.14, 0.75, 0.6, 0.1, 0.19, EffectMetaballs)
	m.start = m.gh.tex.AddText("start", 0.34+offsetX, 0.65, 0.6, 0.05, 0.05, EffectNone)
	m.how = m.gh.tex.AddText("Rules", 0.34+offsetX, 0.58, 0.6, 0.05, 0.05, EffectNone)
	m.scoreBoard = m.gh.tex.AddText("Score", 0.34+offsetX, 0.51, 0.6, 0.05, 0.05, EffectNone)
	m.about = m.gh.tex.AddText("about", 0.34+offsetX, 0.44, 0.6, 0.05, 0.05, EffectNone)
	m.quit = m.gh.tex.AddText("quit", 0.365+offsetX, 0.27, 0.6, 0.05, 0.05, EffectNone)

	for i := range m.logo {
		m.logo[i].ChangeEffect(EffectMetaballs)
	}
	m.hidden = false
}
func (m *Menu) Show() {
	if !m.hidden {
		return
	}
	m.gh.backBg.ChangeEffect(EffectMenu)
	m.gh.bg.Hide()
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
	if x > 0.36 && x < 0.36+(float32(len(m.start))*0.055) && y > 0.65 && y < 0.65+0.05 {
		m.gh.mode.Start(GameModeNormal)
	} else if x > 0.36 && x < 0.36+(float32(len(m.how))*0.055) && y > 0.58 && y < 0.58+0.05 {
		m.gh.menu.Hide()
		m.gh.howto.Show()
		m.gh.backBg.ChangeEffect(EffectMenuStill)
	} else if x > 0.36 && x < 0.36+(float32(len(m.scoreBoard))*0.055) && y > 0.51 && y < 0.51+0.05 {
		m.gh.scoreboard.Show([]Score{})
		m.gh.backBg.ChangeEffect(EffectMenuStill)
		m.gh.menu.Hide()
	} else if x > 0.36 && x < 0.36+(float32(len(m.about))*0.055) && y > 0.44 && y < 0.44+0.05 {
		m.gh.menu.Hide()
		m.gh.backBg.ChangeEffect(EffectMenuStill)
		m.gh.about.Show()
	} else if x > 0.39 && x < 0.39+(float32(len(m.quit))*0.055) && y > 0.27 && y < 0.27+0.05 {
		os.Exit(0)
	}
}

func (m *Menu) Hidden() bool {
	return m.hidden
}
