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

	s.logo = s.gh.tex.AddText("how to play", 0.20, 0.8, 0.6, 0.05, 0.1, EffectMetaballs, SpriteMenu)
	s.back = s.gh.tex.AddText("back", 0.385, 0.1, 0.6, 0.05, 0.05, EffectNone, SpriteMenu)
	s.text = append(s.text, s.gh.tex.AddText("Each slot can have 2 values: 1 OR 0", 0.05, 0.7, 0.6, 0.02, 0.02, EffectNone, SpriteMenu)...)
	s.text = append(s.text, s.gh.tex.AddText("Each row can have values between 1-15", 0.05, 0.65, 0.6, 0.02, 0.02, EffectNone, SpriteMenu)...)

	// TBD: Rewrite this a bit more generic.
	t1 := &Sprite{}
	t1.Init(0.05, 0.45, 0.6, 0.9, 0.15, "tile", g, SpriteMenu)
	t1.Hide()
	s.text = append(s.text, t1)

	t2 := &Sprite{}
	t2.Init(0.08, 0.32, 0.6, 0.4, 0.05, "tile", g, SpriteMenu)
	t2.Hide()
	s.text = append(s.text, t2)
	s.text = append(s.text, s.gh.tex.AddText("0", 0.12, 0.33, 0.6, 0.025, 0.025, EffectNone, SpriteMenu)...)
	s.text = append(s.text, s.gh.tex.AddText("0", 0.22, 0.33, 0.6, 0.025, 0.025, EffectNone, SpriteMenu)...)
	s.text = append(s.text, s.gh.tex.AddText("0", 0.32, 0.33, 0.6, 0.025, 0.025, EffectNone, SpriteMenu)...)
	s.text = append(s.text, s.gh.tex.AddText("1", 0.42, 0.33, 0.6, 0.025, 0.025, EffectNone, SpriteMenu)...)
	s.text = append(s.text, s.gh.tex.AddText("=", 0.50, 0.33, 0.6, 0.025, 0.025, EffectNone, SpriteMenu)...)
	s.text = append(s.text, s.gh.tex.AddText("1", 0.55, 0.33, 0.6, 0.025, 0.025, EffectNone, SpriteMenu)...)

	t3 := &Sprite{}
	t3.Init(0.08, 0.26, 0.6, 0.4, 0.05, "tile", g, SpriteMenu)
	t3.Hide()
	s.text = append(s.text, t3)
	s.text = append(s.text, s.gh.tex.AddText("0", 0.12, 0.27, 0.6, 0.025, 0.025, EffectNone, SpriteMenu)...)
	s.text = append(s.text, s.gh.tex.AddText("1", 0.22, 0.27, 0.6, 0.025, 0.025, EffectNone, SpriteMenu)...)
	s.text = append(s.text, s.gh.tex.AddText("0", 0.32, 0.27, 0.6, 0.025, 0.025, EffectNone, SpriteMenu)...)
	s.text = append(s.text, s.gh.tex.AddText("0", 0.42, 0.27, 0.6, 0.025, 0.025, EffectNone, SpriteMenu)...)
	s.text = append(s.text, s.gh.tex.AddText("=", 0.50, 0.27, 0.6, 0.025, 0.025, EffectNone, SpriteMenu)...)
	s.text = append(s.text, s.gh.tex.AddText("4", 0.55, 0.27, 0.6, 0.025, 0.025, EffectNone, SpriteMenu)...)

	t4 := &Sprite{}
	t4.Init(0.08, 0.20, 0.6, 0.4, 0.05, "tile", g, SpriteMenu)
	t4.Hide()
	s.text = append(s.text, t4)
	s.text = append(s.text, s.gh.tex.AddText("1", 0.12, 0.21, 0.6, 0.025, 0.025, EffectNone, SpriteMenu)...)
	s.text = append(s.text, s.gh.tex.AddText("1", 0.22, 0.21, 0.6, 0.025, 0.025, EffectNone, SpriteMenu)...)
	s.text = append(s.text, s.gh.tex.AddText("0", 0.32, 0.21, 0.6, 0.025, 0.025, EffectNone, SpriteMenu)...)
	s.text = append(s.text, s.gh.tex.AddText("0", 0.42, 0.21, 0.6, 0.025, 0.025, EffectNone, SpriteMenu)...)
	s.text = append(s.text, s.gh.tex.AddText("=", 0.50, 0.21, 0.6, 0.025, 0.025, EffectNone, SpriteMenu)...)
	s.text = append(s.text, s.gh.tex.AddText("12 [8+4]", 0.55, 0.21, 0.6, 0.025, 0.025, EffectNone, SpriteMenu)...)

	s.text = append(s.text, s.gh.tex.AddText("Value 8", 0.07, 0.52, 0.6, 0.02, 0.02, EffectStats, SpriteMenu)...)
	s.text = append(s.text, s.gh.tex.AddText("Value 4", 0.30, 0.52, 0.6, 0.02, 0.02, EffectStats, SpriteMenu)...)
	s.text = append(s.text, s.gh.tex.AddText("Value 2", 0.52, 0.52, 0.6, 0.02, 0.02, EffectStats, SpriteMenu)...)
	s.text = append(s.text, s.gh.tex.AddText("Value 1", 0.76, 0.52, 0.6, 0.02, 0.02, EffectStats, SpriteMenu)...)
	s.text = append(s.text, s.gh.tex.AddText("Examples:", 0.05, 0.40, 0.6, 0.02, 0.02, EffectNone, SpriteMenu)...)

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
