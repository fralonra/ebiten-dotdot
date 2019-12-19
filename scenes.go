package main

import (
	"image/color"
	"strconv"
	"time"

	"github.com/fralonra/dotdot"
	scene "github.com/fralonra/ebiten-scene"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
)

const (
	timeout = 120
)

func jump() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return true
	}
	if len(inpututil.JustPressedTouchIDs()) > 0 {
		return true
	}
	return false
}

const (
	sceneKeyTitle scene.GameSceneKey = iota
	sceneKeyGame
	sceneKeyWin
	sceneKeyLost
)

var (
	sceneTitle *SceneTitle
	sceneGame  *SceneGame
	sceneLost  *SceneLost
	sceneWin   *SceneWin
)

func init() {
	sceneTitle = &SceneTitle{}
	sceneGame = &SceneGame{}
	sceneLost = &SceneLost{}
	sceneWin = &SceneWin{}
}

type SceneTitle struct{}

func (m *SceneTitle) Draw(screen *ebiten.Image) {
	texts := []string{
		"DOTDOT",
		"", "", "", "",
		"PRESS SPACE KEY",
		"",
		"OR TOUCH SCREEN",
	}
	for i, l := range texts {
		x := (screenWidth - len(l)*fontSize) / 2
		text.Draw(screen, l, arcadeFont, x, (i+4)*fontSize, color.White)
	}
}

func (m *SceneTitle) Update() {
	if jump() {
		scene.Switch(sceneKeyGame)
	}
}

func (m *SceneTitle) Start() {}

func (m *SceneTitle) Stop() {}

type GameTimer struct {
	done   chan bool
	ticker *time.Ticker
	time   int
}

func (t *GameTimer) handle() {
	for {
		select {
		case <-t.done:
			return
		case <-t.ticker.C:
			t.time++
		}
	}
}

func (t *GameTimer) stop() {
	t.ticker.Stop()
	t.done <- true
}

type SceneGame struct {
	cursorPos     [2]float64
	game          *dotdot.DotDot
	capturedCount int
	timer         *GameTimer
}

func (m *SceneGame) Draw(screen *ebiten.Image) {
	dots := m.game.GetDots()
	m.capturedCount = 0
	for _, dot := range dots {
		pos := dot.Pos()
		size := dot.Size()
		col := dot.Color()
		ebitenutil.DrawRect(screen, pos[0], pos[1], size, size, col)
		if dot.Captured() {
			m.capturedCount++
			ebitenutil.DrawLine(screen, pos[0], pos[1], m.cursorPos[0], m.cursorPos[1], color.RGBA{
				R: col.(color.RGBA).R,
				G: col.(color.RGBA).G,
				B: col.(color.RGBA).B,
				A: uint8(1 - dot.Distance()/m.game.Distance*200),
			})
		}
	}
	texts := []string{
		"Time Used:" + strconv.Itoa(m.timer.time),
		"Dots Left:" + strconv.Itoa(m.game.Number-m.capturedCount),
	}
	for i, l := range texts {
		x := smallFontSize
		text.Draw(screen, l, smallArcadeFont, x, (i+1)*(smallFontSize+10), color.White)
	}
}

func (m *SceneGame) Update() {
	if m.capturedCount == m.game.Number {
		scene.SwitchWith(sceneKeyWin, func(g scene.GameScene) {
			if scene, ok := g.(*SceneWin); ok {
				scene.time = m.timer.time
			}
		})
		return
	}

	mx, my := ebiten.CursorPosition()
	m.cursorPos = [2]float64{float64(mx), float64(my)}
	m.game.Update(m.cursorPos)

	if m.timer.time > timeout {
		scene.Switch(sceneKeyLost)
		return
	}
}

func (m *SceneGame) Start() {
	m.capturedCount = 0
	m.timer = &GameTimer{
		done:   make(chan bool),
		ticker: time.NewTicker(1 * time.Second),
	}
	go m.timer.handle()
	m.game.Start()
}

func (m *SceneGame) Stop() {
	m.timer.stop()
}

type SceneLost struct{}

func (m *SceneLost) Draw(screen *ebiten.Image) {
	texts := []string{
		"",
		"GAME OVER!",
		"",
		"RUN OUT OF TIME",
	}
	for i, l := range texts {
		x := (screenWidth - len(l)*fontSize) / 2
		text.Draw(screen, l, arcadeFont, x, (i+4)*fontSize, color.White)
	}
}

func (m *SceneLost) Update() {
	if jump() {
		scene.Switch(sceneKeyTitle)
	}
}

func (m *SceneLost) Start() {}

func (m *SceneLost) Stop() {}

type SceneWin struct {
	time int
}

func (m *SceneWin) Draw(screen *ebiten.Image) {
	texts := []string{
		"",
		"GAME WIN!",
		"",
		"YOU USED " + strconv.Itoa(m.time) + "S",
	}
	for i, l := range texts {
		x := (screenWidth - len(l)*fontSize) / 2
		text.Draw(screen, l, arcadeFont, x, (i+4)*fontSize, color.White)
	}
}

func (m *SceneWin) Update() {
	if jump() {
		scene.Switch(sceneKeyTitle)
	}
}

func (m *SceneWin) Start() {}

func (m *SceneWin) Stop() {}
