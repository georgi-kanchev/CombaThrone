package game

// Handles the playground

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color"
)

var View graphics.View
var skyColor = color.TagRGBA("rgb(90, 135, 218)")
var Background graphics.Object

var Layers []assets.TileLayerId
var Tilemaps []graphics.Object

func InitScene() {
	View = graphics.NewView(1)
	Background = graphics.NewSprite(0, 0, 1, assets.LoadImage("data/bgr-field.png"))

	Layers = assets.LoadTileLayersFromTiled("data/map.tmx")
	Tilemaps = make([]graphics.Object, len(Layers))
	for i := range Tilemaps {
		Tilemaps[i] = graphics.NewTilemap(1, Layers[i])
	}

	SpawnUnit(CharacterMan, DutyWalkStraight, TeamAlly)
	SpawnUnit(CharacterWoman, DutyWalkStraight, TeamAlly)
	Units[0].X = 200
}

func UpdateScene() {
	var w, h = View.Size()
	var aspect = w / h
	var _, bottomY = View.PointFromEdge(0.5, 1)

	View.DrawColor(skyColor)
	View.DrawObject(&Background)

	for _, tmap := range Tilemaps {
		if aspect > 1.7778 {
			tmap.ViewFill(&View)
			tmap.Width = w
		} else {
			tmap.ViewFit(&View)
		}

		tmap.Y = bottomY - tmap.Height/2
		Background.Shape = geometry.NewRectangle(tmap.X, tmap.Y, tmap.Width, tmap.Height, 0)
		View.DrawObject(&tmap)
	}
}

func SceneScale() float32 {
	var _, _, _, h = Background.ImageId.CropArea()
	return Background.Height / h
}
