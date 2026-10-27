package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/gui"
	"pure-game-kit/packages/window"
	"pure-game-kit/packages/window/mode"
)

type Settings struct {
	WindowState mode.Enum
}

var CurrentSettings Settings

func InitSettings() {
	settingsLayout = assets.LoadGUILayout("data/gui-settings.xml")

	CurrentSettings = Settings{WindowState: mode.Maximized}
}

func UpdateSettings() {
	gui.Scale = 5

	const scale = 0.2
	var gfx, _, _ = settingsLayout.Box(1, scale)
	var sfx, _, _ = settingsLayout.Box(2, scale)
	gui.Image(PanelId, gfx, false)
	gui.Image(PanelId, sfx, false)

	var titleGfx, _ = settingsLayout.Item(0, 0, 0, scale)
	gui.Text("🔊Graphics", titleGfx, geometry.Area{}, ThemeGUI, false, 0.5)

	for i, txt := range settingsLabels {
		var area, _ = settingsLayout.Item(1+i*2, 0, 0, scale)
		gui.Text(txt, area, geometry.Area{}, ThemeGUI.Variant(1), false, 0.8)
	}

	var fullscreen, _ = settingsLayout.Item(2, 0, 0, scale)
	gui.Button(settingsFullscreen, fullscreen, geometry.Area{}, ThemeGUI.Variant(1), true, 0.8)
	if gui.IsJustClicked() {
		settingsDropdown = !settingsDropdown
	}
	if settingsDropdown {
		const height = 40.0
		fullscreen.Y += fullscreen.Height/2 + height/2
		fullscreen.Width *= 0.8
		fullscreen.Height = height
		gui.List(fullscreen, geometry.Area{}, ThemeGUI.Variant(1), true, 0.6, height/3, "Windowed", "Borderless", "Fullscreen")
		if gui.IsListJustClicked(0) {
			settingsDropdown = false
			CurrentSettings.WindowState = mode.Maximized
			ApplySettings()
		}
		if gui.IsListJustClicked(1) {
			settingsDropdown = false
			CurrentSettings.WindowState = mode.FullscreenBorderless
			ApplySettings()
		}
		if gui.IsListJustClicked(2) {
			settingsDropdown = false
			CurrentSettings.WindowState = mode.Fullscreen
			ApplySettings()
		}
	}

	var titleSfx, _ = settingsLayout.Item(27, 0, 0, scale)
	gui.Text("🔊Audio", titleSfx, geometry.Area{}, ThemeGUI, false, 0.5)
}

func ApplySettings() {
	window.SetMode(CurrentSettings.WindowState)
}

// private ========================================================

var settingsLayout assets.GUILayoutId
var settingsDropdown bool
var settingsFullscreen string

var settingsLabels = []string{
	"Window State", "Monitor", "Quality", "Antialiasing", "Vsync", "Limit FPS", "Show FPS", "",
	"Brightness", "Contrast", "Gamma", "Saturation",
}
