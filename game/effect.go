package game

import "pure-game-kit/packages/utility/number"

type Effect uint8

const (
	EffectMoreRangeOnGarrison Effect = iota
	EffectBowyer
	EffectMiner
	EffectKid
	EffectHorse
	EffectCount
)

var Effects = [EffectCount]Values{
	EffectMoreRangeOnGarrison: Values{ActRange: 2, EffectTimer: number.Infinity(),
		EffectInfo: "🟩" + Tags[IconPlus] + "🌗🟧" + Tags[IconRange] + "2 range ⬜(Garrison)"},
	EffectBowyer: Values{ActRange: 2, EffectTimer: number.Infinity(),
		EffectInfo: "🟩" + Tags[IconPlus] + "🌗🟧" + Tags[IconRange] + "2 range ⬜(Self)"},
	EffectMiner: Values{ActPoints: 8, EffectTimer: number.Infinity(),
		EffectInfo: "🟩" + Tags[IconPlus] + "🟧" + Tags[IconSword] + "8 damage ⬜(Self)"},
	EffectKid: Values{MoveSpeed: 20, EffectTimer: number.Infinity(),
		EffectInfo: "🟩" + Tags[IconPlus] + "🌗🟨" + Tags[IconLeftRight] + "20 speed ⬜(Self)"},
	EffectHorse: Values{MoveSpeed: 20, EffectTimer: number.Infinity(),
		EffectInfo: "🟩" + Tags[IconPlus] + "🌗🟨" + Tags[IconLeftRight] + "20 speed ⬜(Horse)"},
}
