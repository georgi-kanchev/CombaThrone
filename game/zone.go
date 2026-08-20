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
	ZoneDesert
	ZoneRuins
	ZoneCave
	ZoneMine
	ZoneGlacier
	ZoneDocks
	ZoneSwamp
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

var zoneSkyColors = [ZoneCount]uint{
	ZoneField: color.TagRGBA("rgb(98, 171, 212)"), ZoneDesert: color.TagRGBA("rgb(155, 240, 253)"),
	ZoneRuins: color.TagRGBA("rgb(98, 171, 212)"), ZoneCave: color.TagRGBA("rgb(72, 54, 59)"),
	ZoneMine: color.TagRGBA("rgb(61, 36, 59)"), ZoneGlacier: color.TagRGBA("rgb(155, 240, 253)"),
	ZoneDocks: color.TagRGBA("rgb(98, 171, 212)"), ZoneSwamp: color.TagRGBA("rgb(37, 65, 61)"),
	ZoneHell: color.TagRGBA("rgb(227, 177, 109)"),
}

func (e *Zone) UpdateBack() {
	View.DrawColor(e.skyColor)
	View.DrawObject(e.Background)
	View.DrawObject(e.Buildings)
}
func (e *Zone) UpdateFront() {
	if AllyBase.Kind < BaseBarrack {
		e.Ground.X = -TileSize * 4
		View.DrawObject(e.Ground)
	}
	if EnemyBase.Kind < BaseBarrack {
		e.Ground.X = TileSize * 4
		View.DrawObject(e.Ground)
	}
	e.Ground.X = 0
	View.DrawObject(e.Ground)
}
