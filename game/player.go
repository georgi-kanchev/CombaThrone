package game

type PlayerState struct {
	Coins   int
	Pickups []*Pickup
	Units   []*Unit
}

func NewPlayer() *PlayerState {
	return &PlayerState{Units: make([]*Unit, 8)}
}

//=================================================================

var playerLastCoins int

func (p *PlayerState) CoinsJustChanged() bool {
	return playerLastCoins != p.Coins
}

func (p *PlayerState) Update() {
	playerLastCoins = p.Coins
}
