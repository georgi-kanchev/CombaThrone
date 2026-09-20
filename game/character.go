package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/audio"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/utility/text"
)

type Stats struct {
	Name string

	MaxHealth, Speed, Wage int

	ActValue, ActTime, ActRange, RespawnTimer int

	Role Role
}

type CharacterKind uint8
type CharSounds struct{ ActStart, ActTrigger, HitFlesh, HitWood, HitMetal, HitGround []audio.Audio }
type Character struct {
	Stats      Stats
	Hitbox     geometry.Shape
	Animations struct {
		Idle, Walk, ActStart, ActEnd, Hurt, Die []assets.ImageId
	}
	Sounds             CharSounds
	Icon               assets.ImageId
	Origin             ZoneKind
	Info, ActValueName string
	RoleIcon           Icon
	RoleName           string

	Behavior func(self *Unit)
}

const CharMiner, CharCook, CharBowyer, CharDummy, CharCount CharacterKind = 0, 1, 2, 3, 4

var Characters [4]*Character

func NewCharacter(behavior func(self *Unit), origin ZoneKind, stats Stats, info string) *Character {
	var roleIcons = [RoleCount]Icon{
		IconMelee, IconRanged, IconTank, IconCollector, IconSupplier, IconGriefer}
	var roleNames = [RoleCount]string{"Melee", "Ranged", "Tank", "Collector", "Support", "Griefer"}
	var actNames = [RoleCount]string{"damage", "damage", "block", "carry", "buff", "grief"}

	return &Character{
		Behavior: behavior, Stats: stats, Hitbox: geometry.NewRoundedRectangle(0, 7, 18, 35, 0, 1),
		Sounds: CharSounds{HitFlesh: AudioHitFlesh, HitWood: AudioHitWood, HitMetal: AudioHitMetal},
		Info:   info, ActValueName: actNames[stats.Role], RoleIcon: roleIcons[stats.Role], RoleName: roleNames[stats.Role],
	}
}

func InitCharacters() {
	var atlas = assets.LoadAtlas(assets.LoadImage("data/units.png"), "data/units.xml")

	Characters[CharMiner] = NewCharacter(BehaviorMan, ZoneField, Stats{Name: "Miner", Wage: 20, Role: RoleFighter,
		MaxHealth: 20, Speed: 30, ActValue: 2, ActTime: 15, ActRange: 1, RespawnTimer: 100},
		"When close to a Woman:\n🌗🟪"+Tags[IconTimer]+"loses 0.5s rest⬜")

	Characters[CharCook] = NewCharacter(BehaviorWoman, ZoneField, Stats{Name: "Cook", Wage: 10, Role: RoleSupport,
		MaxHealth: 1, Speed: 20, ActValue: 1, ActTime: 18, ActRange: 1, RespawnTimer: 100},
		"When in front of a Man:\n🌗🟨"+Tags[IconMove]+"gains 10 speed⬜")

	Characters[CharDummy] = NewCharacter(BehaviorDummy, ZoneField, Stats{Name: "Dummy", Role: RoleTank,
		MaxHealth: 1, Speed: 0, ActValue: 0, ActTime: 0, ActRange: 0, RespawnTimer: 0},
		"Cannot die.\nAlthough, it would like to.")

	Characters[CharBowyer] = NewCharacter(BehaviorHunter, ZoneField, Stats{Name: "Bowyer", Wage: 40, Role: RoleRanger,
		MaxHealth: 14, Speed: 15, ActValue: 4, ActTime: 20, ActRange: 6, RespawnTimer: 100},
		"When not garrison: \n🟧"+Tags[IconRange]+"gains 2 range⬜")
	Characters[CharBowyer].Sounds = CharSounds{ActTrigger: AudioBow, HitGround: AudioProjectileGround,
		HitFlesh: AudioProjectileFlesh, HitWood: AudioProjectileWood, HitMetal: AudioProjectileMetal}

	for i, c := range Characters {
		var prefix = text.Replace(text.ToLowerCase(c.Stats.Name), " ", "-")
		c.Animations.Idle = atlas.Crops(prefix + "-idle")
		c.Animations.Walk = atlas.Crops(prefix + "-move")
		c.Animations.ActStart = atlas.Crops(prefix + "-prepare")
		c.Animations.ActEnd = atlas.Crops(prefix + "-recover")
		c.Animations.Hurt = atlas.Crops(prefix + "-hurt")
		c.Animations.Die = atlas.Crops(prefix + "-die")
		c.Icon = atlas.Crops(prefix + "-icon")[0]
		Characters[i] = c
	}
}
