package main

import (
	"encoding/binary"
	"encoding/json"
	"image"
	"image/draw"

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
	vert     gl.Attrib
	uv       gl.Attrib
	texID    gl.Texture
	vbo      gl.Buffer
	ubo      gl.Buffer
	vao      gl.VertexArray
	gh       *Game
}

type Texture struct {
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

	//image_draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, image_draw.Src)
	draw.Draw(rgba, b, img, image.Point{0, 0}, draw.Src)

	// x = float64(rgba.Bounds().Max.X - rgba.Bounds().Min.X)
	// y = float64(rgba.Bounds().Max.Y - rgba.Bounds().Min.Y)

	t.texID = t.gh.glc.CreateTexture()
	// t.gh.glc.ActiveTexture(gl.TEXTURE0)
	t.gh.glc.BindTexture(gl.TEXTURE_2D, t.texID)
	t.gh.glc.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	t.gh.glc.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	t.gh.glc.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE) // REPEAT?
	t.gh.glc.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE) // REPEAT?

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
		t.Types[k] = Texture{
			Width:  v.Frame.W,
			Height: v.Frame.H,
			U2:     float32(v.Frame.X+v.Frame.W) / float32(layout.Meta.Size.W),
			V1:     float32(v.Frame.Y) / float32(layout.Meta.Size.H),
			V2:     float32(v.Frame.Y+v.Frame.H) / float32(layout.Meta.Size.H),
			U1:     float32(v.Frame.X) / float32(layout.Meta.Size.W),
		}
	}

	t.verts = make([]float32, 1*20000)
	t.uvs = make([]float32, 1*20000)

	return nil
}

func (t *Textures) Init() {
	//t.vao = t.gh.glc.CreateVertexArray()
	t.vbo = t.gh.glc.CreateBuffer()
	t.ubo = t.gh.glc.CreateBuffer()

	t.vert = t.gh.glc.GetAttribLocation(t.gh.program, "vert")
	t.uv = t.gh.glc.GetAttribLocation(t.gh.program, "uvs")

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
	//t.gh.glc.BindVertexArray(t.vao) // NEEDED?

}

func (t *Textures) Cleanup() {
	t.gh.glc.DeleteVertexArray(t.vao)
	t.gh.glc.DeleteBuffer(t.vbo)
	t.gh.glc.DeleteBuffer(t.ubo)
	t.gh.glc.DeleteTexture(t.texID)
}

func (t *Textures) Draw() {
	t.gh.glc.DrawArrays(gl.TRIANGLES, 0, len(t.Vertices)) // ObjectCount * 6
}

func (t *Textures) Update() {
	for i := range t.gh.objects {
		if t.gh.objects[i].dirty {
			t.UpdateObject(&t.gh.objects[i])
		}
	}
	//t.gh.glc.BufferSubData(gl.ARRAY_BUFFER, 0, t.Vertices)
}

func (t *Textures) UpdateObject(s *Sprite) {
	t.verts[s.id*12] = s.x + s.Texture.Width
	t.verts[s.id*12+1] = s.y

	t.verts[s.id*12+2] = s.x + s.Texture.Width
	t.verts[s.id*12+3] = s.y + s.Texture.Height

	t.verts[s.id*12+4] = s.x
	t.verts[s.id*12+5] = s.y

	t.verts[s.id*12+6] = s.x + s.Texture.Width
	t.verts[s.id*12+7] = s.y + s.Texture.Height

	t.verts[s.id*12+8] = s.x
	t.verts[s.id*12+9] = s.y + s.Texture.Height

	t.verts[s.id*12+10] = s.x
	t.verts[s.id*12+11] = s.y

	t.Vertices = f32.Bytes(binary.LittleEndian, t.verts...)
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

	t.Vertices = f32.Bytes(binary.LittleEndian, t.verts...)
	t.Uvs = f32.Bytes(binary.LittleEndian, t.uvs...)
}
