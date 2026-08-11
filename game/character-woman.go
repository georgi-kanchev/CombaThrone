package game

import "pure-game-kit/packages/geometry"

func CharacterDataWoman() CharacterData {
	return CharacterData{Brain: BehaviorWoman, Hitbox: geometry.NewRoundedRectangle(0, 7, 18, 35, 0, 1),
		Stats: Stats{Name: "Woman", Health: 12, MoveSpeed: 20, HurtTime: 0.5,
			AttackDamage: 1, AttackSpeed: 18, AttackRange: 1}}
}
func BehaviorWoman(self *Unit) {
}
