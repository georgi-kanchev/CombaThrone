package game

import (
	"game/code/constants/music"
	"game/code/constants/popup"
	"game/code/constants/screen"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/audio"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/gui"
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

	VolumeMaster, VolumeSounds, VolumeMusic, VolumeAmbience float32
}

var Settings = SettingsData{Window: mode.Maximized, Quality: 1, VolumeMaster: 1, VolumeSounds: 1, VolumeMusic: 0.0, VolumeAmbience: 1.0}

func InitSettings() {
	settingsLayout = assets.LoadGUILayout("data/gui-settings.xml")
	settingsMonitors, _ = window.Monitors()

	applySettings()
	applyAudioVolume()
}

func UpdateSettings() {
	if Screen != screen.Title {
		DimScreenUI()
	}

	var zero = geometry.Area{}
	gui.Scale = 5.56

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

	for i, txt := range settingsLabelsGfx {
		var area, _ = settingsLayout.Item(1+i*2, 0, 0, scale)
		gui.Text(txt, area, zero, ThemeGUI.Variant(1), false, 0.8)
	}
	for i, txt := range settingsLabelsAudio {
		var area, _ = settingsLayout.Item(28+i*2, 0, 0, scale)
		gui.Text(txt, area, zero, ThemeGUI.Variant(1), false, 0.8)
	}

	var state, _ = settingsLayout.Item(2, 0, 0, scale)
	gui.Button(settingsWindowState, state, zero, ThemeGUI.Variant(1), true, 0.8)
	TryPlayButtonSound()
	if gui.IsJustClicked() {
		settingsDropdown = state
	}

	var monitor, _ = settingsLayout.Item(4, 0, 0, scale)
	gui.Button(settingsMonitor, monitor, zero, ThemeGUI.Variant(1), true, 0.5)
	TryPlayButtonSound()
	if gui.IsJustClicked() {
		settingsDropdown = monitor
	}

	var quality, _ = settingsLayout.Item(8, 0, 0, scale)
	gui.Slider(&Settings.Quality, 0.01, quality, zero, ThemeGUI, true)
	TryPlayButtonSound()
	window.SetQuality(5-Settings.Quality*4, Settings.Filter)
	gui.Text(settingsQuality.Set(number.Round(Settings.Quality*10, 1)), quality, zero, ThemeGUI, false, 0.8)

	var filter, _ = settingsLayout.Item(10, 0, 0, scale)
	gui.Button(settingsFilter, filter, zero, ThemeGUI.Variant(1), Settings.Quality < 1, 0.8)
	TryPlayButtonSound()
	if gui.IsJustClicked() {
		settingsDropdown = filter
	}

	var aa, _ = settingsLayout.Item(12, 0, 0, scale)
	gui.Button(settingsAA, aa, zero, ThemeGUI.Variant(1), true, 0.8)
	TryPlayButtonSound()
	if gui.IsJustClicked() {
		Settings.Antialiasing = !Settings.Antialiasing
		applySettings()
	}

	var vsync, _ = settingsLayout.Item(16, 0, 0, scale)
	gui.Button(settingsVsync, vsync, zero, ThemeGUI.Variant(1), true, 0.8)
	TryPlayButtonSound()
	if gui.IsJustClicked() {
		Settings.Vsync = !Settings.Vsync
		applySettings()
	}

	var limitFPS, _ = settingsLayout.Item(18, 0, 0, scale)
	gui.Slider(&Settings.LimitFPS, 0.01, limitFPS, zero, ThemeGUI, !Settings.Vsync)
	TryPlayButtonSound()
	window.SetTargetFPS(uint8(Settings.LimitFPS * float32(500)))
	gui.Text(settingsLimitFPS.Set(number.Round(Settings.LimitFPS*500)), limitFPS, zero, ThemeGUI, false, 0.8)

	var showFPS, _ = settingsLayout.Item(20, 0, 0, scale)
	gui.Button(settingsShowFPS, showFPS, zero, ThemeGUI.Variant(1), true, 0.8)
	TryPlayButtonSound()
	if gui.IsJustClicked() {
		Settings.ShowFPS = !Settings.ShowFPS
		applySettings()
	}

	var hint, _ = settingsLayout.Item(24, 0, 0, scale)
	gui.Text("* requires restart", hint, zero, ThemeGUI, false, 0.8)

	//=================================================================

	var titleSfx, _ = settingsLayout.Item(27, 0, 0, scale)
	gui.Text("🔊Audio", titleSfx, zero, ThemeGUI, false, 0.5)

	var master, _ = settingsLayout.Item(29, 0, 0, scale)
	gui.Slider(&Settings.VolumeMaster, 0.01, master, zero, ThemeGUI, true)
	TryPlayButtonSound()
	if audio.Volume != Settings.VolumeMaster {
		audio.Volume = Settings.VolumeMaster
		applyAudioVolume()
	}
	gui.Text(settingsMaster.Set(number.Round(audio.Volume*100, 1), "%"), master, zero, ThemeGUI, false, 0.8)

	var snd, _ = settingsLayout.Item(31, 0, 0, scale)
	gui.Slider(&Settings.VolumeSounds, 0.01, snd, zero, ThemeGUI, true)
	TryPlayButtonSound()
	if gui.IsSliderJustChanged() {
		applyAudioVolume()
	}
	gui.Text(settingsSounds.Set(number.Round(Settings.VolumeSounds*100, 1), "%"), snd, zero, ThemeGUI, false, 0.8)

	var mus, _ = settingsLayout.Item(33, 0, 0, scale)
	gui.Slider(&Settings.VolumeMusic, 0.01, mus, zero, ThemeGUI, true)
	TryPlayButtonSound()
	if gui.IsSliderJustChanged() {
		applyAudioVolume()
	}
	gui.Text(settingsMusic.Set(number.Round(Settings.VolumeMusic*100, 1), "%"), mus, zero, ThemeGUI, false, 0.8)

	var amb, _ = settingsLayout.Item(35, 0, 0, scale)
	gui.Slider(&Settings.VolumeAmbience, 0.01, amb, zero, ThemeGUI, true)
	TryPlayButtonSound()
	if gui.IsSliderJustChanged() {
		applyAudioVolume()
	}
	gui.Text(settingsAmbience.Set(number.Round(Settings.VolumeAmbience*100, 1), "%"), amb, zero, ThemeGUI, false, 0.8)

	// menus =================================================================

	var windowResult = menu(state, "Windowed", "Borderless", "Fullscreen")
	if windowResult != -1 {
		Settings.Window = [3]mode.Enum{mode.Maximized, mode.FullscreenBorderless, mode.Fullscreen}[windowResult]
		applySettings()
	}

	var monitorResult = menu(monitor, settingsMonitors...)
	if monitorResult != -1 {
		Settings.Monitor = monitorResult
		applySettings()
	}

	var filterResult = menu(filter, settingsFilters...)
	if filterResult != -1 {
		Settings.Filter = f.Enum(filterResult)
		applySettings()
	}

	// music =================================================================

	if CurrentMusic() != nil {
		ShowCurrentMusic()
	} else {
		PlayMusic(music.Menu)
	}

	//=================================================================

	var panel, _, _ = settingsLayout.Box(0, scale)
	panel = geometry.NewArea(panel.X, panel.Y+panel.Height/2+10, 60, 20)
	gui.Button("Back", panel, zero, ThemeGUI.Variant(1), true, 0.7)
	TryPlayButtonSound()
	if gui.IsJustClicked() {
		Popup = popup.MainMenu
	}
}

// private ========================================================

var settingsLayout assets.GUILayoutId
var settingsDropdown geometry.Area
var settingsWindowState, settingsMonitor, settingsFilter, settingsAA, settingsVsync, settingsShowFPS string
var settingsMonitors []string
var settingsQuality, settingsLimitFPS, settingsMaster, settingsSounds, settingsMusic, settingsAmbience text.Dynamic

var settingsLabelsGfx = []string{"Window", "Monitor", "", "Quality", "Filter", "Antialiasing*", "", "Vsync*", "Limit FPS", "Show FPS"}
var settingsLabelsAudio = []string{"Master", "Game", "Music", "Ambience"}
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

func applySettings() {
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

func applyAudioVolume() {
	for _, v := range allMusic {
		v.Volume = Settings.VolumeMusic
		v.ApplyProperties()
	}
	for _, v := range allSounds {
		v.Volume = Settings.VolumeSounds
		v.ApplyProperties()
	}
	for _, v := range allAmbiences {
		v.Volume = Settings.VolumeAmbience
		v.ApplyProperties()
	}
}
