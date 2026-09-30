package game

import "pure-game-kit/packages/utility/number"

type Effect uint8

const (
	EffectMoreRangeOnGarrison Effect = iota
	EffectBowyer
	EffectMiner
	EffectKid
	EffectHorse
	EffectGoblin
	EffectStabberBandit
	EffectGunnerBandit1
	EffectGunnerBandit2
	EffectVulture
	EffectRaven
	EffectRedSnake
	EffectCount
)

var Effects = [EffectCount]Values{
	EffectMoreRangeOnGarrison: Values{ActRange: 2, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟩" + Tags[IconPlus] + "🌗🟧2 " + Tags[IconRange] + "range ⬜(Garrison)"},
	EffectBowyer: Values{ActRange: 2, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟩" + Tags[IconPlus] + "🌗🟧2 " + Tags[IconRange] + "range ⬜(Self)"},
	EffectMiner: Values{ActPoints: 8, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟩" + Tags[IconPlus] + "🟧8 " + Tags[IconSword] + "damage ⬜(Self)"},
	EffectKid: Values{MoveSpeed: 20, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟩" + Tags[IconPlus] + "🌗🟨20 " + Tags[IconLeftRight] + "speed ⬜(Self)"},
	EffectHorse: Values{MoveSpeed: 20, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟩" + Tags[IconPlus] + "🌗🟨20 " + Tags[IconLeftRight] + "speed ⬜(Horse)"},
	EffectGoblin: Values{ActPoints: 2, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟩" + Tags[IconPlus] + "🟧2 " + Tags[IconSword] + "damage ⬜(Self)"},
	EffectStabberBandit: Values{ActPoints: 1, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "this description is dynamic per unit, see Behavior"},
	EffectGunnerBandit1: Values{ActPoints: 4, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟩" + Tags[IconPlus] + "🟧4 " + Tags[IconSword] + "damage ⬜(Self)"},
	EffectGunnerBandit2: Values{ActTimer: 3, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟥" + Tags[IconPlus] + "🌗🟪3s " + Tags[IconTimer] + "rest ⬜(Self)"},
	EffectVulture: Values{ /*just for tooltip, no stats*/ EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟩" + Tags[IconPlus] + "🌗🟦" + Tags[IconLoop] + "fast revive ⬜(Vulture)"},
	EffectRaven: Values{ActPoints: -10, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟥" + Tags[IconMinus] + "🟧10 " + Tags[IconSword] + "damage ⬜(Raven)"},
	EffectRedSnake: Values{MoveSpeed: -10, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟥" + Tags[IconMinus] + "🌗🟨10 " + Tags[IconLeftRight] + "speed ⬜(Red Snake)"},
}
