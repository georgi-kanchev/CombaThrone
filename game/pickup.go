package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/color/palette"
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

const (
	PickupCoin PickupKind = iota
	PickupGem
	PickupCrystal
	PickupRelic
	PickupRune
	PickupSnowflake
	PickupStar
	PickupKey
	PickupCount
)

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
	if p.SlotUI >= 0 {
		view = GameHUD.View
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
		p.X, p.Y = point.MoveToPointSmooth(p.X, p.Y, hb.X-hb.Width/2-p.Width/2, hb.Y, 0.4)
	}

	view.DrawObject(&p.Object)

	GameHUD.TryShowTooltip(view, p.Object.Shape, false, func(shape geometry.Shape) {
		const width, height = 130.0, 35.0
		var tooltip = geometry.NewRectangle(shape.X, shape.Y-shape.Height/2-height, width, height, 0)
		var col, noMask = palette.White, geometry.Area{}
		GameHUD.Highlight(GameHUD.View, shape, palette.White)
		GameHUD.View.DrawImage(tooltip.X, tooltip.Y, width, height, 0, PanelNinePatchId, col, noMask)

		tooltip.Width -= 12
		tooltip.Height -= 12
		TooltipLabel.Shape = tooltip
		TooltipLabel.Text = pickupDescription[p.Kind]
		TooltipLabel.Effects.TextAlignX, TooltipLabel.Effects.TextAlignY = 0.5, 0.5
		GameHUD.View.DrawObject(TooltipLabel)
	})
}

// private ========================================================

var pickupGroups = [PickupCount]string{"coin", "gem", "crystal", "relic", "rune", "snowflake", "star", "key"}
var pickupDescription = [PickupCount]string{
	PickupCoin:      "A coin with a value of 10.",
	PickupGem:       "gem",
	PickupCrystal:   "crystal",
	PickupRelic:     "Revives all of your 🟥" + Tags[IconDeath] + "dead⬜ units with 🟩" + Tags[IconHealth] + "50% health⬜.",
	PickupRune:      "rune",
	PickupSnowflake: "snowflake",
	PickupStar:      "star",
	PickupKey:       "key",
}
var pickupEffects = [PickupCount]func(pickedUpBy *Unit){
	PickupCoin: func(pickedUpBy *Unit) { Player.Coins += 10 },
}
