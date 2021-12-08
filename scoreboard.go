package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const fileName = "/score"
const noOfRows = 10
const noOfHeaders = 4

type Scoreboard struct {
	hidden bool
	gh     *Game
	logo   []*Sprite
	header []*Sprite
	back   []*Sprite
	rows   [noOfRows][noOfHeaders][]*Sprite
	dir    string
}

type Score struct {
	Pos      string
	Score    int
	Time     int
	Date     string
	NewScore bool
}

func (s *Scoreboard) Init(g *Game) {
	s.gh = g

	// Grab the path to cache dir to find our "files" dir.
	// Note: workaround as we don't have getFilesDir()

	s.dir = os.Getenv("TMPDIR")
	s.dir = "/tmp/" // TBD: REMOVE

	s.dir = strings.ReplaceAll(s.dir, "cache", "files")

	s.logo = s.gh.tex.AddText("scoreboard", 0.25, 0.8, 0.6, 0.05, 0.1, EffectMetaballs)
	s.back = s.gh.tex.AddText("back", 0.385, 0.1, 0.6, 0.05, 0.05, EffectNone)
	s.header = append(s.header, s.gh.tex.AddText("#", 0.05, 0.7, 0.6, 0.03, 0.03, EffectStats)...)
	s.header = append(s.header, s.gh.tex.AddText("Score", 0.15, 0.7, 0.6, 0.03, 0.03, EffectStats)...)
	s.header = append(s.header, s.gh.tex.AddText("Time", 0.35, 0.7, 0.6, 0.03, 0.03, EffectStats)...)
	s.header = append(s.header, s.gh.tex.AddText("Date", 0.55, 0.7, 0.6, 0.03, 0.03, EffectStats)...)

	// TBD: Preload top 10 score positions
	hOffset := float32(0.05)
	sizeX := float32(0.025)
	sizeY := float32(0.035)
	for i := 0; i < noOfRows; i++ {
		s.rows[i][0] = s.gh.tex.AddText("  ", 0.05, 0.7-float32(i+1)*hOffset, 0.6, sizeX, sizeY, EffectNone)
		s.rows[i][1] = s.gh.tex.AddText("      ", 0.15, 0.7-float32(i+1)*hOffset, 0.6, sizeX, sizeY, EffectNone)
		s.rows[i][2] = s.gh.tex.AddText("      ", 0.35, 0.7-float32(i+1)*hOffset, 0.6, sizeX, sizeY, EffectNone)
		s.rows[i][3] = s.gh.tex.AddText("                   ", 0.55, 0.7-float32(i+1)*hOffset, 0.6, sizeX, sizeY, EffectNone)
	}

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
			// We can't really handle errors anyway, so just skip.
			sl, _ := strconv.Atoi(line[1])
			tl, _ := strconv.Atoi(line[2])
			score = append(score, Score{
				Pos:   line[0],
				Date:  line[3],
				Score: sl,
				Time:  tl,
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

	for i, sc := range score {
		if i == noOfRows {
			break
		}
		f.Write([]byte(fmt.Sprintf("%d,%d,%d,%s\n", i+1, sc.Score, sc.Time, sc.Date)))
	}

	f.Close()
}

func (s *Scoreboard) Show(sc []Score) {
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

	for i := 0; i < noOfRows; i++ {
		for n := range s.rows[i] {
			for m := range s.rows[i][n] {
				s.rows[i][n][m].Show()
			}
		}
	}

	if len(sc) == 0 {
		sc = s.LoadFile()
	}

	// Draw score rows
	for i := range sc {
		if i == noOfRows {
			break
		}

		e := Effect(EffectNone)
		if sc[i].NewScore {
			e = EffectStatsBlink
		}

		s.ChangeTexture(fmt.Sprintf("%d", i+1), &s.rows[i][0], e)
		s.ChangeTexture(strconv.Itoa(sc[i].Score), &s.rows[i][1], e)
		s.ChangeTexture(strconv.Itoa(sc[i].Time), &s.rows[i][2], e)
		s.ChangeTexture(sc[i].Date, &s.rows[i][3], e)
	}
}

func (s *Scoreboard) ChangeTexture(val string, sc *[]*Sprite, e Effect) {
	for i := 0; i < len(val); i++ {
		if i < len(*sc) {
			(*sc)[i].ChangeTexture(string(val[i]))
			(*sc)[i].ChangeEffect(e)
		}
	}
}

func (s *Scoreboard) Hide() {
	if s.hidden {
		return
	}
	s.hidden = true

	for i := 0; i < noOfRows; i++ {
		for n := range s.rows[i] {
			for m := range s.rows[i][n] {
				s.rows[i][n][m].Hide()
			}
		}
	}

	for i := range s.logo {
		s.logo[i].Hide()
	}
	for i := range s.back {
		s.back[i].Hide()
	}
	for i := range s.header {
		s.header[i].Hide()
	}
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

func (s *Scoreboard) Add(score, sec int) {
	sc := s.LoadFile()

	sc = append(sc, Score{
		Score:    score,
		Time:     sec,
		Date:     time.Now().Format("20060102 15:04"),
		NewScore: true,
	})

	sort.Slice(sc, func(m, n int) bool {
		return sc[m].Score > sc[n].Score
	})

	s.WriteFile(sc)
	s.Show(sc)
}
