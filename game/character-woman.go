package game

import "pure-game-kit/packages/geometry"

func CharacterDataWoman() CharacterData {
	return CharacterData{Behavior: BehaviorWoman, Hitbox: geometry.NewRoundedRectangle(0, 7, 18, 35, 0, 1),
		Stats: Stats{Name: "Woman", Health: 12, MoveSpeed: 20, HurtTime: 0.5,
			ActionValue: 1, ActionSpeed: 18, ActionRange: 1}}
}
func BehaviorWoman(self *Unit) {
}
