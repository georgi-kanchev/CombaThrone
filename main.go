package main

import (
	game "game/code"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("CombaThrone", game.Settings.Vsync, game.Settings.Antialiasing)

	game.Init()
	for window.KeepOpen() {
		game.Update()
	}
}
