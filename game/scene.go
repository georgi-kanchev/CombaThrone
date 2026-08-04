package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/number"
)

const TileSize, MapCount = 32, 4

var View graphics.View
var Background graphics.Object

var MapLayer assets.TileLayerId
var Map, Grid, AllyBase, EnemyBase, Flags graphics.Object
var TilesetCrops, UserInterfaceCrops assets.AtlasId
var UserInterfaceImg assets.ImageId

var Collisions = make([][]geometry.Shape, 4)
var Masks = [3]geometry.Area{
	LaneLower:  geometry.NewArea(0, 0, 556, 1000),
	LaneMiddle: geometry.NewArea(0, 0, 492, 1000),
	LaneUpper:  geometry.NewArea(0, 0, 428, 1000),
}
var Entrances [6]*EntryData // index is Ally[Lower, Middle, Upper], Enemy[Lower, Middle, Upper]
var Units []*Unit

func InitScene() {
	View = graphics.NewView(5.68)
	Background = graphics.NewSprite(0, 0, 1, assets.LoadImage("data/bgr-field.png"))
	UserInterfaceImg = assets.LoadImage("data/user-interface.png")
	UserInterfaceCrops = assets.LoadAtlas(UserInterfaceImg, "data/user-interface.xml")

	var layers, tileset = assets.LoadTileLayersFromTiled("data/map.tmx")
	var bridge = graphics.NewTilemap(layers[15])
	var upper = graphics.NewTilemap(layers[14])
	var middle = graphics.NewTilemap(layers[13])
	var lower = graphics.NewTilemap(layers[12])
	TilesetCrops = assets.LoadAtlas(tileset, "data/tileset.xml")
	TilesetCrops = assets.LoadAtlas(tileset, "data/tileset.xml")
	Collisions[LaneLower] = lower.TilemapShapes()
	Collisions[LaneMiddle] = middle.TilemapShapes()
	Collisions[LaneUpper] = upper.TilemapShapes()
	Collisions[LaneBridge] = bridge.TilemapShapes()
	MapLayer, Map = layers[9], graphics.NewTilemap(layers[9])
	Flags, Grid = graphics.NewTilemap(layers[10]), graphics.NewTilemap(layers[11])
	AllyBase = graphics.NewTilemap(layers[1])
	EnemyBase = graphics.NewTilemap(layers[1])
	EnemyBase.Width *= -1

	Units = append(Units, NewUnit(CharacterMan, TeamAlly, LaneUpper))
	Units = append(Units, NewUnit(CharacterWoman, TeamEnemy, LaneUpper))
	Units = append(Units, NewUnit(CharacterWoman, TeamEnemy, LaneMiddle))
	Units = append(Units, NewUnit(CharacterWoman, TeamEnemy, LaneLower))

	Entrances = [6]*EntryData{
		LaneLower:      NewEntry(EntryDoor, TeamAlly, LaneLower),
		LaneMiddle:     NewEntry(EntryTallGate, TeamAlly, LaneMiddle),
		LaneUpper:      NewEntry(EntryHole, TeamAlly, LaneUpper),
		LaneLower + 3:  NewEntry(EntryDoor, TeamEnemy, LaneLower),
		LaneMiddle + 3: NewEntry(EntryDoor, TeamEnemy, LaneMiddle),
		LaneUpper + 3:  NewEntry(EntryDoor, TeamEnemy, LaneUpper),
	}
}
func UpdateScene() {
	var _, bly = Background.PointFromEdge(0.5, 1)
	View.FitSize(Background.Width, 0)
	var _, h = View.Size()
	View.Y = (bly - h/2) - 2

	if keyboard.IsKeyJustPressed(key.A) {
		NewUnit(CharacterWoman, TeamEnemy, LaneUpper)
	}

	View.DrawColor(skyColor)
	View.DrawObject(&Background)

	View.DrawObject(&AllyBase)
	View.DrawObject(&EnemyBase)
	View.DrawObject(&Flags)

	for _, g := range Entrances {
		g.Update()
	}
	View.DrawObject(&Map)

	collection.SortByField(Units, func(u *Unit) float32 {
		if u.Stats.Health <= 0 { // dead units go behind all alive units
			return number.NegativeInfinity()
		}
		return u.Y // fall back to Y sort
	})

	for _, u := range Units {
		u.Update()
	}
	for _, g := range Entrances {
		g.HealthBar.Update(g.Tiles[0].Shape, g.Health, g.MaxHealth, geometry.Area{})
	}
	for _, u := range Units { // health bars take the Z order of the units
		u.HealthBar.Update(u.Shape, u.Stats.Health, Characters[u.Character].Stats.Health, u.Mask)
	}

	if keyboard.IsKeyJustPressed(key.A) {
		Entrances[2].TakeDamage(110)
	}
	if keyboard.IsKeyJustPressed(key.S) {
		Units = append(Units, NewUnit(CharacterWoman, TeamEnemy, LaneUpper))
	}
	if keyboard.IsKeyJustPressed(key.W) {
		Units[0].VelocityY = -100
		Units[0].Lane = LaneUpper
	}
	if keyboard.IsKeyJustPressed(key.E) {
		Units = append(Units, NewUnit(CharacterWoman, TeamEnemy, LaneLower))
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
