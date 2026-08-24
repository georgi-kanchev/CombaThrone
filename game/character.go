package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/audio"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/utility/text"
)

type Stats struct {
	Name string

	Health, MoveSpeed, Wage int

	ActValue, ActSpeed, ActRange int

	HurtTime, RespawnTimer float32

	IsOffLaner bool
}

type CharacterKind uint8
type CharSounds struct{ ActionStart, ActionTrigger, HitFlesh, HitWood, HitMetal, HitGround []audio.Audio }
type Character struct {
	Stats      Stats
	Hitbox     geometry.Shape
	Animations struct {
		Idle, Walk, ActionStart, ActionEnd, Hurt, Die []assets.ImageId
	}
	Sounds CharSounds
	Icon   assets.ImageId

	Behavior func(self *Unit)
}

const CharMan, CharWoman, CharHunter CharacterKind = 0, 1, 2

var Characters [3]*Character

func NewCharacter(behavior func(self *Unit), stats Stats) *Character {
	return &Character{Behavior: behavior, Stats: stats, Hitbox: geometry.NewRoundedRectangle(0, 7, 18, 35, 0, 1),
		Sounds: CharSounds{HitFlesh: AudioHitFlesh, HitWood: AudioHitWood, HitMetal: AudioHitMetal}}
}

func InitCharacters() {
	var atlas = assets.LoadAtlas(assets.LoadImage("data/units.png"), "data/units.xml")

	Characters[CharMan] = NewCharacter(BehaviorMan, Stats{Name: "Man", Wage: 20,
		Health: 20, MoveSpeed: 30, HurtTime: 0.5, ActValue: 2, ActSpeed: 15, ActRange: 1, RespawnTimer: 10})

	Characters[CharWoman] = NewCharacter(BehaviorMan, Stats{Name: "Woman", Wage: 10, IsOffLaner: true,
		Health: 12, MoveSpeed: 20, HurtTime: 0.5, ActValue: 1, ActSpeed: 18, ActRange: 1, RespawnTimer: 10})

	Characters[CharHunter] = NewCharacter(BehaviorHunter, Stats{Name: "Hunter", Wage: 40,
		Health: 20, MoveSpeed: 15, HurtTime: 0.5, ActValue: 4, ActSpeed: 20, ActRange: 6, RespawnTimer: 10})
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
