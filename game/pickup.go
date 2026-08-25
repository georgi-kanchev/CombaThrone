package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/mouse/cursor"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/point"
)

type PickupKind uint8
type Pickup struct {
	graphics.Object

	Kind   PickupKind
	Z      float32
	Anim   *motion.Animation[assets.ImageId]
	Target *Unit

	SlotUI int

	lane Lane
}

const PickupCoin, PickupGem, PickupCrystal, PickupRelic PickupKind = 0, 1, 2, 3
const PickupRune, PickupSnowflake, PickupStar, PickupKey PickupKind = 4, 5, 6, 7

var Pickups []*Pickup = make([]*Pickup, 0, 32)

func NewPickup(x float32, kind PickupKind, lane Lane) *Pickup {
	var anim = motion.NewAnimation(6, true, Decor.Crops("pickup-"+pickupGroups[kind])...)
	var data = &Pickup{Object: graphics.NewSprite(x, 0, 1, 0), Z: laneZs[lane], Kind: kind, Anim: &anim, lane: lane}
	var collision = Collisions[lane][0]
	data.SlotUI = -1
	data.Update()
	data.Y = collision.Y - collision.Height/2 - data.Height/2 // taller pickups go below their shadows - pivot bottom
	return data
}

//=================================================================

func (p *Pickup) Update() {
	if p == nil {
		return
	}

	var view = View
	var cur = cursor.Arrow
	if p.SlotUI >= 0 {
		view = GameHUD.View
		cur = cursor.Hand
	}

	p.Anim.TimeScale = TimeScale

	var frame = p.Anim.Frame()
	var crop = frame.CropArea()
	p.ImageId, p.Width, p.Height = frame, crop.Width, crop.Height

	if p.SlotUI >= 0 {
		var x, y = GameHUD.PickupSlotPosition(p.SlotUI)
		p.X, p.Y = point.MoveToPointSmooth(p.X, p.Y, x, y, 0.06)
	} else if p.Target == nil {
		DrawShadow(p.X, p.Z-0.1, p.Width*0.6, p.Height*0.15, 0, p.Mask)
	} else {
		var hb = p.Target.Hitbox()
		p.Mask = p.Target.Mask
		p.X, p.Y = point.MoveToPointSmooth(p.X, p.Y, hb.X, hb.Y-hb.Height/2-p.Height/2-6, 0.4)
	}

	view.DrawObject(&p.Object)
	GameHUD.TryShowTooltip(view, p.Object.Shape, cur)
}

// private ========================================================

var pickupGroups = []string{"coin", "gem", "crystal", "relic", "rune", "snowflake", "star", "key"}
var pickupEffects = []func(pickedUpBy *Unit){
	PickupCoin: func(pickedUpBy *Unit) { Player.Coins += 10 },
}
