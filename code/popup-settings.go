package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
)

func InitSettings() {
	settingsView = new(graphics.NewView(1))
	settingsLayout = assets.LoadGUILayout("data/gui-settings.xml")
}

func UpdateSettings() {
	// var sc = View.Zoom
	// settingsView.Zoom = sc

	// var panel, _, _ = settingsLayout.Box(0)
	// var gfx, _, _ = settingsLayout.Box(0)

	var area, _, _ = settingsLayout.Box(0)
	settingsView.DrawShape(area.X, area.Y, 100, 100, 0, 1, palette.White, geometry.Area{})
	// settingsView.DrawImage(area.X, area.Y, area.Width, area.Height, 0, PanelId, palette.White, geometry.Area{})
}

// private ========================================================

var settingsView *graphics.View
var settingsLayout assets.GUILayoutId
