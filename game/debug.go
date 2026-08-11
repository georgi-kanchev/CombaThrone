package game

import (
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/text"
)

type Debug uint8

const DebugOff, DebugGame, DebugProfile = 0, 1, 2

var DebugMode Debug

var DebugUnitColor = color.TagRGBA("rgba(0, 255, 0, 0.3)")
var DebugHitboxColor = color.TagRGBA("rgba(255, 0, 0, 0.4)")
var DebugGridColor = color.TagRGBA("rgba(0, 0, 0, 0.2)")
var DebugCollisionColor = color.TagRGBA("rgba(0, 255, 255, 0.3)")

var DebugHoveredUnit *Unit

func UpdateDebug() {
	if keyboard.IsKeyJustPressed(key.F3) {
		DebugMode = number.Wrap(DebugMode+1, 0, 3)
	}

	if keyboard.IsKeyJustPressed(key.S) {
		Units[0].TakeDamage(12)
	}

	for i := range 10 {
		if keyboard.IsKeyJustPressed(key.Number0 + i) {
			if keyboard.IsKeyPressed(key.Dot) {
				TimeScale = float32(i) / 10
			} else {
				TimeScale = float32(i)
			}
		}
	}

	switch DebugMode {
	case DebugOff:
		View.DrawDebugInfo(false)
	case DebugGame:
		Grid.Effects.Tint = color.RGBA(0, 0, 0, 50)
		View.DrawObject(&Grid)
		View.DrawDebugInfo(false)
		for _, u := range Units {
			View.DrawShape(u.X, u.Y, u.Width, u.Height, 0, 0, DebugUnitColor, geometry.Area{})
		}
		for _, u := range Units {
			var hb = u.Hitbox()
			View.DrawShape(hb.X, hb.Y, hb.Width, hb.Height, 0, hb.Roundness, DebugHitboxColor, geometry.Area{})
			if u.ContainsPoint(View.MousePosition()) {
				DebugHoveredUnit = u
			}
		}
		var hovered = DebugHoveredUnit
		if hovered != nil {
			var states = [13]string{"idle", "walk", "hurt", "hurt", "dying", "dying", "dying", "dead",
				"attack charge", "attack charge", "attack charge", "attack recover", "attack recover"}
			var info = text.New("state: ", states[hovered.State], "\n",
				"attack timer: ", number.Round(hovered.attackTimer, 1), "\n",
				"hurt timer: ", number.Round(hovered.hurtTimer, 1), "\n",
				"at garrison: ", hovered.IsAtGarrison, "\n",
				"velocity: ", number.Round(hovered.VelocityX, 1), " | ", number.Round(hovered.VelocityY, 1), "\n",
				"move speed x: ", number.Round(hovered.moveSpeedX, 1), "\n",
			)
			View.DrawText(info, hovered.X-hovered.Width/2, hovered.Y-100, 50, 0, palette.White, geometry.Area{})
		}
		for _, cols := range Collisions {
			for _, s := range cols {
				View.DrawShape(s.X, s.Y, s.Width, s.Height, s.Angle, 0, DebugCollisionColor, geometry.Area{})
			}
		}
	case DebugProfile:
		View.DrawDebugInfo(true)
	}
}
