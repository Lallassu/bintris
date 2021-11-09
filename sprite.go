package main

import (
	_ "image/png"
	"strings"

	"golang.org/x/mobile/gl"
)

type Sprite struct {
	gh          *Game
	id          int
	uEffect     gl.Uniform
	fx          float32
	fy          float32
	z           float32
	tx          float32
	ty          float32
	effect      Effect
	hidden      bool
	prevX       float32
	prevY       float32
	Texture     Texture
	dirty       bool
	dirtyUvs    bool
	dirtyEffect bool
}

func (s *Sprite) Init(fx, fy, z, tx, ty float32, tex string, g *Game) {
	s.gh = g
	tex = strings.ToLower(tex)
	s.fx = fx
	s.fy = fy
	s.tx = tx
	s.ty = ty
	s.z = z
	s.id = s.gh.NewID()

	s.Texture = s.gh.tex.Types[tex]
	s.gh.tex.AddSprite(s)
	s.dirty = true
}

func (s *Sprite) Hide() {
	if s.hidden {
		return
	}
	// Offset far away to "hide" it. Since we are already including
	// it in the vertice list.
	s.fx += 10
	s.hidden = true
	s.dirty = true
}

func (s *Sprite) Show() {
	if !s.hidden {
		return
	}
	s.fx -= 10
	s.hidden = false
	s.dirty = true
}

func (s *Sprite) ChangeY(dy float32) {
	s.fy += dy
	s.dirty = true
}

func (s *Sprite) Draw(dt float64) {
	if s.hidden {
		return
	}
}

func (s *Sprite) ChangeTexture(tex string) {
	s.Texture = s.gh.tex.Types[tex]
	s.dirtyUvs = true
	s.dirty = true
}

func (s *Sprite) ChangeEffect(e Effect) {
	s.effect = e
	s.dirtyEffect = true
}
