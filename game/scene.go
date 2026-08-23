package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/cursor"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/text"
	"pure-game-kit/packages/utility/time"
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
	LaneGarrisonPlus1
	LaneGarrisonPlus2
	LaneGarrisonPlus3
	LaneCount
	LaneLayerOffset = 28
)

const TileSize, MapCount = 32.0, 4

var Decor, UserInterface assets.AtlasId

var TimeScale float32 = 1
var View *graphics.View
var CurrentZone *Zone
var Layers []assets.TileLayerId
var SortedY []*graphics.Object // for units & pickups

var Player *PlayerState

var UI *GUI

func InitScene() {
	var view = graphics.NewView(1)
	View = &view
	UserInterface = assets.LoadAtlas(assets.LoadImage("data/user-interface.png"), "data/user-interface.xml")
	UI = NewUI()

	var layers, decor = assets.LoadTileLayersFromTiled("data/map.tmx")
	Layers = layers
	Decor = assets.LoadAtlas(decor, "data/decor.xml")

	assets.FontId(0).EmbedImage(text.At(Tags[IconGlory], 0), UserInterface.Crops("icons")[IconGlory])
	assets.FontId(0).EmbedImage(text.At(Tags[IconHealth], 0), UserInterface.Crops("icons")[IconHealth])
	assets.FontId(0).EmbedImage(text.At(Tags[IconCoin], 0), UserInterface.Crops("icons")[IconCoin])

	for i := range LaneCount {
		var tilemap = graphics.NewTilemap(layers[Lane(LaneLayerOffset)+i])
		laneCollisions[i] = tilemap.TilemapShapes()
	}
	mirrorGarrisonLanes()

	Clouds[CloudsNone] = []assets.ImageId{0}
	Clouds[CloudsWindy] = append(Clouds[CloudsWindy], assets.LoadImage("data/zones/sky-clouds-wind.png"))
	for i := 1; i <= 4; i++ {
		Clouds[CloudsNormal] = append(Clouds[CloudsNormal], assets.LoadImage(text.New("data/zones/sky-clouds", i, ".png")))
	}
	for i, name := range zoneNames {
		ZoneBackgrounds[i] = assets.LoadImage("data/zones/background-" + name + ".png")
		Zones[i] = NewZone(ZoneKind(i))
	}
	CurrentZone = Zones[ZoneField]

	Bases[TeamAlly] = NewBase(TeamAlly, BaseFortress, Garrison3,
		[3]EntranceKind{EntranceDoor, EntranceTallGate, EntranceShortGate})
	Bases[TeamEnemy] = NewBase(TeamEnemy, BaseNone, Garrison3,
		[3]EntranceKind{EntranceNone, EntranceNone, EntranceNone})

	// Units = append(Units, NewUnit(CharWoman, TeamAlly, LaneMiddle))
	// Units = append(Units, NewUnit(CharMan, TeamEnemy, LaneUpper))
	// Units = append(Units, NewUnit(CharMan, TeamEnemy, LaneMiddle))
	// Units = append(Units, NewUnit(CharHunter, TeamEnemy, LaneLower))
	// Units = append(Units, NewUnit(CharHunter, TeamEnemy, LaneLower))
	// Units = append(Units, NewUnit(CharMan, TeamEnemy, LaneMiddle))

	Pickups = append(Pickups, NewPickup(0, PickupRelic, LaneLowerOff))

	PlayAmbience(CurrentZone.kind)

	Player = NewPlayer()
	Player.Units[0] = NewUnit(CharHunter, TeamAlly, LaneMiddle)
	Player.Units[1] = NewUnit(CharWoman, TeamAlly, LaneMiddle)
	Player.Units[2] = NewUnit(CharMan, TeamAlly, LaneMiddle)
}

//=================================================================

func UpdateScene() {
	var _, bly = CurrentZone.Ground.PointFromEdge(0.5, 1)
	View.FitSize(CurrentZone.Ground.Width, 0)
	var _, h = View.Size()
	View.Y = (bly - h/2) - 2

	UI.Tooltip = nil
	mouse.SetCursor(cursor.Default)

	if keyboard.IsKeyJustPressed(key.RightArrow) && CurrentZone.kind < ZoneHell {
		CurrentZone = Zones[CurrentZone.kind+1]
	}
	if keyboard.IsKeyJustPressed(key.LeftArrow) && CurrentZone.kind > ZoneField {
		CurrentZone = Zones[CurrentZone.kind-1]
	}

	CurrentZone.UpdateBack()
	Bases[TeamAlly].UpdateBack()
	Bases[TeamEnemy].UpdateBack()
	iterateRemovable(&Bases[TeamAlly].Entrances, func(e *Entrance) { e.Update() })
	iterateRemovable(&Bases[TeamEnemy].Entrances, func(e *Entrance) { e.Update() })
	CurrentZone.UpdateFront()

	iterateRemovable(&ProjectilesBehind, func(p *Projectile) { p.Update() })
	iterateRemovable(&Pickups, func(p *Pickup) { p.Update() })

	iterateRemovable(&Player.Units, func(u *Unit) {
		if !u.IsSummoned() {
			u.Update()
		}
	})
	collection.SortByField(Units, func(u *Unit) float32 {
		if u.Stats.Health <= 0 { // dead units go behind all alive units
			return number.NegativeInfinity()
		}
		return u.Y // fall back to Y sort
	})
	iterateRemovable(&Units, func(u *Unit) {
		if u.State == StateWaitingToBeSummoned {
			u.State = StateSummoned
		}

		u.Update()
	})

	Bases[TeamAlly].UpdateFront()
	Bases[TeamEnemy].UpdateFront()
	iterateRemovable(&Projectiles, func(p *Projectile) { p.Update() })

	for _, g := range Bases[TeamAlly].Entrances {
		g.HealthBar.Update(g.Tiles[0].Shape, g.Health, g.MaxHealth, geometry.Area{})
	}
	for _, g := range Bases[TeamEnemy].Entrances {
		g.HealthBar.Update(g.Tiles[0].Shape, g.Health, g.MaxHealth, geometry.Area{})
	}
	for _, u := range Units { // health bars take the Z order of the units
		var hb = u.Hitbox()
		hb.Height += 8
		u.HealthBar.Update(hb, u.Stats.Health, Characters[u.Character].Stats.Health, u.Mask)
	}

	UI.Update()
	Player.Update()
}
func DrawShadow(x, z, width, height, angle float32, mask geometry.Area) {
	var lower, upper = laneCollisions[LaneLower][0], laneCollisions[LaneUpper][0]
	var y = number.Map(z, 0, 2, lower.Y-lower.Height/2, upper.Y-upper.Height/2)
	View.DrawShape(x, y, width, height, angle, 1, color.RGBA(0, 0, 0, 100), mask)
}

func PointAtCell(cellX, cellY float32) (x, y float32) {
	var tw, th = Layers[0].TileSize()
	var cols, rows = Layers[0].Size()
	return (cellX-float32(cols)/2)*tw + (tw / 2), (cellY-float32(rows)/2)*th + (th / 2)
}
func CellAtPoint(x, y float32) (cellX, cellY float32) {
	var tw, th = Layers[0].TileSize()
	var cols, rows = Layers[0].Size()
	return x/tw + float32(cols)/2, y/th + float32(rows)/2
}
func TileAtCell(cellX, cellY int, layer assets.TileLayerId) assets.Tile {
	return layer.TileAtCell(cellX, cellY)
}
func DeltaTimeScaled() float32 {
	return time.Delta() * TimeScale
}

// private ========================================================

var highlightCursorColors = map[int]uint{
	cursor.Arrow: palette.LightGray, cursor.Hand: palette.White, cursor.NotAllowed: palette.Red,
}

func mirrorGarrisonLanes() {
	for i := LaneGarrison1; i < LaneGarrisonPlus3+1; i++ {
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
	for i := LaneGarrison1; i < LaneGarrisonPlus3+1; i++ {
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
