package main

import (
	_ "image/png"
	"math/rand"

	"golang.org/x/mobile/gl"
)

type Sprite struct {
	gh      *Game
	id      int
	uEffect gl.Uniform
	x       float32
	y       float32
	z       float32
	scale   float32
	effect  Effect
	hidden  bool
	prevX   float32
	prevY   float32
	Texture Texture
	dirty   bool
}

func (s *Sprite) Init(x, y, z, scale float32, tex string, g *Game) {
	s.gh = g
	//	s.uEffect = s.gh.glc.GetUniformLocation(s.gh.program, "effect")
	s.scale = scale
	s.x = x
	s.y = y
	s.z = z
	s.id = s.gh.NewID()

	s.Texture = s.gh.tex.Types[tex]
	s.gh.tex.AddSprite(s)
	s.dirty = false
}

func (s *Sprite) GetObjectType() ObjectType {
	return ObjectTypeSprite
}

func (s *Sprite) Update(dt float64) {
	s.x += rand.Float32() * 10
	s.y += rand.Float32() * 10
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

func (s *Sprite) Delete() {
	s.hidden = true
}
