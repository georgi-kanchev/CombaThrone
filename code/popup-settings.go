package game

import (
	"game/code/constants/popup"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/gui"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/text"
	"pure-game-kit/packages/window"
	f "pure-game-kit/packages/window/filter"
	"pure-game-kit/packages/window/mode"
)

type SettingsData struct {
	Window            mode.Enum
	Filter            f.Enum
	Monitor           int
	Quality, LimitFPS float32

	Antialiasing, Vsync, ShowFPS bool
}

var Settings = SettingsData{Window: mode.Maximized, Quality: 1}

func InitSettings() {
	settingsLayout = assets.LoadGUILayout("data/gui-settings.xml")
	settingsMonitors, _ = window.Monitors()

	ApplySettings()
}

func UpdateSettings() {
	var zero = geometry.Area{}
	gui.Scale = 5

	gui.Rectangle(0, gui.AreaHUD(0.5, 0.5, 1, 1), zero, true) // menu cancel area (entire screen)
	if gui.IsJustClicked() {
		settingsDropdown = zero
	}

	const scale = 0.2
	var gfx, _, _ = settingsLayout.Box(1, scale)
	var sfx, _, _ = settingsLayout.Box(2, scale)
	gui.Image(PanelId, gfx, false)
	gui.Image(PanelId, sfx, false)

	//=================================================================

	var titleGfx, _ = settingsLayout.Item(0, 0, 0, scale)
	gui.Text("🔊Graphics", titleGfx, zero, ThemeGUI, false, 0.5)

	for i, txt := range settingsLabels {
		var area, _ = settingsLayout.Item(1+i*2, 0, 0, scale)
		gui.Text(txt, area, zero, ThemeGUI.Variant(1), false, 0.8)
	}

	var state, _ = settingsLayout.Item(2, 0, 0, scale)
	gui.Button(settingsWindowState, state, zero, ThemeGUI.Variant(1), true, 0.8)
	if gui.IsJustClicked() {
		settingsDropdown = state
	}

	var monitor, _ = settingsLayout.Item(4, 0, 0, scale)
	gui.Button(settingsMonitor, monitor, zero, ThemeGUI.Variant(1), true, 0.5)
	if gui.IsJustClicked() {
		settingsDropdown = monitor
	}

	var quality, _ = settingsLayout.Item(8, 0, 0, scale)
	gui.Slider(&Settings.Quality, 0.01, quality, zero, ThemeGUI, true)
	window.SetQuality(5-Settings.Quality*4, Settings.Filter)
	gui.Text(settingsQuality.Set(number.Round(Settings.Quality*10, 1)), quality, zero, ThemeGUI, false, 0.8)

	var filter, _ = settingsLayout.Item(10, 0, 0, scale)
	gui.Button(settingsFilter, filter, zero, ThemeGUI.Variant(1), Settings.Quality < 1, 0.8)
	if gui.IsJustClicked() {
		settingsDropdown = filter
	}

	var aa, _ = settingsLayout.Item(12, 0, 0, scale)
	gui.Button(settingsAA, aa, zero, ThemeGUI.Variant(1), true, 0.8)
	if gui.IsJustClicked() {
		Settings.Antialiasing = !Settings.Antialiasing
		ApplySettings()
	}

	var vsync, _ = settingsLayout.Item(16, 0, 0, scale)
	gui.Button(settingsVsync, vsync, zero, ThemeGUI.Variant(1), true, 0.8)
	if gui.IsJustClicked() {
		Settings.Vsync = !Settings.Vsync
		ApplySettings()
	}

	var limitFPS, _ = settingsLayout.Item(18, 0, 0, scale)
	gui.Slider(&Settings.LimitFPS, 0.01, limitFPS, zero, ThemeGUI, !Settings.Vsync)
	window.SetTargetFPS(uint8(Settings.LimitFPS * float32(500)))
	gui.Text(settingsLimitFPS.Set(number.Round(Settings.LimitFPS*500)), limitFPS, zero, ThemeGUI, false, 0.8)

	var showFPS, _ = settingsLayout.Item(20, 0, 0, scale)
	gui.Button(settingsShowFPS, showFPS, zero, ThemeGUI.Variant(1), true, 0.8)
	if gui.IsJustClicked() {
		Settings.ShowFPS = !Settings.ShowFPS
		ApplySettings()
	}

	var hint, _ = settingsLayout.Item(22, 0, 0, scale)
	gui.Text("* requires restart", hint, zero, ThemeGUI, false, 0.8)

	var titleSfx, _ = settingsLayout.Item(27, 0, 0, scale)
	gui.Text("🔊Audio", titleSfx, zero, ThemeGUI, false, 0.5)

	// menus =================================================================

	var windowResult = menu(state, "Windowed", "Borderless", "Fullscreen")
	if windowResult != -1 {
		Settings.Window = [3]mode.Enum{mode.Maximized, mode.FullscreenBorderless, mode.Fullscreen}[windowResult]
		ApplySettings()
	}

	var monitorResult = menu(monitor, settingsMonitors...)
	if monitorResult != -1 {
		Settings.Monitor = monitorResult
		ApplySettings()
	}

	var filterResult = menu(filter, settingsFilters...)
	if filterResult != -1 {
		Settings.Filter = f.Enum(filterResult)
		ApplySettings()
	}

	if keyboard.IsKeyJustPressed(key.Escape) {
		Popup = popup.MainMenu
	}
}

func ApplySettings() {
	window.SetMode(Settings.Window)
	window.MoveToMonitor(Settings.Monitor)
	window.SetQuality(5-Settings.Quality*4, Settings.Filter)
	window.SetTargetFPS(uint8(Settings.LimitFPS * float32(500)))
	settingsWindowState = [4]string{"", "Windowed", "Fullscreen", "Borderless"}[Settings.Window]
	settingsFilter = settingsFilters[Settings.Filter]
	settingsMonitor = settingsMonitors[Settings.Monitor]
	settingsAA = settingsBool[Settings.Antialiasing]
	settingsVsync = settingsBool[Settings.Vsync]
	settingsShowFPS = settingsBool[Settings.ShowFPS]
}

// private ========================================================

var settingsLayout assets.GUILayoutId
var settingsDropdown geometry.Area
var settingsWindowState, settingsMonitor, settingsFilter, settingsAA, settingsVsync, settingsShowFPS string
var settingsMonitors []string
var settingsQuality, settingsLimitFPS text.Dynamic

var settingsLabels = []string{
	"Window", "Monitor", "", "Quality", "Filter", "Antialiasing*", "", "Vsync*", "Limit FPS", "Show FPS",
}
var settingsBool = map[bool]string{false: "Off", true: "On"}
var settingsFilters = []string{"Point", "Bilinear", "Trilinear", "Anisotropic4x", "Anisotropic8x", "Anisotropic16x"}

func menu(area geometry.Area, items ...string) (clicked int) {
	if area == settingsDropdown {
		var count = float32(len(items))
		var itemHeight = area.Height
		area = geometry.NewArea(area.X, area.Y+itemHeight/2+(itemHeight*count)/2, area.Width*1.1, itemHeight*count)
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
