package game

import "pure-game-kit/packages/geometry"

func CharacterDataMan() CharacterData {
	return CharacterData{AnimationPrefix: "man", Brain: BrainMan,
		Stats:  Stats{Name: "Man", Health: 10, Damage: 2, MoveSpeed: 31, HurtTime: 0.5},
		Hitbox: geometry.NewRoundedRectangle(0, 7, 18, 35, 0, 1)}
}
func BrainMan(self *Unit) {
}
