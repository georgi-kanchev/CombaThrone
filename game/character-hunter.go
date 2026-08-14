package game

import (
	"pure-game-kit/packages/geometry"
)

func CharacterDataHunter() CharacterData {
	var data = CharacterData{
		Behavior: BehaviorHunter, Hitbox: geometry.NewRoundedRectangle(0, 7, 18, 35, 0, 1),
		Stats: Stats{
			Name: "Hunter", Health: 20, MoveSpeed: 30, HurtTime: 0.5,
			AttackDamage: 6, AttackSpeed: 20, AttackRange: 6,
		},
	}
	data.Sounds.AttackTrigger = AudioBow
	data.Sounds.HitGround = AudioProjectileGround
	data.Sounds.HitWood = AudioProjectileWood
	data.Sounds.HitMetal = AudioProjectileMetal
	data.Sounds.HitFlesh = AudioProjectileFlesh
	return data
}
func BehaviorHunter(self *Unit) {
}
