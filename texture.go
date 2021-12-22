package main

import (
	"encoding/json"
	"image"
	"image/draw"
	"strings"

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/gl"
)

type Textures struct {
	Types map[string]Texture
	texID gl.Texture
	gh    *Game
}

type Texture struct {
	Name   string
	Width  float32
	Height float32
	U1     float32
	U2     float32
	V1     float32
	V2     float32
}

type Layout struct {
	Meta   Meta
	Frames map[string]Frame
}

type Frame struct {
	Frame Size
}

type Meta struct {
	Size Size
}

type Size struct {
	X float32
	Y float32
	W float32
	H float32
}

// Load sprite sheet into textures
func (t *Textures) Load(texFile, layoutFile string, gh *Game) error {
	t.gh = gh

	imgFile, err := asset.Open(texFile)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return err
	}

	rgba := image.NewRGBA(img.Bounds())
	b := image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{img.Bounds().Max.X, img.Bounds().Max.Y},
	}

	draw.Draw(rgba, b, img, image.Point{0, 0}, draw.Src)

	t.texID = t.gh.glc.CreateTexture()
	t.gh.glc.BindTexture(gl.TEXTURE_2D, t.texID)
	t.gh.glc.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	t.gh.glc.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	t.gh.glc.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	t.gh.glc.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	t.gh.glc.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		rgba.Rect.Size().X,
		rgba.Rect.Size().Y,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		rgba.Pix)

	// Now load the packed layout information and create UVs for the different textures
	l, err := asset.Open(layoutFile)
	if err != nil {
		return err
	}
	layout := Layout{}

	if err := json.NewDecoder(l).Decode(&layout); err != nil {
		return err
	}

	t.Types = make(map[string]Texture)

	for k, v := range layout.Frames {
		k = strings.TrimSuffix(k, ".png")
		t.Types[k] = Texture{
			Name:   k,
			Width:  v.Frame.W,
			Height: v.Frame.H,
			U2:     float32(v.Frame.X+v.Frame.W) / float32(layout.Meta.Size.W),
			V1:     float32(v.Frame.Y) / float32(layout.Meta.Size.H),
			V2:     float32(v.Frame.Y+v.Frame.H) / float32(layout.Meta.Size.H),
			U1:     float32(v.Frame.X) / float32(layout.Meta.Size.W),
		}
	}

	return nil
}

func (t *Textures) Init() {
	t.gh.glc.ActiveTexture(gl.TEXTURE0)
	t.gh.glc.BindTexture(gl.TEXTURE_2D, t.texID)
}

func (t *Textures) Cleanup() {
	t.gh.glc.DeleteTexture(t.texID)
}

func (t *Textures) AddText(txt string, fx, fy, pz, tx, ty float32, effect Effect, sType SpriteType) []*Sprite {
	obj := []*Sprite{}
	if txt == "" {
		return obj
	}
	txt = strings.ToLower(txt)

	for i, ch := range txt {
		c := string(ch)
		if c == "/" {
			c = "slash"
		} else if c == "." {
			c = "dot"
		}
		s := Sprite{}
		s.Init(fx, fy, pz, tx, ty, c, t.gh, sType)

		// 0.005 is for spacing between chars
		s.fx += float32(i) * (tx + 0.005)
		s.dirty = true
		t.gh.AddObjects(&s)
		s.ChangeEffect(effect)
		obj = append(obj, &s)
	}

	return obj
}
