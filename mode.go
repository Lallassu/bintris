package main

import (
	"math/rand"
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
	Type     GameMode
	Speed    float64
	Time     time.Time
	Score    int
	gh       *Game
	lastTile time.Time
}

func (m *Mode) Start(gm GameMode, gh *Game) {
	m.gh = gh
	m.Type = gm
	m.Time = time.Now()
	m.Speed = 0.1
	m.Score = 0
}

func (m *Mode) Reset() {
	for i := range m.gh.tiles {
		m.gh.tiles[i].Hide()
	}
	m.Time = time.Now()
	m.Speed = 0.1
	m.Score = 0
	m.lastTile = time.Now()
}

func (m *Mode) Update(dt float64) {
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
}

func (m *Mode) AddScore(points int) {
	m.Score += points
	// TBD: Update score text
}
