package main

import (
	"encoding/binary"
	"image"
	"image/draw"
	_ "image/png"

	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/gl"
)

type Sprite struct {
	gh         *Game
	id         int
	uModel     gl.Uniform
	uEffect    gl.Uniform
	aPosition  gl.Attrib
	aTexture   gl.Attrib
	texture    gl.Texture
	vbo        gl.Buffer
	vertices   []byte
	x          float64
	y          float64
	z          float64
	scale      float64
	modelf     []float32
	rotation   float64
	width      float32
	height     float32
	sx         float64
	sy         float64
	effect     Effect
	hidden     bool
	objectType ObjectType
}

func (s *Sprite) Init(x, y, z, scale float64, file string, img *image.RGBA, g *Game) {
	s.gh = g
	s.vertices = f32.Bytes(binary.LittleEndian, 0)
	s.uModel = s.gh.glc.GetUniformLocation(s.gh.program, "model")
	s.uEffect = s.gh.glc.GetUniformLocation(s.gh.program, "effect")
	s.aPosition = s.gh.glc.GetAttribLocation(s.gh.program, "position")
	s.aTexture = s.gh.glc.GetAttribLocation(s.gh.program, "texture")
	s.scale = scale
	s.rotation = 180
	s.x = x
	s.y = y
	s.z = z

	s.vbo = s.gh.glc.CreateBuffer()

	s.vertices = f32.Bytes(binary.LittleEndian,
		0.0, 1.0, 0.0, 1.0,
		1.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 0.0,

		0.0, 1.0, 0.0, 1.0,
		1.0, 1.0, 1.0, 1.0,
		1.0, 0.0, 1.0, 0.0,
	)

	var err error
	s.texture, s.sx, s.sy, err = s.LoadTexture(file, img)
	if err != nil {
		panic(err)
	}

	//g.glc.GenerateMipmap(gl.TEXTURE_2D)

	s.id = s.gh.NewID()
}

func (s *Sprite) Copy(x, y, z, scale float64, fs *Sprite) {
	s.gh = fs.gh
	s.id = s.gh.NewID()
	s.vbo = fs.vbo
	s.texture = fs.texture
	s.sx = fs.sx
	s.sy = fs.sy
	s.vertices = fs.vertices
	s.rotation = fs.rotation
	s.x = x
	s.y = y
	s.z = z
	s.scale = scale
	s.uModel = s.gh.glc.GetUniformLocation(s.gh.program, "model")
	s.uEffect = s.gh.glc.GetUniformLocation(s.gh.program, "effect")
	s.aPosition = s.gh.glc.GetAttribLocation(s.gh.program, "position")
	s.aTexture = s.gh.glc.GetAttribLocation(s.gh.program, "texture")
}

func (s *Sprite) GetObjectType() ObjectType {
	return ObjectTypeSprite
}

func (s *Sprite) Update(dt float64) {
}

func (s *Sprite) Draw(dt float64) {
	if s.hidden {
		return
	}

	s.UpdatePosition()

	s.gh.glc.UseProgram(s.gh.program)
	s.gh.glc.UniformMatrix4fv(s.gh.uView, s.gh.viewf)
	s.gh.glc.UniformMatrix4fv(s.gh.uProj, s.gh.projf)

	s.gh.glc.BindBuffer(gl.ARRAY_BUFFER, s.vbo)
	s.gh.glc.BufferData(gl.ARRAY_BUFFER, s.vertices, gl.STATIC_DRAW)
	s.gh.glc.BindTexture(gl.TEXTURE_2D, s.texture)

	s.gh.glc.UniformMatrix4fv(s.uModel, s.modelf)
	s.gh.glc.Uniform1f(s.gh.uTime, float32(s.gh.elapsed))
	s.gh.glc.Uniform1i(s.uEffect, int(s.effect))

	s.gh.glc.EnableVertexAttribArray(s.aPosition)
	s.gh.glc.VertexAttribPointer(s.aPosition, 4, gl.FLOAT, false, 0, 0)

	s.gh.glc.DrawArrays(gl.TRIANGLES, 0, len(s.vertices))
	s.gh.glc.DisableVertexAttribArray(s.aPosition)
}

func (s *Sprite) Delete() {
	s.hidden = true
}
func (s *Sprite) GetY() float64 {
	return s.y
}
func (s *Sprite) GetX() float64 {
	return s.x
}
func (s *Sprite) GetID() int {
	return s.id
}

func (s *Sprite) UpdatePosition() {
	translate := mgl32.Translate3D(float32(s.x), float32(s.y), float32(s.z))
	trans := translate

	scalem4 := mgl32.Scale3D(float32(s.sx*s.scale), float32(s.sy*s.scale), 0.0)

	rot := mgl32.HomogRotate3D(float32(mgl32.DegToRad(float32(s.rotation))), mgl32.Vec3{1.0, 0.0, 0.0})
	trans = translate.Mul4(scalem4).Mul4(rot)

	s.modelf = trans[:]
}

func (s *Sprite) LoadTexture(name string, rgba *image.RGBA) (tex gl.Texture, x float64, y float64, err error) {
	if rgba == nil {
		imgFile, e := asset.Open(name)
		if e != nil {
			err = e
			return
		}

		img, _, e := image.Decode(imgFile)
		if e != nil {
			err = e
			return
		}

		rgba = image.NewRGBA(img.Bounds())
		b := image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{img.Bounds().Max.X, img.Bounds().Max.Y},
		}

		//image_draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, image_draw.Src)
		draw.Draw(rgba, b, img, image.Point{0, 0}, draw.Src)
	}
	x = float64(rgba.Bounds().Max.X - rgba.Bounds().Min.X)
	y = float64(rgba.Bounds().Max.Y - rgba.Bounds().Min.Y)

	tex = s.gh.glc.CreateTexture()
	s.gh.glc.ActiveTexture(gl.TEXTURE0)
	s.gh.glc.BindTexture(gl.TEXTURE_2D, tex)
	s.gh.glc.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	s.gh.glc.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	s.gh.glc.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	s.gh.glc.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	//s.gh.glc.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	//s.gh.glc.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	s.gh.glc.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		rgba.Rect.Size().X,
		rgba.Rect.Size().Y,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		rgba.Pix)

	return
}
