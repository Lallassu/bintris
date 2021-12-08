package main

import (
	"math"
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

var startTimeRelease = time.Second * 3
var maxSpeed = 0.8

// Mode is handling game modes
type Mode struct {
	started      bool
	Type         GameMode
	gameOver     bool
	Speed        float64
	Time         time.Time
	timeSecs     int
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
	timeRelease  time.Duration
}

func (m *Mode) Init(g *Game) {
	m.gameOver = false
	m.gh = g
	m.Time = time.Now()
	m.gameOver1 = m.gh.tex.AddText("GAME", 0.30, 0.55, 0.6, 0.1, 0.19, EffectGameOver1)
	m.gameOver2 = m.gh.tex.AddText("OVER", 0.30, 0.35, 0.6, 0.1, 0.19, EffectGameOver2)

	m.score = g.tex.AddText("Score:", 0.05, 0.91, 0.1, 0.04, 0.055, EffectNone)
	m.currScoreTxt = g.tex.AddText("0     ", 0.33, 0.91, 0.1, 0.04, 0.055, EffectStats)
	m.time = g.tex.AddText("Time:", 0.54, 0.91, 0.1, 0.04, 0.055, EffectNone)
	m.currTimeTxt = g.tex.AddText("0    ", 0.77, 0.91, 0.1, 0.04, 0.055, EffectStats)
	m.Hide()
}

func (m *Mode) Start(gm GameMode) {
	m.started = true
	m.Type = gm
	m.Time = time.Now()
	m.Speed = 0.1
	m.Score = 0
	m.gh.visible = 0
	m.gh.pulse = 2
	m.gh.menu.Hide()
	m.Show()
	m.gameOver = false
	m.timeRelease = startTimeRelease
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
	m.started = false

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

	go func() {
		time.Sleep(3 * time.Second)
		m.gh.scoreboard.Add(m.Score, m.timeSecs)
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

	m.gh.bg.Show()
}

func (m *Mode) Reset() {
	for i := range m.currScoreTxt {
		m.currScoreTxt[i].ChangeTexture(" ")
	}
	m.currScoreTxt[0].ChangeTexture("0")

	for i := range m.currTimeTxt {
		m.currTimeTxt[i].ChangeTexture(" ")
	}
	m.currTimeTxt[0].ChangeTexture("0")

	for i := range m.gh.tiles {
		m.gh.tiles[i].Hide()
	}

	m.timeRelease = startTimeRelease
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
		// Make sure not to pass max speed
		m.Speed = math.Min(maxSpeed, 0.2+time.Since(m.Time).Seconds()/100)
		c := 0
		hidden := []*TileSet{}
		for i := range m.gh.tiles {
			if !m.gh.tiles[i].hidden {
				c++
			} else {
				hidden = append(hidden, &m.gh.tiles[i])
			}
		}
		if len(hidden) > 0 && time.Since(m.lastTile).Seconds() > m.timeRelease.Seconds() || len(hidden) == 15 {
			t := hidden[rand.Intn(len(hidden))]
			t.Reset(1)
			t.SetSpeed(m.Speed)
			m.lastTile = time.Now()
			if m.timeRelease > time.Millisecond*1500 {
				m.timeRelease -= time.Millisecond * 50
			} else {
				m.timeRelease -= time.Millisecond * 10
			}
		}
		m.gh.visible = 14 - float32(len(hidden))
		if len(hidden) == 15-10 {
			m.GameOver()
		}
	case GameModeHard:
	}

	// Update timer
	m.timeSecs = int(time.Since(m.Time).Seconds())
	secs := strconv.Itoa(m.timeSecs)
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
	txt := strconv.Itoa(m.Score)
	for i := 0; i < len(txt); i++ {
		m.currScoreTxt[i].ChangeTexture(string(txt[i]))
	}
}

func (m *Mode) IsGameOver() bool {
	return m.gameOver
}

func (m *Mode) Started() bool {
	return m.started
}
