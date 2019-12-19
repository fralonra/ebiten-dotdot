package main

import (
	"image/color"
	"log"
	"strconv"
	"time"

	"github.com/fralonra/dotdot"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

type Mode int

const (
	ModeTitle Mode = iota
	ModeGame
	ModeGameOver
	ModeGameWin
)

const (
	screenWidth   = 800
	screenHeight  = 600
	timeout       = 20
	tileSize      = 32
	fontSize      = 32
	smallFontSize = fontSize / 3
)

var (
	arcadeFont      font.Face
	smallArcadeFont font.Face
)

func init() {
	tt, err := truetype.Parse(fonts.ArcadeN_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	arcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	smallArcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    smallFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

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

type Game struct {
	cursorPos [2]float64
	game      *dotdot.DotDot
	mode      Mode
	timer     *GameTimer
}

func (g *Game) draw(screen *ebiten.Image) {
	var texts []string
	switch g.mode {
	case ModeTitle:
		texts = []string{
			"DOTDOT",
			"", "", "", "",
			"PRESS SPACE KEY",
			"",
			"OR TOUCH SCREEN",
		}
	case ModeGame:
		if g.timer.time > timeout {
			g.loseGame()
			return
		}
		var capturedCount int
		dots := g.game.GetDots()
		for _, dot := range dots {
			pos := dot.Pos()
			size := dot.Size()
			col := dot.Color()
			ebitenutil.DrawRect(screen, pos[0], pos[1], size, size, col)
			if dot.Captured() {
				capturedCount++
				ebitenutil.DrawLine(screen, pos[0], pos[1], g.cursorPos[0], g.cursorPos[1], color.RGBA{
					R: col.(color.RGBA).R,
					G: col.(color.RGBA).G,
					B: col.(color.RGBA).B,
					A: uint8(1 - dot.Distance()/g.game.Distance*200),
				})
			}
		}
		leftCount := g.game.Number - capturedCount
		if leftCount == 0 {
			g.winGame()
			return
		}
		texts = []string{
			"Time Used:" + strconv.Itoa(g.timer.time),
			"Dots Left:" + strconv.Itoa(leftCount),
		}
		for i, l := range texts {
			x := smallFontSize
			text.Draw(screen, l, smallArcadeFont, x, (i+1)*(smallFontSize+10), color.White)
		}
		return
	case ModeGameOver:
		texts = []string{
			"",
			"GAME OVER!",
			"",
			"RUN OUT OF TIME",
		}
	case ModeGameWin:
		texts = []string{
			"",
			"GAME WIN!",
			"",
			"YOU USED " + strconv.Itoa(g.timer.time) + "S",
		}
	}
	for i, l := range texts {
		x := (screenWidth - len(l)*fontSize) / 2
		text.Draw(screen, l, arcadeFont, x, (i+4)*fontSize, color.White)
	}
}

func (g *Game) handle() {
	switch g.mode {
	case ModeTitle:
		if jump() {
			g.startNewGame()
		}
	case ModeGame:
		g.handleGameUpdate()
	case ModeGameOver, ModeGameWin:
		if jump() {
			g.mode = ModeTitle
		}
	}
}

func (g *Game) handleGameUpdate() {
	go g.timer.handle()

	mx, my := ebiten.CursorPosition()
	g.cursorPos = [2]float64{float64(mx), float64(my)}
	g.game.Update(g.cursorPos)
}

func (g *Game) loseGame() {
	g.timer.stop()
	g.mode = ModeGameOver
}

func (g *Game) startNewGame() {
	g.mode = ModeGame
	g.timer = &GameTimer{
		done:   make(chan bool),
		ticker: time.NewTicker(1 * time.Second),
	}
	g.game.Start()
}

func (g *Game) update(screen *ebiten.Image) error {
	g.handle()
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	g.draw(screen)
	return nil
}

func (g *Game) winGame() {
	g.timer.stop()
	g.mode = ModeGameWin
}

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

func NewGame() {
	game := &Game{}
	game.game = dotdot.New(screenWidth, screenHeight)
	game.mode = ModeTitle
	if err := ebiten.Run(game.update, screenWidth, screenHeight, 1, "DotDot"); err != nil {
		log.Fatal(err)
	}
}
