package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color"
)

const TileSize, MapCount = 32, 4

var View graphics.View
var Background graphics.Object

var LayerGrid assets.TileLayerId
var LayerMaps []assets.TileLayerId
var TilesMap graphics.Object
var CurrentMap = 3

func InitScene() {
	View = graphics.NewView(5.68)
	View.Y = -5
	Background = graphics.NewSprite(0, 0, 1, assets.LoadImage("data/bgr-field.png"))

	var layers = assets.LoadTileLayersFromTiled("data/map.tmx")
	LayerGrid = layers[len(layers)-1]
	LayerMaps = layers[0:MapCount]
	TilesMap = graphics.NewTilemap(1, LayerMaps[CurrentMap])

	SpawnUnit(CharacterMan, TeamAlly)
	SpawnUnit(CharacterWoman, TeamEnemy)
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

	View.DrawObject(&TilesMap)
}

func PointAtCell(cellX, cellY float32) (x, y float32) {
	var tw, th = LayerGrid.TileSize()
	var cols, rows = LayerGrid.Size()
	return (cellX-float32(cols)/2)*tw + (tw / 2), (cellY-float32(rows)/2)*th + (th / 2)
}
func CellAtPoint(x, y float32) (cellX, cellY float32) {
	var tw, th = LayerGrid.TileSize()
	var cols, rows = LayerGrid.Size()
	return x/tw + float32(cols)/2, y/th + float32(rows)/2
}
func TileAtCell(cellX, cellY int, layer assets.TileLayerId) assets.Tile {
	return layer.TileAtCell(cellX, cellY)
}

// private ========================================================

var skyColor = color.TagRGBA("rgb(90, 135, 218)")
