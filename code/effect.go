package game

import (
	"game/code/constants/effect"
	"game/code/constants/icon"
	"pure-game-kit/packages/utility/number"
)

var Effects = [effect.COUNT]Values{
	effect.Dead: Values{DisablesTimer: number.Infinity(), EffectTimer: number.Infinity(), EffectIconText: "🟥" + Tags[icon.Skull],
		EffectInfo: "🌓⬜" + Tags[icon.Skull] + "dead⬜"},
	effect.MoreRangeOnGarrison: Values{ActRange: 2, EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Range],
		EffectInfo: "🟩" + Tags[icon.Plus] + "🌗🟧2 " + Tags[icon.Range] + "range ⬜(Garrison)"},
	effect.Bowyer: Values{ActRange: 2, EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Range],
		EffectInfo: "🟩" + Tags[icon.Plus] + "🌗🟧2 " + Tags[icon.Range] + "range ⬜(Self)"},
	effect.Miner: Values{ActPoints: 8, EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Fighter],
		EffectInfo: "🟩" + Tags[icon.Plus] + "🟧8 " + Tags[icon.Fighter] + "damage ⬜(Self)"},
	effect.Kid: Values{MoveSpeed: 20, EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Speed],
		EffectInfo: "🟩" + Tags[icon.Plus] + "🌗🟨20 " + Tags[icon.Speed] + "speed ⬜(Self)"},
	effect.Horse: Values{MoveSpeed: 20, EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Speed],
		EffectInfo: "🟩" + Tags[icon.Plus] + "🌗🟨20 " + Tags[icon.Speed] + "speed ⬜(Horse)"},
	effect.Goblin: Values{ActPoints: 2, EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Fighter],
		EffectInfo: "🟩" + Tags[icon.Plus] + "🟧2 " + Tags[icon.Fighter] + "damage ⬜(Self)"},
	effect.Stabber: Values{ActPoints: 2, EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Fighter],
		EffectInfo: "this description is dynamic per unit, see Behavior"},
	effect.Gunner1: Values{ActPoints: 4, EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Backliner],
		EffectInfo: "🟩" + Tags[icon.Plus] + "🟧4 " + Tags[icon.Backliner] + "damage ⬜(Self)"},
	effect.Gunner2: Values{ActTimer: 3, EffectTimer: number.Infinity(), EffectIconText: "🌗🟥" + Tags[icon.Hourglass],
		EffectInfo: "🟥" + Tags[icon.Plus] + "🌗🟪3s " + Tags[icon.Hourglass] + "rest ⬜(Self)"},
	effect.Vulture: Values{EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Revive],
		EffectInfo: "🌗🟦" + Tags[icon.Revive] + "faster revive ⬜(Vulture)"},
	effect.RavenFighter: Values{ActPoints: -10, EffectTimer: number.Infinity(), EffectIconText: "🟥" + Tags[icon.Fighter],
		EffectInfo: "🟥" + Tags[icon.Minus] + "🟧10 " + Tags[icon.Fighter] + "damage ⬜(Raven)"},
	effect.RavenBackliner: Values{ActPoints: -10, EffectTimer: number.Infinity(), EffectIconText: "🟥" + Tags[icon.Backliner],
		EffectInfo: "🟥" + Tags[icon.Minus] + "🟧10 " + Tags[icon.Backliner] + "damage ⬜(Raven)"},
	effect.Snake: Values{MoveSpeed: -10, EffectTimer: number.Infinity(), EffectIconText: "🟥" + Tags[icon.Speed],
		EffectInfo: "🟥" + Tags[icon.Minus] + "🌗🟨10 " + Tags[icon.Speed] + "speed ⬜(Snake)"},
	effect.Aquaweaver: Values{ActPoints: 1000, EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Defender],
		EffectInfo: "🟩" + Tags[icon.Plus] + "🟧1000 " + Tags[icon.Defender] + "defence ⬜(Self)"},
	effect.Bear: Values{DisablesTimer: 1.0, EffectTimer: 0.95, EffectIconText: "🟥" + Tags[icon.Fear],
		EffectInfo: "🌓⬜" + Tags[icon.Fear] + "in fear ⬜(Bear)"},
	effect.Boar: Values{DisablesTimer: 4.0, EffectTimer: 3.95, EffectIconText: "🟥" + Tags[icon.Stun],
		EffectInfo: "🌓⬜" + Tags[icon.Stun] + "stunned ⬜(Boar)"},
	effect.Wolf: Values{EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Hourglass],
		EffectInfo: "🌗🟪" + Tags[icon.Hourglass] + "faster rest ⬜(Wolf)"},
}
