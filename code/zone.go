package game

import (
	"game/code/constants/base"
	"game/code/constants/clouds"
	"game/code/constants/team"
	"game/code/constants/zone"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/random"
	"pure-game-kit/packages/utility/text"
)

type Zone struct {
	Name, Info string
	Kind       zone.Enum

	Ground, Clouds, Buildings *graphics.Object

	WindSpeed float32
	SkyColor  uint

	FlagFrames []assets.ImageId
}

var Clouds [clouds.Count][]assets.ImageId
var ZoneBackgrounds [zone.Count]assets.ImageId
var Zones [zone.Count]*Zone

func NewZone(kind zone.Enum) *Zone {
	var names = [zone.Count]string{
		"Field", "Ruins", "Forest", "Swamp", "Desert", "Docks", "Ocean", "Glacier", "Cave", "Mine", "Hell"}
	var skyColors = [zone.Count]uint{
		zone.Field: color.TagRGBA("rgb(98, 171, 212)"), zone.Ruins: color.TagRGBA("rgb(98, 171, 212)"),
		zone.Forest: color.TagRGBA("rgb(31, 55, 54)"), zone.Swamp: color.TagRGBA("rgb(92, 107, 83)"),
		zone.Desert: color.TagRGBA("rgb(155, 240, 253)"), zone.Docks: color.TagRGBA("rgb(98, 171, 212)"),
		zone.Ocean: color.TagRGBA("rgb(98, 171, 212)"), zone.Glacier: color.TagRGBA("rgb(155, 240, 253)"),
		zone.Cave: color.TagRGBA("rgb(72, 54, 59)"), zone.Mine: color.TagRGBA("rgb(61, 36, 59)"),
		zone.Hell: color.TagRGBA("rgb(227, 177, 109)"),
	}
	var infos = [zone.Count]string{"1. The Field of the Vanilla-gers", "2. The Ruins of the Robbing Hoods",
		"3. The Forest of the Cast-aways & Potion-eers", "4. The Swamp of the Abomi Nation",
		"5. The Desert of the Sarcopha-guys", "6. The Docks of the Plank-ton Squid-iots",
		"7. The Cliff of the Balloon-a-ticks", "8. The Glacier of the Satan Claws & Co.",
		"9. The Cave of the Troglo-bites", "10. The Mine of the ???????-????", /*Masquer-oids*/
		"11. The Hell of the Demons-trosities"}
	var cloudPerZone = [zone.Count]clouds.Enum{zone.Field: clouds.Normal, zone.Ruins: clouds.Normal, zone.Forest: clouds.None,
		zone.Swamp: clouds.None, zone.Desert: clouds.None, zone.Docks: clouds.Windy, zone.Ocean: clouds.Windy,
		zone.Glacier: clouds.Windy, zone.Cave: clouds.None, zone.Mine: clouds.None, zone.Hell: clouds.None}
	var ground = graphics.NewTilemap(Layers[zone.LayerOffset+int(kind)*2])
	var buildings = graphics.NewTilemap(Layers[zone.LayerOffset+int(kind)*2+1])
	var randomClouds = Clouds[cloudPerZone[kind]]
	var cloud = random.Pick(randomClouds...)
	var cloudsSprite = graphics.NewSprite(0, 0, 1, cloud)
	var windSpeed float32
	var flagGroup = "flag-" + text.ToLowerCase(names[kind])
	switch cloudPerZone[kind] {
	case clouds.None:
		windSpeed = random.Range[float32](0.2, 0.7)
	case clouds.Normal:
		collection.Remove(randomClouds, cloud) // field and ruins shouldn't have the same clouds
		windSpeed = random.Range[float32](0.5, 2.0)
	}
	if cloudPerZone[kind] == clouds.Windy || kind == zone.Desert {
		windSpeed = random.Range[float32](2.0, 5.0)
	} else if kind == zone.Hell {
		windSpeed = 1
	}
	cloudsSprite.ImageCrop = cloud.CropArea()

	return &Zone{Ground: &ground, Buildings: &buildings, SkyColor: skyColors[kind], Kind: kind, Clouds: &cloudsSprite,
		WindSpeed: windSpeed, FlagFrames: Decor.Crops(flagGroup), Name: names[kind], Info: infos[kind]}
}

//=================================================================

func (z *Zone) UpdateBack() {
	View.DrawColor(z.SkyColor)

	z.Clouds.ImageCrop.X -= z.WindSpeed * DeltaTimeScaled()
	View.DrawObject(z.Clouds)

	View.DrawImage(0, 0, z.Ground.Width, z.Ground.Height, 0, ZoneBackgrounds[z.Kind], palette.White, geometry.Area{})

	var buildingWind = z.WindSpeed
	z.Buildings.Details.TileTimeScale = TimeScale * buildingWind
	View.DrawObject(z.Buildings)
}
func (z *Zone) UpdateFront() {
	if Bases[team.Ally].Kind < base.Barrack {
		z.Ground.X = -TileSize * 4
		View.DrawObject(z.Ground)
	}
	if Bases[team.Enemy].Kind < base.Barrack {
		z.Ground.X = TileSize * 4
		View.DrawObject(z.Ground)
	}
	z.Ground.X = 0
	View.DrawObject(z.Ground)
}
