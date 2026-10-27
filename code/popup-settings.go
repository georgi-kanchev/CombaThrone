package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/gui"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/text"
	"pure-game-kit/packages/window"
	"pure-game-kit/packages/window/filter"
	"pure-game-kit/packages/window/mode"
)

type Settings struct {
	WindowState mode.Enum
	Monitor     int
	Quality     float32
}

var CurrentSettings Settings

func InitSettings() {
	settingsLayout = assets.LoadGUILayout("data/gui-settings.xml")
	settingsMonitors, _ = window.Monitors()

	CurrentSettings = Settings{WindowState: mode.Maximized, Quality: 1}
	ApplySettings()
}

func UpdateSettings() {
	gui.Scale = 5

	gui.Rectangle(0, gui.AreaHUD(0.5, 0.5, 1, 1), geometry.Area{}, true) // menu cancel area (entire screen)
	if gui.IsJustClicked() {
		settingsDropdown = geometry.Area{}
	}

	const scale = 0.2
	var gfx, _, _ = settingsLayout.Box(1, scale)
	var sfx, _, _ = settingsLayout.Box(2, scale)
	gui.Image(PanelId, gfx, false)
	gui.Image(PanelId, sfx, false)

	//=================================================================

	var titleGfx, _ = settingsLayout.Item(0, 0, 0, scale)
	gui.Text("🔊Graphics", titleGfx, geometry.Area{}, ThemeGUI, false, 0.5)

	for i, txt := range settingsLabels {
		var area, _ = settingsLayout.Item(1+i*2, 0, 0, scale)
		gui.Text(txt, area, geometry.Area{}, ThemeGUI.Variant(1), false, 0.8)
	}

	var state, _ = settingsLayout.Item(2, 0, 0, scale)
	gui.Button(settingsWindowState, state, geometry.Area{}, ThemeGUI.Variant(1), true, 0.8)
	if gui.IsJustClicked() {
		settingsDropdown = state
	}
	var monitor, _ = settingsLayout.Item(4, 0, 0, scale)
	gui.Button(settingsMonitor, monitor, geometry.Area{}, ThemeGUI.Variant(1), true, 0.5)
	if gui.IsJustClicked() {
		settingsDropdown = monitor
	}

	var quality, _ = settingsLayout.Item(6, 0, 0, scale)
	gui.Slider(&CurrentSettings.Quality, 0.01, quality, geometry.Area{}, 0, true)
	window.SetQuality(7-CurrentSettings.Quality*6, filter.Point)
	gui.Text(settingsQualityDyn.Set(number.Round(CurrentSettings.Quality*10, 1)), quality, geometry.Area{}, ThemeGUI, false, 0.8)

	var titleSfx, _ = settingsLayout.Item(27, 0, 0, scale)
	gui.Text("🔊Audio", titleSfx, geometry.Area{}, ThemeGUI, false, 0.5)

	// menus =================================================================

	var windowResult = menu(state, "Windowed", "Borderless", "Fullscreen")
	if windowResult != -1 {
		CurrentSettings.WindowState = [3]mode.Enum{mode.Maximized, mode.FullscreenBorderless, mode.Fullscreen}[windowResult]
		ApplySettings()
	}

	var monitorResult = menu(monitor, settingsMonitors...)
	if monitorResult != -1 {
		CurrentSettings.Monitor = monitorResult
		ApplySettings()
	}
}

func ApplySettings() {
	window.SetMode(CurrentSettings.WindowState)
	window.MoveToMonitor(CurrentSettings.Monitor)
	settingsWindowState = [4]string{"", "Windowed", "Fullscreen", "Borderless"}[CurrentSettings.WindowState]
	settingsMonitor = settingsMonitors[CurrentSettings.Monitor]
}

// private ========================================================

var settingsLayout assets.GUILayoutId
var settingsDropdown geometry.Area
var settingsWindowState, settingsMonitor string
var settingsMonitors []string
var settingsQualityDyn text.Dynamic

var settingsLabels = []string{
	"Window State", "Monitor", "Quality", "Antialiasing", "Vsync", "Limit FPS", "Show FPS", "",
	"Brightness", "Contrast", "Gamma", "Saturation",
}

func menu(area geometry.Area, items ...string) (clicked int) {
	if area == settingsDropdown {
		var itemHeight = area.Height
		area = geometry.NewArea(area.X, area.Y+itemHeight/2+(itemHeight*3)/2, area.Width, itemHeight*3)
		gui.List(area, geometry.Area{}, ThemeGUI.Variant(1), true, 0.5, itemHeight, items...)
		for i := range items {
			if gui.IsListJustClicked(i) {
				settingsDropdown = geometry.Area{}
				return i
			}
		}
	}
	return -1
}
