package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/gui"
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
	"pure-game-kit/packages/window"
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

var TitleScreen assets.ImageId
var InGame bool
var InGameTimer float32

var Decor assets.AtlasId

var TimeScale float32 = 1
var View *graphics.View
var CurrentZone *Zone
var Layers []assets.TileLayerId
var SortedY []*graphics.Object // for units & pickups

var Player *PlayerState

func InitScene() {
	TitleScreen = assets.LoadImage("data/title-screen.png")

	var view = graphics.NewView(1)
	View = &view
	GameHUD = NewHUD()

	var layers, decor = assets.LoadTileLayersFromTiled("data/map.tmx")
	Layers = layers
	Decor = assets.LoadAtlas(decor, "data/decor.xml")

	for i := range LaneCount {
		var tilemap = graphics.NewTilemap(layers[Lane(LaneLayerOffset)+i])
		Collisions[i] = tilemap.TilemapShapes()
	}
	mirrorGarrisonLanes()

	Clouds[CloudsNone] = []assets.ImageId{0}
	Clouds[CloudsWindy] = append(Clouds[CloudsWindy], assets.LoadImage("data/zones/sky-clouds-wind.png"))
	for i := 1; i <= 4; i++ {
		Clouds[CloudsNormal] = append(Clouds[CloudsNormal], assets.LoadImage(text.New("data/zones/sky-clouds", i, ".png")))
	}
	for i := range ZoneCount {
		Zones[i] = NewZone(ZoneKind(i))
		ZoneBackgrounds[i] = assets.LoadImage("data/zones/background-" + text.ToLowerCase(Zones[i].Name) + ".png")
	}
	CurrentZone = Zones[ZoneField]

	Bases[TeamAlly] = NewBase(TeamAlly, BaseFortress, Garrison3,
		[3]EntranceKind{EntranceDoor, EntranceShortGate, EntranceTallGate})
	Bases[TeamEnemy] = NewBase(TeamEnemy, BaseBarrack, Garrison3,
		[3]EntranceKind{EntranceDoor, EntranceNone, EntranceNone})

	// Units = append(Units, NewUnit(CharWoman, TeamAlly, LaneMiddle))
	// Units = append(Units, NewUnit(CharMan, TeamEnemy, LaneUpper))
	// Units = append(Units, NewUnit(CharMan, TeamEnemy, LaneMiddle))
	// Units = append(Units, NewUnit(CharHunter, TeamEnemy, LaneLower))
	// Units = append(Units, NewUnit(CharHunter, TeamEnemy, LaneLower))

	Units = append(Units, NewUnit(CharDummy, TeamEnemy, LaneMiddle))
	// Units = append(Units, NewUnit(CharHunter, TeamEnemy, LaneMiddle))

	Pickups = append(Pickups, NewPickup(-240, PickupRelic, LaneLowerOff))
	Pickups = append(Pickups, NewPickup(0, PickupGem, LaneMiddleOff))
	Pickups = append(Pickups, NewPickup(0, PickupCoin, LaneUpperOff))
	Pickups = append(Pickups, NewPickup(100, PickupSnowflake, LaneLowerOff))
	Pickups = append(Pickups, NewPickup(100, PickupCrystal, LaneMiddleOff))
	Pickups = append(Pickups, NewPickup(100, PickupRune, LaneUpperOff))
	Pickups = append(Pickups, NewPickup(150, PickupKey, LaneLowerOff))
	Pickups = append(Pickups, NewPickup(150, PickupStar, LaneMiddleOff))

	Player = NewPlayer()

	Player.Units[0] = NewUnit(CharHunter, TeamAlly, 0)
	Player.Units[1] = NewUnit(CharHunter, TeamAlly, 0)
	Player.Units[2] = NewUnit(CharWoman, TeamAlly, 0)
	Player.Units[3] = NewUnit(CharMan, TeamAlly, 0)
}

//=================================================================

func UpdateTitleScreen() {
	alignView()

	var bgrCrop = TitleScreen.CropArea()
	View.DrawColor(Zones[ZoneDesert].SkyColor)
	View.DrawImage(0, 0, bgrCrop.Width, bgrCrop.Height, 0, TitleScreen, palette.White, geometry.Area{})

	var logo = UserInterface.Crops("logo")[0]
	var logoCrop = logo.CropArea()
	var x, y = PointAtCell(8.5, 1.5)
	View.DrawImage(x, y, logoCrop.Width, logoCrop.Height, 0, logo, palette.White, geometry.Area{})

	gui.Scale = View.Zoom
	var hud = gui.AreaHUD(0.5, 1, 0, 0)
	gui.Button("@ Story Mode", geometry.NewArea(hud.X, hud.Y-TileSize*5.5, 120, 28), geometry.Area{}, ThemeUI, true)
	if gui.IsJustClicked() {
		InGame = true
		PlayAmbience(CurrentZone.Kind)
	}
	gui.Button("* Arena Mode", geometry.NewArea(hud.X, hud.Y-TileSize*4.5, 120, 28), geometry.Area{}, ThemeUI, false)
	if gui.IsFocused() {
		mouse.SetCursor(cursor.NotAllowed)
	}
	gui.Button("Settings", geometry.NewArea(hud.X, hud.Y-TileSize*3, 100, 28), geometry.Area{}, ThemeUI, true)
	gui.Button("Exit", geometry.NewArea(hud.X, hud.Y-TileSize*2, 100, 28), geometry.Area{}, ThemeUI, true)
	if gui.IsJustClicked() {
		window.Close()
	}
}
func UpdateScene() {
	InGameTimer += DeltaTimeScaled()
	alignView()

	mouse.SetCursor(cursor.Default)

	if keyboard.IsKeyJustPressed(key.RightArrow) && CurrentZone.Kind < ZoneHell {
		CurrentZone = Zones[CurrentZone.Kind+1]
	}
	if keyboard.IsKeyJustPressed(key.LeftArrow) && CurrentZone.Kind > ZoneField {
		CurrentZone = Zones[CurrentZone.Kind-1]
	}

	CurrentZone.UpdateBack()
	Bases[TeamAlly].UpdateBack()
	Bases[TeamEnemy].UpdateBack()
	iterateRemovable(&Bases[TeamAlly].Entrances, func(e *Entrance) { e.Update() })
	iterateRemovable(&Bases[TeamEnemy].Entrances, func(e *Entrance) { e.Update() })
	CurrentZone.UpdateFront()

	iterateRemovable(&ProjectilesBehind, func(p *Projectile) { p.Update() })
	iterateRemovable(&Pickups, func(p *Pickup) { p.Update() })

	if mouse.IsAnyButtonJustPressed() {
		PinnedUnit = nil
	}

	iterateRemovable(&Player.Units, func(u *Unit) {
		if !u.IsSummoned() {
			u.Update()
		}
	})
	collection.SortByField(Units, func(u *Unit) float32 {
		if u.Health <= 0 { // dead units go behind all alive units
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

	GameHUD.UpdateBack()
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
		u.HealthBar.Update(hb, u.Health, u.Stats.MaxHealth, u.Mask)
	}

	GameHUD.UpdateFront()
	Player.Update()
}
func DrawShadow(x, z, width, height, angle float32, mask geometry.Area) {
	var lower, upper = Collisions[LaneLower][0], Collisions[LaneUpper][0]
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
func bringGarrisonLanesDown(team Team) {
	for i := LaneGarrison1; i < LaneGarrisonPlus3+1; i++ {
		var length = len(Collisions[i])
		var j = length / 2
		if team == TeamAlly {
			j, length = 0, length/2
		}
		for ; j < length; j++ {
			var shape = Collisions[i][j]
			shape.Y += TileSize
			Collisions[i][j] = shape
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
func alignView() {
	var _, bly = CurrentZone.Ground.PointFromEdge(0.5, 1)
	View.FitSize(CurrentZone.Ground.Width, 0)
	var _, h = View.Size()
	View.Y = (bly - h/2) - 2
}

func statEquation(current, base int, dynamic text.Dynamic) string {
	if current-base > 0 {
		return dynamic.Set("(", base, "+", current-base, ")")
	}
	return dynamic.Set("(", base, current-base, ")")
}
