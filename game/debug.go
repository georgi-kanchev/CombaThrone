package game

import (
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/utility/color"
)

var Debug = true

var DebugUnitColor = color.TagRGBA("rgba(0, 255, 0, 0.5)")
var DebugHitboxColor = color.TagRGBA("rgba(255, 0, 0, 0.5)")
var DebugGridColor = color.TagRGBA("rgba(0, 0, 0, 0.2)")

func UpdateDebug() {
	View.DrawDebugInfo(Debug)

	if keyboard.IsKeyJustPressed(key.F3) {
		Debug = !Debug
	}
}
