package game

import "pure-game-kit/packages/geometry"

func CharacterDataHunter() CharacterData {
	return CharacterData{Brain: BehaviorHunter, Hitbox: geometry.NewRectangle(0, 7, 18, 35, 0),
		Stats: Stats{Name: "Hunter", Health: 20, MoveSpeed: 30, HurtTime: 0.5,
			AttackDamage: 6, AttackSpeed: 20, AttackRange: 6}}
}
func BehaviorHunter(self *Unit) {
}
