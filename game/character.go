package game

// Defines the base stats, animations, behaviors (brain functions) etc of all characters - being a class/template.
// The Unit copies that base data in different points in time and uses/edits it to make the character alive.

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/utility/text"
)

type Animations struct {
	Idle, Walk, AttackStart, AttackEnd, Hurt, Die []assets.ImageId
}

type Stats struct {
	Name string

	Health, MoveSpeed int

	AttackDamage, AttackSpeed, AttackRange int

	HurtTime float32
}

type CharacterData struct {
	Stats      Stats
	Animations Animations
	Hitbox     geometry.Shape

	Brain func(self *Unit)
}

const CharacterMan, CharacterWoman, CharacterHunter Character = 0, 1, 2

var Characters = [3]CharacterData{
	CharacterDataMan(), CharacterDataWoman(), CharacterDataHunter(),
}

func InitCharacters() {
	var atlas = assets.LoadAtlas(assets.LoadImage("data/units.png"), "data/units.xml")

	for i, c := range Characters {
		var prefix = text.ToLowerCase(c.Stats.Name)
		c.Animations.Idle = atlas.Crops(prefix + "_idle")
		c.Animations.Walk = atlas.Crops(prefix + "_walk")
		c.Animations.AttackStart = atlas.Crops(prefix + "_attack_start")
		c.Animations.AttackEnd = atlas.Crops(prefix + "_attack_end")
		c.Animations.Hurt = atlas.Crops(prefix + "_hurt")
		c.Animations.Die = atlas.Crops(prefix + "_death")
		Characters[i] = c
	}
}
