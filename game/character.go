package game

// Defines the base stats, animations, behaviors (brain functions) etc of all characters - being a class/template.
// The Unit copies that base data in different points in time and uses/edits it to make the character alive.

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/audio"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/utility/text"
)

type Stats struct {
	Name string

	Health, MoveSpeed int

	ActionValue, ActionSpeed, ActionRange int

	HurtTime float32
}

type CharacterData struct {
	Stats      Stats
	Hitbox     geometry.Shape
	Animations struct {
		Idle, Walk, ActionStart, ActionEnd, Hurt, Die []assets.ImageId
	}
	Sounds struct {
		ActionTrigger, HitFlesh, HitWood, HitMetal, HitGround []audio.Audio
	}

	Behavior func(self *Unit)
}

const CharacterMan, CharacterWoman, CharacterHunter Character = 0, 1, 2

var Characters [3]CharacterData

func InitCharacters() {
	var atlas = assets.LoadAtlas(assets.LoadImage("data/units.png"), "data/units.xml")
	Characters[CharacterMan] = CharacterDataMan()
	Characters[CharacterWoman] = CharacterDataWoman()
	Characters[CharacterHunter] = CharacterDataHunter()

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
