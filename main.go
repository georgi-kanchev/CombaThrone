package main

import (
	"game/game"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("CombaThrone", true, false)
	// window.SetQuality(2, window.FilterPoint)
	// window.SetMode(window.ModeFullscreenBorderless)
	window.SetTargetFPS(0)

	game.InitCharacters()
	game.InitScene()
	for window.KeepOpen() {
		game.UpdateScene()
		game.UpdateUnits()
		game.UpdateDebug()
	}
}
