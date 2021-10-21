package main

import (
	"fmt"
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
	// s.x += 1 - rand.Float32()
	// s.y += 1 - rand.Float32()
	// s.dirty = true
}

func (s *Sprite) Hidden() bool {
	return s.hidden
}

func (s *Sprite) Draw(dt float64) {
	if s.hidden {
		return
	}
}

func (s *Sprite) Delete() {
	s.hidden = true
}

func (s *Sprite) Resize() {
	if s.scalex != 1.0 {
		s.scalex = float32(s.gh.size.WidthPx) / s.Texture.Width
		s.scaley = float32(s.gh.size.HeightPx) / s.Texture.Height
		fmt.Printf("SCALE_1: %vx%v\n", s.scalex, s.scaley)
		fmt.Printf("SCALE_2: %vx%v\n", 1/s.scalex, 1/s.scaley)
		s.dirty = true
	}
}
