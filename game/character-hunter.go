package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/audio"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/utility/text"
)

func CharacterDataHunter() CharacterData {
	var data = CharacterData{
		Behavior: BehaviorHunter, Hitbox: geometry.NewRoundedRectangle(0, 7, 18, 35, 0, 1),
		Stats: Stats{
			Name: "Hunter", Health: 20, MoveSpeed: 30, HurtTime: 0.5,
			AttackDamage: 6, AttackSpeed: 20, AttackRange: 6,
		},
	}
	data.Sounds.AttackTrigger = append(data.Sounds.AttackTrigger, audio.New(assets.LoadSound("data/audio/bow.mp3", 4)))
	for i := 1; i <= 6; i++ {
		var asset = assets.LoadSound(text.New("data/audio/projectile-flesh", i, ".mp3"), 4)
		data.Sounds.HitFlesh = append(data.Sounds.HitFlesh, audio.New(asset))
	}
	data.Sounds.HitGround = append(data.Sounds.HitGround, audio.New(assets.LoadSound("data/audio/projectile-ground.mp3", 4)))
	return data
}
func BehaviorHunter(self *Unit) {
}
