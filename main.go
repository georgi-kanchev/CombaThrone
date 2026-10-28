package main

import (
	game "game/code"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("CombaThrone", game.CurrentSettings.Vsync, game.CurrentSettings.Antialiasing)

	game.Init()
	for window.KeepOpen() {
		game.Update()
	}
}
