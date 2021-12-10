package main

import (
	"encoding/json"
	"image"
	"image/draw"
	"math"
	"strings"

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/gl"
)

type Textures struct {
	Types    map[string]Texture
	Uvs      []byte
	Vertices []byte
	Effects  []byte
	vert     gl.Attrib
	uv       gl.Attrib
	ef       gl.Attrib
	texID    gl.Texture
	vbo      gl.Buffer
	ubo      gl.Buffer
	ebo      gl.Buffer // Not element buffer, it's effect buffer ;)
	gh       *Game
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

	// TBD: Dynamic sizes?
	t.Vertices = make([]byte, 550000)
	t.Uvs = make([]byte, 550000)
	t.Effects = make([]byte, 55000)

	return nil
}

func (t *Textures) Init() {
	t.vbo = t.gh.glc.CreateBuffer()
	t.ubo = t.gh.glc.CreateBuffer()
	t.ebo = t.gh.glc.CreateBuffer()

	t.vert = t.gh.glc.GetAttribLocation(t.gh.program, "vert")
	t.uv = t.gh.glc.GetAttribLocation(t.gh.program, "uvs")
	t.ef = t.gh.glc.GetAttribLocation(t.gh.program, "effect")

	t.gh.glc.BindBuffer(gl.ARRAY_BUFFER, t.vbo)
	t.gh.glc.BufferData(gl.ARRAY_BUFFER, t.Vertices, gl.DYNAMIC_DRAW)
	t.gh.glc.VertexAttribPointer(t.vert, 2, gl.FLOAT, false, 2*4, 0)

	t.gh.glc.BindBuffer(gl.ARRAY_BUFFER, t.ubo)
	t.gh.glc.BufferData(gl.ARRAY_BUFFER, t.Uvs, gl.STATIC_DRAW)
	t.gh.glc.VertexAttribPointer(t.uv, 2, gl.FLOAT, true, 2*4, 0)

	t.gh.glc.BindBuffer(gl.ARRAY_BUFFER, t.ebo)
	t.gh.glc.BufferData(gl.ARRAY_BUFFER, t.Effects, gl.STATIC_DRAW)
	t.gh.glc.VertexAttribPointer(t.ef, 2, gl.FLOAT, false, 2*4, 0)

	t.gh.glc.EnableVertexAttribArray(t.vert)
	t.gh.glc.EnableVertexAttribArray(t.uv)
	t.gh.glc.EnableVertexAttribArray(t.ef)

	t.gh.glc.UseProgram(t.gh.program)
	t.gh.glc.ActiveTexture(gl.TEXTURE0)
	t.gh.glc.BindTexture(gl.TEXTURE_2D, t.texID)
}

func (t *Textures) Cleanup() {
	t.gh.glc.DeleteBuffer(t.vbo)
	t.gh.glc.DeleteBuffer(t.ubo)
	t.gh.glc.DeleteBuffer(t.ebo)
	t.gh.glc.DeleteTexture(t.texID)
}

func (t *Textures) Draw() {
	t.gh.glc.DrawArrays(gl.TRIANGLES, 0, len(t.Vertices)/6)
}

func (t *Textures) Update() {
	uvUpdated := false
	effectsUpdated := false
	vertsUpdated := false
	for i := range t.gh.objects {
		if t.gh.objects[i].dirty {
			t.UpdateObject(t.gh.objects[i])
			vertsUpdated = true
		}
		if t.gh.objects[i].dirtyUvs {
			t.UpdateUV(t.gh.objects[i])
			uvUpdated = true
		}
		if t.gh.objects[i].dirtyEffect {
			t.UpdateEffect(t.gh.objects[i])
			effectsUpdated = true
		}
	}

	if vertsUpdated {
		t.gh.glc.BindBuffer(gl.ARRAY_BUFFER, t.vbo)
		t.gh.glc.BufferSubData(gl.ARRAY_BUFFER, 0, t.Vertices)
	}

	if uvUpdated {
		t.gh.glc.BindBuffer(gl.ARRAY_BUFFER, t.ubo)
		t.gh.glc.BufferSubData(gl.ARRAY_BUFFER, 0, t.Uvs)
	}

	if effectsUpdated {
		t.gh.glc.BindBuffer(gl.ARRAY_BUFFER, t.ebo)
		t.gh.glc.BufferSubData(gl.ARRAY_BUFFER, 0, t.Effects)
	}
}

func (t *Textures) UpdateEffect(s *Sprite) {
	eff := math.Float32bits(float32(s.effect))
	t.Effects[4*s.id*12] = byte(eff >> 0)
	t.Effects[4*s.id*12+1] = byte(eff >> 8)
	t.Effects[4*s.id*12+2] = byte(eff >> 16)
	t.Effects[4*s.id*12+3] = byte(eff >> 24)

	t.Effects[4*s.id*12+4] = byte(eff >> 0)
	t.Effects[4*s.id*12+5] = byte(eff >> 8)
	t.Effects[4*s.id*12+6] = byte(eff >> 16)
	t.Effects[4*s.id*12+7] = byte(eff >> 24)

	t.Effects[4*s.id*12+8] = byte(eff >> 0)
	t.Effects[4*s.id*12+9] = byte(eff >> 8)
	t.Effects[4*s.id*12+10] = byte(eff >> 16)
	t.Effects[4*s.id*12+11] = byte(eff >> 24)

	t.Effects[4*s.id*12+12] = byte(eff >> 0)
	t.Effects[4*s.id*12+13] = byte(eff >> 8)
	t.Effects[4*s.id*12+14] = byte(eff >> 16)
	t.Effects[4*s.id*12+15] = byte(eff >> 24)

	t.Effects[4*s.id*12+16] = byte(eff >> 0)
	t.Effects[4*s.id*12+17] = byte(eff >> 8)
	t.Effects[4*s.id*12+18] = byte(eff >> 16)
	t.Effects[4*s.id*12+19] = byte(eff >> 24)

	t.Effects[4*s.id*12+20] = byte(eff >> 0)
	t.Effects[4*s.id*12+21] = byte(eff >> 8)
	t.Effects[4*s.id*12+22] = byte(eff >> 16)
	t.Effects[4*s.id*12+23] = byte(eff >> 24)

	t.Effects[4*s.id*12+24] = byte(eff >> 0)
	t.Effects[4*s.id*12+25] = byte(eff >> 8)
	t.Effects[4*s.id*12+26] = byte(eff >> 16)
	t.Effects[4*s.id*12+27] = byte(eff >> 24)

	t.Effects[4*s.id*12+28] = byte(eff >> 0)
	t.Effects[4*s.id*12+29] = byte(eff >> 8)
	t.Effects[4*s.id*12+30] = byte(eff >> 16)
	t.Effects[4*s.id*12+31] = byte(eff >> 24)

	t.Effects[4*s.id*12+32] = byte(eff >> 0)
	t.Effects[4*s.id*12+33] = byte(eff >> 8)
	t.Effects[4*s.id*12+34] = byte(eff >> 16)
	t.Effects[4*s.id*12+35] = byte(eff >> 24)

	t.Effects[4*s.id*12+36] = byte(eff >> 0)
	t.Effects[4*s.id*12+37] = byte(eff >> 8)
	t.Effects[4*s.id*12+38] = byte(eff >> 16)
	t.Effects[4*s.id*12+39] = byte(eff >> 24)

	t.Effects[4*s.id*12+40] = byte(eff >> 0)
	t.Effects[4*s.id*12+41] = byte(eff >> 8)
	t.Effects[4*s.id*12+42] = byte(eff >> 16)
	t.Effects[4*s.id*12+43] = byte(eff >> 24)

	t.Effects[4*s.id*12+44] = byte(eff >> 0)
	t.Effects[4*s.id*12+45] = byte(eff >> 8)
	t.Effects[4*s.id*12+46] = byte(eff >> 16)
	t.Effects[4*s.id*12+47] = byte(eff >> 24)

	s.dirtyEffect = false
}

func (t *Textures) UpdateObject(s *Sprite) {
	sx := math.Float32bits(s.fx)
	sy := math.Float32bits(s.fy)
	sxw := math.Float32bits(s.fx + s.tx)
	syh := math.Float32bits(s.fy + s.ty)

	// Exploding the updates instead of generating a new
	// byte array for each quad increases the performance from 100us
	// down to about 100ns per sprite.

	t.Vertices[4*s.id*12] = byte(sxw >> 0)
	t.Vertices[4*s.id*12+1] = byte(sxw >> 8)
	t.Vertices[4*s.id*12+2] = byte(sxw >> 16)
	t.Vertices[4*s.id*12+3] = byte(sxw >> 24)

	t.Vertices[4*s.id*12+4] = byte(sy >> 0)
	t.Vertices[4*s.id*12+5] = byte(sy >> 8)
	t.Vertices[4*s.id*12+6] = byte(sy >> 16)
	t.Vertices[4*s.id*12+7] = byte(sy >> 24)

	t.Vertices[4*s.id*12+8] = byte(sxw >> 0)
	t.Vertices[4*s.id*12+9] = byte(sxw >> 8)
	t.Vertices[4*s.id*12+10] = byte(sxw >> 16)
	t.Vertices[4*s.id*12+11] = byte(sxw >> 24)

	t.Vertices[4*s.id*12+12] = byte(syh >> 0)
	t.Vertices[4*s.id*12+13] = byte(syh >> 8)
	t.Vertices[4*s.id*12+14] = byte(syh >> 16)
	t.Vertices[4*s.id*12+15] = byte(syh >> 24)

	t.Vertices[4*s.id*12+16] = byte(sx >> 0)
	t.Vertices[4*s.id*12+17] = byte(sx >> 8)
	t.Vertices[4*s.id*12+18] = byte(sx >> 16)
	t.Vertices[4*s.id*12+19] = byte(sx >> 24)

	t.Vertices[4*s.id*12+20] = byte(sy >> 0)
	t.Vertices[4*s.id*12+21] = byte(sy >> 8)
	t.Vertices[4*s.id*12+22] = byte(sy >> 16)
	t.Vertices[4*s.id*12+23] = byte(sy >> 24)

	t.Vertices[4*s.id*12+24] = byte(sxw >> 0)
	t.Vertices[4*s.id*12+25] = byte(sxw >> 8)
	t.Vertices[4*s.id*12+26] = byte(sxw >> 16)
	t.Vertices[4*s.id*12+27] = byte(sxw >> 24)

	t.Vertices[4*s.id*12+28] = byte(syh >> 0)
	t.Vertices[4*s.id*12+29] = byte(syh >> 8)
	t.Vertices[4*s.id*12+30] = byte(syh >> 16)
	t.Vertices[4*s.id*12+31] = byte(syh >> 24)

	t.Vertices[4*s.id*12+32] = byte(sx >> 0)
	t.Vertices[4*s.id*12+33] = byte(sx >> 8)
	t.Vertices[4*s.id*12+34] = byte(sx >> 16)
	t.Vertices[4*s.id*12+35] = byte(sx >> 24)

	t.Vertices[4*s.id*12+36] = byte(syh >> 0)
	t.Vertices[4*s.id*12+37] = byte(syh >> 8)
	t.Vertices[4*s.id*12+38] = byte(syh >> 16)
	t.Vertices[4*s.id*12+39] = byte(syh >> 24)

	t.Vertices[4*s.id*12+40] = byte(sx >> 0)
	t.Vertices[4*s.id*12+41] = byte(sx >> 8)
	t.Vertices[4*s.id*12+42] = byte(sx >> 16)
	t.Vertices[4*s.id*12+43] = byte(sx >> 24)

	t.Vertices[4*s.id*12+44] = byte(sy >> 0)
	t.Vertices[4*s.id*12+45] = byte(sy >> 8)
	t.Vertices[4*s.id*12+46] = byte(sy >> 16)
	t.Vertices[4*s.id*12+47] = byte(sy >> 24)
	s.dirty = false
}

func (t *Textures) AddSprite(s *Sprite) {
	t.UpdateObject(s)
	t.UpdateUV(s)
	t.UpdateEffect(s)
}

func (t *Textures) UpdateUV(s *Sprite) {
	u1 := math.Float32bits(s.Texture.U1)
	u2 := math.Float32bits(s.Texture.U2)
	v1 := math.Float32bits(s.Texture.V1)
	v2 := math.Float32bits(s.Texture.V2)

	t.Uvs[4*s.id*12] = byte(u2 >> 0)
	t.Uvs[4*s.id*12+1] = byte(u2 >> 8)
	t.Uvs[4*s.id*12+2] = byte(u2 >> 16)
	t.Uvs[4*s.id*12+3] = byte(u2 >> 24)

	t.Uvs[4*s.id*12+4] = byte(v2 >> 0)
	t.Uvs[4*s.id*12+5] = byte(v2 >> 8)
	t.Uvs[4*s.id*12+6] = byte(v2 >> 16)
	t.Uvs[4*s.id*12+7] = byte(v2 >> 24)

	t.Uvs[4*s.id*12+8] = byte(u2 >> 0)
	t.Uvs[4*s.id*12+9] = byte(u2 >> 8)
	t.Uvs[4*s.id*12+10] = byte(u2 >> 16)
	t.Uvs[4*s.id*12+11] = byte(u2 >> 24)

	t.Uvs[4*s.id*12+12] = byte(v1 >> 0)
	t.Uvs[4*s.id*12+13] = byte(v1 >> 8)
	t.Uvs[4*s.id*12+14] = byte(v1 >> 16)
	t.Uvs[4*s.id*12+15] = byte(v1 >> 24)

	t.Uvs[4*s.id*12+16] = byte(u1 >> 0)
	t.Uvs[4*s.id*12+17] = byte(u1 >> 8)
	t.Uvs[4*s.id*12+18] = byte(u1 >> 16)
	t.Uvs[4*s.id*12+19] = byte(u1 >> 24)

	t.Uvs[4*s.id*12+20] = byte(v2 >> 0)
	t.Uvs[4*s.id*12+21] = byte(v2 >> 8)
	t.Uvs[4*s.id*12+22] = byte(v2 >> 16)
	t.Uvs[4*s.id*12+23] = byte(v2 >> 24)

	t.Uvs[4*s.id*12+24] = byte(u2 >> 0)
	t.Uvs[4*s.id*12+25] = byte(u2 >> 8)
	t.Uvs[4*s.id*12+26] = byte(u2 >> 16)
	t.Uvs[4*s.id*12+27] = byte(u2 >> 24)

	t.Uvs[4*s.id*12+28] = byte(v1 >> 0)
	t.Uvs[4*s.id*12+29] = byte(v1 >> 8)
	t.Uvs[4*s.id*12+30] = byte(v1 >> 16)
	t.Uvs[4*s.id*12+31] = byte(v1 >> 24)

	t.Uvs[4*s.id*12+32] = byte(u1 >> 0)
	t.Uvs[4*s.id*12+33] = byte(u1 >> 8)
	t.Uvs[4*s.id*12+34] = byte(u1 >> 16)
	t.Uvs[4*s.id*12+35] = byte(u1 >> 24)

	t.Uvs[4*s.id*12+36] = byte(v1 >> 0)
	t.Uvs[4*s.id*12+37] = byte(v1 >> 8)
	t.Uvs[4*s.id*12+38] = byte(v1 >> 16)
	t.Uvs[4*s.id*12+39] = byte(v1 >> 24)

	t.Uvs[4*s.id*12+40] = byte(u1 >> 0)
	t.Uvs[4*s.id*12+41] = byte(u1 >> 8)
	t.Uvs[4*s.id*12+42] = byte(u1 >> 16)
	t.Uvs[4*s.id*12+43] = byte(u1 >> 24)

	t.Uvs[4*s.id*12+44] = byte(v2 >> 0)
	t.Uvs[4*s.id*12+45] = byte(v2 >> 8)
	t.Uvs[4*s.id*12+46] = byte(v2 >> 16)
	t.Uvs[4*s.id*12+47] = byte(v2 >> 24)
	s.dirtyUvs = false
}

func (t *Textures) AddText(txt string, fx, fy, pz, tx, ty float32, effect Effect) []*Sprite {
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
		s.Init(fx, fy, pz, tx, ty, c, t.gh)

		// 0.005 is for spacing between chars
		s.fx += float32(i) * (tx + 0.005)
		s.dirty = true
		t.gh.AddObjects(&s)
		s.ChangeEffect(effect)
		obj = append(obj, &s)
	}

	return obj
}
