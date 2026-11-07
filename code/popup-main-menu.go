package game

import (
	"game/code/constants/icon"
	"game/code/constants/popup"
	"game/code/constants/screen"
	"game/code/constants/sound"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/gui"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/cursor"
	"pure-game-kit/packages/utility/text"
	"pure-game-kit/packages/window"
)

func InitMainMenu() {
	menuContinue.Set(Tags[icon.Glory] + "Continue")
	menuStory.Set(Tags[icon.Book] + "Story Mode")
	menuArena.Set(Tags[icon.Sword] + "Arena Mode")
	menuSettings.Set(Tags[icon.Gear] + "Settings")
	menuCredits.Set(Tags[icon.Glory] + "Credits")
	menuExit.Set(Tags[icon.Door] + "Exit")
	menuToTitle.Set(Tags[icon.Door] + "Give Up")
}

func UpdateMainMenu() {
	var hud = gui.AreaHUD(0.5, 0.5, 0, 0)
	var noMask = geometry.Area{}

	if Screen == screen.Title {
		var logo = GUI.Crops("logo")[0]
		var logoCrop = logo.CropArea()
		gui.Image(logo, geometry.NewArea(hud.X, hud.Y-TileSize*3, logoCrop.Width, logoCrop.Height), false)
	}

	View.FitSize(320, 320) // always fit logo + buttons on screen, regardless of resolution
	gui.Scale = View.Zoom
	gui.ScaleNinePatches = View.Zoom

	gui.Button(menuSettings.Get(), geometry.NewArea(hud.X, hud.Y+TileSize*2, 100, 28), geometry.Area{}, ThemeGUI, true)
	TryPlayButtonSound()
	if gui.IsJustClicked() {
		Popup = popup.Settings
		mouse.SetCursor(cursor.Arrow)
	}

	if Screen == screen.Title {
		gui.Button(menuStory.Get(), geometry.NewArea(hud.X, hud.Y-TileSize*0.5, 120, 28), noMask, ThemeGUI, true)
		TryPlayButtonSound()
		if gui.IsJustClicked() {
			Screen = screen.Scene
			Popup = popup.NONE
			TimeScale = 1
			PauseMusic()
			PauseAmbience()
			PlaySound(sound.Play)
		}
		gui.Button(menuArena.Get(), geometry.NewArea(hud.X, hud.Y+TileSize*0.5, 120, 28), geometry.Area{}, ThemeGUI, false)
		TryPlayButtonSound()
		if gui.IsHovered() {
			mouse.SetCursor(cursor.NotAllowed)
		}

		gui.Button(menuCredits.Get(), geometry.NewArea(hud.X, hud.Y+TileSize*3, 100, 28), geometry.Area{}, ThemeGUI, true)
		TryPlayButtonSound()
		if gui.IsJustClicked() {
			Popup = popup.Credits
			mouse.SetCursor(cursor.Arrow)
		}

		gui.Button(menuExit.Get(), geometry.NewArea(hud.X, hud.Y+TileSize*4, 100, 28), geometry.Area{}, ThemeGUI, true)
		TryPlayButtonSound()
		if gui.IsJustClicked() {
			window.Close()
		}
		return
	}

	gui.Button(menuContinue.Get(), geometry.NewArea(hud.X, hud.Y+TileSize*0.5, 120, 28), noMask, ThemeGUI, true)
	TryPlayButtonSound()
	if gui.IsJustClicked() {
		Popup = popup.NONE
		TimeScale = 1
		ResumeSounds()
	}

	gui.Button(menuToTitle.Get(), geometry.NewArea(hud.X, hud.Y+TileSize*3, 100, 28), geometry.Area{}, ThemeGUI, true)
	TryPlayButtonSound()
	if gui.IsJustClicked() {
		Screen = screen.Title
		Popup = popup.MainMenu
		PauseMusic()
		PauseAmbience()
	}
}

// private ========================================================

var menuStory, menuContinue, menuArena, menuSettings, menuCredits, menuExit, menuToTitle text.Dynamic
