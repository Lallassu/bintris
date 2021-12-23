package main

import (
	"strconv"
)

type TileSet struct {
	id          int
	Size        int // Number of tiles
	Number      int
	NumberStr   string // Textual representation of the bin number.
	Sprites     []*Sprite
	numberSlots []*Sprite
	Speed       float64
	ClickFunc   func(x, y float32)
	sizex       float32
	sizey       float32
	scale       float32
	gh          *Game
	hidden      bool
	tileWidth   float32
	tileHeight  float32
	tile        *Sprite
}

func (t *TileSet) Init(size int, number int, g *Game) {
	t.gh = g
	t.id = g.NewPlayID()
	t.Speed = float64(g.size.HeightPx / 5)
	t.Number = number
	t.Size = size
	t.sizex = t.tileWidth
	t.sizey = t.tileHeight

	t.NumberStr = strconv.Itoa(number)
	c := &Sprite{hidden: false}
	t.Sprites = append(t.Sprites, c)

	c.Init(0.04, 0.1, 0.5, 0.745, 0.1, "tile", g, SpritePlay)
	c.ChangeEffect(EffectTileTop)
	t.tile = c

	// 0.037 is just an offset for the actuall texture that is badly aligned :P
	w := float32(0.732)
	for x := float32(0); x < w; x += w / 4.0 {
		s := g.tex.AddText("0", x+w/4-0.07, 0.125, 0.7, 0.05, 0.05, EffectNumber, SpritePlay)
		t.numberSlots = append(t.numberSlots, s...)
		s = g.tex.AddText("1", x+w/4-0.07, 0.125, 0.7, 0.05, 0.05, EffectNumber, SpritePlay)
		s[0].Hide()
		t.numberSlots = append(t.numberSlots, s...)
	}

	n := g.tex.AddText(t.NumberStr, 0.84, 0.121, 0.7, 0.06, 0.06, EffectNumber, SpritePlay)
	t.Sprites = append(t.Sprites, n...)
}

func (t *TileSet) SetSpeed(s float64) {
	t.Speed = s
}

func (t *TileSet) Click(x, y float32) {
	t.gh.sound.Play("click")
	slot0 := 0
	slot1 := 0
	w := float32(0.732)
	if x > t.tile.fx && x < t.tile.fx+w/4 {
		slot0 = 0
		slot1 = 1
	} else if x > t.tile.fx+w/4 && x < t.tile.fx+(w/4)*2 {
		slot0 = 2
		slot1 = 3
	} else if x > t.tile.fx+(w/4)*2 && x < t.tile.fx+(w/4)*3 {
		slot0 = 4
		slot1 = 5
	} else if x > t.tile.fx+(w/4)*3 && x < t.tile.fx+(w/4)*4 {
		slot0 = 6
		slot1 = 7
	}
	c0 := t.numberSlots[slot0]
	c1 := t.numberSlots[slot1]
	if c0.hidden {
		c0.Show()
		c1.Hide()
	} else {
		c0.Hide()
		c1.Show()
	}

	// Verify number
	t.VerifyNumber()
}

func (t *TileSet) Hide() {
	for i := range t.Sprites {
		t.Sprites[i].Hide()
	}
	for i := range t.numberSlots {
		t.numberSlots[i].Hide()
	}
	t.hidden = true
}

func (t *TileSet) Reset(offset float32) {
	for i := range t.Sprites {
		t.Sprites[i].Show()
		t.Sprites[i].ChangeY(offset)
		t.Sprites[i].ChangeEffect(EffectNumber)
	}
	for i := range t.numberSlots {
		if t.numberSlots[i].Texture.Name != "1" {
			t.numberSlots[i].Show()
		}
		t.numberSlots[i].ChangeY(offset)
		t.numberSlots[i].ChangeEffect(EffectNumber)
	}
	t.hidden = false
}

func (t *TileSet) VerifyNumber() {
	num := 0
	for i := 0; i < t.Size*2; i += 2 {
		if t.numberSlots[i].hidden {
			num |= (1 << (3 - (i / 2)))
		}
	}

	if num == t.Number {
		t.Hide()
		t.gh.mode.AddScore(t.Number)
		t.gh.sound.Play("tile")
	}
}

func (t *TileSet) Update(dt float64) {
	if t.hidden || t.gh.mode.gameOver {
		return
	}

	// Check for collissions
	if t.tile.fy <= 0.036 {
		return
	}

	for _, o := range t.gh.tiles {
		if o.hidden {
			continue
		}
		if o.id == t.id {
			continue
		}

		if t.tile.fy-0.115 < o.tile.fy && o.tile.fy < t.tile.fy {
			t.Speed = o.Speed
			return
		}
	}

	for i := range t.Sprites {
		t.Sprites[i].ChangeY(-float32(t.Speed * dt))
	}
	for i := range t.numberSlots {
		t.numberSlots[i].ChangeY(-float32(t.Speed * dt))
	}
}

func (t *TileSet) GameOver() {
	t.gh.sound.Play("gameover")
	for i := range t.Sprites {
		t.Sprites[i].ChangeEffect(EffectGameOver)
	}
	for i := range t.numberSlots {
		t.numberSlots[i].ChangeEffect(EffectGameOver)
	}
}
