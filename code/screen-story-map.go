package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/gui"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
)

func InitStoryMap() {
	var layers, _ = assets.LoadTileLayersFromTiled("data/story.tmx")
	for i, l := range layers {
		var tmap = new(graphics.NewTilemap(l))
		storyMapLayers = append(storyMapLayers, tmap)

		if i < 3 {
			storyMapFogFade = append(storyMapFogFade, nil)
			continue
		}

		storyMapFogFade = append(storyMapFogFade, motion.NewTween(FOG_OPACITY))
	}
	storyMapLayers[0].Details.Tint = palette.White
}
func Progress() {
	storyMapFogFade[3+Player.Progress].GoTo(1.0, nil, 0)
	Player.Progress++
}
func UpdateStoryMap() {
	View.FillSize(TileSize*18, TileSize*10)

	View.X, View.Y = PointAtCell(-0.5, -5.5)
	for i, l := range storyMapLayers {
		if i >= 3 {
			var alpha = byte(storyMapFogFade[i].CurrentValues()[0])
			l.Details.Tint = color.RGBA(255, 255, 255, alpha)
		}
		View.DrawObject(l)
	}

	if keyboard.IsKeyJustPressed(key.F) {
		Progress()
	}

	gui.Scale = View.Zoom
	gui.Image(PanelId, geometry.NewArea(-TileSize*5-8, -TileSize*3.5-8, TileSize*8, TileSize*3), false)
}

// private ========================================================

const FOG_OPACITY = 128

var storyMapLayers []*graphics.Object
var storyMapFogFade []*motion.Tween
