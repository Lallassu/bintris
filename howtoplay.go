package main

type HowToPlay struct {
	hidden bool
	gh     *Game
	logo   []*Sprite
	text   []*Sprite
	back   []*Sprite
}

func (s *HowToPlay) Init(g *Game) {
	s.gh = g

	s.logo = s.gh.tex.AddText("how to play", 0.20, 0.8, 0.6, 0.05, 0.1, EffectMetaballs)
	s.back = s.gh.tex.AddText("back", 0.385, 0.1, 0.6, 0.05, 0.05, EffectNone)
	s.text = append(s.text, s.gh.tex.AddText("Binary - It is just about adding numbers.", 0.00, 0.75, 0.6, 0.02, 0.02, EffectStats)...)
	s.text = append(s.text, s.gh.tex.AddText("Each slot can have 2 values: 1 OR 0.", 0.05, 0.7, 0.6, 0.02, 0.02, EffectNone)...)
	s.text = append(s.text, s.gh.tex.AddText("Each row can have values between 1-15.", 0.05, 0.65, 0.6, 0.02, 0.02, EffectNone)...)
	s.text = append(s.text, s.gh.tex.AddText("Left slot is the most significant.", 0.05, 0.60, 0.6, 0.02, 0.02, EffectNone)...)
	s.text = append(s.text, s.gh.tex.AddText("Slot 1: Value 8 - left", 0.05, 0.55, 0.6, 0.02, 0.02, EffectNone)...)
	s.text = append(s.text, s.gh.tex.AddText("Slot 2: Value 4", 0.05, 0.50, 0.6, 0.02, 0.02, EffectNone)...)
	s.text = append(s.text, s.gh.tex.AddText("Slot 3: Value 2", 0.05, 0.45, 0.6, 0.02, 0.02, EffectNone)...)
	s.text = append(s.text, s.gh.tex.AddText("Slot 4: Value 1 - right", 0.05, 0.40, 0.6, 0.02, 0.02, EffectNone)...)
	s.text = append(s.text, s.gh.tex.AddText("Examples:", 0.05, 0.35, 0.6, 0.02, 0.02, EffectNone)...)
	s.text = append(s.text, s.gh.tex.AddText("0 0 0 1 = 1   = 0 + 0 + 0 + 1", 0.1, 0.30, 0.6, 0.025, 0.025, EffectStats)...)
	s.text = append(s.text, s.gh.tex.AddText("1 0 0 0 = 8   = 8 + 0 + 0 + 0", 0.1, 0.27, 0.6, 0.025, 0.025, EffectStats)...)
	s.text = append(s.text, s.gh.tex.AddText("1 1 0 0 = 12  = 8 + 4 + 0 + 0", 0.1, 0.23, 0.6, 0.025, 0.025, EffectStats)...)
	s.text = append(s.text, s.gh.tex.AddText("1 1 1 1 = 15  = 8 + 4 + 2 + 1", 0.1, 0.20, 0.6, 0.025, 0.025, EffectStats)...)

	s.Hide()
}

func (s *HowToPlay) Show() {
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

func (s *HowToPlay) Hide() {
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

func (s *HowToPlay) KeyDown(x, y float32) {
	if x > 0.385 && x < 0.385+(float32(len(s.back))*0.055) && y > 0.1 && y < 0.1+0.05 {
		s.Hide()
		s.gh.menu.Show()
	}
}

func (s *HowToPlay) Hidden() bool {
	return s.hidden
}
