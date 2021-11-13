package main

import (
	"math/rand"
	"strconv"
	"time"
)

type GameMode int

const (
	GameModeEasy GameMode = iota
	GameModeNormal
	GameModeHard
)

// Mode is handling game modes
type Mode struct {
	Type         GameMode
	gameOver     bool
	Speed        float64
	Time         time.Time
	Score        int
	gh           *Game
	lastTile     time.Time
	secsPrev     string
	score        []*Sprite
	time         []*Sprite
	currTimeTxt  []*Sprite
	currScoreTxt []*Sprite
	gameOver1    []*Sprite
	gameOver2    []*Sprite
	hidden       bool
}

func (m *Mode) Init(g *Game) {
	m.gh = g
	m.gameOver1 = m.gh.tex.AddText("GAME", 0.30, 0.75, 0.6, 0.1, 0.19, EffectGameOver1)
	m.gameOver2 = m.gh.tex.AddText("OVER", 0.30, 0.55, 0.6, 0.1, 0.19, EffectGameOver2)

	m.score = g.tex.AddText("Score:", 0.05, 0.91, 0.1, 0.04, 0.055, EffectNone)
	m.currScoreTxt = g.tex.AddText("0000", 0.33, 0.91, 0.1, 0.04, 0.055, EffectStats)
	m.time = g.tex.AddText("Time:", 0.54, 0.91, 0.1, 0.04, 0.055, EffectNone)
	m.currTimeTxt = g.tex.AddText("0000", 0.77, 0.91, 0.1, 0.04, 0.055, EffectStats)
	m.Hide()
}

func (m *Mode) Start(gm GameMode) {
	m.Type = gm
	m.Time = time.Now()
	m.Speed = 0.1
	m.Score = 0
	m.gh.menu.Hide()
	m.Show()
	m.gameOver = false
}

func (m *Mode) Hide() {
	if m.hidden {
		return
	}
	m.hidden = true
	for i := range m.score {
		m.score[i].Hide()
	}

	for i := range m.time {
		m.time[i].Hide()
		m.time[i].ChangeEffect(EffectNone)
	}
	for i := range m.currTimeTxt {
		m.currTimeTxt[i].Hide()
		m.currTimeTxt[i].ChangeEffect(EffectStats)
	}
	for i := range m.currScoreTxt {
		m.currScoreTxt[i].Hide()
		m.currScoreTxt[i].ChangeEffect(EffectStats)
	}

	// Same len
	for i := range m.gameOver1 {
		m.gameOver1[i].Hide()
		m.gameOver2[i].Hide()
	}
}

func (m *Mode) GameOver() {
	m.gameOver = true
	for i := range m.gameOver1 {
		m.gameOver1[i].Show()
		m.gameOver2[i].Show()
	}

	for i := range m.time {
		m.time[i].ChangeEffect(EffectGameOver)
	}
	for i := range m.currTimeTxt {
		m.currTimeTxt[i].ChangeEffect(EffectGameOver)
	}
	for i := range m.currScoreTxt {
		m.currScoreTxt[i].ChangeEffect(EffectGameOver)
	}
	for i := range m.score {
		m.score[i].ChangeEffect(EffectGameOver)
	}
	m.gh.GameOver()

	// TBD: Print score + time

	go func() {
		time.Sleep(5 * time.Second)
		m.gh.Reset()
	}()
}

func (m *Mode) Show() {
	if !m.hidden {
		return
	}
	m.hidden = false
	for i := range m.score {
		m.score[i].Show()
		m.score[i].ChangeEffect(EffectNone)
	}

	for i := range m.time {
		m.time[i].Show()
	}
	for i := range m.currTimeTxt {
		m.currTimeTxt[i].Show()
	}
	for i := range m.currScoreTxt {
		m.currScoreTxt[i].Show()
	}

}

func (m *Mode) Reset() {
	for i := range m.gh.tiles {
		m.gh.tiles[i].Hide()
	}
	m.Time = time.Now()
	m.Speed = 0.1
	m.Score = 0
	m.lastTile = time.Now()
	m.Show()
}

func (m *Mode) Update(dt float64) {
	if m.gh == nil || m.gameOver {
		return
	}

	switch m.Type {
	case GameModeEasy:
	case GameModeNormal:
		m.Speed = 0.9 + time.Since(m.Time).Seconds()/100
		c := 0
		hidden := []*TileSet{}
		for i := range m.gh.tiles {
			if !m.gh.tiles[i].hidden {
				c++
			} else {
				hidden = append(hidden, &m.gh.tiles[i])
			}
		}
		if len(hidden) > 0 && time.Since(m.lastTile).Seconds() > 2 || len(hidden) == 15 {
			t := hidden[rand.Intn(len(hidden))]
			t.Reset(1)
			t.SetSpeed(m.Speed)
			m.lastTile = time.Now()
		}
		if len(hidden) == 13 { //15-9 {
			m.GameOver()
		}
	case GameModeHard:
	}

	// Update timer
	secs := strconv.Itoa(int(time.Since(m.Time).Seconds()))
	// Only update if changed.
	if secs != m.secsPrev {
		for i := 0; i < len(secs); i++ {
			if i < 4 {
				m.currTimeTxt[i].ChangeTexture(string(secs[i]))
			}
		}
	}
	m.secsPrev = secs
}

func (m *Mode) AddScore(points int) {
	m.Score += points
	// TBD: Update score text
	txt := strconv.Itoa(m.Score)
	for i := 0; i < len(txt); i++ {
		if i < 4 {
			m.currScoreTxt[i].ChangeTexture(string(txt[i]))
		}
	}
}
