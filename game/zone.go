package game

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color"
)

type ZoneKind uint8
type Zone struct {
	Background, Ground, Buildings *graphics.Object

	kind     ZoneKind
	skyColor uint
}

const (
	ZoneField ZoneKind = iota
	ZoneRuins
	ZoneSwamp
	ZoneDesert
	ZoneDocks
	ZoneGlacier
	ZoneCave
	ZoneMine
	ZoneHell
	ZoneCount
	ZoneLayerOffset = 10
)

var Zones [ZoneCount]*Zone

func NewZone(kind ZoneKind) *Zone {
	var bgr = graphics.NewSprite(0, 0, 1, Backgrounds[kind])
	var ground = graphics.NewTilemap(Layers[ZoneLayerOffset+int(kind)*2])
	var buildings = graphics.NewTilemap(Layers[ZoneLayerOffset+int(kind)*2+1])
	return &Zone{Background: &bgr, Ground: &ground, Buildings: &buildings, skyColor: zoneSkyColors[kind], kind: kind}
}

//=================================================================

var zoneNames = [ZoneCount]string{
	ZoneField: "field", ZoneRuins: "ruins", ZoneSwamp: "swamp", ZoneDesert: "desert", ZoneDocks: "docks",
	ZoneGlacier: "glacier", ZoneCave: "cave", ZoneMine: "mine", ZoneHell: "hell",
}
var zoneSkyColors = [ZoneCount]uint{
	ZoneField: color.TagRGBA("rgb(98, 171, 212)"), ZoneDesert: color.TagRGBA("rgb(155, 240, 253)"),
	ZoneRuins: color.TagRGBA("rgb(98, 171, 212)"), ZoneCave: color.TagRGBA("rgb(72, 54, 59)"),
	ZoneMine: color.TagRGBA("rgb(61, 36, 59)"), ZoneGlacier: color.TagRGBA("rgb(155, 240, 253)"),
	ZoneDocks: color.TagRGBA("rgb(98, 171, 212)"), ZoneSwamp: color.TagRGBA("rgb(37, 65, 61)"),
	ZoneHell: color.TagRGBA("rgb(227, 177, 109)"),
}
var zoneInfos = [ZoneCount]string{
	ZoneField: "1. The Field of the Vanilla-gers", ZoneRuins: "2. The Ruins of the Robbing Hoods",
	ZoneSwamp: "3. The Swamp of the Abomi Nation", ZoneDesert: "4. The Desert of the Sarcopha-guys",
	ZoneDocks: "5. The Docks of the Plank-ton Pirates", ZoneGlacier: "6. The Glacier of the Satan Claws & Co.",
	ZoneCave: "7. The Cave of the Troglo-bites", ZoneMine: "8. The Mine of the Avant Guards",
	ZoneHell: "9. The Hell of the Demons-trosities",
}

func (z *Zone) UpdateBack() {
	View.DrawColor(z.skyColor)
	View.DrawObject(z.Background)
	z.Buildings.Effects.TileTimeScale = TimeScale
	View.DrawObject(z.Buildings)
}
func (z *Zone) UpdateFront() {
	if AllyBase.Kind < BaseBarrack {
		z.Ground.X = -TileSize * 4
		View.DrawObject(z.Ground)
	}
	if EnemyBase.Kind < BaseBarrack {
		z.Ground.X = TileSize * 4
		View.DrawObject(z.Ground)
	}
	z.Ground.X = 0
	View.DrawObject(z.Ground)
}
