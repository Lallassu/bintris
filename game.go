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
	ids     int
	idLock  sync.Mutex
	images  *glutil.Images
	fps     *debug.FPS
	glc     gl.Context
	lastTS  time.Time
	frameDt float64
	elapsed float64
	uEffect gl.Uniform
	uTime   gl.Uniform
	uModel  gl.Uniform
	uView   gl.Uniform
	uProj   gl.Uniform
	touchX  float32
	touchY  float32
	lastX   int
	lastY   int
	size    size.Event
	//objects map[int]Object
	objects []Sprite
	program gl.Program
	projf   []float32
	viewf   []float32
	tex     Textures
	//font    *Font
}

func (g *Game) Init(glctx gl.Context) {
	g.glc = glctx

	var err error
	g.program, err = glutil.CreateProgram(g.glc, vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}

	g.glc.Viewport(0, 0, g.size.WidthPx, g.size.HeightPx)

	g.uEffect = g.glc.GetUniformLocation(g.program, "effect")
	g.uTime = g.glc.GetUniformLocation(g.program, "uTime")

	rand.Seed(time.Now().Unix())

	//g.objects = make(map[int]Object)

	//  g.font = &Font{}
	//  g.font.Init(g)

	//  for i := 1; i < 8; i++ {
	//  	tSet := &TileSet{}
	//  	tSet.Init(1.0, 4, i, 12, float64(100+i*30), g)
	//  	g.AddObjects(tSet)
	//  }
	//g.AddObjects(g.font.AddText("bintris", 210, 311, 0.5, 0.8, EffectMetaballsBlue)...)

	g.glc.BlendFunc(gl.BLEND_SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	g.glc.Enable(gl.CULL_FACE)
	g.glc.FrontFace(gl.CCW)
	g.glc.Enable(gl.BLEND)
	g.glc.Disable(gl.DEPTH_TEST)
	g.glc.Disable(gl.SCISSOR_TEST)
	g.glc.CullFace(gl.BACK)
	//g.glc.DepthFunc(gl.LESS)

	g.tex = Textures{}
	if err = g.tex.Load("packed.png", "packed.json", g); err != nil {
		panic(err)
	}

	s2 := Sprite{}
	s2.Init(0, 0, 0, 1.0, 1.0, "bg3", g)
	s2.scalex = float32(g.size.WidthPx) / s2.Texture.Width
	s2.scaley = float32(g.size.HeightPx) / s2.Texture.Height
	fmt.Printf("X: %v Y: %v\n", s2.scalex, s2.scaley)
	s2.dirty = true
	fmt.Printf("s2.scale: %v, height: %v text: %v\n", s2.scalex, g.size.HeightPx, s2.Texture.Height)
	fmt.Printf("s2.scale: %v, width: %v text: %v\n", s2.scaley, g.size.WidthPx, s2.Texture.Width)
	g.AddObjects(s2)

	for i := 0; i < 10; i++ {
		s2 := Sprite{}
		s2.Init(float32(i*20), float32(i*20), 0, 1.0, 1.0, "4", g)
		g.AddObjects(s2)
	}
	g.tex.Init()

	//  g.images = glutil.NewImages(g.glc)
	//  g.fps = debug.NewFPS(g.images)
	g.lastTS = time.Now()

}

func (g *Game) Stop() {
	g.glc.DeleteProgram(g.program)
	//g.fps.Release()
	//  g.images.Release()
}

func (g *Game) Draw() {
	dt := time.Since(g.lastTS).Seconds()
	g.frameDt += dt
	g.lastTS = time.Now()

	g.glc.ClearColor(0.1, 0.1, 0.1, 1.0)
	g.glc.Clear(gl.COLOR_BUFFER_BIT) //| gl.DEPTH_BUFFER_BIT)

	for {
		if g.frameDt >= wMaxInvFPS {
			g.elapsed += wMaxInvFPS
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
	g.tex.SetResolution()
	g.tex.Update()
	g.tex.Draw()
	//g.fps.Draw(g.size)
}

func (g *Game) Click(x, y float32) {
	// Make sure we don't generate too many clicks.
	if g.lastX == int(x) && g.lastY == int(y) {
		return
	}
	g.lastX = int(x)
	g.lastY = int(y)
	//g.AddObjects(g.font.AddText(fmt.Sprintf("%vx%v", g.size.WidthPx, g.size.HeightPx), 100, 100, 1.0, 0.3, EffectMetaballsBlue)...)

	ah := float32(g.size.HeightPx) / 320
	aw := float32(g.size.WidthPx) / 320
	y = 320 - (y / ah)
	x = x / aw

	// for i := range g.objects {
	// 	if !c.deleted {
	// 		if float64(x) > c.X && float64(x) < c.X+(float64(c.sizex)*float64(c.scale)) &&
	// 			float64(y) < c.Y && float64(y) > c.Y-(float64(c.sizey)*float64(c.scale)) {
	// 			c.Click(x, y)
	// 		}
	// 	}
	// }
}

func (g *Game) Resize(e size.Event) {
	g.size = e
	g.tex.SetResolution()
	fmt.Printf("REsize: %v\n", e)
}

func (g *Game) AddObjects(obj ...Sprite) {
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
