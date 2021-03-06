package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/audio/al"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

const (
	version    = "1.0"
	wMaxInvFPS = 1 / 60.0
)

type Game struct {
	visible     float32
	playIds     int
	menuIds     int
	pulse       float32
	idLock      sync.Mutex
	glc         gl.Context
	lastTS      time.Time
	frameDt     float64
	elapsed     float64
	uTime       gl.Uniform
	uPulse      gl.Uniform
	uTouchX     gl.Uniform
	uTouchY     gl.Uniform
	touchX      float32
	touchY      float32
	lastX       int
	lastY       int
	size        size.Event
	sizePrev    size.Event
	objectsMenu []*Sprite
	objectsPlay []*Sprite
	program     gl.Program
	projf       []float32
	viewf       []float32
	tex         Textures
	tiles       []TileSet
	mode        Mode
	about       About
	howto       HowToPlay
	menu        Menu
	scoreboard  Scoreboard
	bg          *Sprite
	backBg      *Sprite
	menuBg      *Sprite
	clicked     time.Time
	sound       Sound
	glData      *GLData
	tileIds     int
}

func (g *Game) Init(glctx gl.Context) {
	g.glc = glctx

	var err error
	g.program, err = glutil.CreateProgram(g.glc, vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}
	g.glc.UseProgram(g.program)

	g.uTime = g.glc.GetUniformLocation(g.program, "uTime")
	g.uPulse = g.glc.GetUniformLocation(g.program, "uPulse")
	g.uTouchX = g.glc.GetUniformLocation(g.program, "uTouchX")
	g.uTouchY = g.glc.GetUniformLocation(g.program, "uTouchY")

	rand.Seed(time.Now().Unix())

	g.playIds = 707
	g.glData = &GLData{}
	g.glData.Init(g, 480000)

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

	g.backBg = &Sprite{}
	g.backBg.Init(0.0, 0.0, -0.0001, 1.0, 1.0, "blank", g, SpriteMenu)
	g.backBg.ChangeEffect(EffectMenu)

	g.menuBg = &Sprite{}
	g.menuBg.Init(0.0, 0.0, -0.001, 1.0, 1.0, "menubg", g, SpriteMenu)
	g.menuBg.ChangeEffect(EffectNone)

	g.bg = &Sprite{}
	g.bg.Init(0.0, 0.0, 0, 1.0, 1.0, "bg", g, SpritePlay)
	g.bg.ChangeEffect(EffectBg)

	for i := 1; i <= 15; i++ {
		ts := TileSet{}
		ts.Init(4, i, g)
		g.tiles = append(g.tiles, ts)
	}

	g.sound = Sound{}
	g.sound.Init()
	g.sound.Load("main", "sounds/main.wav", al.FormatMono16, 9000)
	g.sound.Load("gameover", "sounds/gameover.wav", al.FormatMono8, 11025)
	g.sound.Load("click", "sounds/click.wav", al.FormatMono8, 11025)
	g.sound.Load("tile", "sounds/tile.wav", al.FormatMono8, 11025)
	g.sound.Load("bitrot", "sounds/bitrot.wav", al.FormatMono16, 44100)
	g.sound.Load("win", "sounds/win.wav", al.FormatMono16, 9000)

	g.scoreboard.Init(g)
	g.menu.Init(g)
	g.mode.Init(g)
	g.about.Init(g)
	g.howto.Init(g)
	g.tex.Init()

	g.glData.Enable(SpriteMenu)
	g.menu.Show()
	g.lastTS = time.Now()

	fmt.Printf("Menu: %v, Play: %v\n", g.menuIds, g.playIds)
	g.sound.Play("main")
}

func (g *Game) Stop() {
	g.glc.DeleteProgram(g.program)
	g.sound.Close()
	g.tex.Cleanup()
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
			g.mode.Update(wMaxInvFPS)
			if g.mode.started {
				div := float32(4.5)
				if g.visible/div != g.pulse {
					if g.pulse < g.visible/div {
						g.pulse += 0.007
					} else if g.pulse > g.visible/div {
						g.pulse -= 0.007
					}
				}
			}
		} else {
			break
		}

		g.frameDt -= wMaxInvFPS
	}

	// Avoid flickering tiles
	for i := range g.tiles {
		if !g.tiles[i].hidden {
			g.tiles[i].Update(dt) //wMaxInvFPS)
		}
	}

	g.glc.Uniform1f(g.uTime, float32(g.elapsed))
	g.glc.Uniform1f(g.uPulse, g.pulse)
	g.glc.Uniform1f(g.uTouchX, g.touchX)
	g.glc.Uniform1f(g.uTouchY, g.touchY)
	g.glData.Draw()
	g.glData.Update(g.glData.sType)
}

func (g *Game) GameOver() {
	for i := range g.tiles {
		g.tiles[i].GameOver()
	}
	g.bg.ChangeEffect(EffectGameOver)
}

func (g *Game) Reset() {
	g.bg.ChangeEffect(EffectBg)
	g.backBg.Show()
	g.menuBg.Show()
	for i := range g.tiles {
		g.tiles[i].Hide()
	}

	g.mode.Hide()
	g.bg.Hide()
}

func (g *Game) Click(sz size.Event, x, y float32) {
	x /= float32(sz.WidthPx)
	y /= float32(sz.HeightPx)
	y = 1 - y
	g.touchX = x
	g.touchY = y

	if time.Since(g.clicked) < time.Duration(150*time.Millisecond) {
		return
	}
	g.clicked = time.Now()

	if !g.mode.IsGameOver() && g.mode.Started() {
		for i, c := range g.tiles {
			if !c.hidden {
				// Offset Y a bit to have a bit off click-free area between tiles
				if float32(x) > c.tile.fx && float32(x) < c.tile.fx+0.822 &&
					float32(y) > c.tile.fy+0.01 && float32(y) < c.tile.fy+0.09 {
					g.tiles[i].Click(x, y)
					break
				}
			}
		}
	} else {
		if !g.menu.Hidden() {
			g.menu.KeyDown(x, y)
		}
		if !g.scoreboard.Hidden() {
			g.scoreboard.KeyDown(x, y)
		}
		if !g.about.Hidden() {
			g.about.KeyDown(x, y)
		}
		if !g.howto.Hidden() {
			g.howto.KeyDown(x, y)
		}
	}
}

func (g *Game) AddObjects(sType SpriteType, obj ...*Sprite) {
	for i := range obj {
		if sType == SpritePlay {
			g.objectsPlay = append(g.objectsPlay, obj[i])
		} else {
			g.objectsMenu = append(g.objectsMenu, obj[i])
		}
	}
}

func (g *Game) HideAll() {
	for i := range g.objectsPlay {
		g.objectsPlay[i].Hide()
	}
	for i := range g.objectsMenu {
		g.objectsMenu[i].Hide()
	}
}

func (g *Game) NewMenuID() int {
	g.idLock.Lock()
	defer g.idLock.Unlock()

	g.menuIds++
	return g.menuIds - 1
}

func (g *Game) NewPlayID() int {
	g.idLock.Lock()
	defer g.idLock.Unlock()

	g.playIds++
	return g.playIds - 1
}

func (g *Game) NewTileID() int {
	g.idLock.Lock()
	defer g.idLock.Unlock()

	g.tileIds++
	return g.tileIds - 1
}
