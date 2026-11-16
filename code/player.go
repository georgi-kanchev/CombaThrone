package game

import "game/code/constants/progress"

type PlayerState struct {
	Progress     progress.Enum
	Glory, Coins int
	Pickups      []*Pickup
	Units        []*Unit
}

func NewPlayer() *PlayerState {
	return &PlayerState{Units: make([]*Unit, 8)}
}

//=================================================================

var playerLastCoins, playerLastGlory int

func (p *PlayerState) CoinsJustChanged() bool {
	return playerLastCoins != p.Coins
}
func (p *PlayerState) GloryJustChanged() bool {
	return playerLastGlory != p.Glory
}

func (p *PlayerState) Update() {
	playerLastCoins = p.Coins
}
