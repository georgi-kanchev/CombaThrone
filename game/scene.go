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
	"pure-game-kit/packages/utility/time"
)

type Environment uint8

const TileSize, MapCount = 32, 4

const LaneLower, LaneMiddle, LaneUpper Lane = 0, 1, 2
const LaneGarrison1, LaneGarrison2, LaneGarrison3, LaneGarrison4, LaneGarrison5 Lane = 3, 4, 5, 6, 7
const LaneGarrisonPlus1, LaneGarrisonPlus2, LaneGarrisonPlus3, LaneGarrisonPlus4, LaneGarrisonPlus5 Lane = 8, 9, 10, 11, 12
const EnvironmentPlains Environment = 0

var TimeScale float32 = 1

var View graphics.View
var Background graphics.Object

var MapLayer assets.TileLayerId
var Map, Grid, Flags graphics.Object
var TilesetCrops, UserInterfaceCrops assets.AtlasId
var UserInterfaceImg assets.ImageId

var AllyBase, EnemyBase BaseData
var Layers []assets.TileLayerId
var Collisions = make([][]geometry.Shape, 13)
var Masks = [13]geometry.Area{
	LaneLower:  geometry.NewArea(0, 0, 556, 1000),
	LaneMiddle: geometry.NewArea(0, 0, 492, 1000),
	LaneUpper:  geometry.NewArea(0, 0, 428, 1000),
}
var Entrances [6]*Entrance // index is Ally[Lower, Middle, Upper], Enemy[Lower, Middle, Upper]
var Units []*Unit
var Projectiles []*Projectile

func InitScene() {
	View = graphics.NewView(1)
	Background = graphics.NewSprite(0, 0, 1, assets.LoadImage("data/bgr-field.png"))
	UserInterfaceImg = assets.LoadImage("data/user-interface.png")
	UserInterfaceCrops = assets.LoadAtlas(UserInterfaceImg, "data/user-interface.xml")

	var layers, tileset = assets.LoadTileLayersFromTiled("data/map.tmx")
	var low, mid, up = graphics.NewTilemap(layers[12]), graphics.NewTilemap(layers[13]), graphics.NewTilemap(layers[14])
	var g1, g2, g3 = graphics.NewTilemap(layers[15]), graphics.NewTilemap(layers[16]), graphics.NewTilemap(layers[17])
	var g4, g5 = graphics.NewTilemap(layers[18]), graphics.NewTilemap(layers[19])
	var p1, p2, p3 = graphics.NewTilemap(layers[20]), graphics.NewTilemap(layers[21]), graphics.NewTilemap(layers[22])
	var p4, p5 = graphics.NewTilemap(layers[23]), graphics.NewTilemap(layers[24])

	Layers = layers
	TilesetCrops = assets.LoadAtlas(tileset, "data/tileset.xml")
	TilesetCrops = assets.LoadAtlas(tileset, "data/tileset.xml")
	Collisions[LaneLower] = low.TilemapShapes()
	Collisions[LaneMiddle] = mid.TilemapShapes()
	Collisions[LaneUpper] = up.TilemapShapes()
	Collisions[LaneGarrison1] = g1.TilemapShapes()
	Collisions[LaneGarrison2] = g2.TilemapShapes()
	Collisions[LaneGarrison3] = g3.TilemapShapes()
	Collisions[LaneGarrison4] = g4.TilemapShapes()
	Collisions[LaneGarrison5] = g5.TilemapShapes()
	Collisions[LaneGarrisonPlus1] = p1.TilemapShapes()
	Collisions[LaneGarrisonPlus2] = p2.TilemapShapes()
	Collisions[LaneGarrisonPlus3] = p3.TilemapShapes()
	Collisions[LaneGarrisonPlus4] = p4.TilemapShapes()
	Collisions[LaneGarrisonPlus5] = p5.TilemapShapes()

	mirrorGarrisonLanes()

	MapLayer, Map = layers[9], graphics.NewTilemap(layers[9])
	Flags, Grid = graphics.NewTilemap(layers[10]), graphics.NewTilemap(layers[11])

	AllyBase = NewBase(BaseFortress, Garrison3, true)
	EnemyBase = NewBase(BaseFort, Garrison0, false)

	Units = append(Units, NewUnit(CharacterMan, TeamAlly, LaneUpper))
	Units = append(Units, NewUnit(CharacterMan, TeamEnemy, LaneUpper))
	Units = append(Units, NewUnit(CharacterWoman, TeamEnemy, LaneMiddle))
	Units = append(Units, NewUnit(CharacterHunter, TeamAlly, LaneGarrisonPlus1))
	Units = append(Units, NewUnit(CharacterHunter, TeamAlly, LaneGarrisonPlus2))
	Units = append(Units, NewUnit(CharacterHunter, TeamAlly, LaneGarrisonPlus3))
	Units = append(Units, NewUnit(CharacterHunter, TeamAlly, LaneGarrisonPlus4))
	Units = append(Units, NewUnit(CharacterHunter, TeamAlly, LaneGarrisonPlus5))
	Units = append(Units, NewUnit(CharacterHunter, TeamAlly, LaneGarrison1))
	Units = append(Units, NewUnit(CharacterHunter, TeamAlly, LaneGarrison2))
	Units = append(Units, NewUnit(CharacterHunter, TeamAlly, LaneGarrison3))
	Units = append(Units, NewUnit(CharacterHunter, TeamAlly, LaneGarrison4))
	Units = append(Units, NewUnit(CharacterHunter, TeamAlly, LaneGarrison5))

	Entrances = [6]*Entrance{
		LaneLower:      NewEntrance(EntranceTallGate, TeamAlly, LaneLower),
		LaneMiddle:     NewEntrance(EntranceTallGate, TeamAlly, LaneMiddle),
		LaneUpper:      NewEntrance(EntranceTallGate, TeamAlly, LaneUpper),
		LaneLower + 3:  NewEntrance(EntranceDoor, TeamEnemy, LaneLower),
		LaneMiddle + 3: NewEntrance(EntranceDoor, TeamEnemy, LaneMiddle),
		LaneUpper + 3:  NewEntrance(EntranceDoor, TeamEnemy, LaneUpper),
	}

	PlayAmbience(EnvironmentPlains)
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

	View.DrawObject(AllyBase.Back)
	View.DrawObject(AllyBase.GarrisonBack)
	View.DrawObject(EnemyBase.Back)
	View.DrawObject(EnemyBase.GarrisonBack)

	for _, g := range Entrances {
		g.Update()
	}
	View.DrawObject(&Map)
	View.DrawObject(&Flags)

	collection.SortByField(Units, func(u *Unit) float32 {
		if u.Stats.Health <= 0 { // dead units go behind all alive units
			return number.NegativeInfinity()
		}
		return u.Y // fall back to Y sort
	})
	for _, u := range Units {
		if u != nil { // unit may have died, faded out & removed during an update
			u.Update()
		}
	}

	View.DrawObject(AllyBase.Front)
	View.DrawObject(AllyBase.GarrisonFront)
	View.DrawObject(EnemyBase.Front)
	View.DrawObject(EnemyBase.GarrisonFront)

	for _, p := range Projectiles {
		if p != nil { // projectile may have been destroyed, faded out & removed during an update
			p.Update()
		}
	}

	for _, g := range Entrances {
		g.HealthBar.Update(g.Tiles[0].Shape, g.Health, g.MaxHealth, geometry.Area{})
	}
	for _, u := range Units { // health bars take the Z order of the units
		var hb = u.Hitbox()
		hb.Height += 8
		u.HealthBar.Update(hb, u.Stats.Health, Characters[u.Character].Stats.Health, u.Mask)
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

func DeltaTimeScaled() float32 {
	return time.Delta() * TimeScale
}

func DrawShadow(x, z, width, height, angle float32, mask geometry.Area) {
	var lower, upper = Collisions[LaneLower][0], Collisions[LaneUpper][0]
	var y = number.Map(z, 0, 2, lower.Y-lower.Height/2, upper.Y-upper.Height/2)
	View.DrawShape(x, y, width, height, angle, 1, color.RGBA(0, 0, 0, 100), mask)
}

// private ========================================================

var skyColor = color.TagRGBA("rgb(98, 171, 212)")

func mirrorGarrisonLanes() {
	for i := LaneGarrison1; i < LaneGarrisonPlus5+1; i++ {
		var length = len(Collisions[i])
		for j := range length {
			var shape = Collisions[i][j]
			shape.X *= -1
			if shape.Angle != 0 {
				shape.Angle += 90
			}
			Collisions[i] = append(Collisions[i], shape)
		}
	}
}
func bringGarrisonLanesDown(ally bool) {
	for i := LaneGarrison1; i < LaneGarrisonPlus5+1; i++ {
		var length = len(Collisions[i])
		var j = length / 2
		if ally {
			j, length = 0, length/2
		}
		for ; j < length; j++ {
			var shape = Collisions[i][j]
			shape.Y += TileSize
			Collisions[i][j] = shape
		}
	}
}
