package main

import (
	"math"

	"golang.org/x/mobile/gl"
)

type GLData struct {
	Uvs      []byte
	Vertices []byte
	Effects  []byte
	vert     gl.Attrib
	uv       gl.Attrib
	ef       gl.Attrib
	vbo      gl.Buffer
	ubo      gl.Buffer
	ebo      gl.Buffer // Not element buffer, it's effect buffer ;)
	gh       *Game
	sType    SpriteType
	vertMenu int
}

func (g *GLData) Init(gh *Game, bufferSize int) {
	g.gh = gh

	// TBD: Dynamic sizes?
	g.Vertices = make([]byte, bufferSize)
	g.Uvs = make([]byte, bufferSize)
	g.Effects = make([]byte, bufferSize)

	g.vbo = g.gh.glc.CreateBuffer()
	g.ubo = g.gh.glc.CreateBuffer()
	g.ebo = g.gh.glc.CreateBuffer()

	g.vert = g.gh.glc.GetAttribLocation(g.gh.program, "vert")
	g.uv = g.gh.glc.GetAttribLocation(g.gh.program, "uvs")
	g.ef = g.gh.glc.GetAttribLocation(g.gh.program, "effect")

	g.gh.glc.BindBuffer(gl.ARRAY_BUFFER, g.vbo)
	g.gh.glc.BufferData(gl.ARRAY_BUFFER, g.Vertices, gl.DYNAMIC_DRAW)
	g.gh.glc.VertexAttribPointer(g.vert, 2, gl.FLOAT, false, 2*4, 0)

	g.gh.glc.BindBuffer(gl.ARRAY_BUFFER, g.ubo)
	g.gh.glc.BufferData(gl.ARRAY_BUFFER, g.Uvs, gl.STATIC_DRAW)
	g.gh.glc.VertexAttribPointer(g.uv, 2, gl.FLOAT, true, 2*4, 0)

	g.gh.glc.BindBuffer(gl.ARRAY_BUFFER, g.ebo)
	g.gh.glc.BufferData(gl.ARRAY_BUFFER, g.Effects, gl.STATIC_DRAW)
	g.gh.glc.VertexAttribPointer(g.ef, 2, gl.FLOAT, false, 2*4, 0)

	g.gh.glc.EnableVertexAttribArray(g.vert)
	g.gh.glc.EnableVertexAttribArray(g.uv)
	g.gh.glc.EnableVertexAttribArray(g.ef)
}

func (g *GLData) Enable(sType SpriteType) {
	g.gh.glc.BindBuffer(gl.ARRAY_BUFFER, g.vbo)
	g.gh.glc.BufferData(gl.ARRAY_BUFFER, g.Vertices, gl.DYNAMIC_DRAW)
	g.gh.glc.BindBuffer(gl.ARRAY_BUFFER, g.ubo)
	g.gh.glc.BufferData(gl.ARRAY_BUFFER, g.Uvs, gl.STATIC_DRAW)
	g.gh.glc.BindBuffer(gl.ARRAY_BUFFER, g.ebo)
	g.gh.glc.BufferData(gl.ARRAY_BUFFER, g.Effects, gl.STATIC_DRAW)
	g.sType = sType
}

func (g *GLData) Cleanup() {
	g.gh.glc.DeleteBuffer(g.vbo)
	g.gh.glc.DeleteBuffer(g.ubo)
	g.gh.glc.DeleteBuffer(g.ebo)
}

func (g *GLData) AddSprite(s *Sprite) {
	if s.sType == SpriteMenu {
		g.vertMenu += 6
	}

	g.UpdateObject(s)
	g.UpdateUV(s)
	g.UpdateEffect(s)
}

func (g *GLData) Draw() {
	if g.sType == SpritePlay {
		g.gh.glc.DrawArrays(gl.TRIANGLES, g.vertMenu, len(g.Vertices))
	} else {
		g.gh.glc.DrawArrays(gl.TRIANGLES, 0, g.vertMenu)
	}
}

func (g *GLData) Update(sType SpriteType) {
	uvUpdated := false
	effectsUpdated := false
	vertsUpdated := false
	var objects []*Sprite
	if sType == SpritePlay {
		objects = g.gh.objectsPlay
	} else {
		objects = g.gh.objectsMenu
	}

	for i := range objects {
		if objects[i].dirty && objects[i].sType == sType {
			g.UpdateObject(objects[i])
			vertsUpdated = true
		}
		if objects[i].dirtyUvs && objects[i].sType == sType {
			g.UpdateUV(objects[i])
			uvUpdated = true
		}
		if objects[i].dirtyEffect && objects[i].sType == sType {
			g.UpdateEffect(objects[i])
			effectsUpdated = true
		}
	}

	if vertsUpdated {
		g.gh.glc.BindBuffer(gl.ARRAY_BUFFER, g.vbo)
		g.gh.glc.BufferSubData(gl.ARRAY_BUFFER, 0, g.Vertices)
	}

	if uvUpdated {
		g.gh.glc.BindBuffer(gl.ARRAY_BUFFER, g.ubo)
		g.gh.glc.BufferSubData(gl.ARRAY_BUFFER, 0, g.Uvs)
	}

	if effectsUpdated {
		g.gh.glc.BindBuffer(gl.ARRAY_BUFFER, g.ebo)
		g.gh.glc.BufferSubData(gl.ARRAY_BUFFER, 0, g.Effects)
	}
}

func (g *GLData) UpdateEffect(s *Sprite) {
	eff := math.Float32bits(float32(s.effect))
	g.Effects[4*s.id*12] = byte(eff >> 0)
	g.Effects[4*s.id*12+1] = byte(eff >> 8)
	g.Effects[4*s.id*12+2] = byte(eff >> 16)
	g.Effects[4*s.id*12+3] = byte(eff >> 24)

	g.Effects[4*s.id*12+4] = byte(eff >> 0)
	g.Effects[4*s.id*12+5] = byte(eff >> 8)
	g.Effects[4*s.id*12+6] = byte(eff >> 16)
	g.Effects[4*s.id*12+7] = byte(eff >> 24)

	g.Effects[4*s.id*12+8] = byte(eff >> 0)
	g.Effects[4*s.id*12+9] = byte(eff >> 8)
	g.Effects[4*s.id*12+10] = byte(eff >> 16)
	g.Effects[4*s.id*12+11] = byte(eff >> 24)

	g.Effects[4*s.id*12+12] = byte(eff >> 0)
	g.Effects[4*s.id*12+13] = byte(eff >> 8)
	g.Effects[4*s.id*12+14] = byte(eff >> 16)
	g.Effects[4*s.id*12+15] = byte(eff >> 24)

	g.Effects[4*s.id*12+16] = byte(eff >> 0)
	g.Effects[4*s.id*12+17] = byte(eff >> 8)
	g.Effects[4*s.id*12+18] = byte(eff >> 16)
	g.Effects[4*s.id*12+19] = byte(eff >> 24)

	g.Effects[4*s.id*12+20] = byte(eff >> 0)
	g.Effects[4*s.id*12+21] = byte(eff >> 8)
	g.Effects[4*s.id*12+22] = byte(eff >> 16)
	g.Effects[4*s.id*12+23] = byte(eff >> 24)

	g.Effects[4*s.id*12+24] = byte(eff >> 0)
	g.Effects[4*s.id*12+25] = byte(eff >> 8)
	g.Effects[4*s.id*12+26] = byte(eff >> 16)
	g.Effects[4*s.id*12+27] = byte(eff >> 24)

	g.Effects[4*s.id*12+28] = byte(eff >> 0)
	g.Effects[4*s.id*12+29] = byte(eff >> 8)
	g.Effects[4*s.id*12+30] = byte(eff >> 16)
	g.Effects[4*s.id*12+31] = byte(eff >> 24)

	g.Effects[4*s.id*12+32] = byte(eff >> 0)
	g.Effects[4*s.id*12+33] = byte(eff >> 8)
	g.Effects[4*s.id*12+34] = byte(eff >> 16)
	g.Effects[4*s.id*12+35] = byte(eff >> 24)

	g.Effects[4*s.id*12+36] = byte(eff >> 0)
	g.Effects[4*s.id*12+37] = byte(eff >> 8)
	g.Effects[4*s.id*12+38] = byte(eff >> 16)
	g.Effects[4*s.id*12+39] = byte(eff >> 24)

	g.Effects[4*s.id*12+40] = byte(eff >> 0)
	g.Effects[4*s.id*12+41] = byte(eff >> 8)
	g.Effects[4*s.id*12+42] = byte(eff >> 16)
	g.Effects[4*s.id*12+43] = byte(eff >> 24)

	g.Effects[4*s.id*12+44] = byte(eff >> 0)
	g.Effects[4*s.id*12+45] = byte(eff >> 8)
	g.Effects[4*s.id*12+46] = byte(eff >> 16)
	g.Effects[4*s.id*12+47] = byte(eff >> 24)

	s.dirtyEffect = false
}

func (g *GLData) UpdateObject(s *Sprite) {
	sx := math.Float32bits(s.fx)
	sy := math.Float32bits(s.fy)
	sxw := math.Float32bits(s.fx + s.tx)
	syh := math.Float32bits(s.fy + s.ty)

	// Exploding the updates instead of generating a new
	// byte array for each quad increases the performance from 100us
	// down to about 100ns per sprite.

	g.Vertices[4*s.id*12] = byte(sxw >> 0)
	g.Vertices[4*s.id*12+1] = byte(sxw >> 8)
	g.Vertices[4*s.id*12+2] = byte(sxw >> 16)
	g.Vertices[4*s.id*12+3] = byte(sxw >> 24)

	g.Vertices[4*s.id*12+4] = byte(sy >> 0)
	g.Vertices[4*s.id*12+5] = byte(sy >> 8)
	g.Vertices[4*s.id*12+6] = byte(sy >> 16)
	g.Vertices[4*s.id*12+7] = byte(sy >> 24)

	g.Vertices[4*s.id*12+8] = byte(sxw >> 0)
	g.Vertices[4*s.id*12+9] = byte(sxw >> 8)
	g.Vertices[4*s.id*12+10] = byte(sxw >> 16)
	g.Vertices[4*s.id*12+11] = byte(sxw >> 24)

	g.Vertices[4*s.id*12+12] = byte(syh >> 0)
	g.Vertices[4*s.id*12+13] = byte(syh >> 8)
	g.Vertices[4*s.id*12+14] = byte(syh >> 16)
	g.Vertices[4*s.id*12+15] = byte(syh >> 24)

	g.Vertices[4*s.id*12+16] = byte(sx >> 0)
	g.Vertices[4*s.id*12+17] = byte(sx >> 8)
	g.Vertices[4*s.id*12+18] = byte(sx >> 16)
	g.Vertices[4*s.id*12+19] = byte(sx >> 24)

	g.Vertices[4*s.id*12+20] = byte(sy >> 0)
	g.Vertices[4*s.id*12+21] = byte(sy >> 8)
	g.Vertices[4*s.id*12+22] = byte(sy >> 16)
	g.Vertices[4*s.id*12+23] = byte(sy >> 24)

	g.Vertices[4*s.id*12+24] = byte(sxw >> 0)
	g.Vertices[4*s.id*12+25] = byte(sxw >> 8)
	g.Vertices[4*s.id*12+26] = byte(sxw >> 16)
	g.Vertices[4*s.id*12+27] = byte(sxw >> 24)

	g.Vertices[4*s.id*12+28] = byte(syh >> 0)
	g.Vertices[4*s.id*12+29] = byte(syh >> 8)
	g.Vertices[4*s.id*12+30] = byte(syh >> 16)
	g.Vertices[4*s.id*12+31] = byte(syh >> 24)

	g.Vertices[4*s.id*12+32] = byte(sx >> 0)
	g.Vertices[4*s.id*12+33] = byte(sx >> 8)
	g.Vertices[4*s.id*12+34] = byte(sx >> 16)
	g.Vertices[4*s.id*12+35] = byte(sx >> 24)

	g.Vertices[4*s.id*12+36] = byte(syh >> 0)
	g.Vertices[4*s.id*12+37] = byte(syh >> 8)
	g.Vertices[4*s.id*12+38] = byte(syh >> 16)
	g.Vertices[4*s.id*12+39] = byte(syh >> 24)

	g.Vertices[4*s.id*12+40] = byte(sx >> 0)
	g.Vertices[4*s.id*12+41] = byte(sx >> 8)
	g.Vertices[4*s.id*12+42] = byte(sx >> 16)
	g.Vertices[4*s.id*12+43] = byte(sx >> 24)

	g.Vertices[4*s.id*12+44] = byte(sy >> 0)
	g.Vertices[4*s.id*12+45] = byte(sy >> 8)
	g.Vertices[4*s.id*12+46] = byte(sy >> 16)
	g.Vertices[4*s.id*12+47] = byte(sy >> 24)
	s.dirty = false
}

func (g *GLData) UpdateUV(s *Sprite) {
	u1 := math.Float32bits(s.Texture.U1)
	u2 := math.Float32bits(s.Texture.U2)
	v1 := math.Float32bits(s.Texture.V1)
	v2 := math.Float32bits(s.Texture.V2)

	g.Uvs[4*s.id*12] = byte(u2 >> 0)
	g.Uvs[4*s.id*12+1] = byte(u2 >> 8)
	g.Uvs[4*s.id*12+2] = byte(u2 >> 16)
	g.Uvs[4*s.id*12+3] = byte(u2 >> 24)

	g.Uvs[4*s.id*12+4] = byte(v2 >> 0)
	g.Uvs[4*s.id*12+5] = byte(v2 >> 8)
	g.Uvs[4*s.id*12+6] = byte(v2 >> 16)
	g.Uvs[4*s.id*12+7] = byte(v2 >> 24)

	g.Uvs[4*s.id*12+8] = byte(u2 >> 0)
	g.Uvs[4*s.id*12+9] = byte(u2 >> 8)
	g.Uvs[4*s.id*12+10] = byte(u2 >> 16)
	g.Uvs[4*s.id*12+11] = byte(u2 >> 24)

	g.Uvs[4*s.id*12+12] = byte(v1 >> 0)
	g.Uvs[4*s.id*12+13] = byte(v1 >> 8)
	g.Uvs[4*s.id*12+14] = byte(v1 >> 16)
	g.Uvs[4*s.id*12+15] = byte(v1 >> 24)

	g.Uvs[4*s.id*12+16] = byte(u1 >> 0)
	g.Uvs[4*s.id*12+17] = byte(u1 >> 8)
	g.Uvs[4*s.id*12+18] = byte(u1 >> 16)
	g.Uvs[4*s.id*12+19] = byte(u1 >> 24)

	g.Uvs[4*s.id*12+20] = byte(v2 >> 0)
	g.Uvs[4*s.id*12+21] = byte(v2 >> 8)
	g.Uvs[4*s.id*12+22] = byte(v2 >> 16)
	g.Uvs[4*s.id*12+23] = byte(v2 >> 24)

	g.Uvs[4*s.id*12+24] = byte(u2 >> 0)
	g.Uvs[4*s.id*12+25] = byte(u2 >> 8)
	g.Uvs[4*s.id*12+26] = byte(u2 >> 16)
	g.Uvs[4*s.id*12+27] = byte(u2 >> 24)

	g.Uvs[4*s.id*12+28] = byte(v1 >> 0)
	g.Uvs[4*s.id*12+29] = byte(v1 >> 8)
	g.Uvs[4*s.id*12+30] = byte(v1 >> 16)
	g.Uvs[4*s.id*12+31] = byte(v1 >> 24)

	g.Uvs[4*s.id*12+32] = byte(u1 >> 0)
	g.Uvs[4*s.id*12+33] = byte(u1 >> 8)
	g.Uvs[4*s.id*12+34] = byte(u1 >> 16)
	g.Uvs[4*s.id*12+35] = byte(u1 >> 24)

	g.Uvs[4*s.id*12+36] = byte(v1 >> 0)
	g.Uvs[4*s.id*12+37] = byte(v1 >> 8)
	g.Uvs[4*s.id*12+38] = byte(v1 >> 16)
	g.Uvs[4*s.id*12+39] = byte(v1 >> 24)

	g.Uvs[4*s.id*12+40] = byte(u1 >> 0)
	g.Uvs[4*s.id*12+41] = byte(u1 >> 8)
	g.Uvs[4*s.id*12+42] = byte(u1 >> 16)
	g.Uvs[4*s.id*12+43] = byte(u1 >> 24)

	g.Uvs[4*s.id*12+44] = byte(v2 >> 0)
	g.Uvs[4*s.id*12+45] = byte(v2 >> 8)
	g.Uvs[4*s.id*12+46] = byte(v2 >> 16)
	g.Uvs[4*s.id*12+47] = byte(v2 >> 24)
	s.dirtyUvs = false
}
