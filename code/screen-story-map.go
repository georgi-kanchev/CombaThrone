package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
)

func InitStoryMap() {
	var layers, _ = assets.LoadTileLayersFromTiled("data/story.tmx")
	for _, l := range layers {
		storyMapLayers = append(storyMapLayers, new(graphics.NewTilemap(l)))
	}
}
func UpdateStoryMap() {
	View.FillSize(TileSize*18, TileSize*10)

	View.X, View.Y = PointAtCell(-0.5, -5.5)
	for _, l := range storyMapLayers {
		View.DrawObject(l)
	}
}

// private ========================================================

var storyMapLayers []*graphics.Object
