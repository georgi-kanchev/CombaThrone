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

var Characters map[Character]CharacterData = map[Character]CharacterData{
	CharacterMan: CharacterDataMan(), CharacterWoman: CharacterDataWoman(),
}

func InitCharacters() {
	var animations = assets.LoadAnimations(assets.LoadImage("data/units.png"), "data/units.xml")

	for i, c := range Characters {
		c.Animations.Idle = animations.Frames(c.AnimationPrefix + "-idle")
		c.Animations.Walk = animations.Frames(c.AnimationPrefix + "-walk")
		c.Animations.AttackStart = animations.Frames(c.AnimationPrefix + "-attack-start")
		c.Animations.AttackEnd = animations.Frames(c.AnimationPrefix + "-attack-end")
		c.Animations.Hurt = animations.Frames(c.AnimationPrefix + "-hurt")
		c.Animations.Die = animations.Frames(c.AnimationPrefix + "-death")
		Characters[i] = c
	}
}
