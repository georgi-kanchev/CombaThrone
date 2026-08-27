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

	HurtTime float32

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

	Behavior func(self *Unit)
}

const CharMan, CharWoman, CharHunter CharacterKind = 0, 1, 2

var Characters [3]*Character

func NewCharacter(behavior func(self *Unit), stats Stats, info, actValueName string, origin ZoneKind) *Character {
	return &Character{
		Behavior: behavior, Stats: stats, Hitbox: geometry.NewRoundedRectangle(0, 7, 18, 35, 0, 1),
		Sounds: CharSounds{HitFlesh: AudioHitFlesh, HitWood: AudioHitWood, HitMetal: AudioHitMetal},
		Info:   info, ActValueName: actValueName,
	}
}

func InitCharacters() {
	var atlas = assets.LoadAtlas(assets.LoadImage("data/units.png"), "data/units.xml")

	Characters[CharMan] = NewCharacter(BehaviorMan, Stats{Name: "Man", Wage: 20, Role: RoleMelee,
		MaxHealth: 20, Speed: 30, HurtTime: 0.5, ActValue: 2, ActTime: 15, ActRange: 1, RespawnTimer: 100},
		"\nPunches.", "damage", ZoneField)

	Characters[CharWoman] = NewCharacter(BehaviorMan, Stats{Name: "Woman", Wage: 10, Role: RoleHealer,
		MaxHealth: 12, Speed: 20, HurtTime: 0.5, ActValue: 1, ActTime: 18, ActRange: 1, RespawnTimer: 100},
		"\nPunches.", "heal", ZoneField)

	Characters[CharHunter] = NewCharacter(BehaviorHunter, Stats{Name: "Hunter", Wage: 40, Role: RoleRanged,
		MaxHealth: 20, Speed: 15, HurtTime: 0.5, ActValue: 4, ActTime: 20, ActRange: 6, RespawnTimer: 100},
		"\nShoots arrows.", "damage", ZoneField)
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

// private ========================================================

var roleNames = [RoleCount]string{"Melee", "Ranged", "Tank", "Mage", "Healer", "Collector", "Supplier", "Trapper"}
var roleIcons = [RoleCount]Icon{
	IconMelee, IconRanged, IconTank, IconMage, IconHealer, IconCollector, IconSupplier, IconTrapper}
