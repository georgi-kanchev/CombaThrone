package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/random"
)

type CloudsKind uint8
type ZoneKind uint8
type Zone struct {
	Ground, Clouds, Buildings *graphics.Object

	WindSpeed float32

	kind     ZoneKind
	skyColor uint
}

const CloudsNone, CloudsNormal, CloudsWindy, CloudsCount CloudsKind = 0, 1, 2, 3

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

var Clouds [CloudsCount][]assets.ImageId
var ZoneBackgrounds [ZoneCount]assets.ImageId
var Zones [ZoneCount]*Zone

func NewZone(kind ZoneKind) *Zone {
	var ground = graphics.NewTilemap(Layers[ZoneLayerOffset+int(kind)*2])
	var buildings = graphics.NewTilemap(Layers[ZoneLayerOffset+int(kind)*2+1])
	var cloudsKind = zoneClouds[kind]
	var randomClouds = Clouds[cloudsKind]
	var cloud = random.Pick(randomClouds...)
	var clouds = graphics.NewSprite(0, 0, 1, cloud)
	var windSpeed float32
	switch cloudsKind {
	case CloudsNone:
		windSpeed = random.Range[float32](0.2, 0.7)
	case CloudsNormal:
		collection.Remove(randomClouds, cloud) // field and ruins shouldn't have the same clouds
		windSpeed = random.Range[float32](0.5, 2.0)
	case CloudsWindy:
		windSpeed = random.Range[float32](2.0, 5.0)
	}
	clouds.ImageCrop = cloud.CropArea()
	return &Zone{Ground: &ground, Buildings: &buildings, skyColor: zoneSkyColors[kind], kind: kind, Clouds: &clouds,
		WindSpeed: windSpeed}
}

//=================================================================

var zoneNames = [ZoneCount]string{
	ZoneField: "field", ZoneRuins: "ruins", ZoneSwamp: "swamp", ZoneDesert: "desert", ZoneDocks: "docks",
	ZoneGlacier: "glacier", ZoneCave: "cave", ZoneMine: "mine", ZoneHell: "hell",
}
var zoneSkyColors = [ZoneCount]uint{
	ZoneField: color.TagRGBA("rgb(98, 171, 212)"), ZoneRuins: color.TagRGBA("rgb(98, 171, 212)"),
	ZoneSwamp: color.TagRGBA("rgb(37, 65, 61)"), ZoneDesert: color.TagRGBA("rgb(155, 240, 253)"),
	ZoneDocks: color.TagRGBA("rgb(98, 171, 212)"), ZoneGlacier: color.TagRGBA("rgb(155, 240, 253)"),
	ZoneCave: color.TagRGBA("rgb(72, 54, 59)"), ZoneMine: color.TagRGBA("rgb(61, 36, 59)"),
	ZoneHell: color.TagRGBA("rgb(227, 177, 109)"),
}
var zoneInfos = [ZoneCount]string{
	ZoneField: "1. The Field of the Vanilla-gers", ZoneRuins: "2. The Ruins of the Robbing Hoods",
	ZoneSwamp: "3. The Swamp of the Abomi Nation", ZoneDesert: "4. The Desert of the Sarcopha-guys",
	ZoneDocks: "5. The Docks of the Plank-ton Pirates", ZoneGlacier: "6. The Glacier of the Satan Claws & Co.",
	ZoneCave: "7. The Cave of the Troglo-bites", ZoneMine: "8. The Mine of the Avant Guards",
	ZoneHell: "9. The Hell of the Demons-trosities",
}
var zoneClouds = [ZoneCount]CloudsKind{
	ZoneField: CloudsNormal, ZoneRuins: CloudsNormal, ZoneSwamp: CloudsNone, ZoneDesert: CloudsNone,
	ZoneDocks: CloudsWindy, ZoneGlacier: CloudsWindy, ZoneCave: CloudsNone, ZoneMine: CloudsNone,
	ZoneHell: CloudsNone,
}

func (z *Zone) UpdateBack() {
	View.DrawColor(z.skyColor)

	z.Clouds.ImageCrop.X -= z.WindSpeed * DeltaTimeScaled()
	View.DrawObject(z.Clouds)

	View.DrawImage(0, 0, z.Ground.Width, z.Ground.Height, 0, ZoneBackgrounds[z.kind], palette.White, geometry.Area{})

	var buildingWind = z.WindSpeed
	if z.kind == ZoneHell {
		buildingWind = 0
	}
	z.Buildings.Effects.TileTimeScale = TimeScale + buildingWind
	View.DrawObject(z.Buildings)
}
func (z *Zone) UpdateFront() {
	if Bases[TeamAlly].Kind < BaseBarrack {
		z.Ground.X = -TileSize * 4
		View.DrawObject(z.Ground)
	}
	if Bases[TeamEnemy].Kind < BaseBarrack {
		z.Ground.X = TileSize * 4
		View.DrawObject(z.Ground)
	}
	z.Ground.X = 0
	View.DrawObject(z.Ground)
}
