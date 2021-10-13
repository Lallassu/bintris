package main

import (
	"image"
	"image/draw"
	"strings"

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/gl"
)

type Font struct {
	chars  map[string]*Sprite
	width  int
	height int
	g      *Game
}

func (f *Font) Init(g *Game) {
	f.g = g

	imgFile, err := asset.Open("font.png")
	if err != nil {
		return
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return
	}

	f.width = 18
	f.height = 32
	bx := img.Bounds().Max.X

	f.chars = make(map[string]*Sprite)

	chars := []string{
		" ", "!", "\"", "#", "$", "%", "", "'", "(", ")", " ", "+", ",", "-", ".", "/", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", ":", ";", "<", "=", ">", "?",
		"@", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "[", "\\", "]", "^", "_",
		"'", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "{", "|", "}", "~", " ",
	}

	y := 0
	x := 0
	for i, c := range chars {
		if i%(int(bx)/int(f.width)+1) == 0 && i > 0 {
			x = 0
			y++
		}

		if f.chars[c] == nil {
			f.chars[c] = &Sprite{}

			b := image.Rectangle{
				Min: image.Point{
					X: x * f.width,
					Y: y * f.height,
				},
				Max: image.Point{
					X: x*f.width + f.width - 4,
					Y: y*f.height + f.height - 2,
				},
			}

			rgba := image.NewRGBA(b)
			draw.Draw(rgba, b, img, b.Bounds().Min, draw.Over)
			f.chars[c].Init(0, 0, 0, 0, "", rgba, g)
		}
		x++
	}
}

func (f *Font) AddText(txt string, px, py, pz, scale float64, effect Effect) []Object {
	obj := []Object{}
	if txt == "" {
		return obj
	}
	txt = strings.ToLower(txt)

	vbos := make(map[string]gl.Buffer)

	for i, ch := range txt {
		if _, ok := f.chars[string(ch)]; !ok {
			continue
		}
		c := *f.chars[string(ch)]
		if _, ok := vbos[string(ch)]; !ok {
			vbos[string(ch)] = f.chars[string(ch)].vbo
		}

		c.vbo = vbos[string(ch)]
		c.id = f.g.NewID()
		c.x = px + float64(i*(f.width))*scale
		c.effect = effect
		c.y = py
		c.z = pz
		c.scale = scale
		obj = append(obj, &c)
	}
	return obj
}
