package game

import (
	"game/code/constants/effect"
	"game/code/constants/icon"
	"pure-game-kit/packages/utility/number"
)

var Effects = [effect.COUNT]Values{
	effect.Dead: {DisablesTimer: number.Infinity(), EffectTimer: number.Infinity(), EffectIconText: "🟥" + Tags[icon.Skull],
		EffectInfo: "🌓⬜" + Tags[icon.Skull] + "dead⬜"},
	effect.MoreRangeOnGarrison: {ActRange: 2, EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Range],
		EffectInfo: "🟩" + Tags[icon.Plus] + "🌗🟧2 " + Tags[icon.Range] + "range ⬜(Garrison)"},
	effect.Bowyer: {ActRange: 2, EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Range],
		EffectInfo: "🟩" + Tags[icon.Plus] + "🌗🟧2 " + Tags[icon.Range] + "range ⬜(Self)"},
	effect.Miner: {ActPoints: 8, EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Fighter],
		EffectInfo: "🟩" + Tags[icon.Plus] + "🟧8 " + Tags[icon.Fighter] + "damage ⬜(Self)"},
	effect.Kid: {MoveSpeed: 20, EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Speed],
		EffectInfo: "🟩" + Tags[icon.Plus] + "🌗🟨20 " + Tags[icon.Speed] + "speed ⬜(Self)"},
	effect.Horse: {MoveSpeed: 20, EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Speed],
		EffectInfo: "🟩" + Tags[icon.Plus] + "🌗🟨20 " + Tags[icon.Speed] + "speed ⬜(Horse)"},
	effect.Goblin: {ActPoints: 2, EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Fighter],
		EffectInfo: "🟩" + Tags[icon.Plus] + "🟧2 " + Tags[icon.Fighter] + "damage ⬜(Self)"},
	effect.Stabber: {ActPoints: 2, EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Fighter],
		EffectInfo: "this description is dynamic per unit, see Behavior"},
	effect.Gunner1: {ActPoints: 4, EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Backliner],
		EffectInfo: "🟩" + Tags[icon.Plus] + "🟧4 " + Tags[icon.Backliner] + "damage ⬜(Self)"},
	effect.Gunner2: {ActTimer: 3, EffectTimer: number.Infinity(), EffectIconText: "🌗🟥" + Tags[icon.Hourglass],
		EffectInfo: "🟥" + Tags[icon.Plus] + "🌗🟪3s " + Tags[icon.Hourglass] + "rest ⬜(Self)"},
	effect.Vulture: {EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Revive],
		EffectInfo: "🌗🟦" + Tags[icon.Revive] + "faster revive ⬜(Vulture)"},
	effect.RavenFighter: {ActPoints: -10, EffectTimer: number.Infinity(), EffectIconText: "🟥" + Tags[icon.Fighter],
		EffectInfo: "🟥" + Tags[icon.Minus] + "🟧10 " + Tags[icon.Fighter] + "damage ⬜(Raven)"},
	effect.RavenBackliner: {ActPoints: -10, EffectTimer: number.Infinity(), EffectIconText: "🟥" + Tags[icon.Backliner],
		EffectInfo: "🟥" + Tags[icon.Minus] + "🟧10 " + Tags[icon.Backliner] + "damage ⬜(Raven)"},
	effect.Snake: {MoveSpeed: -10, EffectTimer: number.Infinity(), EffectIconText: "🟥" + Tags[icon.Speed],
		EffectInfo: "🟥" + Tags[icon.Minus] + "🌗🟨10 " + Tags[icon.Speed] + "speed ⬜(Snake)"},
	effect.Aquaweaver: {ActPoints: 1000, EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Defender],
		EffectInfo: "🟩" + Tags[icon.Plus] + "🟧1000 " + Tags[icon.Defender] + "block ⬜(Self)"},
	effect.Bear: {DisablesTimer: 1.0, EffectTimer: 0.95, EffectIconText: "🟥" + Tags[icon.Fear],
		EffectInfo: "🌓⬜" + Tags[icon.Fear] + "in fear ⬜(Bear)"},
	effect.Boar: {DisablesTimer: 4.0, EffectTimer: 3.95, EffectIconText: "🟥" + Tags[icon.Stun],
		EffectInfo: "🌓⬜" + Tags[icon.Stun] + "stunned ⬜(Boar)"},
	effect.Wolf: {EffectTimer: number.Infinity(), EffectIconText: "🟩" + Tags[icon.Hourglass],
		EffectInfo: "🌗🟪" + Tags[icon.Hourglass] + "faster rest ⬜(Wolf)"},
	effect.Shaman: {ActPoints: -10, EffectTimer: 5.0, EffectIconText: "🟥" + Tags[icon.Fighter], // TODO: find damage icon
		EffectInfo: "🟥" + Tags[icon.Minus] + "🟧10 " + Tags[icon.Fighter] + "damage ⬜(Shaman)"},
}
