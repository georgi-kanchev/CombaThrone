package game

// Defines the base stats, animations, behaviors (brain functions) etc of all characters - being a class/template.
// The Unit copies that base data in different points in time and uses/edits it to make the character alive.

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
)

type Animations struct {
	Idle, Walk, Attack, Hurt, Die []assets.ImageId
}

type Stats struct {
	Name           string
	Health, Damage int
}

type CharacterData struct {
	Stats           Stats
	Animations      Animations
	AnimationPrefix string
	Hitbox          geometry.Area

	Brain func(self *Unit)
}

const CharacterMan, CharacterWoman Character = 0, 1

var Characters map[Character]CharacterData = make(map[Character]CharacterData)

func InitCharacters() {
	var animations = assets.LoadAnimations(assets.LoadImage("data/units.png"), "data/animations.xml")

	Characters[CharacterMan] = CharacterData{AnimationPrefix: "man", Brain: BrainMan,
		Stats:  Stats{Name: "Man", Health: 10, Damage: 2},
		Hitbox: geometry.NewArea(0, 6.5, 18, 35)}
	Characters[CharacterWoman] = CharacterData{AnimationPrefix: "woman", Brain: BrainWoman,
		Stats:  Stats{Name: "Woman", Health: 5, Damage: 1},
		Hitbox: geometry.NewArea(0, 6.5, 18, 35)}

	for i, c := range Characters {
		c.Animations.Idle = animations.Frames(c.AnimationPrefix + "-idle")
		c.Animations.Walk = animations.Frames(c.AnimationPrefix + "-walk")
		c.Animations.Attack = animations.Frames(c.AnimationPrefix + "-attack")
		c.Animations.Hurt = animations.Frames(c.AnimationPrefix + "-hurt")
		c.Animations.Die = animations.Frames(c.AnimationPrefix + "-die")
		Characters[i] = c
	}
}
