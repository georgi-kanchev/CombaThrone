package game

import "pure-game-kit/packages/utility/color"

var Debug = true

var DebugUnitColor = color.TagRGBA("rgba(0, 255, 0, 0.5)")
var DebugHitboxColor = color.TagRGBA("rgba(255, 0, 0, 0.5)")

func UpdateDebug() {
	if !Debug {
		return
	}

	View.DrawDebugInfo(true)
}
