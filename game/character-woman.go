package game

import "pure-game-kit/packages/geometry"

func CharacterDataWoman() CharacterData {
	return CharacterData{AnimationPrefix: "woman", Brain: BrainWoman,
		Stats:  Stats{Name: "Woman", Health: 5, Damage: 1, MoveSpeed: 20},
		Hitbox: geometry.NewRoundedRectangle(0, 7, 18, 35, 0, 1)}
}
func BrainWoman(self *Unit) {
}
