package main

import (
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
	uTime    gl.Uniform
	touchX   float32
	touchY   float32
	lastX    int
	lastY    int
	size     size.Event
	sizePrev size.Event
	objects  []*Sprite
	program  gl.Program
	projf    []float32
	viewf    []float32
	tex      Textures
	tiles    []TileSet
	mode     Mode
	menu     Menu
}

func (g *Game) Init(glctx gl.Context) {
	g.glc = glctx

	var err error
	g.program, err = glutil.CreateProgram(g.glc, vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}

	g.uTime = g.glc.GetUniformLocation(g.program, "uTime")

	rand.Seed(time.Now().Unix())

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
	s2.Init(0.0, 0.0, 0, 1.0, 1.0, "bg", g)
	s2.dirty = true
	g.AddObjects(s2)

	for i := 1; i <= 15; i++ {
		ts := TileSet{}
		ts.Init(4, i, g)
		ts.Hide()
		g.tiles = append(g.tiles, ts)
	}

	g.menu.Init(g)
	g.mode.Init(g)
	g.tex.Init()

	g.images = glutil.NewImages(g.glc)
	g.fps = debug.NewFPS(g.images)
	g.lastTS = time.Now()
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

	for {
		if g.frameDt >= wMaxInvFPS {
			g.elapsed += wMaxInvFPS
			for i := range g.tiles {
				if !g.tiles[i].hidden {
					g.tiles[i].Update(wMaxInvFPS)
				}
			}
			g.mode.Update(wMaxInvFPS)
		} else {
			break
		}

		g.frameDt -= wMaxInvFPS
	}

	g.glc.Uniform1f(g.uTime, float32(g.elapsed))
	g.tex.Draw()
	g.tex.Update()

	//	g.fps.Draw(g.size)
}

func (g *Game) Click(sz size.Event, x, y float32) {
	x /= float32(sz.WidthPx)
	y /= float32(sz.HeightPx)
	y = 1 - y

	if g.menu.hidden {
		for i, c := range g.tiles {
			if !c.hidden {
				// Offset Y a bit to have a bit off click-free area between tiles
				if float32(x) > c.tile.fx && float32(x) < c.tile.fx+0.822 &&
					float32(y) > c.tile.fy+0.01 && float32(y) < c.tile.fy+0.09 {
					g.tiles[i].Click(x, y)
				}
			}
		}
	} else {
		g.menu.KeyDown(x, y)
	}
}

func (g *Game) AddObjects(obj ...*Sprite) {
	for i := range obj {
		g.objects = append(g.objects, obj[i])
	}
}

func (g *Game) NewID() int {
	g.idLock.Lock()
	defer g.idLock.Unlock()

	g.ids++
	return g.ids - 1
}
