package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/gui"
	"pure-game-kit/packages/window"
	"pure-game-kit/packages/window/mode"
)

type Settings struct {
	Fullscreen bool
}

var CurrentSettings Settings

func InitSettings() {
	settingsLayout = assets.LoadGUILayout("data/gui-settings.xml")

	CurrentSettings = Settings{Fullscreen: false}
}

func UpdateSettings() {
	const scale = 0.2
	var gfx, _, _ = settingsLayout.Box(1, scale)
	var sfx, _, _ = settingsLayout.Box(2, scale)
	gui.Image(PanelId, gfx, false)
	gui.Image(PanelId, sfx, false)

	var titleGfx, _ = settingsLayout.Item(0, 0, 0, scale)
	gui.Text("🔊Graphics", titleGfx, geometry.Area{}, ThemeGUI, false)

	for i, txt := range settingsLabels {
		var area, _ = settingsLayout.Item(1+i*2, 0, 0, scale)
		gui.Text(txt, area, geometry.Area{}, ThemeGUI.Variant(1), false)
	}

	var fullscreen, _ = settingsLayout.Item(2, 0, 0, scale)
	fullscreen.X -= fullscreen.Width/2 - fullscreen.Height
	fullscreen.Width = fullscreen.Height
	gui.Button(settingsFullscreen, fullscreen, geometry.Area{}, ThemeGUI.Variant(1), true)
	if gui.IsJustClicked() {
		CurrentSettings.Fullscreen = !CurrentSettings.Fullscreen
		settingsFullscreen = ""
		if CurrentSettings.Fullscreen {
			settingsFullscreen = "X"
		}
		ApplySettings()
	}

	var titleSfx, _ = settingsLayout.Item(25, 0, 0, scale)
	gui.Text("🔊Audio", titleSfx, geometry.Area{}, ThemeGUI, false)
}

func ApplySettings() {
	if CurrentSettings.Fullscreen {
		window.SetMode(mode.FullscreenBorderless)
	} else {
		window.SetMode(mode.Maximized)
	}
}

// private ========================================================

var settingsLayout assets.GUILayoutId

var settingsLabels = []string{
	"Fullscreen", "Monitor", "Quality", "Antialiasing", "Vsync", "Limit FPS", "Show FPS", "",
	"Brightness", "Contrast", "Gamma", "Saturation",
}
var settingsFullscreen string
