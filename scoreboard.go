package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

const fileName = "/score"

type Scoreboard struct {
	hidden bool
	gh     *Game
	logo   []*Sprite
	header []*Sprite
	back   []*Sprite
	dir    string
}

type Score struct {
	Pos   string
	Score string
	Time  string
	Date  string
}

func (s *Scoreboard) Init(g *Game) {
	s.gh = g

	// Grab the path to cache dir to find our "files" dir.
	// Note: workaround as we don't have getFilesDir()

	// s.dir = os.Getenv("TMPDIR")
	s.dir = "/tmp/" // TBD: REMOVE

	s.dir = strings.ReplaceAll(s.dir, "cache", "files")

	s.logo = s.gh.tex.AddText("scoreboard", 0.25, 0.8, 0.6, 0.05, 0.1, EffectMetaballs)
	s.back = s.gh.tex.AddText("back", 0.385, 0.1, 0.6, 0.05, 0.05, EffectNone)
	s.header = append(s.header, s.gh.tex.AddText("#", 0.1, 0.7, 0.6, 0.03, 0.03, EffectNone)...)
	s.header = append(s.header, s.gh.tex.AddText("Score", 0.2, 0.7, 0.6, 0.03, 0.03, EffectNone)...)
	s.header = append(s.header, s.gh.tex.AddText("Time", 0.5, 0.7, 0.6, 0.03, 0.03, EffectNone)...)
	s.header = append(s.header, s.gh.tex.AddText("Date", 0.75, 0.7, 0.6, 0.03, 0.03, EffectNone)...)

	// TBD: Preload top 10 score positions

	s.Hide()
}

func (s *Scoreboard) LoadFile() []Score {
	content, err := ioutil.ReadFile(s.dir + fileName)
	if err != nil {
		return []Score{}
	}

	score := []Score{}
	rows := strings.Split(string(content), "\n")
	for _, v := range rows {
		line := strings.Split(v, ",")
		if len(line) == 4 {
			score = append(score, Score{
				Pos:   line[0],
				Score: line[1],
				Time:  line[2],
				Date:  line[3],
			})
		}
	}

	return score
}

func (s *Scoreboard) WriteFile(score []Score) {
	f, err := os.Create(s.dir + "/score")
	if err != nil {
		return
	}

	sort.Slice(score, func(m, n int) bool {
		return score[m].Score > score[n].Score
	})

	for i, sc := range score {
		// Only save top 10
		if i == 9 {
			break
		}
		f.Write([]byte(fmt.Sprintf("%d,%s,%s,%s\n", i+1, sc.Score, sc.Time, sc.Date)))
	}

	f.Close()
}

func (s *Scoreboard) Show() {
	if !s.hidden {
		return
	}
	s.hidden = false

	for i := range s.logo {
		s.logo[i].Show()
	}
	for i := range s.back {
		s.back[i].Show()
	}

	for i := range s.header {
		s.header[i].Show()
	}

	// TBD: Load from file and add to list.
	// sc := s.LoadFile()
}

func (s *Scoreboard) Hide() {
	if s.hidden {
		return
	}
	s.hidden = true

	for i := range s.logo {
		s.logo[i].Hide()
	}
	for i := range s.back {
		s.back[i].Hide()
	}
	for i := range s.header {
		s.header[i].Hide()
	}

	// TBD: Hide all entries
}

func (s *Scoreboard) KeyDown(x, y float32) {
	if x > 0.385 && x < 0.385+(float32(len(s.back))*0.055) && y > 0.1 && y < 0.1+0.05 {
		s.Hide()
		s.gh.menu.Show()
	}
}

func (s *Scoreboard) Hidden() bool {
	return s.hidden
}

func (s *Scoreboard) Add(score, sec string) {
	sc := s.LoadFile()

	sc = append(sc, Score{
		Score: score,
		Time:  sec,
		Date:  time.Now().String(),
	})

	s.WriteFile(sc)
}
