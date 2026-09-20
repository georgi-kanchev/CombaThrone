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

	game.LoadAudio()
	game.InitUI() // needs to be before all other img loads since it relies on IDs
	game.InitCharacters()
	game.InitScene()
	for window.KeepOpen() {
		game.UpdateAudio()

		if !game.InGame {
			game.UpdateTitleScreen()
			continue
		}

		game.UpdateScene()
		game.UpdateDebug()
	}
}
