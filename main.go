package main

import (
	"game/game"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("CombaThrone", true, false)
	// window.SetQuality(1.75, window.FilterPoint)
	// window.SetMode(window.ModeFullscreenBorderless)
	window.SetTargetFPS(0)

	game.LoadAudio()
	game.InitCharacters()
	game.InitScene()
	for window.KeepOpen() {
		game.UpdateAudio()
		game.UpdateScene()
		game.UpdateDebug()
	}
}
