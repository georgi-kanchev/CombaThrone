package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/motion"
)

type PickupKind uint8
type PickupData struct {
	graphics.Object

	Kind   PickupKind
	Z      float32
	Anim   *motion.Animation[assets.ImageId]
	Effect func(pickedUpBy *Unit)
}

const PickupCoin, PickupGem, PickupCrystal, PickupRelic PickupKind = 0, 1, 2, 3
const PickupRune, PickupSnowflake, PickupStar, PickupKey PickupKind = 4, 5, 6, 7

var Pickups []*PickupData = make([]*PickupData, 0, 32)

func NewPickup(x, y, z float32, kind PickupKind) *PickupData {
	var anim = motion.NewAnimation(6, true, Decor.Crops("pickup-"+pickupGroups[kind])...)
	var data = &PickupData{Object: graphics.NewSprite(x, y, 1, 0), Z: z, Kind: kind, Anim: &anim}
	data.Update()
	data.Y -= data.Height / 2 // taller pickups go below their shadows - accounting for that by pivoting bottom
	return data
}

//=================================================================

func (p *PickupData) Pickup(by *Unit) {
	p.Effect(by)

}

func (p *PickupData) Update() {
	var frame = p.Anim.Frame()
	var crop = frame.CropArea()
	p.ImageId, p.Width, p.Height = frame, crop.Width, crop.Height

	DrawShadow(p.X, p.Z, p.Width*0.6, p.Height*0.15, 0, p.Mask)
	View.DrawObject(&p.Object)
}

// private ========================================================

var pickupGroups = []string{"coin", "gem", "crystal", "relic", "rune", "snowflake", "star", "key"}
var pickupEffects = []func(pickedUpBy *Unit){
	PickupCoin: func(pickedUpBy *Unit) { AllyBase.Gold += 10 },
}
