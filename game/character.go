package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/audio"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/utility/text"
)

type Values struct {
	Name, EffectInfo string

	MaxHealth, MoveSpeed, Wage int
	ActPoints, ActRange        int

	ActTimer, RespawnTimer float32
	EffectTimer            float32 // remaining Effect time (used only for Unit.Effects - not in Unit/Character)

	Role Role
}

type CharacterKind uint8
type CharSounds struct{ ActStart, ActTrigger, HitFlesh, HitWood, HitMetal, HitGround []audio.Audio }
type Character struct {
	Values              Values
	Hitbox              geometry.Shape
	Animations          struct{ Idle, Walk, ActStart, ActEnd, Hurt, Die []assets.ImageId }
	Sounds              CharSounds
	Icon                assets.ImageId
	Origin              ZoneKind
	Info, ActPointsName string
	RoleIcon            Icon
	RoleName            string
}

const CharMiner, CharCook, CharBowyer, CharDummy, CharCount CharacterKind = 0, 1, 2, 3, 4

var Characters [4]*Character

func NewCharacter(origin ZoneKind, stats Values, info string) *Character {
	var roleIcons = [RoleCount]Icon{
		IconSword, IconBow, IconShield, IconBag, IconHand, IconDebuff}
	var roleNames = [RoleCount]string{"Fighter", "Ranger", "Defender", "Collector", "Supplier", "Griefer"}
	var actNames = [RoleCount]string{"damage", "damage", "block", "carry", "buff", "grief"}

	return &Character{
		Values: stats, Hitbox: geometry.NewRoundedRectangle(0, 7, 18, 35, 0, 1),
		Sounds: CharSounds{HitFlesh: AudioHitFlesh, HitWood: AudioHitWood, HitMetal: AudioHitMetal},
		Info:   info, ActPointsName: actNames[stats.Role], RoleIcon: roleIcons[stats.Role], RoleName: roleNames[stats.Role],
	}
}

func InitCharacters() {
	var atlas = assets.LoadAtlas(assets.LoadImage("data/units.png"), "data/units.xml")

	Characters[CharDummy] = NewCharacter(ZoneField, Values{Name: "Dummy", Role: RoleDefender,
		MaxHealth: 1, MoveSpeed: 0, ActPoints: 0, ActTimer: 0, ActRange: 0, RespawnTimer: 0},
		"Cannot 🟥"+Tags[IconSkull]+"die⬜. Even though it really wants to.")
	Characters[CharDummy].Hitbox.Width += 4
	Characters[CharMiner] = NewCharacter(ZoneField, Values{Name: "Miner", Wage: 20, Role: RoleFighter,
		MaxHealth: 20, MoveSpeed: 20, ActPoints: 2, ActTimer: 2.0, ActRange: 1, RespawnTimer: 10.0},
		"🟩"+Tags[IconPlus]+"🟧"+Tags[IconSword]+"8 damage ⬜against "+Tags[IconDoor]+"entrances.")
	Characters[CharCook] = NewCharacter(ZoneField, Values{Name: "Cook", Wage: 10, Role: RoleSupplier,
		MaxHealth: 1, MoveSpeed: 15, ActPoints: 4, ActTimer: 10.0, ActRange: 1, RespawnTimer: 10.0},
		"🟩"+Tags[IconPlus]+"🌗🟩"+Tags[IconHealth]+"4 health ⬜to an 🟩ally ⬜passing by.")
	Characters[CharBowyer] = NewCharacter(ZoneField, Values{Name: "Bowyer", Wage: 40, Role: RoleRanger,
		MaxHealth: 14, MoveSpeed: 15, ActPoints: 4, ActTimer: 2.0, ActRange: 6, RespawnTimer: 10.0},
		"🟩"+Tags[IconPlus]+"🌗🟧"+Tags[IconRange]+"2 range ⬜when grounded.")
	Characters[CharBowyer].Sounds = CharSounds{ActTrigger: AudioBow, HitGround: AudioProjectileGround,
		HitFlesh: AudioProjectileFlesh, HitWood: AudioProjectileWood, HitMetal: AudioProjectileMetal}

	for i, c := range Characters {
		var prefix = text.Replace(text.ToLowerCase(c.Values.Name), " ", "-")
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
