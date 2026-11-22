package game

import (
	"game/code/constants/base"
	"game/code/constants/clouds"
	"game/code/constants/entrance"
	"game/code/constants/garrison"
	"game/code/constants/lane"
	"game/code/constants/state"
	"game/code/constants/team"
	"game/code/constants/zone"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/cursor"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/text"
	"pure-game-kit/packages/utility/time"
)

const TileSize, MapCount = 32.0, 4

var Decor assets.AtlasId
var SplashAtlas assets.AtlasId

var TimeScale float32 = 1
var SceneHUD *HUD
var CurrentZone *Zone
var Layers []assets.TileLayerId

var Player *PlayerState

var InGameTimer float32

func InitScene() {
	SceneHUD = NewHUD()
	InitCharacters()

	var layers, decor = assets.LoadTileLayersFromTiled("data/decor.tmx")
	Layers = layers
	Decor = assets.LoadAtlas(decor, "data/decor.xml")
	SplashAtlas = assets.LoadAtlas(assets.LoadImage("data/splash.png"), "data/splash.xml")

	for i := range lane.COUNT {
		var tilemap = graphics.NewTilemap(layers[lane.Enum(lane.LayerOffset)+i])
		Collisions[i] = tilemap.TilemapShapes()
	}
	mirrorGarrisonLanes()

	Clouds[clouds.NONE] = []assets.ImageId{0}
	Clouds[clouds.Windy] = append(Clouds[clouds.Windy], assets.LoadImage("data/zones/sky-clouds-wind.png"))
	for i := 1; i <= 4; i++ {
		Clouds[clouds.Normal] = append(Clouds[clouds.Normal], assets.LoadImage(text.New("data/zones/sky-clouds", i, ".png")))
	}
	for i := range zone.COUNT {
		Zones[i] = NewZone(zone.Enum(i))
		ZoneBackgrounds[i] = assets.LoadImage("data/zones/background-" + text.ToLowerCase(Zones[i].Name) + ".png")
	}
	CurrentZone = Zones[zone.Field]

	Bases[team.Ally] = NewBase(team.Ally, base.NONE, garrison.NONE, entrance.NONE, entrance.NONE, entrance.NONE)
	Bases[team.Enemy] = NewBase(team.Enemy, base.Barrack, garrison.NONE, entrance.NONE, entrance.ShortGate, entrance.Door)

	// Units = append(Units, NewUnit(char.Troll, team.Ally, lane.Middle))
	// Units = append(Units, NewUnit(char.Miner, team.Enemy, lane.Middle))

	// Pickups = append(Pickups, NewPickup(-240, pickup.Relic, lane.LowerOff))
	// Pickups = append(Pickups, NewPickup(0, pickup.Gem, lane.MiddleOff))
	// Pickups = append(Pickups, NewPickup(0, pickup.Coin, lane.UpperOff))
	// Pickups = append(Pickups, NewPickup(100, pickup.Snowflake, lane.LowerOff))
	// Pickups = append(Pickups, NewPickup(100, pickup.Crystal, lane.MiddleOff))
	// Pickups = append(Pickups, NewPickup(100, pickup.Rune, lane.UpperOff))
	// Pickups = append(Pickups, NewPickup(150, pickup.Key, lane.LowerOff))
	// Pickups = append(Pickups, NewPickup(150, pickup.Star, lane.MiddleOff))

	// Player.Units[0] = NewUnit(char.Adventurer, team.Ally, 0)
	// Player.Units[1] = NewUnit(char.Miner, team.Ally, 0)
	// Player.Units[2] = NewUnit(char.Cook, team.Ally, 0)
	// Player.Units[3] = NewUnit(char.Bowyer, team.Ally, 0)
	// Player.Units[4] = NewUnit(char.Smith, team.Ally, 0)
	// Player.Units[5] = NewUnit(char.Kid, team.Ally, 0)
	// Player.Units[6] = NewUnit(char.Fisherman, team.Ally, 0)
	// Player.Units[7] = NewUnit(char.Horse, team.Ally, 0)
}

//=================================================================

func UpdateScene() {
	InGameTimer += DeltaTimeScaled()
	AlignView()

	mouse.SetCursor(cursor.Arrow)

	CurrentZone.UpdateBack()
	Bases[team.Ally].UpdateBack()
	Bases[team.Enemy].UpdateBack()
	iterateRemovable(&Bases[team.Ally].Entrances, func(e *Entrance) { e.Update() })
	iterateRemovable(&Bases[team.Enemy].Entrances, func(e *Entrance) { e.Update() })
	CurrentZone.UpdateFront()

	iterateRemovable(&ShotsBehind, func(p *Shot) { p.Update() })
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
		if u.Health <= 0 && u.Values.MaxHealth > 0 { // dead units go behind all alive units
			return number.NegativeInfinity()
		}
		return -u.Z
	})
	iterateRemovable(&Units, func(u *Unit) {
		if u.State == state.WaitingToBeSummoned {
			u.State = state.Summoned
		}
		u.Update()
	})

	Bases[team.Ally].UpdateFront()
	Bases[team.Enemy].UpdateFront()

	const scale = 2.65
	SceneHUD.View.FitSize(1920/scale, 1080/scale) // always fit the HUD based on 16:9, regardless of resolution

	if Story > 100 {
		SceneHUD.UpdateBack()
	}

	iterateRemovable(&Splashes, func(s *Splash) { s.Update() })

	if Story > 100 {
		iterateRemovable(&Shots, func(p *Shot) { p.Update() })

		for _, g := range Bases[team.Ally].Entrances {
			g.HealthBar.Update(g.Tiles[0].Shape, g.Health, g.MaxHealth, geometry.Area{})
		}
		for _, g := range Bases[team.Enemy].Entrances {
			g.HealthBar.Update(g.Tiles[0].Shape, g.Health, g.MaxHealth, geometry.Area{})
		}
		for _, u := range Units { // health bars take the Z order of the units
			if u.IsSummoned() { // this check prevents healthbar flash the very first summon frame
				var hb = u.Hitbox()
				hb.Height += 8
				u.HealthBar.Update(hb, u.Health, u.Values.MaxHealth, u.Mask)
				u.drawDisables()
			}
		}
		SceneHUD.UpdateFront()
		Player.Update()
	}

	UpdateStory()
}

func DrawShadow(x, z, width, height, angle float32, mask geometry.Area) {
	var lower, upper = Collisions[lane.Lower][0], Collisions[lane.Upper][0]
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
	for i := lane.Garrison1; i < lane.GarrisonPlus3+1; i++ {
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
func bringGarrisonLanesDown(forTeam team.Enum) {
	for i := lane.Garrison1; i < lane.GarrisonPlus3+1; i++ {
		var length = len(Collisions[i])
		var j = length / 2
		if forTeam == team.Ally {
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
