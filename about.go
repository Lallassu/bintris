package main

type About struct {
	hidden bool
	gh     *Game
	logo   []*Sprite
	text   []*Sprite
	back   []*Sprite
}

func (s *About) Init(g *Game) {
	s.gh = g

	s.logo = s.gh.tex.AddText("about", 0.37, 0.8, 0.6, 0.05, 0.1, EffectMetaballs)
	s.back = s.gh.tex.AddText("back", 0.385, 0.1, 0.6, 0.05, 0.05, EffectNone)
	s.text = append(s.text, s.gh.tex.AddText(" Author:", 0.05, 0.7, 0.6, 0.02, 0.02, EffectStats)...)
	s.text = append(s.text, s.gh.tex.AddText("Magnus Persson", 0.25, 0.7, 0.6, 0.02, 0.02, EffectNone)...)
	s.text = append(s.text, s.gh.tex.AddText("Version:", 0.05, 0.65, 0.6, 0.02, 0.02, EffectStats)...)
	s.text = append(s.text, s.gh.tex.AddText(version, 0.25, 0.65, 0.6, 0.02, 0.02, EffectNone)...)
	s.text = append(s.text, s.gh.tex.AddText(" Source:", 0.05, 0.60, 0.6, 0.02, 0.02, EffectStats)...)
	s.text = append(s.text, s.gh.tex.AddText("github.com/lallassu/bintris", 0.25, 0.60, 0.6, 0.02, 0.02, EffectNone)...)
	s.text = append(s.text, s.gh.tex.AddText("Credits:", 0.05, 0.55, 0.6, 0.02, 0.02, EffectStats)...)
	s.text = append(s.text, s.gh.tex.AddText("https://thebookofshaders.com", 0.1, 0.50, 0.6, 0.017, 0.017, EffectNone)...)
	s.text = append(s.text, s.gh.tex.AddText("it is written in go!", 0.25, 0.20, 0.6, 0.02, 0.02, EffectStatsBlink)...)

	s.Hide()
}

func (s *About) Show() {
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

	for i := range s.text {
		s.text[i].Show()
	}
}

func (s *About) Hide() {
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
	for i := range s.text {
		s.text[i].Hide()
	}
}

func (s *About) KeyDown(x, y float32) {
	if x > 0.385 && x < 0.385+(float32(len(s.back))*0.055) && y > 0.1 && y < 0.1+0.05 {
		s.Hide()
		s.gh.menu.Show()
	}
}

func (s *About) Hidden() bool {
	return s.hidden
}
