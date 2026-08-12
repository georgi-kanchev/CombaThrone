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

const TileSize, MapCount = 32, 4

const LaneLower, LaneMiddle, LaneUpper Lane = 0, 1, 2
const LaneLowerGarrison, LaneMiddleGarrison, LaneUpperGarrison Lane = 3, 4, 5
const LaneGarrison1, LaneGarrison2, LaneGarrison3 Lane = 6, 7, 8

var TimeScale float32 = 1

var View graphics.View
var Background graphics.Object

var MapLayer assets.TileLayerId
var Map, Grid, AllyBase, EnemyBase, Flags graphics.Object
var TilesetCrops, UserInterfaceCrops assets.AtlasId
var UserInterfaceImg assets.ImageId

var Collisions = make([][]geometry.Shape, 10)
var Masks = [9]geometry.Area{
	LaneLower:  geometry.NewArea(0, 0, 556, 1000),
	LaneMiddle: geometry.NewArea(0, 0, 492, 1000),
	LaneUpper:  geometry.NewArea(0, 0, 428, 1000),
}
var Entrances [6]*EntranceData // index is Ally[Lower, Middle, Upper], Enemy[Lower, Middle, Upper]
var Units []*Unit
var Projectiles []*Projectile

func InitScene() {
	View = graphics.NewView(5.68)
	Background = graphics.NewSprite(0, 0, 1, assets.LoadImage("data/bgr-field.png"))
	UserInterfaceImg = assets.LoadImage("data/user-interface.png")
	UserInterfaceCrops = assets.LoadAtlas(UserInterfaceImg, "data/user-interface.xml")

	var layers, tileset = assets.LoadTileLayersFromTiled("data/map.tmx")
	var low, mid, up = graphics.NewTilemap(layers[12]), graphics.NewTilemap(layers[13]), graphics.NewTilemap(layers[14])
	var lowG, midG, upG = graphics.NewTilemap(layers[15]), graphics.NewTilemap(layers[16]), graphics.NewTilemap(layers[17])
	var g1, g2, g3 = graphics.NewTilemap(layers[18]), graphics.NewTilemap(layers[19]), graphics.NewTilemap(layers[20])

	TilesetCrops = assets.LoadAtlas(tileset, "data/tileset.xml")
	TilesetCrops = assets.LoadAtlas(tileset, "data/tileset.xml")
	Collisions[LaneLower] = low.TilemapShapes()
	Collisions[LaneMiddle] = mid.TilemapShapes()
	Collisions[LaneUpper] = up.TilemapShapes()
	Collisions[LaneLowerGarrison] = lowG.TilemapShapes()
	Collisions[LaneMiddleGarrison] = midG.TilemapShapes()
	Collisions[LaneUpperGarrison] = upG.TilemapShapes()
	Collisions[LaneGarrison1] = g1.TilemapShapes()
	Collisions[LaneGarrison2] = g2.TilemapShapes()
	Collisions[LaneGarrison3] = g3.TilemapShapes()

	mirrorGarrisonLanes()
	liftGarrisonLanesUp(true)

	MapLayer, Map = layers[9], graphics.NewTilemap(layers[9])
	Flags, Grid = graphics.NewTilemap(layers[10]), graphics.NewTilemap(layers[11])
	AllyBase = graphics.NewTilemap(layers[5])
	EnemyBase = graphics.NewTilemap(layers[1])
	EnemyBase.Width *= -1

	Units = append(Units, NewUnit(CharacterHunter, TeamAlly, LaneUpper))
	// Units = append(Units, NewUnit(CharacterMan, TeamEnemy, LaneUpper))
	// Units = append(Units, NewUnit(CharacterWoman, TeamEnemy, LaneMiddle))
	// Units = append(Units, NewUnit(CharacterHunter, TeamEnemy, LaneLower))
	Units = append(Units, NewUnit(CharacterHunter, TeamAlly, LaneUpperGarrison))
	Units = append(Units, NewUnit(CharacterHunter, TeamAlly, LaneMiddleGarrison))
	Units = append(Units, NewUnit(CharacterHunter, TeamAlly, LaneLowerGarrison))

	Entrances = [6]*EntranceData{
		LaneLower:      NewEntrance(EntranceDoor, TeamAlly, LaneLower),
		LaneMiddle:     NewEntrance(EntranceTallGate, TeamAlly, LaneMiddle),
		LaneUpper:      NewEntrance(EntranceHole, TeamAlly, LaneUpper),
		LaneLower + 3:  NewEntrance(EntranceDoor, TeamEnemy, LaneLower),
		LaneMiddle + 3: NewEntrance(EntranceDoor, TeamEnemy, LaneMiddle),
		LaneUpper + 3:  NewEntrance(EntranceShortGate, TeamEnemy, LaneUpper),
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
		if u != nil { // unit may have died, faded out & removed during an update
			u.Update()
		}
	}

	for _, p := range Projectiles {
		if p != nil { // projectile may have been destroyed, faded out & removed during an update
			p.Update()
		}
	}

	for _, g := range Entrances {
		g.HealthBar.Update(g.Tiles[0].Shape, g.Health, g.MaxHealth, geometry.Area{})
	}
	for _, u := range Units { // health bars take the Z order of the units
		u.HealthBar.Update(u.Shape, u.Stats.Health, Characters[u.Character].Stats.Health, u.Mask)
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
	var l = []Lane{LaneLowerGarrison, LaneMiddleGarrison, LaneUpperGarrison, LaneGarrison1, LaneGarrison2, LaneGarrison3}
	for _, lane := range l {
		var length = len(Collisions[lane])
		for i := range length {
			var shape = Collisions[lane][i]
			shape.X *= -1
			if shape.Angle != 0 {
				shape.Angle += 90
			}
			Collisions[lane] = append(Collisions[lane], shape)
		}
	}
}
func liftGarrisonLanesUp(ally bool) {
	var l = []Lane{LaneLowerGarrison, LaneMiddleGarrison, LaneUpperGarrison, LaneGarrison1, LaneGarrison2, LaneGarrison3}
	for _, lane := range l {
		var length = len(Collisions[lane])
		var i = length / 2
		if ally {
			i, length = 0, length/2
		}
		for ; i < length; i++ {
			var shape = Collisions[lane][i]
			shape.Y -= TileSize
			Collisions[lane][i] = shape
		}
	}
}
