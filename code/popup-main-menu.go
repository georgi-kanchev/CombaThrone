package game

import (
	"game/code/constants/icon"
	"game/code/constants/popup"
	"game/code/constants/screen"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/gui"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/cursor"
	"pure-game-kit/packages/utility/text"
	"pure-game-kit/packages/window"
)

func InitMainMenu() {
	menuStory.Set(Tags[icon.Book] + "Story Mode")
	menuArena.Set(Tags[icon.Sword] + "Arena Mode")
	menuSettings.Set(Tags[icon.Gear] + "Settings")
	menuCredits.Set(Tags[icon.Glory] + "Credits")
	menuExit.Set(Tags[icon.Door] + "Exit")
}

func UpdateMainMenu() {
	var logo = GUI.Crops("logo")[0]
	var logoCrop = logo.CropArea()
	var hud = gui.AreaHUD(0.5, 0.5, 0, 0)
	gui.Image(logo, geometry.NewArea(hud.X, hud.Y-TileSize*3, logoCrop.Width, logoCrop.Height), false)

	View.FitSize(320, 320) // always fit logo + buttons on screen, regardless of resolution
	gui.Scale = View.Zoom

	var noMask = geometry.Area{}
	gui.Button(menuStory.Get(), geometry.NewArea(hud.X, hud.Y-TileSize*0.5, 120, 28), noMask, ThemeGUI, true)
	if gui.IsJustClicked() {
		PlayAmbience(CurrentZone.Kind)
		Screen = screen.Scene
		Popup = popup.NONE
	}

	gui.Button(menuArena.Get(), geometry.NewArea(hud.X, hud.Y+TileSize*0.5, 120, 28), geometry.Area{}, ThemeGUI, false)
	if gui.IsFocused() {
		mouse.SetCursor(cursor.NotAllowed)
	}

	gui.Button(menuSettings.Get(), geometry.NewArea(hud.X, hud.Y+TileSize*2, 100, 28), geometry.Area{}, ThemeGUI, true)
	if gui.IsJustClicked() {
		Popup = popup.Settings
	}

	gui.Button(menuCredits.Get(), geometry.NewArea(hud.X, hud.Y+TileSize*3, 100, 28), geometry.Area{}, ThemeGUI, true)
	if gui.IsJustClicked() {
		Popup = popup.Credits
	}

	gui.Button(menuExit.Get(), geometry.NewArea(hud.X, hud.Y+TileSize*4, 100, 28), geometry.Area{}, ThemeGUI, true)
	if gui.IsJustClicked() {
		window.Close()
	}
}

// private ========================================================

var menuStory, menuArena, menuSettings, menuCredits, menuExit text.Dynamic
