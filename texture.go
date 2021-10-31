package main

import (
	"encoding/binary"
	"encoding/json"
	"image"
	"image/draw"
	"math"
	"strings"

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/gl"
)

type Textures struct {
	Types    map[string]Texture
	Uvs      []byte
	Vertices []byte
	verts    []float32
	uvs      []float32
	res      gl.Uniform
	vert     gl.Attrib
	uv       gl.Attrib
	texID    gl.Texture
	vbo      gl.Buffer
	ubo      gl.Buffer
	vao      gl.VertexArray
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
	t.verts = make([]float32, 1*20000)
	t.Vertices = make([]byte, 1*200000)
	t.uvs = make([]float32, 1*20000)

	return nil
}

func (t *Textures) Init() {
	//t.vao = t.gh.glc.CreateVertexArray()
	t.vbo = t.gh.glc.CreateBuffer()
	t.ubo = t.gh.glc.CreateBuffer()

	t.vert = t.gh.glc.GetAttribLocation(t.gh.program, "vert")
	t.uv = t.gh.glc.GetAttribLocation(t.gh.program, "uvs")
	t.res = t.gh.glc.GetUniformLocation(t.gh.program, "res")
	t.gh.glc.Uniform2fv(t.res, []float32{float32(t.gh.size.WidthPx), float32(t.gh.size.HeightPx)})

	//t.gh.glc.BindVertexArray(t.vao)
	t.gh.glc.BindBuffer(gl.ARRAY_BUFFER, t.vbo)
	t.gh.glc.BufferData(gl.ARRAY_BUFFER, t.Vertices, gl.DYNAMIC_DRAW)
	t.gh.glc.VertexAttribPointer(t.vert, 2, gl.FLOAT, false, 2*4, 0)

	t.gh.glc.BindBuffer(gl.ARRAY_BUFFER, t.ubo)
	t.gh.glc.BufferData(gl.ARRAY_BUFFER, t.Uvs, gl.STATIC_DRAW)
	t.gh.glc.VertexAttribPointer(t.uv, 2, gl.FLOAT, true, 2*4, 0)

	t.gh.glc.EnableVertexAttribArray(t.vert)
	t.gh.glc.EnableVertexAttribArray(t.uv)

	t.gh.glc.UseProgram(t.gh.program)
	t.gh.glc.ActiveTexture(gl.TEXTURE0)
	t.gh.glc.BindTexture(gl.TEXTURE_2D, t.texID)
}

func (t *Textures) SetResolution() {
	if t.gh != nil {
		t.gh.glc.Uniform2fv(t.res, []float32{float32(t.gh.size.WidthPx), float32(t.gh.size.HeightPx)})
	}
}

func (t *Textures) Cleanup() {
	t.gh.glc.DeleteVertexArray(t.vao)
	t.gh.glc.DeleteBuffer(t.vbo)
	t.gh.glc.DeleteBuffer(t.ubo)
	t.gh.glc.DeleteTexture(t.texID)
}

func (t *Textures) Draw() {
	t.gh.glc.DrawArrays(gl.TRIANGLES, 0, len(t.Vertices)/6)
}

func (t *Textures) Update() {
	for i := range t.gh.objects {
		if t.gh.objects[i].dirty {
			t.UpdateObject(t.gh.objects[i])
		}
	}
	t.gh.glc.BindBuffer(gl.ARRAY_BUFFER, t.vbo)
	t.gh.glc.BufferData(gl.ARRAY_BUFFER, t.Vertices, gl.DYNAMIC_DRAW)
}

func (t *Textures) UpdateObject(s *Sprite) {
	sx := math.Float32bits(s.x)
	sy := math.Float32bits(s.y)
	sxw := math.Float32bits(s.x + s.Texture.Width*s.scalex)
	syh := math.Float32bits(s.y + s.Texture.Height*s.scaley)

	// Exploding the updates instead of generating a new
	// byte array for each quad increases the performance from 100us
	// down to about 100ns per sprite.

	//t.verts[s.id*12] = s.x + s.Texture.Width*s.scale
	t.Vertices[4*s.id*12] = byte(sxw >> 0)
	t.Vertices[4*s.id*12+1] = byte(sxw >> 8)
	t.Vertices[4*s.id*12+2] = byte(sxw >> 16)
	t.Vertices[4*s.id*12+3] = byte(sxw >> 24)

	//t.verts[s.id*12+1] = s.y
	t.Vertices[4*s.id*12+4] = byte(sy >> 0)
	t.Vertices[4*s.id*12+5] = byte(sy >> 8)
	t.Vertices[4*s.id*12+6] = byte(sy >> 16)
	t.Vertices[4*s.id*12+7] = byte(sy >> 24)

	//t.verts[s.id*12+2] = s.x + s.Texture.Width*s.scale
	t.Vertices[4*s.id*12+8] = byte(sxw >> 0)
	t.Vertices[4*s.id*12+9] = byte(sxw >> 8)
	t.Vertices[4*s.id*12+10] = byte(sxw >> 16)
	t.Vertices[4*s.id*12+11] = byte(sxw >> 24)

	//t.verts[s.id*12+3] = s.y + s.Texture.Height*s.scale
	t.Vertices[4*s.id*12+12] = byte(syh >> 0)
	t.Vertices[4*s.id*12+13] = byte(syh >> 8)
	t.Vertices[4*s.id*12+14] = byte(syh >> 16)
	t.Vertices[4*s.id*12+15] = byte(syh >> 24)

	// t.verts[s.id*12+4] = s.x
	t.Vertices[4*s.id*12+16] = byte(sx >> 0)
	t.Vertices[4*s.id*12+17] = byte(sx >> 8)
	t.Vertices[4*s.id*12+18] = byte(sx >> 16)
	t.Vertices[4*s.id*12+19] = byte(sx >> 24)

	// t.verts[s.id*12+5] = s.y
	t.Vertices[4*s.id*12+20] = byte(sy >> 0)
	t.Vertices[4*s.id*12+21] = byte(sy >> 8)
	t.Vertices[4*s.id*12+22] = byte(sy >> 16)
	t.Vertices[4*s.id*12+23] = byte(sy >> 24)

	//t.verts[s.id*12+6] = s.x + s.Texture.Width*s.scale
	t.Vertices[4*s.id*12+24] = byte(sxw >> 0)
	t.Vertices[4*s.id*12+25] = byte(sxw >> 8)
	t.Vertices[4*s.id*12+26] = byte(sxw >> 16)
	t.Vertices[4*s.id*12+27] = byte(sxw >> 24)

	// t.verts[s.id*12+7] = s.y + s.Texture.Height*s.scale
	t.Vertices[4*s.id*12+28] = byte(syh >> 0)
	t.Vertices[4*s.id*12+29] = byte(syh >> 8)
	t.Vertices[4*s.id*12+30] = byte(syh >> 16)
	t.Vertices[4*s.id*12+31] = byte(syh >> 24)

	//t.verts[s.id*12+8] = s.x
	t.Vertices[4*s.id*12+32] = byte(sx >> 0)
	t.Vertices[4*s.id*12+33] = byte(sx >> 8)
	t.Vertices[4*s.id*12+34] = byte(sx >> 16)
	t.Vertices[4*s.id*12+35] = byte(sx >> 24)

	//t.verts[s.id*12+9] = s.y + s.Texture.Height*s.scale
	t.Vertices[4*s.id*12+36] = byte(syh >> 0)
	t.Vertices[4*s.id*12+37] = byte(syh >> 8)
	t.Vertices[4*s.id*12+38] = byte(syh >> 16)
	t.Vertices[4*s.id*12+39] = byte(syh >> 24)

	//t.verts[s.id*12+10] = s.x
	t.Vertices[4*s.id*12+40] = byte(sx >> 0)
	t.Vertices[4*s.id*12+41] = byte(sx >> 8)
	t.Vertices[4*s.id*12+42] = byte(sx >> 16)
	t.Vertices[4*s.id*12+43] = byte(sx >> 24)

	//t.verts[s.id*12+11] = s.y
	t.Vertices[4*s.id*12+44] = byte(sy >> 0)
	t.Vertices[4*s.id*12+45] = byte(sy >> 8)
	t.Vertices[4*s.id*12+46] = byte(sy >> 16)
	t.Vertices[4*s.id*12+47] = byte(sy >> 24)
	s.dirty = false
}

func (t *Textures) AddSprite(s *Sprite) {
	t.UpdateObject(s)
	t.uvs[s.id*12] = s.Texture.U2
	t.uvs[s.id*12+1] = s.Texture.V2

	t.uvs[s.id*12+2] = s.Texture.U2
	t.uvs[s.id*12+3] = s.Texture.V1

	t.uvs[s.id*12+4] = s.Texture.U1
	t.uvs[s.id*12+5] = s.Texture.V2

	t.uvs[s.id*12+6] = s.Texture.U2
	t.uvs[s.id*12+7] = s.Texture.V1

	t.uvs[s.id*12+8] = s.Texture.U1
	t.uvs[s.id*12+9] = s.Texture.V1

	t.uvs[s.id*12+10] = s.Texture.U1
	t.uvs[s.id*12+11] = s.Texture.V2

	// UVs are just updated once per sprite so we never need to generate this
	// more than once, compared to vertices.
	t.Vertices = f32.Bytes(binary.LittleEndian, t.verts...)
	t.Uvs = f32.Bytes(binary.LittleEndian, t.uvs...)
}

func (t *Textures) AddText(txt string, px, py, pz, scalex, scaley float32, effect Effect) []*Sprite {
	obj := []*Sprite{}
	if txt == "" {
		return obj
	}
	txt = strings.ToLower(txt)

	for i, ch := range txt {
		c := string(ch)
		if ch == ':' {
			c = "colon"
		}
		s := Sprite{}
		s.Init(px, py, pz, scalex, scaley, c, t.gh)

		// 10 is just an arbitrary offset between characters
		s.x += float32(i)*s.Texture.Width*scalex + float32(i)*10
		t.gh.AddObjects(&s)
		obj = append(obj, &s)
	}

	return obj
}