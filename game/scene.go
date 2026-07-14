package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color"
)

const LayerMap, LayerDoors, LayerFlags, LayerGrid = 0, 1, 2, 3

var View graphics.View
var skyColor = color.TagRGBA("rgb(90, 135, 218)")
var Background graphics.Object

var Layers []assets.TileLayerId
var Tilemaps []graphics.Object

func InitScene() {
	View = graphics.NewView(5.68)
	View.Y = -5
	Background = graphics.NewSprite(0, 0, 1, assets.LoadImage("data/bgr-field.png"))

	Layers = assets.LoadTileLayersFromTiled("data/map.tmx")
	Tilemaps = make([]graphics.Object, len(Layers))
	for i := range Tilemaps {
		Tilemaps[i] = graphics.NewTilemap(1, Layers[i])
	}
	Tilemaps[LayerGrid].Effects.Tint = DebugGridColor

	SpawnUnit(CharacterMan, DutyWalkStraight, TeamAlly)
	SpawnUnit(CharacterWoman, DutyWalkStraight, TeamAlly)
	Units[0].CellX = 3
}
func UpdateScene() {
	View.DrawColor(skyColor)
	View.DrawObject(&Background)

	for i, tmap := range Tilemaps {
		if !Debug && i == LayerGrid {
			continue
		}
		View.DrawObject(&tmap)
	}
}

func PointAt(cellX, cellY int) (x, y float32) {
	var tw, th = Layers[LayerMap].TileSize()
	var cols, rows = Layers[LayerMap].Size()
	return (float32(cellX)-float32(cols)/2)*tw + (tw / 2), (float32(cellY)-float32(rows)/2)*th + (th / 2)
}
func CellAt(x, y float32) (cellX, cellY int) {
	var tw, th = Layers[LayerMap].TileSize()
	var cols, rows = Layers[LayerMap].Size()
	return int(x/tw + float32(cols)/2), int(y/th + float32(rows)/2)
}
func TileAt(cellX, cellY, layer int) assets.Tile {
	return Layers[layer].TileAtCell(cellX, cellY)
}
