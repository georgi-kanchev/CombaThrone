package game

import (
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/utility/color"
)

var Debug = false

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

	var ux, uy, uw, uh = LaneUpper[0].X, LaneUpper[0].Y, LaneUpper[0].Width, LaneUpper[0].Height
	View.DrawShape(ux, uy, uw, uh, 0, 0, DebugCollisionColor, geometry.Area{})
	var mx, my, mw, mh = LaneMiddle[0].X, LaneMiddle[0].Y, LaneMiddle[0].Width, LaneMiddle[0].Height
	View.DrawShape(mx, my, mw, mh, 0, 0, DebugCollisionColor, geometry.Area{})
	var lx, ly, lw, lh = LaneLower[0].X, LaneLower[0].Y, LaneLower[0].Width, LaneLower[0].Height
	View.DrawShape(lx, ly, lw, lh, 0, 0, DebugCollisionColor, geometry.Area{})
}
