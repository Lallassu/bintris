package main

import (
	_ "image/png"
	"strings"

	"golang.org/x/mobile/gl"
)

type Sprite struct {
	gh      *Game
	id      int
	uEffect gl.Uniform
	x       float32
	y       float32
	z       float32
	scalex  float32
	scaley  float32
	effect  Effect
	hidden  bool
	prevX   float32
	prevY   float32
	Texture Texture
	dirty   bool
}

func (s *Sprite) Init(x, y, z, scalex, scaley float32, tex string, g *Game) {
	s.gh = g
	tex = strings.ToLower(tex)
	//	s.uEffect = s.gh.glc.GetUniformLocation(s.gh.program, "effect")
	s.scalex = scalex
	s.scaley = scaley
	s.x = x
	s.y = y
	s.z = z
	s.id = s.gh.NewID()

	s.Texture = s.gh.tex.Types[tex]
	s.gh.tex.AddSprite(s)
	s.dirty = true
}

func (s *Sprite) GetObjectType() ObjectType {
	return ObjectTypeSprite
}

func (s *Sprite) Update(dt float64) {
}

func (s *Sprite) Hide() {
	// Offset far away to "hide" it. Since we are already including
	// it in the vertice list.
	s.x += 10000
	s.hidden = true
	s.dirty = true
}

func (s *Sprite) Show() {
	s.x -= 10000
	s.dirty = true
}

func (s *Sprite) ChangeY(dy float32) {
	s.y += dy
	s.dirty = true
}

func (s *Sprite) Hidden() bool {
	return s.hidden
}

func (s *Sprite) Draw(dt float64) {
	if s.hidden {
		return
	}
}

func (s *Sprite) GetID() int {
	return s.id
}

func (s *Sprite) GetX() float32 {
	return s.x
}

func (s *Sprite) GetY() float32 {
	return s.y
}

func (s *Sprite) Delete() {
	s.hidden = true
}

func (s *Sprite) Resize() {
	//a := float32(s.gh.size.WidthPx) / float32(s.gh.size.HeightPx)
	sx := float32(s.gh.size.WidthPx) / float32(s.gh.sizePrev.WidthPx)
	sy := float32(s.gh.size.HeightPx) / float32(s.gh.sizePrev.HeightPx)

	s.scalex *= sx // float32(s.gh.size.WidthPx) / float32(s.gh.sizePrev.WidthPx)
	s.scaley *= sy // float32(s.gh.size.HeightPx) / float32(s.gh.sizePrev.HeightPx)
	s.x *= sx      //float32(s.gh.size.WidthPx) / float32(s.gh.sizePrev.WidthPx)
	s.y *= sy      //float32(s.gh.size.HeightPx) / float32(s.gh.sizePrev.HeightPx)
	s.dirty = true
}
