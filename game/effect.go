package game

import "pure-game-kit/packages/utility/number"

type Effect uint8

const (
	EffectMoreRangeOnGarrison Effect = iota
	EffectMoreRangeOnGround
	EffectMoreDmgVsEntrances
	EffectCount
)

var Effects = [EffectCount]Values{
	EffectMoreRangeOnGarrison: Values{ActRange: 2, EffectTimer: number.Infinity(),
		EffectInfo: "🟩" + Tags[IconPlus] + "🌗🟧" + Tags[IconRange] + "2 range ⬜(Garrison)"},
	EffectMoreRangeOnGround: Values{ActRange: 2, EffectTimer: number.Infinity(),
		EffectInfo: "🟩" + Tags[IconPlus] + "🌗🟧" + Tags[IconRange] + "2 range ⬜(Self)"},
	EffectMoreDmgVsEntrances: Values{ActPoints: 8, EffectTimer: number.Infinity(),
		EffectInfo: "🟩" + Tags[IconPlus] + "🟧" + Tags[IconSword] + "8 damage ⬜(Self)"},
}
