package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color"
)

const TileSize, MapCount = 32, 4

var View graphics.View
var Background graphics.Object

var MapLayer assets.TileLayerId
var Map, Grid, AllyBase, EnemyBase, Flags graphics.Object
var AllyGates, EnemyGates [2]graphics.Object

var LaneUpper, LaneMiddle, LaneLower [2]geometry.Shape

func InitScene() {
	View = graphics.NewView(5.68)
	View.Y = -5
	Background = graphics.NewSprite(0, 0, 1, assets.LoadImage("data/bgr-field.png"))

	var layers = assets.LoadTileLayersFromTiled("data/map.tmx")
	var upper = graphics.NewTilemap(1, layers[20])
	var middle = graphics.NewTilemap(1, layers[19])
	var lower = graphics.NewTilemap(1, layers[18])
	LaneUpper[0] = upper.TilemapShapes()[0]
	LaneMiddle[0] = middle.TilemapShapes()[0]
	LaneLower[0] = lower.TilemapShapes()[0]
	MapLayer = layers[0]
	Grid, Map = graphics.NewTilemap(1, layers[17]), graphics.NewTilemap(1, MapLayer)
	Flags = graphics.NewTilemap(1, layers[16])
	AllyBase = graphics.NewTilemap(1, layers[1])
	AllyGates[0], AllyGates[1] = graphics.NewTilemap(1, layers[11]), graphics.NewTilemap(1, layers[11])
	EnemyBase = graphics.NewTilemap(1, layers[1])
	EnemyGates[0], EnemyGates[1] = graphics.NewTilemap(1, layers[10]), graphics.NewTilemap(1, layers[10])
	EnemyBase.Width *= -1
	EnemyGates[0].Width *= -1
	EnemyGates[1].Width *= -1

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

	View.DrawObject(&Map)
	View.DrawObject(&AllyBase)
	View.DrawObject(&AllyGates[0])
	View.DrawObject(&AllyGates[1])
	View.DrawObject(&EnemyBase)
	View.DrawObject(&EnemyGates[0])
	View.DrawObject(&EnemyGates[1])
	View.DrawObject(&Flags)
}

func PointAtCell(cellX, cellY float32) (x, y float32) {
	var tw, th = MapLayer.TileSize()
	var cols, rows = MapLayer.Size()
	return (cellX-float32(cols)/2)*tw + (tw / 2), (cellY-float32(rows)/2)*th + (th / 2)
}
func CellAtPoint(x, y float32) (cellX, cellY float32) {
	var tw, th = MapLayer.TileSize()
	var cols, rows = MapLayer.Size()
	return x/tw + float32(cols)/2, y/th + float32(rows)/2
}
func TileAtCell(cellX, cellY int, layer assets.TileLayerId) assets.Tile {
	return layer.TileAtCell(cellX, cellY)
}

// private ========================================================

var skyColor = color.TagRGBA("rgb(90, 135, 218)")
