package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/point"
)

type PickupKind uint8
type Pickup struct {
	graphics.Object

	Description string
	Lane        Lane
	Kind        PickupKind
	Z           float32

	Anim      *motion.Animation[assets.ImageId]
	HasShadow bool

	SlotUI int
	Effect func()
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
	var pickupGroups = [PickupCount]string{"coin", "gem", "crystal", "relic", "rune", "snowflake", "star", "key"}
	var anim = motion.NewAnimation(6, true, Decor.Crops("pickup-"+pickupGroups[kind])...)
	var data = &Pickup{
		Object: graphics.NewSprite(x, 0, 1, 0), Z: laneZs[lane], Kind: kind, Anim: &anim, Lane: lane, HasShadow: true}
	var collision = Collisions[lane][0]
	data.SlotUI = -1
	data.Update()
	data.Y = collision.Y - collision.Height/2 - data.Height/2 // taller pickups go below their shadows - pivot bottom

	switch kind {
	case PickupCoin:
		data.Description = "Gives you 🟨" + Tags[IconCoin] + "10 coins⬜. Not bad for a single coin, eh?"
		data.Effect = func() { Player.Coins += 10 }
	case PickupGem:
		data.Description = "🟩" + Tags[IconPlus] + "🌗🟩10 " + Tags[IconHeart] + "health⬜ to all 🟩" +
			Tags[IconUnit] + "allies⬜."
	case PickupCrystal:
		data.Description = "🟩" + Tags[IconPlus] + "🟧10 damage⬜ to all 🟩" + Tags[IconUnit] + "ally⬜ " +
			Tags[IconSword] + "Fighters and " + Tags[IconBow] + "Rangers."
	case PickupRelic:
		data.Description = "🌗🟦" + Tags[IconLoop] + "Revives⬜ all 🟥" + Tags[IconSkull] +
			"dead🟩 " + Tags[IconUnit] + "allies⬜ and gives them 🌗🟩full " + Tags[IconHeart] + "health⬜."
	case PickupRune:
		data.Description = "Prevents 🟥" + Tags[IconUnit] + "enemies⬜ from appearing for 20s."
	case PickupSnowflake:
		data.Description = "Prevents all 🟥" + Tags[IconUnit] + "enemies⬜ from moving. They can still act."
	case PickupStar:
		data.Description = "Gives you 🟩" + Tags[IconGlory] + "100 Glory⬜. So glorious!"
	case PickupKey:
		data.Description = "Unlocks a treasure chest (if owned)."
	}
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

		if p.Kind == PickupCoin {
			x, y = GameHUD.Coins.X, GameHUD.Coins.Y
			p.X, p.Y = point.MoveToPointSmooth(p.X, p.Y, x, y, 0.06)

			if p.SlotUI >= 0 && number.IsWithin(p.X, x, 3) && number.IsWithin(p.Y, y, 3) {
				p.Effect()
				GameHUD.Pickups = collection.Remove(GameHUD.Pickups, p)
			}
		}

		p.X, p.Y = point.MoveToPointSmooth(p.X, p.Y, x, y, 0.06)
	}
	if p.HasShadow {
		DrawShadow(p.X, p.Z-0.1, p.Width*0.6, p.Height*0.15, 0, p.Mask)
	}

	view.DrawObject(&p.Object)
}
