package main

import (
	"game/game"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("CombaThrone", true, true)

	game.InitCharacters()

	game.InitScene()
	for window.KeepOpen() {
		game.UpdateScene()
		game.UpdateUnits()
	}
}
