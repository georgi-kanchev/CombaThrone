package game

import (
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/utility/color"
)

var Debug = true

var DebugUnitColor = color.TagRGBA("rgba(0, 255, 0, 0.3)")
var DebugHitboxColor = color.TagRGBA("rgba(255, 0, 0, 0.3)")
var DebugGridColor = color.TagRGBA("rgba(0, 0, 0, 0.2)")
var DebugCollisionColor = color.TagRGBA("rgba(0, 255, 255, 0.15)")
var DebugAttackColor = color.TagRGBA("rgb(255, 255, 255)")

func UpdateDebug() {
	View.DrawDebugInfo(Debug)

	if keyboard.IsKeyJustPressed(key.F3) {
		Debug = !Debug
	}

	if !Debug {
		return
	}
	for _, u := range Units {
		View.DrawShape(u.X, u.Y, u.Width, u.Height, 0, 0, DebugUnitColor, geometry.Area{})
	}
	for _, u := range Units {
		var x, y = u.AttackPoint()
		View.DrawShape(x, y, 2, 2, 0, 1, DebugAttackColor, geometry.Area{})
	}
	Grid.Effects.Tint = color.RGBA(0, 0, 0, 50)
	View.DrawObject(&Grid)

	for _, cols := range Collisions {
		for _, s := range cols {
			View.DrawShape(s.X, s.Y, s.Width, s.Height, 0, 0, DebugCollisionColor, geometry.Area{})
		}
	}
}
