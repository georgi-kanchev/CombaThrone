package game

import (
	"game/code/constants/effect"
	"game/code/constants/icon"
	"pure-game-kit/packages/utility/number"
)

var Effects = [effect.Count]Values{
	effect.MoreRangeOnGarrison: Values{ActRange: 2, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟩" + Tags[icon.Plus] + "🌗🟧2 " + Tags[icon.Range] + "range ⬜(Garrison)"},
	effect.Bowyer: Values{ActRange: 2, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟩" + Tags[icon.Plus] + "🌗🟧2 " + Tags[icon.Range] + "range ⬜(Self)"},
	effect.Miner: Values{ActPoints: 8, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟩" + Tags[icon.Plus] + "🟧8 " + Tags[icon.Sword] + "damage ⬜(Self)"},
	effect.Kid: Values{MoveSpeed: 20, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟩" + Tags[icon.Plus] + "🌗🟨20 " + Tags[icon.LeftRight] + "speed ⬜(Self)"},
	effect.Horse: Values{MoveSpeed: 20, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟩" + Tags[icon.Plus] + "🌗🟨20 " + Tags[icon.LeftRight] + "speed ⬜(Horse)"},
	effect.Goblin: Values{ActPoints: 2, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟩" + Tags[icon.Plus] + "🟧2 " + Tags[icon.Sword] + "damage ⬜(Self)"},
	effect.Stabber: Values{ActPoints: 1, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "this description is dynamic per unit, see Behavior"},
	effect.Gunner1: Values{ActPoints: 4, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟩" + Tags[icon.Plus] + "🟧4 " + Tags[icon.Sword] + "damage ⬜(Self)"},
	effect.Gunner2: Values{ActTimer: 3, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟥" + Tags[icon.Plus] + "🌗🟪3s " + Tags[icon.Timer] + "rest ⬜(Self)"},
	effect.Vulture: Values{ /*just for tooltip, no stats*/ EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟩" + Tags[icon.Plus] + "🌗🟦" + Tags[icon.Loop] + "fast revive ⬜(Vulture)"},
	effect.Raven: Values{ActPoints: -10, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟥" + Tags[icon.Minus] + "🟧10 " + Tags[icon.Sword] + "damage ⬜(Raven)"},
	effect.Snake: Values{MoveSpeed: -10, EffectTimer: number.Infinity(), EffectTickInterval: number.Infinity(),
		EffectInfo: "🟥" + Tags[icon.Minus] + "🌗🟨10 " + Tags[icon.LeftRight] + "speed ⬜(Red Snake)"},
}
