package game

import "pure-game-kit/packages/geometry"

func CharacterDataMan() CharacterData {
	return CharacterData{AnimationPrefix: "man", Brain: BehaviorMan,
		Stats: Stats{Name: "Man", Health: 20, MoveSpeed: 31, HurtTime: 0.5,
			AttackDamage: 2, AttackSpeed: 15, AttackRange: 1},
		Hitbox: geometry.NewRoundedRectangle(0, 7, 18, 35, 0, 1)}
}
func BehaviorMan(self *Unit) {
}
