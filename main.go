package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 800
	screenHeight = 600
	windowTitle  = "DotDot"
)

func main() {
	game := NewGame()
	if err := ebiten.Run(game.update, screenWidth, screenHeight, 1, windowTitle); err != nil {
		log.Fatal(err)
	}
}
