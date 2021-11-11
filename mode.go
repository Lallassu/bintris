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
	hidden       bool
}

func (m *Mode) Init(g *Game) {
	m.gh = g
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
	}
	for i := range m.currTimeTxt {
		m.currTimeTxt[i].Hide()
	}
	for i := range m.currScoreTxt {
		m.currScoreTxt[i].Hide()
	}
}

func (m *Mode) Show() {
	if !m.hidden {
		return
	}
	m.hidden = false
	for i := range m.score {
		m.score[i].Show()
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
	if m.gh == nil {
		return
	}

	switch m.Type {
	case GameModeEasy:
	case GameModeNormal:
		m.Speed = 0.1 + time.Since(m.Time).Seconds()/100
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
		// if len(hidden) == 15-9 {
		// 	m.Reset()
		// }
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
