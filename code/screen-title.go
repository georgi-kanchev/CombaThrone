package game

import (
	"game/code/constants/zone"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/utility/color/palette"
)

func InitTitleScreen() {
	titleScreen = assets.LoadImage("data/title-screen.png")
}

func UpdateTitleScreen() {
	AlignView()

	var bgrCrop = titleScreen.CropArea()
	View.DrawColor(Zones[zone.Desert].SkyColor)
	View.DrawImage(0, 0, bgrCrop.Width, bgrCrop.Height, 0, titleScreen, palette.White, geometry.Area{})
}

// private ========================================================

var titleScreen assets.ImageId
