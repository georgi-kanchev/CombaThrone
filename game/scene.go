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

	SpawnUnit(CharacterMan, DutyUseStairs, TeamAlly)
	SpawnUnit(CharacterWoman, DutyWalkStraight, TeamEnemy)
	Units[0].X = -32*8 + 16
	Units[1].X = 32*7 + 16
}
func UpdateScene() {
	var _, bly = Background.PointFromEdge(0.5, 1)
	View.FitSize(Background.Width, 0)
	var _, h = View.Size()
	View.Y = (bly - h/2) - 8

	View.DrawColor(skyColor)
	View.DrawObject(&Background)

	for i, tmap := range Tilemaps {
		if !Debug && i == LayerGrid {
			continue
		}
		View.DrawObject(&tmap)
	}
}

func PointAtCell(cellX, cellY float32) (x, y float32) {
	var tw, th = Layers[LayerMap].TileSize()
	var cols, rows = Layers[LayerMap].Size()
	return (cellX-float32(cols)/2)*tw + (tw / 2), (cellY-float32(rows)/2)*th + (th / 2)
}
func CellAtPoint(x, y float32) (cellX, cellY float32) {
	var tw, th = Layers[LayerMap].TileSize()
	var cols, rows = Layers[LayerMap].Size()
	return x/tw + float32(cols)/2, y/th + float32(rows)/2
}
func TileAtCell(cellX, cellY, layer int) assets.Tile {
	return Layers[layer].TileAtCell(cellX, cellY)
}
