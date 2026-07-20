package game

import "pure-game-kit/packages/geometry"

func CharacterDataWoman() CharacterData {
	return CharacterData{AnimationPrefix: "woman", Brain: BrainWoman,
		Stats:  Stats{Name: "Woman", Health: 5, Damage: 1, Speed: 20, Duty: DutyMiddle},
		Hitbox: geometry.NewRoundedRectangle(0, 7, 18, 30, 0, 1)}
}
func BrainWoman(self *Unit) {
}
