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
type CharSounds struct{ ActionStart, ActionTrigger, HitFlesh, HitWood, HitMetal, HitGround []audio.Audio }
type Character struct {
	Stats      Stats
	Hitbox     geometry.Shape
	Animations struct {
		Idle, Walk, ActionStart, ActionEnd, Hurt, Die []assets.ImageId
	}
	Sounds             CharSounds
	Icon               assets.ImageId
	Origin             ZoneKind
	Info, ActValueName string
	RoleIcon           Icon
	RoleName           string

	Behavior func(self *Unit)
}

const CharMan, CharWoman, CharHunter, CharDummy, CharCount CharacterKind = 0, 1, 2, 3, 4

var Characters [4]*Character

func NewCharacter(behavior func(self *Unit), origin ZoneKind, stats Stats, info string) *Character {
	var roleIcons = [RoleCount]Icon{
		IconMelee, IconRanged, IconTank, IconMage, IconHealer, IconCollector, IconSupplier, IconTrapper}
	var roleNames = [RoleCount]string{"Melee", "Ranged", "Tank", "Mage", "Healer", "Collector", "Supplier", "Trapper"}
	var actNames = [RoleCount]string{"damage", "damage", "block", "damage", "heal", "carry", "buff", "debuff"}

	return &Character{
		Behavior: behavior, Stats: stats, Hitbox: geometry.NewRoundedRectangle(0, 7, 18, 35, 0, 1),
		Sounds: CharSounds{HitFlesh: AudioHitFlesh, HitWood: AudioHitWood, HitMetal: AudioHitMetal},
		Info:   info, ActValueName: actNames[stats.Role], RoleIcon: roleIcons[stats.Role], RoleName: roleNames[stats.Role],
	}
}

func InitCharacters() {
	var atlas = assets.LoadAtlas(assets.LoadImage("data/units.png"), "data/units.xml")

	Characters[CharMan] = NewCharacter(BehaviorMan, ZoneField, Stats{Name: "Man", Wage: 20, Role: RoleMelee,
		MaxHealth: 20, Speed: 30, ActValue: 2, ActTime: 15, ActRange: 1, RespawnTimer: 100},
		"When close to a Woman:\n🌗🟪"+Tags[IconTimer]+"loses 0.5s rest⬜")

	Characters[CharWoman] = NewCharacter(BehaviorWoman, ZoneField, Stats{Name: "Woman", Wage: 10, Role: RoleHealer,
		MaxHealth: 1, Speed: 20, ActValue: 1, ActTime: 18, ActRange: 1, RespawnTimer: 100},
		"When in front of a Man:\n🌗🟨"+Tags[IconMove]+"gains 10 speed⬜")

	Characters[CharDummy] = NewCharacter(BehaviorDummy, ZoneField, Stats{Name: "Dummy", Role: RoleTank,
		MaxHealth: 1, Speed: 0, ActValue: 0, ActTime: 0, ActRange: 0, RespawnTimer: 0},
		"Cannot die.\nAlthough, it would like to.")

	Characters[CharHunter] = NewCharacter(BehaviorHunter, ZoneField, Stats{Name: "Hunter", Wage: 40, Role: RoleRanged,
		MaxHealth: 14, Speed: 15, ActValue: 4, ActTime: 20, ActRange: 6, RespawnTimer: 100},
		"When not garrison: \n🟧"+Tags[IconRange]+"gains 2 range⬜")
	Characters[CharHunter].Sounds = CharSounds{ActionTrigger: AudioBow, HitGround: AudioProjectileGround,
		HitFlesh: AudioProjectileFlesh, HitWood: AudioProjectileWood, HitMetal: AudioProjectileMetal}

	for i, c := range Characters {
		var prefix = text.ToLowerCase(c.Stats.Name)
		c.Animations.Idle = atlas.Crops(prefix + "-idle")
		c.Animations.Walk = atlas.Crops(prefix + "-walk")
		c.Animations.ActionStart = atlas.Crops(prefix + "-action-start")
		c.Animations.ActionEnd = atlas.Crops(prefix + "-action-end")
		c.Animations.Hurt = atlas.Crops(prefix + "-hurt")
		c.Animations.Die = atlas.Crops(prefix + "-death")
		c.Icon = atlas.Crops(prefix + "-icon")[0]
		Characters[i] = c
	}
}
