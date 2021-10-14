package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/go-gl/mathgl/mgl32"
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
	objects []Object
	program gl.Program
	projf   []float32
	viewf   []float32
	font    *Font
}

func (g *Game) Init(glctx gl.Context) {
	g.glc = glctx

	var err error
	g.program, err = glutil.CreateProgram(g.glc, vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}

	g.UpdateView()
	g.glc.Viewport(0, 0, g.size.WidthPx, g.size.HeightPx)

	g.uModel = g.glc.GetUniformLocation(g.program, "model")
	g.uView = g.glc.GetUniformLocation(g.program, "view")
	g.uProj = g.glc.GetUniformLocation(g.program, "projection")
	g.uEffect = g.glc.GetUniformLocation(g.program, "effect")
	g.uTime = g.glc.GetUniformLocation(g.program, "uTime")

	rand.Seed(time.Now().Unix())

	//g.objects = make(map[int]Object)

	g.font = &Font{}
	g.font.Init(g)

	for i := 1; i < 8; i++ {
		tSet := &TileSet{}
		tSet.Init(1.0, 4, i, 12, float64(100+i*30), g)
		g.AddObjects(tSet)
	}
	g.AddObjects(g.font.AddText("bintris", 210, 311, 0.5, 0.8, EffectMetaballsBlue)...)

	s2 := &Sprite{}
	s2.Init(0, 320, 0, 1.0, "bg3.png", nil, g)
	s2.effect = EffectMetaballs
	g.AddObjects(s2)

	g.images = glutil.NewImages(g.glc)
	g.fps = debug.NewFPS(g.images)
	g.lastTS = time.Now()
}

func (g *Game) Stop() {
	g.glc.DeleteProgram(g.program)
	g.fps.Release()
	g.images.Release()
}

func (g *Game) Draw() {
	dt := time.Since(g.lastTS).Seconds()
	g.frameDt += dt
	g.lastTS = time.Now()

	g.glc.FrontFace(gl.CCW)
	g.glc.CullFace(gl.BACK)
	g.glc.Enable(gl.CULL_FACE)
	g.glc.BlendFunc(gl.BLEND_SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	g.glc.Enable(gl.BLEND)
	g.glc.Enable(gl.DEPTH_TEST)
	g.glc.DepthFunc(gl.LESS)
	g.glc.ClearColor(0, 0, 0, 0)
	g.glc.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	for {
		if g.frameDt >= wMaxInvFPS {
			g.elapsed += wMaxInvFPS
			for k := range g.objects {
				g.objects[k].Update(float64(wMaxInvFPS))
			}
		} else {
			break
		}

		g.frameDt -= wMaxInvFPS
	}

	for k := range g.objects {
		g.objects[k].Draw(float64(wMaxInvFPS))
	}

	g.glc.FrontFace(gl.CW)
	g.glc.Disable(gl.DEPTH_TEST)
	g.fps.Draw(g.size)
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

	for i := range g.objects {
		switch g.objects[i].(type) {
		case *TileSet:
			c := g.objects[i].(*TileSet)
			if !c.deleted {
				if float64(x) > c.X && float64(x) < c.X+(float64(c.sizex)*float64(c.scale)) &&
					float64(y) < c.Y && float64(y) > c.Y-(float64(c.sizey)*float64(c.scale)) {
					c.Click(x, y)
				}
			}
		}
	}
}

func (g *Game) Resize(e size.Event) {
	g.size = e
	g.UpdateProjection()
	g.UpdateView()
}

func (g *Game) UpdateProjection() {
	ww := float32(320)
	wh := float32(320)

	//a := ww / wh
	//v := float32(g.size.WidthPx) / float32(g.size.HeightPx)

	//projection := mgl32.Mat4{}
	//if v >= a {
	//	projection = mgl32.Ortho2D(0, 2*v/a*ww/2.0, 0, 2*(wh/2.0))
	//} else {
	//	projection = mgl32.Ortho2D(0, 2*ww/2.0, 0, (a/v*wh/2.0)*2)
	//}

	//aspect := float32(g.size.WidthPx) / float32(g.size.HeightPx)
	//projection := mgl32.Ortho2D(0, ww, -(ww/g.size.PixelsPerPt)/2, wh+(ww/g.size.PixelsPerPt)/2)
	fmt.Printf("> %vx%v\n", g.size.WidthPx, g.size.HeightPx)
	h := float32(320)
	if h > float32(g.size.HeightPx) {
		h = float32(g.size.HeightPx)
	}
	projection := mgl32.Ortho2D(0, ww, 0, h+10) //float32(g.size.HeightPx))
	fmt.Printf("Left: %v Right: %v Top: %v Bottom: %v\n", 0, ww, -(ww/g.size.PixelsPerPt)/2, wh+(ww/g.size.PixelsPerPt)/2)

	g.projf = projection[:]
}

func (g *Game) UpdateView() {
	view := mgl32.Translate3D(float32(0), float32(0), float32(0))
	g.viewf = view[:]
}

func (g *Game) AddObjects(obj ...Object) {
	for i := range obj {
		//if _, ok := g.objects[obj[i].GetID()]; !ok {
		//g.objects[obj[i].GetID()] = obj[i]
		g.objects = append(g.objects, obj[i])
		//	}
	}
}

func (g *Game) DeleteObject(obj Object) {
	for i := range g.objects {
		if obj.GetID() == g.objects[i].GetID() {
			g.objects[i].Delete()
			//		delete(g.objects, i)
		}
	}
}

func (g *Game) NewID() int {
	g.idLock.Lock()
	defer g.idLock.Unlock()

	g.ids++
	return g.ids
}
