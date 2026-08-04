package game

// Defines the base stats, animations, behaviors (brain functions) etc of all characters - being a class/template.
// The Unit copies that base data in different points in time and uses/edits it to make the character alive.

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
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
	Stats           Stats
	Animations      Animations
	AnimationPrefix string
	Hitbox          geometry.Shape

	Brain func(self *Unit)
}

const CharacterMan, CharacterWoman Character = 0, 1

var Characters = [2]CharacterData{
	CharacterDataMan(), CharacterDataWoman(),
}

func InitCharacters() {
	var atlas = assets.LoadAtlas(assets.LoadImage("data/units.png"), "data/units.xml")

	for i, c := range Characters {
		c.Animations.Idle = atlas.Crops(c.AnimationPrefix + "-idle")
		c.Animations.Walk = atlas.Crops(c.AnimationPrefix + "-walk")
		c.Animations.AttackStart = atlas.Crops(c.AnimationPrefix + "-attack-start")
		c.Animations.AttackEnd = atlas.Crops(c.AnimationPrefix + "-attack-end")
		c.Animations.Hurt = atlas.Crops(c.AnimationPrefix + "-hurt")
		c.Animations.Die = atlas.Crops(c.AnimationPrefix + "-death")
		Characters[i] = c
	}
}
