package main

import (
	"strconv"
	"time"
)

type TileSet struct {
	Size        int // Number of tiles
	Number      int
	NumberStr   string // Textual representation of the bin number.
	X           float64
	Y           float64
	Objects     []Object
	number      []Object
	numberSlots []Object
	Speed       float64
	id          int
	ClickFunc   func(x, y float32)
	sizex       int
	sizey       int
	scale       float64
	clicked     time.Time
	deleted     bool
	objectType  ObjectType
	gh          *Game
}

const tileWidth = 241
const tileHeight = 30

var tileSprite *Sprite
var tileInited bool

func (t *TileSet) Init(scale float64, size int, number int, x, y float64, g *Game) {
	t.gh = g
	t.id = g.NewID()
	t.Speed = 50
	t.Number = number
	t.Size = size
	t.Y = y
	t.X = x
	t.sizex = tileWidth
	t.sizey = tileHeight
	t.scale = scale

	t.NumberStr = strconv.Itoa(number)
	c := &Sprite{hidden: false}
	t.Objects = append(t.Objects, c)

	// TBD: rewrite this
	if tileInited {
		c.Copy(13, t.Y, 0.6, 1.0, tileSprite)
	} else {
		c.Init(13, t.Y, 0.6, 1.0, "tile.png", nil, g)
		tileSprite = c
		tileInited = true
	}

	for x := tileWidth / 4; x < tileWidth; x += tileWidth / 4 {
		txt := g.font.AddText("0", float64(x-(tileWidth/8)+8), float64(y)-3, 0.7, 0.7, 0)
		t.numberSlots = append(t.numberSlots, txt...)
		txt = g.font.AddText("1", float64(x-(tileWidth/8)+8), float64(y)-3, 0.7, 0.7, 0)
		txt[0].(*Sprite).hidden = true
		t.numberSlots = append(t.numberSlots, txt...)
	}

	t.number = g.font.AddText(t.NumberStr, 280, t.Y-2, 0.7, 0.8, 0)
	t.Objects = append(t.Objects, t.number...)
	// Create tileset based on size

	//t.ClickFunc = func(x, y float32) {
	//	c.Click(x, y)
	//}

	g.AddObjects(t.Objects...)
	g.AddObjects(t.numberSlots...)
}

func (t *TileSet) Delete() {
	t.deleted = true
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
	if float64(x) > t.X && float64(x) < t.X+tileWidth/float64(t.Size) {
		slot0 = 0
		slot1 = 1
	} else if float64(x) > t.X+tileWidth/float64(t.Size) && float64(x) < t.X+tileWidth/float64(t.Size)*2 {
		slot0 = 2
		slot1 = 3
	} else if float64(x) > t.X+tileWidth/float64(t.Size)*2 && float64(x) < t.X+tileWidth/float64(t.Size)*3 {
		slot0 = 4
		slot1 = 5
	} else if float64(x) > t.X+tileWidth/float64(t.Size)*3 && float64(x) < t.X+tileWidth/float64(t.Size)*4 {
		slot0 = 6
		slot1 = 7
	}
	c0 := t.numberSlots[slot0].(*Sprite)
	c1 := t.numberSlots[slot1].(*Sprite)
	if c0.hidden {
		c0.hidden = false
		c1.hidden = true
	} else {
		c0.hidden = true
		c1.hidden = false
	}

	// Verify number
	t.VerifyNumber()
}

func (t *TileSet) VerifyNumber() {
	num := 0
	for i := 0; i < t.Size*2; i += 2 {
		if t.numberSlots[i].(*Sprite).hidden {
			num |= (1 << (3 - (i / 2)))
		}
	}

	if num == t.Number {
		for i := range t.Objects {
			t.gh.DeleteObject(t.Objects[i])
		}
		for i := range t.numberSlots {
			t.gh.DeleteObject(t.numberSlots[i])
		}
		t.gh.DeleteObject(t)
	}
}

func (t *TileSet) GetObjectType() ObjectType {
	return ObjectTypeTileSet
}

func (t *TileSet) Draw(dt float64) {
}

func (c *TileSet) GetX() float64 {
	return c.X
}

func (c *TileSet) GetY() float64 {
	return c.Y
}

func (t *TileSet) Update(dt float64) {
	if t.deleted {
		return
	}

	// Check for collissions
	if t.Y < 40 {
		return
	}

	for _, o := range t.gh.objects {
		if o.GetObjectType() == ObjectTypeTileSet {
			if o.(*TileSet).deleted {
				continue
			}
			if o.GetID() == t.id {
				continue
			}

			if int(t.Y-tileHeight-4) < int(o.GetY()) && int(o.GetY()) < int(t.Y) {
				return
			}
		}
	}

	t.Y -= t.Speed * dt
	for i := range t.Objects {
		t.Objects[i].(*Sprite).y -= t.Speed * dt
	}
	for i := range t.numberSlots {
		t.numberSlots[i].(*Sprite).y -= t.Speed * dt
	}
}
