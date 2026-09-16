package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/random"
	"pure-game-kit/packages/utility/text"
)

type CloudsKind uint8
type ZoneKind uint8
type Zone struct {
	Name, Info string
	Kind       ZoneKind

	Ground, Clouds, Buildings *graphics.Object

	WindSpeed float32
	SkyColor  uint

	FlagFrames []assets.ImageId
}

const CloudsNone, CloudsNormal, CloudsWindy, CloudsCount CloudsKind = 0, 1, 2, 3

const (
	ZoneField ZoneKind = iota
	ZoneRuins
	ZoneForest
	ZoneSwamp
	ZoneDesert
	ZoneDocks
	ZoneOcean
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
	var names = [ZoneCount]string{
		"Field", "Ruins", "Forest", "Swamp", "Desert", "Docks", "Ocean", "Glacier", "Cave", "Mine", "Hell"}
	var skyColors = [ZoneCount]uint{
		ZoneField: color.TagRGBA("rgb(98, 171, 212)"), ZoneRuins: color.TagRGBA("rgb(98, 171, 212)"),
		ZoneForest: color.TagRGBA("rgb(31, 55, 54)"), ZoneSwamp: color.TagRGBA("rgb(92, 107, 83)"),
		ZoneDesert: color.TagRGBA("rgb(155, 240, 253)"), ZoneDocks: color.TagRGBA("rgb(98, 171, 212)"),
		ZoneOcean: color.TagRGBA("rgb(98, 171, 212)"), ZoneGlacier: color.TagRGBA("rgb(155, 240, 253)"),
		ZoneCave: color.TagRGBA("rgb(72, 54, 59)"), ZoneMine: color.TagRGBA("rgb(61, 36, 59)"),
		ZoneHell: color.TagRGBA("rgb(227, 177, 109)"),
	}
	var infos = [ZoneCount]string{"1. The Field of the Vanilla-gers", "2. The Ruins of the Robbing Hoods",
		"3. The Forest of the Cast-aways & Potion-eers", "4. The Swamp of the Abomi Nation",
		"5. The Desert of the Sarcopha-guys", "6. The Docks of the Plank-ton Squid-iots",
		"7. The Cliff of the Balloon-a-ticks", "8. The Glacier of the Satan Claws & Co.",
		"9. The Cave of the Troglo-bites", "10. The Mine of the ???????-????", /*Masquer-oids*/
		"11. The Hell of the Demons-trosities"}
	var cloudPerZone = [ZoneCount]CloudsKind{ZoneField: CloudsNormal, ZoneRuins: CloudsNormal, ZoneForest: CloudsNone,
		ZoneSwamp: CloudsNone, ZoneDesert: CloudsNone, ZoneDocks: CloudsWindy, ZoneOcean: CloudsWindy,
		ZoneGlacier: CloudsWindy, ZoneCave: CloudsNone, ZoneMine: CloudsNone, ZoneHell: CloudsNone}
	var ground = graphics.NewTilemap(Layers[ZoneLayerOffset+int(kind)*2])
	var buildings = graphics.NewTilemap(Layers[ZoneLayerOffset+int(kind)*2+1])
	var randomClouds = Clouds[cloudPerZone[kind]]
	var cloud = random.Pick(randomClouds...)
	var clouds = graphics.NewSprite(0, 0, 1, cloud)
	var windSpeed float32
	var flagGroup = "flag-" + text.ToLowerCase(names[kind])
	switch cloudPerZone[kind] {
	case CloudsNone:
		windSpeed = random.Range[float32](0.2, 0.7)
	case CloudsNormal:
		collection.Remove(randomClouds, cloud) // field and ruins shouldn't have the same clouds
		windSpeed = random.Range[float32](0.5, 2.0)
	}
	if cloudPerZone[kind] == CloudsWindy || kind == ZoneDesert {
		windSpeed = random.Range[float32](2.0, 5.0)
	} else if kind == ZoneHell {
		windSpeed = 1
	}
	clouds.ImageCrop = cloud.CropArea()

	return &Zone{Ground: &ground, Buildings: &buildings, SkyColor: skyColors[kind], Kind: kind, Clouds: &clouds,
		WindSpeed: windSpeed, FlagFrames: Decor.Crops(flagGroup), Name: names[kind], Info: infos[kind]}
}

//=================================================================

func (z *Zone) UpdateBack() {
	View.DrawColor(z.SkyColor)

	z.Clouds.ImageCrop.X -= z.WindSpeed * DeltaTimeScaled()
	View.DrawObject(z.Clouds)

	View.DrawImage(0, 0, z.Ground.Width, z.Ground.Height, 0, ZoneBackgrounds[z.Kind], palette.White, geometry.Area{})

	var buildingWind = z.WindSpeed
	z.Buildings.Effects.TileTimeScale = TimeScale * buildingWind
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
