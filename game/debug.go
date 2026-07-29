package game

import (
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/number"
)

type Debug uint8

const DebugOff, DebugUnits, DebugGrid, DebugStats = 0, 1, 2, 3

var DebugMode Debug

var DebugUnitColor = color.TagRGBA("rgba(0, 255, 0, 0.3)")
var DebugHitboxColor = color.TagRGBA("rgba(255, 0, 0, 0.3)")
var DebugGridColor = color.TagRGBA("rgba(0, 0, 0, 0.2)")
var DebugCollisionColor = color.TagRGBA("rgba(0, 255, 255, 0.15)")
var DebugAttackColor = color.TagRGBA("rgb(255, 255, 255)")

func UpdateDebug() {
	if keyboard.IsKeyJustPressed(key.F3) {
		DebugMode = number.Wrap(DebugMode+1, 0, 4)
	}

	switch DebugMode {
	case DebugOff:
		View.DrawDebugInfo(false)
	case DebugUnits:
		View.DrawDebugInfo(false)
		for _, u := range Units {
			View.DrawShape(u.X, u.Y, u.Width, u.Height, 0, 0, DebugUnitColor, geometry.Area{})
		}
		for _, u := range Units {
			var hb = u.Hitbox()
			var x, y = u.AttackPoint()
			View.DrawShape(x, y, 2, 2, 0, 1, DebugAttackColor, geometry.Area{})
			View.DrawShape(hb.X, hb.Y, hb.Width, hb.Height, 0, hb.Roundness, DebugHitboxColor, geometry.Area{})
		}
		for _, cols := range Collisions {
			for _, s := range cols {
				View.DrawShape(s.X, s.Y, s.Width, s.Height, 0, 0, DebugCollisionColor, geometry.Area{})
			}
		}
	case DebugGrid:
		View.DrawDebugInfo(false)
		Grid.Effects.Tint = color.RGBA(0, 0, 0, 50)
		View.DrawObject(&Grid)
	case DebugStats:
		View.DrawDebugInfo(true)
	}

}
