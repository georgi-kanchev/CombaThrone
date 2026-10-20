package game

import (
	"game/code/constants/icon"
	"game/code/constants/zone"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/gui"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/text"
)

func InitTitleScreen() {
	titleScreen = assets.LoadImage("data/title-screen.png")
}

func UpdateTitleScreen() {
	alignView()

	const baseAspect = 16.0 / 9.0
	var bgrCrop = titleScreen.CropArea()
	var logo = UserInterface.Crops("logo")[0]
	var logoCrop = logo.CropArea()
	var hud = gui.AreaHUD(0.5, 0.5, 0, 0)

	View.DrawColor(Zones[zone.Desert].SkyColor)
	View.DrawImage(0, 0, bgrCrop.Width, bgrCrop.Height, 0, titleScreen, palette.White, geometry.Area{})

	gui.Image(logo, geometry.NewArea(hud.X, hud.Y-TileSize*3, logoCrop.Width, logoCrop.Height), false)

	View.FitSize(320, 320) // always fit logo + buttons on screen, regardless of resolution
	gui.Scale = View.Zoom

	menuStory.Set(Tags[icon.Book] + "Story Mode")
	menuArena.Set(Tags[icon.Sword] + "Arena Mode")
	menuSettings.Set(Tags[icon.Gear] + "Settings")
	menuCredits.Set(Tags[icon.Glory] + "Credits")
	menuExit.Set(Tags[icon.Door] + "Exit")

	var noMask = geometry.Area{}
	gui.Button(menuStory.Get(), geometry.NewArea(hud.X, hud.Y-TileSize*0.5, 120, 28), noMask, ThemeUI, true)
	// if gui.IsJustClicked() {
	// 	PlayAmbience(CurrentZone.Kind)
	// 	Screen = screen.Scene
	// }

	// gui.Button(menuArena.Get(), geometry.NewArea(hud.X, hud.Y+TileSize*0.5, 120, 28), geometry.Area{}, ThemeUI, false)
	// if gui.IsFocused() {
	// 	mouse.SetCursor(cursor.NotAllowed)
	// }

	// gui.Button(menuSettings.Get(), geometry.NewArea(hud.X, hud.Y+TileSize*2, 100, 28), geometry.Area{}, ThemeUI, true)
	// if gui.IsJustClicked() {
	// 	Screen = screen.Settings
	// }

	// gui.Button(menuCredits.Get(), geometry.NewArea(hud.X, hud.Y+TileSize*3, 100, 28), geometry.Area{}, ThemeUI, true)
	// if gui.IsJustClicked() {
	// 	Screen = screen.Credits
	// }

	// gui.Button(menuExit.Get(), geometry.NewArea(hud.X, hud.Y+TileSize*4, 100, 28), geometry.Area{}, ThemeUI, true)
	// if gui.IsJustClicked() {
	// 	window.Close()
	// }
}

// private ========================================================

var titleScreen assets.ImageId
var menuStory, menuArena, menuSettings, menuCredits, menuExit text.Dynamic
