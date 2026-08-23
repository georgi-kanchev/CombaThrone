package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/mouse/cursor"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/easing"
	"pure-game-kit/packages/utility/number"
)

type PickupKind uint8
type Pickup struct {
	graphics.Object

	Kind PickupKind
	Z    float32
	Anim *motion.Animation[assets.ImageId]

	lane                  Lane
	startX, startY, timer float32
	target                *Unit
}

const PickupCoin, PickupGem, PickupCrystal, PickupRelic PickupKind = 0, 1, 2, 3
const PickupRune, PickupSnowflake, PickupStar, PickupKey PickupKind = 4, 5, 6, 7

var Pickups []*Pickup = make([]*Pickup, 0, 32)

func NewPickup(x float32, kind PickupKind, lane Lane) *Pickup {
	var anim = motion.NewAnimation(6, true, Decor.Crops("pickup-"+pickupGroups[kind])...)
	var data = &Pickup{Object: graphics.NewSprite(x, 0, 1, 0), Z: laneZs[lane], Kind: kind, Anim: &anim, lane: lane}
	var collision = laneCollisions[lane][0]
	data.Update()
	data.Y = collision.Y - collision.Height/2 - data.Height/2 // taller pickups go below their shadows - pivot bottom
	return data
}

//=================================================================

func (p *Pickup) Pickup(by *Unit) {
	p.timer = pickupDuration
	p.target = by
	p.startX, p.startY = p.X, p.Y
}

func (p *Pickup) Update() {
	p.Anim.TimeScale = TimeScale
	p.timer -= DeltaTimeScaled()

	if p.timer > 0 {
		var t = number.Map(p.timer, pickupDuration, 0, 0, 1)
		p.X = number.Map(easing.BackInOut(t), 0, 1, p.startX, p.target.X)
		p.Y = number.Map(easing.CubicIn(t), 0, 1, p.startY, p.target.Y)
	}

	var frame = p.Anim.Frame()
	var crop = frame.CropArea()
	p.ImageId, p.Width, p.Height = frame, crop.Width, crop.Height

	DrawShadow(p.X, p.Z-0.1, p.Width*0.6, p.Height*0.15, 0, p.Mask)
	View.DrawObject(&p.Object)
	UI.TryShowTooltip(View, p.Object.Shape, cursor.Arrow)
}

// private ========================================================

const pickupDuration = 1.0

var pickupGroups = []string{"coin", "gem", "crystal", "relic", "rune", "snowflake", "star", "key"}
var pickupEffects = []func(pickedUpBy *Unit){
	PickupCoin: func(pickedUpBy *Unit) { Player.Coins += 10 },
}
