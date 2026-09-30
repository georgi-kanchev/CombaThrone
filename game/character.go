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

	ActTimer, ReviveTimer, SleepTimer float32

	EffectTimer        float32 // remaining time before removal (used only for Unit.Effects - not in Unit/Character.Values)
	EffectTickInterval float32 // used only for Unit.Effects - not in Unit/Character.Values
	EffectTickDamage   int     // used only for Unit.Effects - not in Unit/Character.Values

	Role Role
}

type CharacterKind uint8
type CharSounds struct{ ActStart, ActTrigger, HitFlesh, HitWood, HitMetal, HitGround []audio.Audio }
type CharProjectile struct {
	Kind ProjectileKind
	ParabolaMultiplier,
	Speed float32
}
type Character struct {
	Values              Values
	Hitbox              geometry.Shape
	Animations          struct{ Idle, Move, Prepare, Recover, Hurt, Die []assets.ImageId }
	Sounds              CharSounds
	Icon                assets.ImageId
	Origin              ZoneKind
	Info, ActPointsName string
	RoleIcon            Icon
	RoleName            string
	IsFlying            bool

	Projectile CharProjectile
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
var CharactersFlying = map[CharacterKind]struct{}{CharVulture: {}, CharRaven: {}}

func NewCharacter(origin ZoneKind, hitbox geometry.Shape, stats Values, info string, proj ...CharProjectile) *Character {
	var roleIcons = [RoleCount]Icon{IconSword, IconBow, IconShield, IconDebuff, IconHand, IconBag}
	var roleNames = [RoleCount]string{"Fighter", "Ranger", "Defender", "Griefer", "Supplier", "Collector"}
	var actNames = [RoleCount]string{"damage", "damage", "block", "grief", "buff", "pickup"}

	if len(proj) == 0 {
		proj = []CharProjectile{CharProjectile{ParabolaMultiplier: 1, Speed: 1}}
	}

	return &Character{Values: stats, Hitbox: hitbox,
		Sounds: CharSounds{HitFlesh: AudioHitFlesh, HitWood: AudioHitWood, HitMetal: AudioHitMetal},
		Info:   info, ActPointsName: actNames[stats.Role], RoleIcon: roleIcons[stats.Role], RoleName: roleNames[stats.Role],
		Projectile: proj[0],
	}
}

func InitCharacters() {
	var atlas = assets.LoadAtlas(assets.LoadImage("data/units.png"), "data/units.xml")

	Characters[CharDummy] = NewCharacter(ZoneField, geometry.NewRoundedRectangle(0, 6, 24, 34, 0, 1),
		Values{Name: "Dummy", Role: RoleDefender},
		"Cannot 🟥"+Tags[IconSkull]+"die⬜. Even though it really wants to.")
	Characters[CharMiner] = NewCharacter(ZoneField, geometry.NewRoundedRectangle(0, 6, 18, 35, 0, 1),
		Values{Name: "Miner", Wage: 20, Role: RoleFighter, MaxHealth: 15, MoveSpeed: 20,
			ActPoints: 2, ActTimer: 2.0, ActRange: 1, ReviveTimer: 10.0},
		"🟩"+Tags[IconPlus]+"🟧"+"8 "+Tags[IconSword]+"damage ⬜against "+Tags[IconDoor]+"entrances.")
	Characters[CharCook] = NewCharacter(ZoneField, geometry.NewRoundedRectangle(0, 0, 18, 35, 0, 1),
		Values{Name: "Cook", Wage: 10, Role: RoleSupplier, MoveSpeed: 15,
			ActPoints: 4, ActTimer: 10.0, ActRange: 1, ReviveTimer: 10.0},
		"🟧"+Tags[IconHand]+"Heals⬜ a 🌗🟨"+Tags[IconLeftRight]+"passing by 🟩"+Tags[IconUnit]+"ally⬜.")
	Characters[CharBowyer] = NewCharacter(ZoneField, geometry.NewRoundedRectangle(0, 7, 18, 35, 0, 1),
		Values{Name: "Bowyer", Wage: 40, Role: RoleRanger, MaxHealth: 14, MoveSpeed: 15,
			ActPoints: 4, ActTimer: 2.0, ActRange: 6, ReviveTimer: 10.0},
		"🟩"+Tags[IconPlus]+"🌗🟧2 "+Tags[IconRange]+"range ⬜when 🌗🟨"+Tags[IconLeftRight]+"grounded⬜.")
	Characters[CharBowyer].Sounds = CharSounds{ActTrigger: AudioBow, HitGround: AudioProjectileGround,
		HitFlesh: AudioProjectileFlesh, HitWood: AudioProjectileWood, HitMetal: AudioProjectileMetal}
	Characters[CharSmith] = NewCharacter(ZoneField, geometry.NewRoundedRectangle(0, 0, 24, 48, 0, 1),
		Values{Name: "Smith", Wage: 20, Role: RoleDefender, MaxHealth: 40, MoveSpeed: 12,
			ActPoints: 1, ActTimer: 5.0, ActRange: 1, ReviveTimer: 20.0},
		"Pushes the 🟥"+Tags[IconUnit]+"enemy ⬜in front.")
	Characters[CharKid] = NewCharacter(ZoneField, geometry.NewRoundedRectangle(0, 4, 12, 24, 0, 1),
		Values{Name: "Kid", Wage: 20, Role: RoleCollector, MoveSpeed: 20, ActPoints: 1, ReviveTimer: 20.0},
		"🟩"+Tags[IconPlus]+"🌗🟨20 "+Tags[IconLeftRight]+"speed⬜ when 🟧"+Tags[IconBag]+"empty handed⬜.")
	Characters[CharFisherman] = NewCharacter(ZoneField, geometry.NewRoundedRectangle(-4, 6, 20, 35, 0, 1),
		Values{Name: "Fisherman", Wage: 20, Role: RoleGriefer, MoveSpeed: 15, ActTimer: 15.0, ReviveTimer: 20.0},
		"Hooks up ⬜a 🌗🟨"+Tags[IconLeftRight]+"passing by 🟥"+Tags[IconUnit]+"enemy⬜.")
	Characters[CharHorse] = NewCharacter(ZoneField, geometry.NewRoundedRectangle(0, 8, 48, 32, 0, 0.5),
		Values{Name: "Horse", Wage: 20, Role: RoleSupplier, MoveSpeed: 60, ActPoints: 20, ActRange: 3, ReviveTimer: 20.0},
		"🟧"+Tags[IconHand]+"Speeds up 🌗🟧"+Tags[IconRange]+"nearby 🟩"+Tags[IconUnit]+"allies⬜.")
	Characters[CharBrownBunny] = NewCharacter(ZoneField, geometry.NewRoundedRectangle(0, 0, 16, 16, 0, 1),
		Values{Name: "Brown\nBunny", Wage: 20, Role: RoleCollector, MoveSpeed: 30, ActPoints: 1, ReviveTimer: 5.0},
		"Very cute and soft.")
	Characters[CharWhiteBunny] = NewCharacter(ZoneField, geometry.NewRoundedRectangle(0, 0, 16, 16, 0, 1),
		Values{Name: "White\nBunny", Wage: 20, Role: RoleCollector, MoveSpeed: 30, ActPoints: 1, ReviveTimer: 5.0},
		"Very cute and soft.")

	Characters[CharTroll] = NewCharacter(ZoneRuins, geometry.NewRoundedRectangle(0, 6, 32, 52, 0, 1),
		Values{Name: "Troll", Wage: 20, Role: RoleDefender, MaxHealth: 60, MoveSpeed: 10,
			ActPoints: 2, ActTimer: 8.0, ActRange: 1, ReviveTimer: number.NaN()},
		"Cannot 🌗🟦"+Tags[IconLoop]+"revive. Pushes 🌗🟧"+Tags[IconRange]+"nearby 🟥"+Tags[IconUnit]+"enemies⬜ away.")
	Characters[CharSkirmisherGoblin] = NewCharacter(ZoneRuins, geometry.NewRoundedRectangle(-2, 0, 18, 30, 0, 1),
		Values{Name: "Skirmisher\nGoblin", Wage: 20, Role: RoleFighter,
			MaxHealth: 12, MoveSpeed: 25, ActPoints: 2, ActTimer: 1.5, ActRange: 1, ReviveTimer: 20.0},
		"🟩"+Tags[IconPlus]+"🟧"+"2 "+Tags[IconSword]+"damage ⬜against "+Tags[IconShield]+"Defenders.")
	Characters[CharArbalestierGoblin] = NewCharacter(ZoneRuins, geometry.NewRoundedRectangle(-2, 0, 18, 30, 0, 1),
		Values{Name: "Arbalestier\nGoblin", Wage: 20, Role: RoleRanger,
			MaxHealth: 12, MoveSpeed: 25, ActPoints: 6, ActTimer: 3.0, ActRange: 4, ReviveTimer: 20.0},
		"🟩"+Tags[IconPlus]+"🟧"+"2 "+Tags[IconSword]+"damage ⬜against "+Tags[IconShield]+"Defenders.",
		CharProjectile{ParabolaMultiplier: 0.2, Speed: 2})
	Characters[CharStabberBandit] = NewCharacter(ZoneRuins, geometry.NewRoundedRectangle(-4, 7, 20, 33, 0, 1),
		Values{Name: "Stabber\nBandit", Wage: 20, Role: RoleFighter,
			MaxHealth: 16, MoveSpeed: 20, ActPoints: 2, ActTimer: 1.0, ActRange: 1, ReviveTimer: 30.0},
		"🟩"+Tags[IconPlus]+"🟧"+"1 "+Tags[IconSword]+"damage ⬜upon killing an 🟥"+Tags[IconUnit]+"enemy⬜.")
	Characters[CharGunnerBandit] = NewCharacter(ZoneRuins, geometry.NewRoundedRectangle(-1, 7, 20, 33, 0, 1),
		Values{Name: "Gunner\nBandit", Wage: 20, Role: RoleRanger,
			MaxHealth: 14, MoveSpeed: 20, ActPoints: 8, ActTimer: 5.0, ActRange: 8, ReviveTimer: 30.0},
		"🟩"+Tags[IconPlus]+"🟧"+"4 "+Tags[IconSword]+"damage⬜ 🟥"+Tags[IconPlus]+"🌗🟪3s "+Tags[IconTimer]+"rest⬜"+
			" when melee vs 🟥"+Tags[IconUnit]+"enemy⬜.",
		CharProjectile{ParabolaMultiplier: 0.05, Speed: 3.5, Kind: ProjectileBullet})
	Characters[CharVulture] = NewCharacter(ZoneRuins, geometry.NewRoundedRectangle(0, 0, 32, 32, 0, 1),
		Values{Name: "Vulture", Wage: 20, Role: RoleSupplier, MoveSpeed: 40, ActRange: 3, ReviveTimer: 20.0},
		"🟧"+Tags[IconHand]+"Faster 🌗🟦"+Tags[IconLoop]+"revive⬜ for 🌗🟧"+Tags[IconRange]+"nearby 🟥"+
			Tags[IconSkull]+"dead 🟩"+Tags[IconUnit]+"allies⬜.")
	Characters[CharRaven] = NewCharacter(ZoneRuins, geometry.NewRoundedRectangle(0, 4, 30, 24, 0, 1),
		Values{Name: "Raven", Wage: 20, Role: RoleGriefer, MoveSpeed: 40, ActPoints: 10, ActRange: 3, ReviveTimer: 20.0},
		"🟧"+Tags[IconDebuff]+"Weakens⬜ 🌗🟧"+Tags[IconRange]+"nearby 🟥"+Tags[IconUnit]+"enemies⬜.")
	Characters[CharRedSnake] = NewCharacter(ZoneRuins, geometry.NewRoundedRectangle(0, 0, 30, 18, 0, 1),
		Values{Name: "Red Snake", Wage: 20, Role: RoleGriefer, MoveSpeed: 10, ActPoints: 10, ActRange: 3, ReviveTimer: 10.0},
		"🟧"+Tags[IconDebuff]+"Slows down 🌗🟧"+Tags[IconRange]+"nearby 🟥"+Tags[IconUnit]+"enemies⬜.")

	for i, c := range Characters {
		if c == nil {
			continue
		}
		var prefix = text.Replace(text.Replace(text.ToLowerCase(c.Values.Name), " ", "-"), "\n", "-")
		c.Animations.Idle = atlas.Crops(prefix + "-idle")
		c.Animations.Move = atlas.Crops(prefix + "-move")
		c.Animations.Prepare = atlas.Crops(prefix + "-prepare")
		c.Animations.Recover = atlas.Crops(prefix + "-recover")
		c.Animations.Hurt = atlas.Crops(prefix + "-hurt")
		c.Animations.Die = atlas.Crops(prefix + "-die")

		var _, flying = CharactersFlying[CharacterKind(i)]
		c.IsFlying = flying
		c.Icon = atlas.Crops(prefix + "-icon")[0]

		if c.IsFlying {
			c.Hitbox.Height += TileSize
			c.Hitbox.Y += TileSize / 2
		}

		Characters[i] = c
	}
}
