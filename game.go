package main

import (
	"github.com/fralonra/dotdot"
	scene "github.com/fralonra/ebiten-scene"
	"github.com/hajimehoshi/ebiten"
)

type Game struct {
	core *dotdot.DotDot
}

func (g *Game) update(screen *ebiten.Image) error {
	return scene.Update(screen)
}

func NewGame() (g *Game) {
	g = &Game{}
	g.core = dotdot.New(screenWidth, screenHeight)
	sceneGame.game = g.core
	scene.Set(sceneKeyTitle, sceneTitle)
	scene.Set(sceneKeyGame, sceneGame)
	scene.Set(sceneKeyLost, sceneLost)
	scene.Set(sceneKeyWin, sceneWin)
	scene.Switch(sceneKeyTitle)
	return
}
