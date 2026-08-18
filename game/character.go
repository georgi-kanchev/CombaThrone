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

	HurtTime float32
}

type CharSounds struct{ ActionStart, ActionTrigger, HitFlesh, HitWood, HitMetal, HitGround []audio.Audio }

type CharacterData struct {
	Stats      Stats
	Hitbox     geometry.Shape
	Animations struct {
		Idle, Walk, ActionStart, ActionEnd, Hurt, Die []assets.ImageId
	}
	Sounds CharSounds

	Behavior func(self *Unit)
}

const CharMan, CharWoman, CharHunter Character = 0, 1, 2

var Characters [3]CharacterData

func NewCharacter(behavior func(self *Unit), stats Stats) CharacterData {
	return CharacterData{Behavior: behavior, Stats: stats, Hitbox: geometry.NewRoundedRectangle(0, 7, 18, 35, 0, 1),
		Sounds: CharSounds{HitFlesh: AudioHitFlesh, HitWood: AudioHitWood, HitMetal: AudioHitMetal}}
}

func InitCharacters() {
	var atlas = assets.LoadAtlas(assets.LoadImage("data/units.png"), "data/units.xml")

	Characters[CharMan] = NewCharacter(BehaviorMan, Stats{Name: "Man", Wage: 20,
		Health: 20, MoveSpeed: 30, HurtTime: 0.5, ActValue: 2, ActSpeed: 15, ActRange: 1})

	Characters[CharWoman] = NewCharacter(BehaviorMan, Stats{Name: "Woman", Wage: 10,
		Health: 12, MoveSpeed: 20, HurtTime: 0.5, ActValue: 1, ActSpeed: 18, ActRange: 1})

	Characters[CharHunter] = NewCharacter(BehaviorHunter, Stats{Name: "Hunter", Wage: 40,
		Health: 20, MoveSpeed: 30, HurtTime: 0.5, ActValue: 4, ActSpeed: 20, ActRange: 6})
	Characters[CharHunter].Sounds = CharSounds{ActionTrigger: AudioBow, HitGround: AudioProjectileGround,
		HitFlesh: AudioProjectileFlesh, HitWood: AudioProjectileWood, HitMetal: AudioProjectileMetal}

	for i, c := range Characters {
		var prefix = text.ToLowerCase(c.Stats.Name)
		c.Animations.Idle = atlas.Crops(prefix + "_idle")
		c.Animations.Walk = atlas.Crops(prefix + "_walk")
		c.Animations.ActionStart = atlas.Crops(prefix + "_action_start")
		c.Animations.ActionEnd = atlas.Crops(prefix + "_action_end")
		c.Animations.Hurt = atlas.Crops(prefix + "_hurt")
		c.Animations.Die = atlas.Crops(prefix + "_death")
		Characters[i] = c
	}
}
