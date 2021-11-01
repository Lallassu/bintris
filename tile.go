package main

import (
	"fmt"
	"strconv"
	"time"
)

type TileSet struct {
	id          int
	Size        int // Number of tiles
	Number      int
	NumberStr   string // Textual representation of the bin number.
	X           float32
	Y           float32
	Sprites     []*Sprite
	numberSlots []*Sprite
	Speed       float64
	ClickFunc   func(x, y float32)
	sizex       float32
	sizey       float32
	scale       float32
	clicked     time.Time
	objectType  ObjectType
	gh          *Game
	hidden      bool
	tileWidth   float32
	tileHeight  float32
	tile        *Sprite
}

func (t *TileSet) Init(scale float32, size int, number int, x, y float32, g *Game) {
	t.gh = g
	t.id = g.NewID()
	t.Speed = float64(g.size.HeightPx / 5)
	t.Number = number
	t.Size = size
	t.Y = g.Y(y)
	fmt.Printf("==> %v\n", t.Y)
	t.X = g.X(x)
	t.sizex = t.tileWidth
	t.sizey = t.tileHeight
	t.scale = scale
	t.tileWidth = g.X(241)
	t.tileHeight = g.Y(30)

	t.NumberStr = strconv.Itoa(number)
	c := &Sprite{hidden: false}
	t.Sprites = append(t.Sprites, c)

	c.Init(g.X(13), t.Y, 0.6, 1.0, 1.0, "tile", g)
	c.scalex = c.gh.SX(3.2)
	c.scaley = c.gh.SY(2)
	t.tile = c
	g.AddObjects(c)

	offset := float32(t.gh.size.HeightPx / 25)
	for x := 241 / 4; x < 241; x += 241 / 4 {
		s := g.tex.AddText("0", g.X(float32(x)-(241/8)+8), t.Y+offset, 0.7, c.gh.SX(3.8), c.gh.SY(4.0), EffectMetaballsBlue)

		t.numberSlots = append(t.numberSlots, s...)
		s = g.tex.AddText("1", g.X(float32(x)-(241/8)+8), t.Y+offset, 0.7, c.gh.SX(3.8), c.gh.SY(4.0), EffectMetaballsBlue)
		s[0].Hide()
		t.numberSlots = append(t.numberSlots, s...)
	}

	n := g.tex.AddText(t.NumberStr, g.X(270), g.Y(y+10), 0.7, c.gh.SX(3.8), c.gh.SY(4.0), EffectMetaballsBlue)
	t.Sprites = append(t.Sprites, n...)

	//t.ClickFunc = func(x, y float32) {
	//	c.Click(x, y)
	//}

}

func (t *TileSet) SetSpeed(s int) {
	t.Speed = float64(t.gh.size.HeightPx / s)
}

func (s *TileSet) Hidden() bool {
	return s.hidden
}

func (t *TileSet) GetID() int {
	return t.id
}

func (t *TileSet) Click(x, y float32) {
	if time.Since(t.clicked) < time.Duration(100*time.Millisecond) {
		return
	}

	t.clicked = time.Now()
	slot0 := 0
	slot1 := 0
	if x > t.X && x < t.X+t.tileWidth/float32(t.Size) {
		slot0 = 0
		slot1 = 1
	} else if x > t.X+t.tileWidth/float32(t.Size) && x < t.X+t.tileWidth/float32(t.Size)*2 {
		slot0 = 2
		slot1 = 3
	} else if x > t.X+t.tileWidth/float32(t.Size)*2 && x < t.X+t.tileWidth/float32(t.Size)*3 {
		slot0 = 4
		slot1 = 5
	} else if x > t.X+t.tileWidth/float32(t.Size)*3 && x < t.X+t.tileWidth/float32(t.Size)*4 {
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

func (t *TileSet) Reset() {
	t.hidden = false
	t.tile.Show()
	t.Y = t.gh.Y(320)

	offset := float32(t.gh.size.HeightPx / 25)
	for i := range t.numberSlots {
		t.numberSlots[i].Show()
		t.numberSlots[i].x = t.X + float32(i)*float32(t.tileWidth/4) + t.gh.X(15)
		t.numberSlots[i].y = t.Y + offset
		t.numberSlots[i].Show()
	}

	for i := range t.Sprites {
		t.Sprites[i].Show()
	}

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
	}
}

func (t *TileSet) Hide() {
	for i := range t.Sprites {
		t.Sprites[i].Hide()
	}
	for i := range t.numberSlots {
		t.numberSlots[i].Hide()
	}
	t.tile.Hide()
	t.hidden = true
}

func (t *TileSet) GetObjectType() ObjectType {
	return ObjectTypeTileSet
}

func (t *TileSet) Draw(dt float64) {
}

func (c *TileSet) GetX() float32 {
	return c.X
}

func (c *TileSet) GetY() float32 {
	return c.Y
}

func (t *TileSet) Update(dt float64) {
	if t.hidden {
		return
	}

	// Check for collissions
	if t.Y < t.gh.Y(9) {
		return
	}

	for _, o := range t.gh.tiles {
		if o.hidden {
			continue
		}
		if o.GetID() == t.id {
			continue
		}

		if int(t.Y-t.tileHeight-t.gh.Y(10)) < int(o.GetY()) && int(o.GetY()) < int(t.Y) {
			return
		}
	}

	t.Y -= float32(t.Speed * dt)
	for i := range t.Sprites {
		t.Sprites[i].ChangeY(-float32(t.Speed * dt))
	}
	for i := range t.numberSlots {
		t.numberSlots[i].ChangeY(-float32(t.Speed * dt))
	}
}

func (t *TileSet) Resize() {
	//a := float32(s.gh.size.WidthPx) / float32(s.gh.size.HeightPx)
	t.tileWidth *= float32(t.gh.size.WidthPx) / float32(t.gh.sizePrev.WidthPx)
	t.tileHeight *= float32(t.gh.size.HeightPx) / float32(t.gh.sizePrev.HeightPx)
	t.X *= float32(t.gh.size.WidthPx) / float32(t.gh.sizePrev.WidthPx)
	t.Y *= float32(t.gh.size.HeightPx) / float32(t.gh.sizePrev.HeightPx)
	t.Speed *= float64(float32(t.gh.size.HeightPx) / float32(t.gh.sizePrev.HeightPx))
}
