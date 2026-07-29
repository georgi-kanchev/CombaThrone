package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/utility/color"
)

const TileSize, MapCount = 32, 4

var View graphics.View
var Background graphics.Object

var MapLayer assets.TileLayerId
var Map, Grid, AllyBase, EnemyBase, Flags graphics.Object
var TilesetCrops, UserInterfaceCrops assets.AnimationsId
var UserInterfaceImg assets.ImageId

var Collisions = make(map[Duty][]geometry.Shape)
var Masks = map[Duty]geometry.Area{
	DutyLower:  geometry.NewArea(0, 0, 560, 1000),
	DutyMiddle: geometry.NewArea(0, 0, 500, 1000),
	DutyUpper:  geometry.NewArea(0, 0, 432, 1000),
}

var Units []*Unit
var Gates []*EntryData

func InitScene() {
	View = graphics.NewView(5.68)
	Background = graphics.NewSprite(0, 0, 1, assets.LoadImage("data/bgr-field.png"))
	UserInterfaceImg = assets.LoadImage("data/user-interface.png")
	UserInterfaceCrops = assets.LoadAnimations(UserInterfaceImg, "data/user-interface.xml")

	var layers, tileset = assets.LoadTileLayersFromTiled("data/map.tmx")
	var upper = graphics.NewTilemap(1, layers[20])
	var middle = graphics.NewTilemap(1, layers[19])
	var lower = graphics.NewTilemap(1, layers[18])
	TilesetCrops = assets.LoadAnimations(tileset, "data/tileset.xml")
	TilesetCrops = assets.LoadAnimations(tileset, "data/tileset.xml")
	Collisions[DutyLower] = lower.TilemapShapes()
	Collisions[DutyMiddle] = middle.TilemapShapes()
	Collisions[DutyUpper] = upper.TilemapShapes()
	MapLayer = layers[15]
	Grid, Map = graphics.NewTilemap(1, layers[17]), graphics.NewTilemap(1, MapLayer)
	Flags = graphics.NewTilemap(1, layers[16])
	AllyBase = graphics.NewTilemap(1, layers[0])
	EnemyBase = graphics.NewTilemap(1, layers[0])
	EnemyBase.Width *= -1
	Units = append(Units, NewUnit(CharacterMan, TeamAlly, DutyLower))
	Units = append(Units, NewUnit(CharacterWoman, TeamEnemy, DutyMiddle))
	Units = append(Units, NewUnit(CharacterWoman, TeamEnemy, DutyLower))

	Gates = append(Gates, NewEntry(EntryHole, TeamAlly, DutyUpper))
	Gates = append(Gates, NewEntry(EntryTallGate, TeamAlly, DutyMiddle))
	Gates = append(Gates, NewEntry(EntryDoor, TeamAlly, DutyLower))

	Gates = append(Gates, NewEntry(EntryDoor, TeamEnemy, DutyUpper))
	Gates = append(Gates, NewEntry(EntryShortGate, TeamEnemy, DutyMiddle))
	Gates = append(Gates, NewEntry(EntryHole, TeamEnemy, DutyLower))
}
func UpdateScene() {
	var _, bly = Background.PointFromEdge(0.5, 1)
	View.FitSize(Background.Width, 0)
	var _, h = View.Size()
	View.Y = (bly - h/2) - 2

	if keyboard.IsKeyJustPressed(key.A) {
		NewUnit(CharacterWoman, TeamEnemy, DutyUpper)
	}

	View.DrawColor(skyColor)
	View.DrawObject(&Background)

	View.DrawObject(&AllyBase)
	View.DrawObject(&EnemyBase)
	View.DrawObject(&Flags)

	for _, g := range Gates {
		g.Update()
	}
	View.DrawObject(&Map)

	for _, u := range Units {
		u.Update()
	}
	for _, g := range Gates {
		g.HealthBar.Update(g.Tiles[0].Shape, g.Health, g.MaxHealth)
	}
	for _, u := range Units {
		u.HealthBar.Update(u.Shape, u.Stats.Health, Characters[u.Character].Stats.Health)
	}

	if keyboard.IsKeyPressed(key.A) {
		Gates[1].ApplyDamage(1)
	}
	if keyboard.IsKeyJustPressed(key.S) {
		Units[2].Stats.Health--
	}
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

var skyColor = color.TagRGBA("rgb(98, 171, 212)")
