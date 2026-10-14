package game

import (
	char "game/code/constants/character"
	"game/code/constants/icon"
	proj "game/code/constants/projectile"
	"game/code/constants/role"
	"game/code/constants/zone"
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

	EffectTimer      float32 // remaining time before removal (used only for Unit.Effects - not in Unit/Character.Values)
	EffectTickDamage int     // used only for Unit.Effects - not in Unit/Character.Values

	Role role.Enum
}

type CharSounds struct{ ActStart, ActTrigger, HitFlesh, HitWood, HitMetal, HitGround []audio.Audio }
type CharProjectile struct {
	Kind proj.Enum
	ParabolaMultiplier,
	Speed float32
}
type Character struct {
	Values              Values
	Hitbox              geometry.Shape
	Animations          struct{ Idle, Move, Prepare, Recover, Hurt, Die []assets.ImageId }
	Sounds              CharSounds
	Icon                assets.ImageId
	Origin              zone.Enum
	Info, ActPointsName string
	RoleIcon            icon.Enum
	RoleName            string
	IsFlying            bool

	Projectile CharProjectile
}

var Characters [char.Count]*Character
var CharactersFlying = map[char.Enum]struct{}{char.Vulture: {}, char.Raven: {}}

func NewCharacter(origin zone.Enum, hitbox geometry.Shape, stats Values, info string, proj ...CharProjectile) *Character {
	var roleIcons = [role.Count]icon.Enum{icon.Sword, icon.Bow, icon.Shield, icon.Debuff, icon.Hand, icon.Bag}
	var roleNames = [role.Count]string{"Fighter", "Ranger", "Defender", "Griefer", "Supplier", "Collector"}
	var actNames = [role.Count]string{"damage", "damage", "block", "grief", "buff", "pickups"}

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

	Characters[char.Dummy] = NewCharacter(zone.Field, geometry.NewRoundedRectangle(0, 6, 24, 34, 0, 1),
		Values{Name: "Dummy", Role: role.Defender},
		"Cannot 🟥"+Tags[icon.Skull]+"die⬜. Even though it really wants to.")
	Characters[char.Miner] = NewCharacter(zone.Field, geometry.NewRoundedRectangle(0, 6, 18, 35, 0, 1),
		Values{Name: "Miner", Wage: 20, Role: role.Fighter, MaxHealth: 15, MoveSpeed: 20,
			ActPoints: 2, ActTimer: 2.0, ActRange: 1, ReviveTimer: 10.0},
		"🟩"+Tags[icon.Plus]+"🟧"+"8 "+Tags[icon.Sword]+"damage ⬜against "+Tags[icon.Door]+"entrances.")
	Characters[char.Cook] = NewCharacter(zone.Field, geometry.NewRoundedRectangle(0, 0, 18, 35, 0, 1),
		Values{Name: "Cook", Wage: 10, Role: role.Supplier, MoveSpeed: 15,
			ActPoints: 4, ActTimer: 10.0, ActRange: 1, ReviveTimer: 10.0},
		"🟧"+Tags[icon.Hand]+"Heals⬜ a 🌗🟨"+Tags[icon.LeftRight]+"passing by 🟩"+Tags[icon.Unit]+"ally⬜.")
	Characters[char.Bowyer] = NewCharacter(zone.Field, geometry.NewRoundedRectangle(0, 7, 18, 35, 0, 1),
		Values{Name: "Bowyer", Wage: 40, Role: role.Ranger, MaxHealth: 10, MoveSpeed: 15,
			ActPoints: 4, ActTimer: 2.0, ActRange: 6, ReviveTimer: 10.0},
		"🟩"+Tags[icon.Plus]+"🌗🟧2 "+Tags[icon.Range]+"range ⬜when 🌗🟨"+Tags[icon.LeftRight]+"grounded⬜.")
	Characters[char.Bowyer].Sounds = CharSounds{ActTrigger: AudioBow, HitGround: AudioProjectileGround,
		HitFlesh: AudioProjectileFlesh, HitWood: AudioProjectileWood, HitMetal: AudioProjectileMetal}
	Characters[char.Smith] = NewCharacter(zone.Field, geometry.NewRoundedRectangle(0, 0, 24, 48, 0, 1),
		Values{Name: "Smith", Wage: 20, Role: role.Defender, MaxHealth: 40, MoveSpeed: 12,
			ActPoints: 1, ActTimer: 5.0, ActRange: 1, ReviveTimer: 20.0},
		"Pushes the 🟥"+Tags[icon.Unit]+"enemy⬜ in front.")
	Characters[char.Kid] = NewCharacter(zone.Field, geometry.NewRoundedRectangle(0, 4, 12, 24, 0, 1),
		Values{Name: "Kid", Wage: 20, Role: role.Collector, MoveSpeed: 20, ActPoints: 1, ReviveTimer: 20.0},
		"🟩"+Tags[icon.Plus]+"🌗🟨20 "+Tags[icon.LeftRight]+"speed⬜ when 🟧"+Tags[icon.Bag]+"empty handed⬜.")
	Characters[char.Fisherman] = NewCharacter(zone.Field, geometry.NewRoundedRectangle(-4, 6, 20, 35, 0, 1),
		Values{Name: "Fisherman", Wage: 20, Role: role.Griefer, MoveSpeed: 15, ActTimer: 15.0, ReviveTimer: 20.0},
		"Hooks up ⬜a 🌗🟨"+Tags[icon.LeftRight]+"passing by 🟥"+Tags[icon.Unit]+"enemy⬜.")
	Characters[char.Horse] = NewCharacter(zone.Field, geometry.NewRoundedRectangle(0, 8, 48, 32, 0, 0.5),
		Values{Name: "Horse", Wage: 20, Role: role.Supplier, MoveSpeed: 60, ActPoints: 20, ActRange: 3, ReviveTimer: 20.0},
		"🟧"+Tags[icon.Hand]+"Speeds up 🌗🟧"+Tags[icon.Range]+"nearby 🟩"+Tags[icon.Unit]+"allies⬜.")
	Characters[char.Hare] = NewCharacter(zone.Field, geometry.NewRoundedRectangle(0, 0, 16, 16, 0, 1),
		Values{Name: "Hare", Wage: 20, Role: role.Collector, MoveSpeed: 30, ActPoints: 1, ReviveTimer: 5.0},
		"Very cute and soft.")
	Characters[char.Bunny] = NewCharacter(zone.Field, geometry.NewRoundedRectangle(0, 0, 16, 16, 0, 1),
		Values{Name: "Bunny", Wage: 20, Role: role.Collector, MoveSpeed: 30, ActPoints: 1, ReviveTimer: 5.0},
		"Very cute and soft.")

	Characters[char.Troll] = NewCharacter(zone.Ruins, geometry.NewRoundedRectangle(0, 6, 32, 52, 0, 1),
		Values{Name: "Troll", Wage: 20, Role: role.Defender, MaxHealth: 60, MoveSpeed: 10,
			ActPoints: 2, ActTimer: 8.0, ActRange: 1, ReviveTimer: number.NaN()},
		"Cannot 🌗🟦"+Tags[icon.Loop]+"revive.⬜ Pushes 🌗🟧"+Tags[icon.Range]+"nearby 🟥"+Tags[icon.Unit]+"enemies⬜ away.")
	Characters[char.Spearlin] = NewCharacter(zone.Ruins, geometry.NewRoundedRectangle(-2, 0, 18, 30, 0, 1),
		Values{Name: "Spearlin", Wage: 20, Role: role.Fighter,
			MaxHealth: 12, MoveSpeed: 25, ActPoints: 2, ActTimer: 1.5, ActRange: 1, ReviveTimer: 20.0},
		"🟩"+Tags[icon.Plus]+"🟧"+"2 "+Tags[icon.Sword]+"damage ⬜against "+Tags[icon.Shield]+"Defenders.")
	Characters[char.Arbalestier] = NewCharacter(zone.Ruins, geometry.NewRoundedRectangle(-2, 0, 18, 30, 0, 1),
		Values{Name: "Arbalestier", Wage: 20, Role: role.Ranger,
			MaxHealth: 12, MoveSpeed: 25, ActPoints: 6, ActTimer: 3.0, ActRange: 4, ReviveTimer: 20.0},
		"🟩"+Tags[icon.Plus]+"🟧"+"2 "+Tags[icon.Sword]+"damage ⬜against "+Tags[icon.Shield]+"Defenders.",
		CharProjectile{ParabolaMultiplier: 0.2, Speed: 2})
	Characters[char.Stabber] = NewCharacter(zone.Ruins, geometry.NewRoundedRectangle(-4, 7, 20, 33, 0, 1),
		Values{Name: "Stabber", Wage: 20, Role: role.Fighter,
			MaxHealth: 16, MoveSpeed: 20, ActPoints: 2, ActTimer: 1.0, ActRange: 1, ReviveTimer: 30.0},
		"🟩"+Tags[icon.Plus]+"🟧"+"1 "+Tags[icon.Sword]+"damage ⬜upon killing an 🟥"+Tags[icon.Unit]+"enemy⬜.")
	Characters[char.Gunner] = NewCharacter(zone.Ruins, geometry.NewRoundedRectangle(-1, 7, 20, 33, 0, 1),
		Values{Name: "Gunner", Wage: 20, Role: role.Ranger,
			MaxHealth: 14, MoveSpeed: 20, ActPoints: 8, ActTimer: 8.0, ActRange: 8, ReviveTimer: 30.0},
		"🟩"+Tags[icon.Plus]+"🟧"+"4 "+Tags[icon.Sword]+"damage⬜ 🟥"+Tags[icon.Plus]+"🌗🟪3s "+Tags[icon.Timer]+"rest⬜"+
			" when melee vs 🟥"+Tags[icon.Unit]+"enemy⬜.",
		CharProjectile{ParabolaMultiplier: 0.05, Speed: 3.5, Kind: proj.SmallBullet})
	Characters[char.Vulture] = NewCharacter(zone.Ruins, geometry.NewRoundedRectangle(0, 0, 32, 32, 0, 1),
		Values{Name: "Vulture", Wage: 20, Role: role.Supplier, MoveSpeed: 40, ActRange: 3, ReviveTimer: 20.0},
		"🟧"+Tags[icon.Hand]+"Faster 🌗🟦"+Tags[icon.Loop]+"revive⬜ for 🌗🟧"+Tags[icon.Range]+"nearby 🟥"+
			Tags[icon.Skull]+"dead 🟩"+Tags[icon.Unit]+"allies⬜.")
	Characters[char.Raven] = NewCharacter(zone.Ruins, geometry.NewRoundedRectangle(0, 4, 30, 24, 0, 1),
		Values{Name: "Raven", Wage: 20, Role: role.Griefer, MoveSpeed: 40, ActPoints: 10, ActRange: 3, ReviveTimer: 20.0},
		"🟧"+Tags[icon.Debuff]+"Weakens⬜ 🌗🟧"+Tags[icon.Range]+"nearby 🟥"+Tags[icon.Unit]+"enemies⬜.")
	Characters[char.Snake] = NewCharacter(zone.Ruins, geometry.NewRoundedRectangle(0, 0, 30, 18, 0, 1),
		Values{Name: "Snake", Wage: 20, Role: role.Griefer, MoveSpeed: 10, ActPoints: 10, ActRange: 3, ReviveTimer: 10.0},
		"🟧"+Tags[icon.Debuff]+"Slows down 🌗🟧"+Tags[icon.Range]+"nearby 🟥"+Tags[icon.Unit]+"enemies⬜.")

	Characters[char.Pyromancer] = NewCharacter(zone.Ruins, geometry.NewRoundedRectangle(0, 6, 18, 35, 0, 1),
		Values{Name: "Pyromancer", Wage: 20, Role: role.Ranger, MoveSpeed: 20, MaxHealth: 8,
			ActPoints: 0, ActTimer: 8.0, ActRange: 6, ReviveTimer: 30.0},
		"Fire "+Tags[icon.Splash]+"splash. Better against "+Tags[icon.Door]+"wooden doors.",
		CharProjectile{ParabolaMultiplier: 1, Speed: 0.6, Kind: proj.RedPotion})
	Characters[char.Healer] = NewCharacter(zone.Ruins, geometry.NewRoundedRectangle(0, 6, 18, 35, 0, 1),
		Values{Name: "Healer", Wage: 20, Role: role.Supplier, MoveSpeed: 20,
			ActPoints: 2, ActTimer: 5.0, ActRange: 6, ReviveTimer: 30.0},
		"Heal "+Tags[icon.Splash]+"splash.",
		CharProjectile{ParabolaMultiplier: 1, Speed: 0.6, Kind: proj.GreenPotion})
	Characters[char.Toxicologist] = NewCharacter(zone.Ruins, geometry.NewRoundedRectangle(0, 6, 18, 35, 0, 1),
		Values{Name: "Toxicologist", Wage: 20, Role: role.Ranger, MoveSpeed: 20, MaxHealth: 8,
			ActPoints: 2, ActTimer: 5.0, ActRange: 6, ReviveTimer: 30.0},
		"Poison "+Tags[icon.Splash]+"splash. Better against "+Tags[icon.Door]+"metal gates.",
		CharProjectile{ParabolaMultiplier: 1, Speed: 0.6, Kind: proj.YellowPotion})

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
