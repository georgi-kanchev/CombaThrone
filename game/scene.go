package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/text"
	"pure-game-kit/packages/utility/time"
)

type Layer uint8
type Environment uint8

const ( // layer tilemaps
	LayerCamp Layer = iota
	LayerBarrack
	LayerFortBack0
	LayerFortFront0
	LayerFortBack1
	LayerFortFront1
	LayerFortBack2
	LayerFortFront2
	LayerFortBack3
	LayerFortFront3

	LayerField
	LayerFieldNoBase
	LayerFieldBuildings
	LayerDesert
	LayerDesertNoBase
	LayerDesertBuildings

	LayerFlags
	LayerGrid
	LayerCount
)

const ( // lanes (collision layers)
	LaneLower Lane = iota
	LaneLowerOff
	LaneMiddle
	LaneMiddleOff
	LaneUpper
	LaneUpperOff
	LaneGarrison1
	LaneGarrison2
	LaneGarrison3
	LaneGarrison4
	LaneGarrison5
	LaneGarrisonPlus1
	LaneGarrisonPlus2
	LaneGarrisonPlus3
	LaneGarrisonPlus4
	LaneGarrisonPlus5
	LaneCount
)

const EnvironmentField, EnvironmentDesert Environment = 0, 1

const TileSize, MapCount = 32, 4

var TimeScale float32 = 1
var View graphics.View
var Background, Map, MapNoBase, MapBuildings, Grid, Flags *graphics.Object
var Decor, UserInterface assets.AtlasId
var Layers [LayerCount]assets.TileLayerId

var SortedY []*graphics.Object // for units & pickups

func InitScene() {
	View = graphics.NewView(1)
	UserInterface = assets.LoadAtlas(assets.LoadImage("data/user-interface.png"), "data/user-interface.xml")

	assets.FontId(0).EmbedImage(text.At(Tags[IconGlory], 0), UserInterface.Crops("icons")[IconGlory])
	assets.FontId(0).EmbedImage(text.At(Tags[IconHealth], 0), UserInterface.Crops("icons")[IconHealth])

	var layers, decor = assets.LoadTileLayersFromTiled("data/map.tmx")
	var low = graphics.NewTilemap(layers[Lane(LayerCount)+LaneLower])
	var lowOff = graphics.NewTilemap(layers[Lane(LayerCount)+LaneLowerOff])
	var mid = graphics.NewTilemap(layers[Lane(LayerCount)+LaneMiddle])
	var midOff = graphics.NewTilemap(layers[Lane(LayerCount)+LaneMiddleOff])
	var up = graphics.NewTilemap(layers[Lane(LayerCount)+LaneUpper])
	var upOff = graphics.NewTilemap(layers[Lane(LayerCount)+LaneUpperOff])
	var g1 = graphics.NewTilemap(layers[Lane(LayerCount)+LaneGarrison1])
	var g2 = graphics.NewTilemap(layers[Lane(LayerCount)+LaneGarrison2])
	var g3 = graphics.NewTilemap(layers[Lane(LayerCount)+LaneGarrison3])
	var g4 = graphics.NewTilemap(layers[Lane(LayerCount)+LaneGarrison4])
	var g5 = graphics.NewTilemap(layers[Lane(LayerCount)+LaneGarrison5])
	var p1 = graphics.NewTilemap(layers[Lane(LayerCount)+LaneGarrisonPlus1])
	var p2 = graphics.NewTilemap(layers[Lane(LayerCount)+LaneGarrisonPlus2])
	var p3 = graphics.NewTilemap(layers[Lane(LayerCount)+LaneGarrisonPlus3])
	var p4 = graphics.NewTilemap(layers[Lane(LayerCount)+LaneGarrisonPlus4])
	var p5 = graphics.NewTilemap(layers[Lane(LayerCount)+LaneGarrisonPlus5])

	for i := range LayerCount {
		Layers[i] = layers[i]
	}

	Decor = assets.LoadAtlas(decor, "data/decor.xml")
	laneCollisions[LaneLower] = low.TilemapShapes()
	laneCollisions[LaneLowerOff] = lowOff.TilemapShapes()
	laneCollisions[LaneMiddle] = mid.TilemapShapes()
	laneCollisions[LaneMiddleOff] = midOff.TilemapShapes()
	laneCollisions[LaneUpper] = up.TilemapShapes()
	laneCollisions[LaneUpperOff] = upOff.TilemapShapes()
	laneCollisions[LaneGarrison1] = g1.TilemapShapes()
	laneCollisions[LaneGarrison2] = g2.TilemapShapes()
	laneCollisions[LaneGarrison3] = g3.TilemapShapes()
	laneCollisions[LaneGarrison4] = g4.TilemapShapes()
	laneCollisions[LaneGarrison5] = g5.TilemapShapes()
	laneCollisions[LaneGarrisonPlus1] = p1.TilemapShapes()
	laneCollisions[LaneGarrisonPlus2] = p2.TilemapShapes()
	laneCollisions[LaneGarrisonPlus3] = p3.TilemapShapes()
	laneCollisions[LaneGarrisonPlus4] = p4.TilemapShapes()
	laneCollisions[LaneGarrisonPlus5] = p5.TilemapShapes()

	mirrorGarrisonLanes()

	var bgr = graphics.NewSprite(0, 0, 1, assets.LoadImage("data/environments/background-desert.png"))
	var _map, mapNoBase = graphics.NewTilemap(layers[LayerDesert]), graphics.NewTilemap(layers[LayerDesertNoBase])
	var mapBuildings = graphics.NewTilemap(layers[LayerDesertBuildings])
	Background, Map, MapNoBase, MapBuildings = &bgr, &_map, &mapNoBase, &mapBuildings

	var flag, grid = graphics.NewTilemap(layers[LayerFlags]), graphics.NewTilemap(layers[LayerGrid])
	Flags, Grid = &flag, &grid

	AllyBase = NewBase(TeamAlly, SaveState{
		Kind: BaseCamp, Garrison: Garrison3, EntranceKinds: [3]EntranceKind{
			EntranceNone, EntranceNone, EntranceNone,
		},
	})
	EnemyBase = NewBase(TeamEnemy, SaveState{
		Kind: BaseNone, Garrison: Garrison3, EntranceKinds: [3]EntranceKind{
			EntranceNone, EntranceNone, EntranceNone,
		},
	})

	// Units = append(Units, NewUnit(CharWoman, TeamAlly, LaneUpperOff))
	// Units = append(Units, NewUnit(CharWoman, TeamEnemy, LaneUpperOff))
	// Units = append(Units, NewUnit(CharMan, TeamAlly, LaneUpper))
	// Units = append(Units, NewUnit(CharHunter, TeamAlly, LaneLower))
	// Units = append(Units, NewUnit(CharMan, TeamEnemy, LaneMiddle))
	// Units = append(Units, NewUnit(CharHunter, TeamEnemy, LaneGarrisonPlus1))
	// Units = append(Units, NewUnit(CharHunter, TeamEnemy, LaneGarrisonPlus3))
	// Units = append(Units, NewUnit(CharHunter, TeamEnemy, LaneGarrisonPlus4))
	// Units = append(Units, NewUnit(CharHunter, TeamEnemy, LaneGarrisonPlus5))
	// Units = append(Units, NewUnit(CharHunter, TeamEnemy, LaneGarrison1))
	// Units = append(Units, NewUnit(CharHunter, TeamEnemy, LaneGarrison2))
	// Units = append(Units, NewUnit(CharHunter, TeamEnemy, LaneGarrison3))
	// Units = append(Units, NewUnit(CharHunter, TeamEnemy, LaneGarrison4))
	// Units = append(Units, NewUnit(CharHunter, TeamEnemy, LaneGarrison5))

	Pickups = append(Pickups, NewPickup(0, PickupRelic, LaneLowerOff))

	PlayAmbience(EnvironmentField)
}

//=================================================================

func UpdateScene() {
	var _, bly = Background.PointFromEdge(0.5, 1)
	View.FitSize(Background.Width, 0)
	var _, h = View.Size()
	View.Y = (bly - h/2) - 2

	View.DrawColor(skyColor)
	View.DrawObject(Background)
	View.DrawObject(MapBuildings)

	AllyBase.UpdateBack()
	EnemyBase.UpdateBack()

	View.DrawObject(Flags)
	AllyBase.UpdateMiddle()
	EnemyBase.UpdateMiddle()

	iterateRemovable(&AllyBase.Entrances, func(e *Entrance) { e.Update() })
	iterateRemovable(&EnemyBase.Entrances, func(e *Entrance) { e.Update() })

	if AllyBase.Kind < BaseBarrack {
		MapNoBase.Width = number.Absolute(MapNoBase.Width)
		View.DrawObject(MapNoBase)
	}
	if EnemyBase.Kind < BaseBarrack {
		MapNoBase.Width = -number.Absolute(MapNoBase.Width)
		View.DrawObject(MapNoBase)
	}
	View.DrawObject(Map)

	iterateRemovable(&ProjectilesBehind, func(p *Projectile) { p.Update() })
	iterateRemovable(&Pickups, func(p *Pickup) { p.Update() })

	collection.SortByField(Units, func(u *Unit) float32 {
		if u.Stats.Health <= 0 { // dead units go behind all alive units
			return number.NegativeInfinity()
		}
		return u.Y // fall back to Y sort
	})
	iterateRemovable(&Units, func(u *Unit) { u.Update() })
	AllyBase.UpdateFront()
	EnemyBase.UpdateFront()
	iterateRemovable(&Projectiles, func(p *Projectile) { p.Update() })

	for _, g := range AllyBase.Entrances {
		g.HealthBar.Update(g.Tiles[0].Shape, g.Health, g.MaxHealth, geometry.Area{})
	}
	for _, g := range EnemyBase.Entrances {
		g.HealthBar.Update(g.Tiles[0].Shape, g.Health, g.MaxHealth, geometry.Area{})
	}
	for _, u := range Units { // health bars take the Z order of the units
		var hb = u.Hitbox()
		hb.Height += 8
		u.HealthBar.Update(hb, u.Stats.Health, Characters[u.Character].Stats.Health, u.Mask)
	}
}

func PointAtCell(cellX, cellY float32) (x, y float32) {
	var tw, th = Layers[LayerGrid].TileSize()
	var cols, rows = Layers[LayerGrid].Size()
	return (cellX-float32(cols)/2)*tw + (tw / 2), (cellY-float32(rows)/2)*th + (th / 2)
}
func CellAtPoint(x, y float32) (cellX, cellY float32) {
	var tw, th = Layers[LayerGrid].TileSize()
	var cols, rows = Layers[LayerGrid].Size()
	return x/tw + float32(cols)/2, y/th + float32(rows)/2
}
func TileAtCell(cellX, cellY int, layer assets.TileLayerId) assets.Tile {
	return layer.TileAtCell(cellX, cellY)
}

func DeltaTimeScaled() float32 {
	return time.Delta() * TimeScale
}

func DrawShadow(x, z, width, height, angle float32, mask geometry.Area) {
	var lower, upper = laneCollisions[LaneLower][0], laneCollisions[LaneUpper][0]
	var y = number.Map(z, 0, 2, lower.Y-lower.Height/2, upper.Y-upper.Height/2)
	View.DrawShape(x, y, width, height, angle, 1, color.RGBA(0, 0, 0, 100), mask)
}

// private ========================================================

var skyColor = color.TagRGBA("rgb(98, 171, 212)")

// var skyColor = color.TagRGBA("rgb(227, 177, 109)")

func mirrorGarrisonLanes() {
	for i := LaneGarrison1; i < LaneGarrisonPlus5+1; i++ {
		var length = len(laneCollisions[i])
		for j := range length {
			var shape = laneCollisions[i][j]
			shape.X *= -1
			if shape.Angle != 0 {
				shape.Angle += 90
			}
			laneCollisions[i] = append(laneCollisions[i], shape)
		}
	}
}
func bringGarrisonLanesDown(team Team) {
	for i := LaneGarrison1; i < LaneGarrisonPlus5+1; i++ {
		var length = len(laneCollisions[i])
		var j = length / 2
		if team == TeamAlly {
			j, length = 0, length/2
		}
		for ; j < length; j++ {
			var shape = laneCollisions[i][j]
			shape.Y += TileSize
			laneCollisions[i][j] = shape
		}
	}
}
func iterateRemovable[T any](things *[]*T, iteration func(*T)) {
	for i := 0; i < len(*things); i++ {
		var lastLength = len(*things)
		iteration((*things)[i])
		if lastLength != len(*things) {
			i-- // was removed during update
		}
	}
}
