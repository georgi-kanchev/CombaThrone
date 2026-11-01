package game

import (
	char "game/code/constants/character"
	"game/code/constants/icon"
	"game/code/constants/role"
	proj "game/code/constants/shot"
	"game/code/constants/sound"
	"game/code/constants/splash"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/text"
)

type Values struct {
	Name, EffectInfo string

	MaxHealth, MoveSpeed, Wage int
	ActPoints, ActRange        int

	ActTimer, ReviveTimer float32

	DisablesTimer float32 // cannot move/act when disabled (always 0 or above because effects add accumulatevly to it)

	EffectTimer    float32 // remaining time before removal (used only for Unit.Effects - not in Unit/Character.Values)
	EffectIconText string

	CannotActVsEntrances bool

	Role role.Enum
}

type CharSounds struct{ ActStart, ActTrigger, HitFlesh, HitWood, HitMetal, HitGround sound.Enum }
type ShotInfo struct {
	Kind                      proj.Enum
	ParabolaMultiplier, Speed float32
	CastCustomSplashInstead   splash.Enum
	Tint                      uint
}
type Character struct {
	Values              Values
	Hitbox              geometry.Shape
	Animations          struct{ Idle, Move, Prepare, Recover, Hurt, Die []assets.ImageId }
	Sounds              CharSounds
	Icon                assets.ImageId
	Info, ActPointsName string
	RoleIcon            icon.Enum
	RoleName            string
	IsFlying            bool

	ShotInfo ShotInfo
}

var Characters [char.COUNT]*Character
var CharactersFlying = map[char.Enum]struct{}{char.Vulture: {}, char.Raven: {}}

func NewCharacter(hitbox geometry.Shape, stats Values, info string, shot *ShotInfo, sounds *CharSounds) *Character {
	var roleIcons = [role.COUNT]icon.Enum{icon.Fighter, icon.Backliner, icon.Defender, icon.Griefer, icon.Supplier, icon.Collector}
	var roleNames = [role.COUNT]string{"Fighter", "Backliner", "Defender", "Griefer", "Supplier", "Collector"}
	var actNames = [role.COUNT]string{"damage", "damage", "block", "grief", "buff", "pickups"}

	if shot == nil {
		shot = new(ShotInfo{Kind: proj.Arrow, ParabolaMultiplier: 1, Speed: 1})
	}
	if sounds == nil {
		sounds = new(CharSounds{HitFlesh: sound.HitFlesh, HitWood: sound.HitWood, HitMetal: sound.HitMetal})
	}

	return &Character{Values: stats, Hitbox: hitbox, Sounds: *sounds, Info: info, ActPointsName: actNames[stats.Role],
		RoleIcon: roleIcons[stats.Role], RoleName: roleNames[stats.Role], ShotInfo: *shot}
}

func InitCharacters() {
	var atlas = assets.LoadAtlas(assets.LoadImage("data/units.png"), "data/units.xml")

	Characters[char.Dummy] = NewCharacter(geometry.NewRoundedRectangle(0, 6, 24, 34, 0, 1),
		Values{Name: "Dummy", Role: role.Defender},
		"Cannot 🟥"+Tags[icon.Skull]+"die⬜. Even though it really wants to.", nil, nil)
	Characters[char.Miner] = NewCharacter(geometry.NewRoundedRectangle(0, 6, 18, 35, 0, 1),
		Values{Name: "Miner", Wage: 20, Role: role.Fighter, MaxHealth: 15, MoveSpeed: 20,
			ActPoints: 2, ActTimer: 2.0, ActRange: 1, ReviveTimer: 10.0},
		"🟩"+Tags[icon.Plus]+"🟧"+"8 "+Tags[icon.Wound]+"damage ⬜against "+Tags[icon.Door]+"entrances.", nil, nil)
	Characters[char.Cook] = NewCharacter(geometry.NewRoundedRectangle(0, 0, 18, 35, 0, 1),
		Values{Name: "Cook", Wage: 10, Role: role.Supplier, MoveSpeed: 15,
			ActPoints: 4, ActTimer: 10.0, ActRange: 1, ReviveTimer: 10.0},
		"🟧"+Tags[icon.Supplier]+"Heals⬜ a 🌗🟨"+Tags[icon.Speed]+"passing by 🟩"+Tags[icon.Unit]+"ally⬜.", nil, nil)
	Characters[char.Bowyer] = NewCharacter(geometry.NewRoundedRectangle(0, 7, 18, 35, 0, 1),
		Values{Name: "Bowyer", Wage: 40, Role: role.Backliner, MaxHealth: 10, MoveSpeed: 15,
			ActPoints: 4, ActTimer: 2.0, ActRange: 6, ReviveTimer: 10.0},
		"🟩"+Tags[icon.Plus]+"🌗🟧2 "+Tags[icon.Range]+"range ⬜when 🌗🟨"+Tags[icon.Speed]+"grounded⬜.",
		new(ShotInfo{Kind: proj.Arrow, ParabolaMultiplier: 1, Speed: 1}),
		new(CharSounds{ActTrigger: sound.Bow, HitGround: sound.ShotGround,
			HitFlesh: sound.ShotFlesh, HitWood: sound.ShotWood, HitMetal: sound.ShotMetal}))
	Characters[char.Smith] = NewCharacter(geometry.NewRoundedRectangle(0, 0, 24, 48, 0, 1),
		Values{Name: "Smith", Wage: 20, Role: role.Defender, MaxHealth: 40, MoveSpeed: 12, ActPoints: 1,
			ActTimer: 5.0, ActRange: 1, ReviveTimer: 20.0}, "Pushes the 🟥"+Tags[icon.Unit]+"enemy⬜ in front.", nil, nil)
	Characters[char.Kid] = NewCharacter(geometry.NewRoundedRectangle(0, 4, 12, 24, 0, 1),
		Values{Name: "Kid", Wage: 20, Role: role.Collector, MoveSpeed: 20, ActPoints: 1, ReviveTimer: 20.0},
		"🟩"+Tags[icon.Plus]+"🌗🟨20 "+Tags[icon.Speed]+"speed⬜ when 🟧"+Tags[icon.Collector]+"empty handed⬜.", nil, nil)
	Characters[char.Fisherman] = NewCharacter(geometry.NewRoundedRectangle(-4, 6, 20, 35, 0, 1),
		Values{Name: "Fisherman", Wage: 20, Role: role.Griefer, MoveSpeed: 15, ActTimer: 15.0, ReviveTimer: 20.0},
		"Hooks up ⬜a 🌗🟨"+Tags[icon.Speed]+"passing by 🟥"+Tags[icon.Unit]+"enemy⬜.", nil, nil)
	Characters[char.Horse] = NewCharacter(geometry.NewRoundedRectangle(0, 8, 48, 32, 0, 0.5),
		Values{Name: "Horse", Wage: 20, Role: role.Supplier, MoveSpeed: 60, ActPoints: 20, ActRange: 3, ReviveTimer: 20.0},
		"🟧"+Tags[icon.Supplier]+"Speeds up 🌗🟧"+Tags[icon.Range]+"nearby 🟩"+Tags[icon.Unit]+"allies⬜.", nil, nil)
	Characters[char.Hare] = NewCharacter(geometry.NewRoundedRectangle(0, 0, 16, 16, 0, 1),
		Values{Name: "Hare", Wage: 20, Role: role.Collector, MoveSpeed: 30, ActPoints: 1, ReviveTimer: 5.0},
		"Very cute and soft.", nil, nil)
	Characters[char.Bunny] = NewCharacter(geometry.NewRoundedRectangle(0, 0, 16, 16, 0, 1),
		Values{Name: "Bunny", Wage: 20, Role: role.Collector, MoveSpeed: 30, ActPoints: 1, ReviveTimer: 5.0},
		"Very cute and soft.", nil, nil)

	Characters[char.Troll] = NewCharacter(geometry.NewRoundedRectangle(0, 6, 32, 52, 0, 1),
		Values{Name: "Troll", Wage: 20, Role: role.Defender, MaxHealth: 60, MoveSpeed: 10,
			ActPoints: 2, ActTimer: 8.0, ActRange: 1, ReviveTimer: number.NaN()},
		"Cannot 🌗🟦"+Tags[icon.Revive]+"revive.⬜ Pushes 🌗🟧"+Tags[icon.Range]+"nearby 🟥"+Tags[icon.Unit]+"enemies⬜ away.",
		nil, nil)
	Characters[char.Spearlin] = NewCharacter(geometry.NewRoundedRectangle(-2, 0, 18, 30, 0, 1),
		Values{Name: "Spearlin", Wage: 20, Role: role.Fighter,
			MaxHealth: 12, MoveSpeed: 25, ActPoints: 2, ActTimer: 1.5, ActRange: 1, ReviveTimer: 20.0},
		"🟩"+Tags[icon.Plus]+"🟧"+"2 "+Tags[icon.Wound]+"damage ⬜against "+Tags[icon.Defender]+"Defenders.", nil,
		new(CharSounds{HitFlesh: sound.Stab, HitWood: sound.HitWood, HitMetal: sound.HitMetal}))
	Characters[char.Arbalestier] = NewCharacter(geometry.NewRoundedRectangle(-2, 0, 18, 30, 0, 1),
		Values{Name: "Arbalestier", Wage: 20, Role: role.Backliner,
			MaxHealth: 12, MoveSpeed: 25, ActPoints: 6, ActTimer: 3.0, ActRange: 4, ReviveTimer: 20.0},
		"🟩"+Tags[icon.Plus]+"🟧"+"2 "+Tags[icon.Wound]+"damage ⬜against "+Tags[icon.Defender]+"Defenders.",
		new(ShotInfo{ParabolaMultiplier: 0.2, Speed: 2}), nil)
	Characters[char.Stabber] = NewCharacter(geometry.NewRoundedRectangle(-4, 7, 20, 33, 0, 1),
		Values{Name: "Stabber", Wage: 20, Role: role.Fighter,
			MaxHealth: 16, MoveSpeed: 20, ActPoints: 2, ActTimer: 1.0, ActRange: 1, ReviveTimer: 30.0},
		"🟩"+Tags[icon.Plus]+"🟧"+"2 "+Tags[icon.Wound]+"damage ⬜upon killing an 🟥"+Tags[icon.Unit]+"enemy⬜.", nil, nil)
	Characters[char.Gunner] = NewCharacter(geometry.NewRoundedRectangle(-1, 7, 20, 33, 0, 1),
		Values{Name: "Gunner", Wage: 20, Role: role.Backliner,
			MaxHealth: 14, MoveSpeed: 20, ActPoints: 8, ActTimer: 8.0, ActRange: 8, ReviveTimer: 30.0},
		"🟩"+Tags[icon.Plus]+"🟧"+"4 "+Tags[icon.Wound]+"damage⬜ 🟥"+Tags[icon.Plus]+"🌗🟪3s "+Tags[icon.Hourglass]+"rest⬜"+
			" when melee vs 🟥"+Tags[icon.Unit]+"enemy⬜.",
		new(ShotInfo{ParabolaMultiplier: 0.05, Speed: 3.5, Kind: proj.SmallBullet}), nil)
	Characters[char.Vulture] = NewCharacter(geometry.NewRoundedRectangle(0, 0, 32, 32, 0, 1),
		Values{Name: "Vulture", Wage: 20, Role: role.Supplier, MoveSpeed: 40, ActRange: 3, ReviveTimer: 20.0},
		"🌓⬜"+Tags[icon.Skull]+"Dead 🟩"+Tags[icon.Unit]+"allies 🌗🟧"+Tags[icon.Range]+
			"nearby 🌗🟦"+Tags[icon.Revive]+"revive⬜ twice as fast.", nil, nil)
	Characters[char.Raven] = NewCharacter(geometry.NewRoundedRectangle(0, 4, 30, 24, 0, 1),
		Values{Name: "Raven", Wage: 20, Role: role.Griefer, MoveSpeed: 40, ActPoints: 10, ActRange: 3, ReviveTimer: 20.0},
		"🟥"+Tags[icon.Minus]+"🟧10 "+Tags[icon.Wound]+"damage⬜ to 🌗🟧"+Tags[icon.Range]+"nearby 🟥"+
			Tags[icon.Unit]+"enemies⬜.", nil, nil)
	Characters[char.Snake] = NewCharacter(geometry.NewRoundedRectangle(0, 0, 30, 18, 0, 1),
		Values{Name: "Snake", Wage: 20, Role: role.Griefer, MoveSpeed: 10, ActPoints: 10, ActRange: 3, ReviveTimer: 10.0},
		"🟧"+Tags[icon.Griefer]+"Slows down 🌗🟧"+Tags[icon.Range]+"nearby 🟥"+Tags[icon.Unit]+"enemies⬜.", nil, nil)

	Characters[char.Pyromancer] = NewCharacter(geometry.NewRoundedRectangle(0, 6, 18, 35, 0, 1),
		Values{Name: "Pyromancer", Wage: 20, Role: role.Backliner, MoveSpeed: 20, MaxHealth: 8,
			ActPoints: 1, ActTimer: 8.0, ActRange: 6, ReviveTimer: 30.0},
		"Fire "+Tags[icon.Splash]+"splash. Better against "+Tags[icon.Door]+"wooden doors.",
		new(ShotInfo{ParabolaMultiplier: 1, Speed: 0.6, Kind: proj.RedPotion}), nil)
	Characters[char.Healer] = NewCharacter(geometry.NewRoundedRectangle(0, 6, 18, 35, 0, 1),
		Values{Name: "Healer", Wage: 20, Role: role.Supplier, MoveSpeed: 20,
			ActPoints: 2, ActTimer: 5.0, ActRange: 6, ReviveTimer: 30.0},
		"Heal "+Tags[icon.Splash]+"splash.",
		new(ShotInfo{ParabolaMultiplier: 1, Speed: 0.6, Kind: proj.GreenPotion}), nil)
	Characters[char.Toxicologist] = NewCharacter(geometry.NewRoundedRectangle(0, 6, 18, 35, 0, 1),
		Values{Name: "Toxicologist", Wage: 20, Role: role.Backliner, MoveSpeed: 20, MaxHealth: 8,
			ActPoints: 2, ActTimer: 5.0, ActRange: 6, ReviveTimer: 30.0},
		"Poison "+Tags[icon.Splash]+"splash. Better against "+Tags[icon.Door]+"metal gates.",
		new(ShotInfo{ParabolaMultiplier: 1, Speed: 0.6, Kind: proj.YellowPotion}), nil)
	Characters[char.Windbender] = NewCharacter(geometry.NewRoundedRectangle(0, 4, 18, 40, 0, 1),
		Values{Name: "Windbender", Wage: 20, Role: role.Backliner, MoveSpeed: 20, MaxHealth: 8,
			ActPoints: 6, ActTimer: 12.0, ActRange: 3, ReviveTimer: 30.0, CannotActVsEntrances: true},
		"Casts a whirlwind. Useless against "+Tags[icon.Door]+"entrances.",
		new(ShotInfo{CastCustomSplashInstead: splash.Whirlwind}), nil)
	Characters[char.Aquaweaver] = NewCharacter(geometry.NewRoundedRectangle(0, 4, 18, 40, 0, 1),
		Values{Name: "Aquaweaver", Wage: 20, Role: role.Defender, MoveSpeed: 20, MaxHealth: 16,
			ActPoints: 3, ActTimer: 12.0, ActRange: 1, ReviveTimer: 30.0},
		"Casts a water shield, blocking all damage.",
		new(ShotInfo{CastCustomSplashInstead: splash.WaterShield}), nil)
	Characters[char.Stormcaller] = NewCharacter(geometry.NewRoundedRectangle(0, 4, 18, 40, 0, 1),
		Values{Name: "Stormcaller", Wage: 20, Role: role.Backliner, MoveSpeed: 20, MaxHealth: 8,
			ActPoints: 6, ActTimer: 12.0, ActRange: 2, ReviveTimer: 30.0, CannotActVsEntrances: true},
		"Casts a zap. Useless against "+Tags[icon.Door]+"entrances.",
		new(ShotInfo{CastCustomSplashInstead: splash.LightningZap}), nil)

	Characters[char.Bear] = NewCharacter(geometry.NewRoundedRectangle(0, 4, 40, 25, 0, 1),
		Values{Name: "Bear", Wage: 20, Role: role.Griefer, MoveSpeed: 15, ActRange: 2, ReviveTimer: 30.0},
		"Inflicts 🌓⬜"+Tags[icon.Fear]+"fear⬜ to 🌗🟧"+Tags[icon.Range]+"nearby 🟥"+Tags[icon.Unit]+"enemies⬜.", nil, nil)
	Characters[char.Wolf] = NewCharacter(geometry.NewRoundedRectangle(0, 4, 40, 25, 0, 1),
		Values{Name: "Wolf", Wage: 20, Role: role.Supplier, MoveSpeed: 35, ReviveTimer: 30.0},
		"Follows an 🟩"+Tags[icon.Unit]+"ally⬜. Makes them 🌗🟪"+Tags[icon.Hourglass]+"rest⬜ 50% faster.", nil, nil)
	Characters[char.Boar] = NewCharacter(geometry.NewRoundedRectangle(0, 4, 40, 25, 0, 1),
		Values{Name: "Boar", Wage: 20, Role: role.Griefer, MoveSpeed: 60, ReviveTimer: 30.0},
		"Charges, then 🟧"+Tags[icon.Griefer]+"pushes⬜ and 🌓⬜"+Tags[icon.Stun]+"stuns⬜ for 4s.", nil, nil)

	Characters[char.Shaman] = NewCharacter(geometry.NewRoundedRectangle(-1, 9, 20, 30, 0, 1),
		Values{Name: "Shaman", Wage: 20, Role: role.Backliner,
			MaxHealth: 14, MoveSpeed: 20, ActPoints: 5, ActTimer: 5.0, ActRange: 5, ReviveTimer: 30.0},
		"Hitting an 🟥"+Tags[icon.Unit]+"enemy⬜ causes 🟥"+Tags[icon.Minus]+"🟧5 "+Tags[icon.Wound]+"damage⬜ for 3s.",
		new(ShotInfo{ParabolaMultiplier: 0.4, Speed: 2.0, Kind: proj.MediumSpit, Tint: color.TagHex("#ccf61f")}), nil)

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

		var _, flying = CharactersFlying[char.Enum(i)]
		c.IsFlying = flying
		c.Icon = atlas.Crops(prefix + "-icon")[0]

		if c.IsFlying {
			c.Hitbox.Height += TileSize
			c.Hitbox.Y += TileSize / 2
		}

		Characters[i] = c
	}
}
