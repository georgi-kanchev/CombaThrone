package game

import "pure-game-kit/packages/geometry"

func CharacterDataMan() CharacterData {
	return CharacterData{Behavior: BehaviorMan, Hitbox: geometry.NewRoundedRectangle(0, 7, 18, 35, 0, 1),
		Stats: Stats{Name: "Man", Health: 20, MoveSpeed: 30, HurtTime: 0.5,
			ActionValue: 2, ActionSpeed: 15, ActionRange: 1}}
}
func BehaviorMan(self *Unit) {
}
