package game

import (
	"pure-game-kit/packages/execution/condition"
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
		View.DrawGrid(1, TileSize, TileSize, color.RGBA(0, 0, 0, 100))
		View.DrawGrid(1, TileSize*100, TileSize*100, palette.Black)
		// View.DrawObject(Grid)
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
		for _, cols := range laneCollisions {
			for _, s := range cols {
				View.DrawShape(s.X, s.Y, s.Width, s.Height, s.Angle, 0, DebugCollisionColor, geometry.Area{})
			}
		}

		var hovered = DebugHoveredUnit
		if hovered != nil {
			if condition.TrueEvery(0.2, -111) {
				debugInfo = text.New("state: ", debugStates[hovered.State], "\n",
					"action timer: ", number.Round(hovered.actionTimer, 1), "\n",
					"hurt timer: ", number.Round(hovered.hurtTimer, 1), "\n",
					"velocity: ", number.Round(hovered.VelocityX, 1), " | ", number.Round(hovered.VelocityY, 1), "\n",
					"move speed: ", number.Round(hovered.moveSpeedX, 1), "\n",
					"returning: ", hovered.IsReturning, "\n",
				)
			}

			var topX, topY = View.PointFromEdge(0.5, 0)
			View.DrawText(debugInfo, topX, topY, 50, 0, palette.White, geometry.Area{})
		}
	case DebugProfile:
		View.DrawDebugInfo(true)
	}
}

// private ========================================================

var debugStates = []string{
	StateIdling: "idle", StateWalking: "walk",
	StateHurtStart: "hurt", StateHurting: "hurt",
	StateDyingStart: "dying", StateDying: "dying", StateDyingEnd: "dying", StateDecaying: "dead",
	StateActionStart: "action charge", StateActionCharging: "action charge", StateActionTrigger: "action charge",
	StateActionRecovering: "action recover", StateActionEnd: "action recover"}
var debugInfo = ""
