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
	LayerPlains
	LayerPlainsNoBase
	LayerFlags
	LayerGrid
	LayerCount
)

const ( // lanes (collision layers)
	LaneLower Lane = iota
	LaneMiddle
	LaneUpper
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

const EnvironmentPlains Environment = 0

const TileSize, MapCount = 32, 4

var TimeScale float32 = 1
var View graphics.View
var Background, MapNoBase, Map, Grid, Flags *graphics.Object
var DecorCrops, UserInterfaceCrops assets.AtlasId
var Layers [LayerCount]assets.TileLayerId

func InitScene() {
	View = graphics.NewView(1)
	UserInterfaceCrops = assets.LoadAtlas(assets.LoadImage("data/user-interface.png"), "data/user-interface.xml")

	assets.FontId(0).EmbedImage(text.At(Tags[IconGlory], 0), UserInterfaceCrops.Crops("icons")[IconGlory])
	assets.FontId(0).EmbedImage(text.At(Tags[IconHealth], 0), UserInterfaceCrops.Crops("icons")[IconHealth])

	var layers, decor = assets.LoadTileLayersFromTiled("data/map.tmx")
	var low = graphics.NewTilemap(layers[Lane(LayerCount)+LaneLower])
	var mid = graphics.NewTilemap(layers[Lane(LayerCount)+LaneMiddle])
	var up = graphics.NewTilemap(layers[Lane(LayerCount)+LaneUpper])
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

	DecorCrops = assets.LoadAtlas(decor, "data/decor.xml")
	DecorCrops = assets.LoadAtlas(decor, "data/decor.xml")
	laneCollisions[LaneLower] = low.TilemapShapes()
	laneCollisions[LaneMiddle] = mid.TilemapShapes()
	laneCollisions[LaneUpper] = up.TilemapShapes()
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

	var bgr = graphics.NewSprite(0, 0, 1, assets.LoadImage("data/bgr-field.png"))
	var _map, mapNoBase = graphics.NewTilemap(layers[LayerPlains]), graphics.NewTilemap(layers[LayerPlainsNoBase])
	Background, Map, MapNoBase = &bgr, &_map, &mapNoBase

	var flag, grid = graphics.NewTilemap(layers[LayerFlags]), graphics.NewTilemap(layers[LayerGrid])
	Flags, Grid = &flag, &grid
	Flags.Y += TileSize * 0.1

	AllyBase = NewBase(SaveState{Kind: BaseBarrack}, true)
	EnemyBase = NewBase(SaveState{Kind: BaseFort, Garrison: Garrison1}, false)

	Units = append(Units, NewUnit(CharHunter, TeamAlly, LaneMiddle))
	Units = append(Units, NewUnit(CharMan, TeamAlly, LaneUpper))
	Units = append(Units, NewUnit(CharMan, TeamEnemy, LaneUpper))
	Units = append(Units, NewUnit(CharWoman, TeamEnemy, LaneMiddle))
	// Units = append(Units, NewUnit(CharHunter, TeamEnemy, LaneGarrison5))
	// Units = append(Units, NewUnit(CharHunter, TeamAlly, LaneGarrisonPlus2))
	// Units = append(Units, NewUnit(CharHunter, TeamAlly, LaneGarrisonPlus3))
	// Units = append(Units, NewUnit(CharHunter, TeamAlly, LaneGarrisonPlus4))
	// Units = append(Units, NewUnit(CharHunter, TeamAlly, LaneGarrisonPlus5))
	// Units = append(Units, NewUnit(CharHunter, TeamAlly, LaneGarrison1))
	// Units = append(Units, NewUnit(CharHunter, TeamAlly, LaneGarrison2))
	// Units = append(Units, NewUnit(CharHunter, TeamAlly, LaneGarrison3))
	// Units = append(Units, NewUnit(CharHunter, TeamAlly, LaneGarrison4))
	// Units = append(Units, NewUnit(CharHunter, TeamAlly, LaneGarrison5))

	Entrances = [6]*Entrance{
		LaneUpper:      NewEntrance(EntranceDoor, TeamAlly, LaneUpper),
		LaneMiddle:     NewEntrance(EntranceDoor, TeamAlly, LaneMiddle),
		LaneLower:      NewEntrance(EntranceNone, TeamAlly, LaneLower),
		LaneUpper + 3:  NewEntrance(EntranceNone, TeamEnemy, LaneUpper),
		LaneMiddle + 3: NewEntrance(EntranceNone, TeamEnemy, LaneMiddle),
		LaneLower + 3:  NewEntrance(EntranceNone, TeamEnemy, LaneLower),
	}

	Pickups = append(Pickups, NewPickup(0, 58, 2.75, PickupRelic))

	PlayAmbience(EnvironmentPlains)
}

//=================================================================

func UpdateScene() {
	var _, bly = Background.PointFromEdge(0.5, 1)
	View.FitSize(Background.Width, 0)
	var _, h = View.Size()
	View.Y = (bly - h/2) - 2

	View.DrawColor(skyColor)
	View.DrawObject(Background)

	AllyBase.UpdateBack()
	EnemyBase.UpdateBack()

	for _, g := range Entrances {
		g.Update()
	}

	if AllyBase.Kind < BaseBarrack {
		MapNoBase.Width = number.Absolute(MapNoBase.Width)
		View.DrawObject(MapNoBase)
	}
	if EnemyBase.Kind < BaseBarrack {
		MapNoBase.Width = -number.Absolute(MapNoBase.Width)
		View.DrawObject(MapNoBase)
	}
	View.DrawObject(Map)
	View.DrawObject(Flags)

	for _, p := range ProjectilesBehind {
		if p != nil { // may have been faded out & removed during an update
			p.Update()
		}
	}

	for _, p := range Pickups {
		if p != nil { // may have been removed during an update
			p.Update()
		}
	}

	collection.SortByField(Units, func(u *Unit) float32 {
		if u.Stats.Health <= 0 { // dead units go behind all alive units
			return number.NegativeInfinity()
		}
		return u.Y // fall back to Y sort
	})
	for _, u := range Units {
		if u != nil { // may have died, faded out & removed during an update
			u.Update()
		}
	}

	AllyBase.UpdateFront()
	EnemyBase.UpdateFront()

	for _, p := range Projectiles {
		if p != nil { // may have been faded out & removed during an update
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
func bringGarrisonLanesDown(ally bool) {
	for i := LaneGarrison1; i < LaneGarrisonPlus5+1; i++ {
		var length = len(laneCollisions[i])
		var j = length / 2
		if ally {
			j, length = 0, length/2
		}
		for ; j < length; j++ {
			var shape = laneCollisions[i][j]
			shape.Y += TileSize
			laneCollisions[i][j] = shape
		}
	}
}
