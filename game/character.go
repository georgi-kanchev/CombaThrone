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

	ActTimer, RespawnTimer, SleepTimer float32

	EffectTimer float32 // remaining Effect time (used only for Unit.Effects - not in Unit/Character)

	Role Role
}

type CharacterKind uint8
type CharSounds struct{ ActStart, ActTrigger, HitFlesh, HitWood, HitMetal, HitGround []audio.Audio }
type Character struct {
	Values              Values
	Hitbox              geometry.Shape
	Animations          struct{ Idle, Walk, Prepare, Recover, Hurt, Die []assets.ImageId }
	Sounds              CharSounds
	Icon                assets.ImageId
	Origin              ZoneKind
	Info, ActPointsName string
	RoleIcon            Icon
	RoleName            string
}

const (
	CharDummy CharacterKind = iota
	CharMiner
	CharCook
	CharBowyer
	CharSmith
	CharKid
	CharFisherman
	CharHorse
	CharBrownBunny
	CharWhiteBunny
	CharCount
)

var Characters [CharCount]*Character

func NewCharacter(origin ZoneKind, stats Values, info string) *Character {
	var roleIcons = [RoleCount]Icon{IconSword, IconBow, IconShield, IconDebuff, IconHand, IconBag}
	var roleNames = [RoleCount]string{"Fighter", "Ranger", "Defender", "Griefer", "Supplier", "Collector"}
	var actNames = [RoleCount]string{"damage", "damage", "block", "grief", "buff", "pickup"}

	return &Character{
		Values: stats, Hitbox: geometry.NewRoundedRectangle(0, 0, 18, 35, 0, 1),
		Sounds: CharSounds{HitFlesh: AudioHitFlesh, HitWood: AudioHitWood, HitMetal: AudioHitMetal},
		Info:   info, ActPointsName: actNames[stats.Role], RoleIcon: roleIcons[stats.Role], RoleName: roleNames[stats.Role],
	}
}

func InitCharacters() {
	var atlas = assets.LoadAtlas(assets.LoadImage("data/units.png"), "data/units.xml")

	Characters[CharDummy] = NewCharacter(ZoneField, Values{Name: "Dummy", Role: RoleDefender,
		MaxHealth: 0, MoveSpeed: 0, ActPoints: 0, ActTimer: 0, ActRange: 0, RespawnTimer: 0},
		"Cannot 🟥"+Tags[IconSkull]+"die⬜. Even though it really wants to.")
	Characters[CharDummy].Hitbox = geometry.NewRoundedRectangle(0, 6, 24, 34, 0, 1)

	Characters[CharMiner] = NewCharacter(ZoneField, Values{Name: "Miner", Wage: 20, Role: RoleFighter,
		MaxHealth: 20, MoveSpeed: 20, ActPoints: 2, ActTimer: 2.0, ActRange: 1, RespawnTimer: 10.0},
		"🟩"+Tags[IconPlus]+"🟧"+"8 "+Tags[IconSword]+"damage ⬜against "+Tags[IconDoor]+"entrances.")
	Characters[CharMiner].Hitbox.Y += 6

	Characters[CharCook] = NewCharacter(ZoneField, Values{Name: "Cook", Wage: 10, Role: RoleSupplier,
		MaxHealth: 0, MoveSpeed: 15, ActPoints: 4, ActTimer: 10.0, ActRange: 1, RespawnTimer: 10.0},
		"🟧"+Tags[IconHand]+"Heals⬜ a 🌗🟨"+Tags[IconLeftRight]+"passing by 🟩"+Tags[IconUnit]+"ally⬜.")

	Characters[CharBowyer] = NewCharacter(ZoneField, Values{Name: "Bowyer", Wage: 40, Role: RoleRanger,
		MaxHealth: 14, MoveSpeed: 15, ActPoints: 4, ActTimer: 2.0, ActRange: 6, RespawnTimer: 10.0},
		"🟩"+Tags[IconPlus]+"🌗🟧2 "+Tags[IconRange]+"range ⬜when 🌗🟨"+Tags[IconLeftRight]+"grounded⬜.")
	Characters[CharBowyer].Hitbox.Y += 7
	Characters[CharBowyer].Sounds = CharSounds{ActTrigger: AudioBow, HitGround: AudioProjectileGround,
		HitFlesh: AudioProjectileFlesh, HitWood: AudioProjectileWood, HitMetal: AudioProjectileMetal}

	Characters[CharSmith] = NewCharacter(ZoneField, Values{Name: "Smith", Wage: 20, Role: RoleDefender,
		MaxHealth: 40, MoveSpeed: 12, ActPoints: 1, ActTimer: 5.0, ActRange: 1, RespawnTimer: 20.0},
		"Pushes the 🟥"+Tags[IconUnit]+"enemy ⬜in front.")
	Characters[CharSmith].Hitbox = geometry.NewRoundedRectangle(0, 0, 24, 48, 0, 1)

	Characters[CharKid] = NewCharacter(ZoneField, Values{Name: "Kid", Wage: 20, Role: RoleCollector,
		MaxHealth: 0, MoveSpeed: 20, ActPoints: 1, ActTimer: 0.0, ActRange: 0, RespawnTimer: 20.0},
		"🟩"+Tags[IconPlus]+"🌗🟨20 "+Tags[IconLeftRight]+"speed⬜ when 🟧"+Tags[IconBag]+"empty handed⬜.")
	Characters[CharKid].Hitbox = geometry.NewRoundedRectangle(0, 4, 12, 24, 0, 1)

	Characters[CharFisherman] = NewCharacter(ZoneField, Values{Name: "Fisherman", Wage: 20, Role: RoleGriefer,
		MaxHealth: 0, MoveSpeed: 15, ActPoints: 0, ActTimer: 15.0, ActRange: 0, RespawnTimer: 20.0},
		"🟧"+Tags[IconHand]+"Hooks up ⬜a 🌗🟨"+Tags[IconLeftRight]+"passing by 🟥"+Tags[IconUnit]+"enemy⬜.")
	Characters[CharFisherman].Hitbox.Y = 6

	Characters[CharHorse] = NewCharacter(ZoneField, Values{Name: "Horse", Wage: 20, Role: RoleSupplier,
		MaxHealth: 0, MoveSpeed: 60, ActPoints: 20, ActTimer: 0.0, ActRange: 3, RespawnTimer: 20.0},
		"🟧"+Tags[IconHand]+"Speeds up 🌗🟧"+Tags[IconRange]+"closeby 🟩"+Tags[IconUnit]+"allies⬜.")
	Characters[CharHorse].Hitbox = geometry.NewRoundedRectangle(0, 8, 48, 32, 0, 0.5)

	Characters[CharBrownBunny] = NewCharacter(ZoneField, Values{Name: "Brown Bunny", Wage: 20, Role: RoleCollector,
		MaxHealth: 0, MoveSpeed: 30, ActPoints: 1, ActTimer: 0.0, ActRange: 0, RespawnTimer: 5.0},
		"Very cute and soft.")
	Characters[CharBrownBunny].Hitbox = geometry.NewRoundedRectangle(0, 0, 16, 16, 0, 1)

	Characters[CharWhiteBunny] = NewCharacter(ZoneField, Values{Name: "White Bunny", Wage: 20, Role: RoleCollector,
		MaxHealth: 0, MoveSpeed: 30, ActPoints: 1, ActTimer: 0.0, ActRange: 0, RespawnTimer: 5.0},
		"Very cute and soft.")
	Characters[CharWhiteBunny].Hitbox = geometry.NewRoundedRectangle(0, 0, 16, 16, 0, 1)

	for i, c := range Characters {
		var prefix = text.Replace(text.ToLowerCase(c.Values.Name), " ", "-")
		c.Animations.Idle = atlas.Crops(prefix + "-idle")
		c.Animations.Walk = atlas.Crops(prefix + "-move")
		c.Animations.Prepare = atlas.Crops(prefix + "-prepare")
		c.Animations.Recover = atlas.Crops(prefix + "-recover")
		c.Animations.Hurt = atlas.Crops(prefix + "-hurt")
		c.Animations.Die = atlas.Crops(prefix + "-die")
		c.Icon = atlas.Crops(prefix + "-icon")[0]
		Characters[i] = c
	}
}
