package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

const (
	wMaxInvFPS = 1 / 60.0
)

type Game struct {
	ids      int
	idLock   sync.Mutex
	images   *glutil.Images
	fps      *debug.FPS
	glc      gl.Context
	lastTS   time.Time
	frameDt  float64
	elapsed  float64
	uEffect  gl.Uniform
	uTime    gl.Uniform
	uModel   gl.Uniform
	uView    gl.Uniform
	uProj    gl.Uniform
	touchX   float32
	touchY   float32
	lastX    int
	lastY    int
	size     size.Event
	sizePrev size.Event
	//objects map[int]Object
	objects  []*Sprite
	program  gl.Program
	projf    []float32
	viewf    []float32
	tex      Textures
	tiles    []TileSet
	initDone bool
}

func (g *Game) Init(glctx gl.Context) {
	g.glc = glctx

	var err error
	g.program, err = glutil.CreateProgram(g.glc, vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}

	if g.size.HeightPx == 0 {
		g.size.WidthPx = 1080
		g.size.HeightPx = 2000

	}
	g.glc.Viewport(0, 0, g.size.WidthPx, g.size.HeightPx)

	// g.uEffect = g.glc.GetUniformLocation(g.program, "effect")
	// g.uTime = g.glc.GetUniformLocation(g.program, "uTime")

	rand.Seed(time.Now().Unix())

	//g.objects = make(map[int]Object)

	//  g.font = &Font{}
	//  g.font.Init(g)

	//  for i := 1; i < 8; i++ {
	//  	tSet := &TileSet{}
	//  	tSet.Init(1.0, 4, i, 12, float64(100+i*30), g)
	//  	g.AddObjects(tSet)
	//  }

	g.glc.BlendFuncSeparate(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA, gl.ONE, gl.ONE)
	g.glc.FrontFace(gl.CCW)
	g.glc.Enable(gl.BLEND)
	g.glc.Disable(gl.DEPTH_TEST)
	g.glc.Disable(gl.SCISSOR_TEST)
	g.glc.Enable(gl.CULL_FACE)
	g.glc.CullFace(gl.BACK)

	g.tex = Textures{}
	if err = g.tex.Load("packed.png", "packed.json", g); err != nil {
		panic(err)
	}

	s2 := &Sprite{}
	s2.Init(0, 0, 0, 1.0, 1.0, "bg", g)
	s2.scalex = float32(g.size.WidthPx) / s2.Texture.Width
	s2.scaley = float32(g.size.HeightPx) / s2.Texture.Height
	s2.dirty = true
	g.AddObjects(s2)

	for i := 1; i <= 15; i++ {
		ts := TileSet{}
		ts.Init(1.0, 4, i, 20, 320, g)
		ts.SetSpeed(4)
		g.tiles = append(g.tiles, ts)
	}

	g.tex.AddText("bintris", g.X(225), g.Y(290), 0.0,
		g.SX(6),
		g.SY(4.6),
		EffectMetaballsBlue)

	g.tex.AddText("Score:", g.X(15), g.Y(295), 0.0,
		g.SX(8),
		g.SY(10),
		EffectMetaballsBlue)
	g.tex.AddText("0000", g.X(70), g.Y(295), 0.0,
		g.SX(8),
		g.SY(10),
		EffectMetaballsBlue)

	g.tex.AddText("Time:", g.X(120), g.Y(295), 0.0,
		g.SX(8),
		g.SY(10),
		EffectMetaballsBlue)
	g.tex.AddText("0000", g.X(165), g.Y(295), 0.0,
		g.SX(8),
		g.SY(10),
		EffectMetaballsBlue)

	g.tex.Init()
	g.tex.SetResolution()

	//g.images = glutil.NewImages(g.glc)
	//	g.fps = debug.NewFPS(g.images)
	g.lastTS = time.Now()
	g.initDone = true
}

func (g *Game) Stop() {
	g.glc.DeleteProgram(g.program)
	//g.tex.Cleanup()
	// g.fps.Release()
	// g.images.Release()
}

func (g *Game) Draw() {
	dt := time.Since(g.lastTS).Seconds()
	g.frameDt += dt
	g.lastTS = time.Now()

	g.glc.ClearColor(0.0, 0.0, 0.0, 0.0)
	g.glc.Clear(gl.COLOR_BUFFER_BIT)

	c := 0
	hidden := []*TileSet{}
	for i := range g.tiles {
		if !g.tiles[i].hidden {
			c++
		} else {
			hidden = append(hidden, &g.tiles[i])
		}
	}
	x := 0
	for i := c; i <= 8; i++ {
		fmt.Printf("ADD: %v\n", i)
		hidden[x].Reset()
		x++
	}

	for {
		if g.frameDt >= wMaxInvFPS {
			g.elapsed += wMaxInvFPS
			for i := range g.tiles {
				// TBD: Only draw visible tiles
				if !g.tiles[i].Hidden() {
					g.tiles[i].Update(wMaxInvFPS)
				}
			}
			for k := range g.objects {
				if g.objects[k].Hidden() {
					continue
				}
				g.objects[k].Update(float64(wMaxInvFPS))
			}
		} else {
			break
		}

		g.frameDt -= wMaxInvFPS
	}

	for k := range g.objects {
		if g.objects[k].Hidden() {
			continue
		}
		//g.objects[k].Draw(float64(wMaxInvFPS))
	}

	g.tex.Draw()
	g.tex.Update()
	//g.fps.Draw(g.size)
}

func (g *Game) Click(x, y float32) {
	// Make sure we don't generate too many clicks.
	if g.lastX == int(x) && g.lastY == int(y) {
		return
	}
	g.lastX = int(x)
	g.lastY = int(y)

	y = float32(g.size.HeightPx) - y

	for i, c := range g.tiles {
		if !c.hidden {
			if float32(x) > c.X && float32(x) < c.X+(float32(c.tileWidth)) &&
				float32(y) > c.Y && float32(y) < c.Y+(float32(c.tileHeight)) {
				g.tiles[i].Click(x, y)
			}
		}
	}
}

func (g *Game) Resize(e size.Event) {
	if g.sizePrev.WidthPx == e.WidthPx &&
		g.sizePrev.HeightPx == e.HeightPx {
		return
	}

	if g.glc != nil {
		g.glc.Viewport(0, 0, g.size.WidthPx, g.size.HeightPx)
	}

	g.sizePrev = g.size
	g.size = e

	if g.initDone {
		g.tex.SetResolution()
	}

	// Resize objects.
	for i := range g.objects {
		g.objects[i].Resize()
	}
	for i := range g.tiles {
		g.tiles[i].Resize()
	}
}

func (g *Game) AddObjects(obj ...*Sprite) {
	for i := range obj {
		//if _, ok := g.objects[obj[i].GetID()]; !ok {
		//g.objects[obj[i].GetID()] = obj[i]
		g.objects = append(g.objects, obj[i])
		//	}
	}
}

func (g *Game) DeleteObject(obj Sprite) {
	for i := range g.objects {
		if obj.id == g.objects[i].id {
			g.objects[i].Delete()
			//		delete(g.objects, i)
		}
	}
}

func (g *Game) NewID() int {
	g.idLock.Lock()
	defer g.idLock.Unlock()

	g.ids++
	return g.ids - 1
}

// Y calculates absolute position on a virtual sized viewport
func (g *Game) Y(y float32) float32 {
	vy := float32(320)
	pp := float32(g.size.HeightPx) / vy
	return y * pp
}

// VY returns the virtual position for a given Y
func (g *Game) VY(y float32) float32 {
	vy := float32(320)
	pp := float32(g.size.HeightPx) / vy
	return y / pp
}

// SY converts to relative scale of window size
// TBD: These are inverted meaning 3 > 2
func (g *Game) SY(y float32) float32 {
	return float32(g.size.HeightPx) / (y * 100)
}

// SX converts to relative scale of window size
func (g *Game) SX(x float32) float32 {
	return float32(g.size.WidthPx) / (x * 100)
}

// VX returns the virtual position for a given Y
func (g *Game) VX(x float32) float32 {
	vx := float32(320)
	pp := float32(g.size.WidthPx) / vx
	return x / pp
}

// X calculates absolute position on a virtual sized viewport
func (g *Game) X(x float32) float32 {
	vx := float32(320)
	pp := float32(g.size.WidthPx) / vx
	return x * pp
}
