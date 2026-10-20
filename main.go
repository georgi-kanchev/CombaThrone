package main

import (
	game "game/code"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("CombaThrone", true, false)

	game.Init()
	for window.KeepOpen() {
		game.Update()
	}
}
