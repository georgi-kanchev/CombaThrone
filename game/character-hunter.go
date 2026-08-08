package game

import "pure-game-kit/packages/geometry"

func CharacterDataHunter() CharacterData {
	return CharacterData{Brain: BehaviorHunter, Hitbox: geometry.NewRoundedRectangle(0, 7, 18, 35, 0, 1),
		Stats: Stats{Name: "Hunter", Health: 20, MoveSpeed: 30, HurtTime: 0.5,
			AttackDamage: 2, AttackSpeed: 15, AttackRange: 1}}
}
func BehaviorHunter(self *Unit) {
}
