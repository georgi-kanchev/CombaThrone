package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/audio"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/utility/number"
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

	Projectile ProjectileKind
	ProjectileParabolaMultiplier,
	ProjectileSpeed float32
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

	CharTroll
	CharArbalestierGoblin
	CharSkirmisherGoblin
	CharStabberBandit
	CharGunnerBandit
	CharVulture
	CharRaven
	CharRedSnake

	CharCount
)

var Characters [CharCount]*Character

func NewCharacter(origin ZoneKind, hitbox geometry.Shape, stats Values, info string, projParabSpeed ...float32) *Character {
	var roleIcons = [RoleCount]Icon{IconSword, IconBow, IconShield, IconDebuff, IconHand, IconBag}
	var roleNames = [RoleCount]string{"Fighter", "Ranger", "Defender", "Griefer", "Supplier", "Collector"}
	var actNames = [RoleCount]string{"damage", "damage", "block", "grief", "buff", "pickup"}

	if len(projParabSpeed) == 0 {
		projParabSpeed = []float32{1, 1}
	}

	return &Character{Values: stats, Hitbox: hitbox,
		Sounds: CharSounds{HitFlesh: AudioHitFlesh, HitWood: AudioHitWood, HitMetal: AudioHitMetal},
		Info:   info, ActPointsName: actNames[stats.Role], RoleIcon: roleIcons[stats.Role], RoleName: roleNames[stats.Role],
		ProjectileParabolaMultiplier: projParabSpeed[0], ProjectileSpeed: projParabSpeed[1],
	}
}

func InitCharacters() {
	var atlas = assets.LoadAtlas(assets.LoadImage("data/units.png"), "data/units.xml")

	Characters[CharDummy] = NewCharacter(ZoneField, geometry.NewRoundedRectangle(0, 6, 24, 34, 0, 1),
		Values{Name: "Dummy", Role: RoleDefender},
		"Cannot 🟥"+Tags[IconSkull]+"die⬜. Even though it really wants to.")

	Characters[CharMiner] = NewCharacter(ZoneField, geometry.NewRoundedRectangle(0, 6, 18, 35, 0, 1),
		Values{Name: "Miner", Wage: 20, Role: RoleFighter, MaxHealth: 20, MoveSpeed: 20,
			ActPoints: 2, ActTimer: 2.0, ActRange: 1, RespawnTimer: 10.0},
		"🟩"+Tags[IconPlus]+"🟧"+"8 "+Tags[IconSword]+"damage ⬜against "+Tags[IconDoor]+"entrances.")

	Characters[CharCook] = NewCharacter(ZoneField, geometry.NewRoundedRectangle(0, 0, 18, 35, 0, 1),
		Values{Name: "Cook", Wage: 10, Role: RoleSupplier, MoveSpeed: 15,
			ActPoints: 4, ActTimer: 10.0, ActRange: 1, RespawnTimer: 10.0},
		"🟧"+Tags[IconHand]+"Heals⬜ a 🌗🟨"+Tags[IconLeftRight]+"passing by 🟩"+Tags[IconUnit]+"ally⬜.")

	Characters[CharBowyer] = NewCharacter(ZoneField, geometry.NewRoundedRectangle(0, 7, 18, 35, 0, 1),
		Values{Name: "Bowyer", Wage: 40, Role: RoleRanger, MaxHealth: 14, MoveSpeed: 15,
			ActPoints: 4, ActTimer: 2.0, ActRange: 6, RespawnTimer: 10.0},
		"🟩"+Tags[IconPlus]+"🌗🟧2 "+Tags[IconRange]+"range ⬜when 🌗🟨"+Tags[IconLeftRight]+"grounded⬜.")
	Characters[CharBowyer].Sounds = CharSounds{ActTrigger: AudioBow, HitGround: AudioProjectileGround,
		HitFlesh: AudioProjectileFlesh, HitWood: AudioProjectileWood, HitMetal: AudioProjectileMetal}

	Characters[CharSmith] = NewCharacter(ZoneField, geometry.NewRoundedRectangle(0, 0, 24, 48, 0, 1),
		Values{Name: "Smith", Wage: 20, Role: RoleDefender, MaxHealth: 40, MoveSpeed: 12,
			ActPoints: 1, ActTimer: 5.0, ActRange: 1, RespawnTimer: 20.0},
		"Pushes the 🟥"+Tags[IconUnit]+"enemy ⬜in front.")

	Characters[CharKid] = NewCharacter(ZoneField, geometry.NewRoundedRectangle(0, 4, 12, 24, 0, 1),
		Values{Name: "Kid", Wage: 20, Role: RoleCollector, MoveSpeed: 20, ActPoints: 1, RespawnTimer: 20.0},
		"🟩"+Tags[IconPlus]+"🌗🟨20 "+Tags[IconLeftRight]+"speed⬜ when 🟧"+Tags[IconBag]+"empty handed⬜.")

	Characters[CharFisherman] = NewCharacter(ZoneField, geometry.NewRoundedRectangle(0, 6, 18, 35, 0, 1),
		Values{Name: "Fisherman", Wage: 20, Role: RoleGriefer, MoveSpeed: 15, ActTimer: 15.0, RespawnTimer: 20.0},
		"🟧"+Tags[IconHand]+"Hooks up ⬜a 🌗🟨"+Tags[IconLeftRight]+"passing by 🟥"+Tags[IconUnit]+"enemy⬜.")

	Characters[CharHorse] = NewCharacter(ZoneField, geometry.NewRoundedRectangle(0, 8, 48, 32, 0, 0.5),
		Values{Name: "Horse", Wage: 20, Role: RoleSupplier, MoveSpeed: 60, ActPoints: 20, ActRange: 3, RespawnTimer: 20.0},
		"🟧"+Tags[IconHand]+"Speeds up 🌗🟧"+Tags[IconRange]+"nearby 🟩"+Tags[IconUnit]+"allies⬜.")

	Characters[CharBrownBunny] = NewCharacter(ZoneField, geometry.NewRoundedRectangle(0, 0, 16, 16, 0, 1),
		Values{Name: "Brown Bunny", Wage: 20, Role: RoleCollector, MoveSpeed: 30, ActPoints: 1, RespawnTimer: 5.0},
		"Very cute and soft.")

	Characters[CharWhiteBunny] = NewCharacter(ZoneField, geometry.NewRoundedRectangle(0, 0, 16, 16, 0, 1),
		Values{Name: "White Bunny", Wage: 20, Role: RoleCollector, MoveSpeed: 30, ActPoints: 1, RespawnTimer: 5.0},
		"Very cute and soft.")

	//=================================================================

	Characters[CharTroll] = NewCharacter(ZoneRuins, geometry.NewRoundedRectangle(0, 6, 32, 52, 0, 1),
		Values{Name: "Troll", Wage: 20, Role: RoleDefender, MaxHealth: 60, MoveSpeed: 10,
			ActPoints: 2, ActTimer: 8, ActRange: 1, RespawnTimer: number.NaN()},
		"Cannot 🌗🟦"+Tags[IconLoop]+"respawn⬜. Pushes 🌗🟧"+Tags[IconRange]+"nearby 🟥"+Tags[IconUnit]+"enemies⬜ away.")

	Characters[CharSkirmisherGoblin] = NewCharacter(ZoneRuins, geometry.NewRoundedRectangle(-2, 0, 18, 30, 0, 1),
		Values{Name: "Skirmisher\nGoblin", Wage: 20, Role: RoleFighter,
			MaxHealth: 15, MoveSpeed: 25, ActPoints: 2, ActTimer: 1.5, ActRange: 1, RespawnTimer: 20.0},
		"🟩"+Tags[IconPlus]+"🟧"+"2 "+Tags[IconSword]+"damage ⬜against "+Tags[IconShield]+"Defenders.")

	Characters[CharArbalestierGoblin] = NewCharacter(ZoneRuins, geometry.NewRoundedRectangle(-2, 0, 18, 30, 0, 1),
		Values{Name: "Arbalestier\nGoblin", Wage: 20, Role: RoleRanger,
			MaxHealth: 15, MoveSpeed: 25, ActPoints: 4, ActTimer: 3.0, ActRange: 4, RespawnTimer: 20.0},
		"🟩"+Tags[IconPlus]+"🟧"+"2 "+Tags[IconSword]+"damage ⬜against "+Tags[IconShield]+"Defenders.", 0.2, 2)

	for i, c := range Characters {
		if c == nil {
			continue
		}
		var prefix = text.Replace(text.Replace(text.ToLowerCase(c.Values.Name), " ", "-"), "\n", "-")
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
