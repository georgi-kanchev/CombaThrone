package main

import (
	"game/game"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("CombaThrone", false, true)

	window.SetTargetFPS(0)

	game.InitCharacters()

	game.InitScene()
	for window.KeepOpen() {
		game.UpdateScene()
		game.UpdateUnits()
		game.UpdateDebug()
	}
}
