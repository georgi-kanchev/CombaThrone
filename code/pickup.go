package game

import (
	"game/code/constants/icon"
	"game/code/constants/lane"
	"game/code/constants/pickup"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/point"
)

type Pickup struct {
	graphics.Object

	Description string
	Lane        lane.Enum
	Kind        pickup.Enum
	Z           float32

	Anim      *motion.Animation[assets.ImageId]
	HasShadow bool

	SlotUI int
	Effect func()
}

var Pickups []*Pickup = make([]*Pickup, 0, 32)

func NewPickup(x float32, kind pickup.Enum, atLane lane.Enum) *Pickup {
	var pickupGroups = [pickup.Count]string{"coin", "gem", "crystal", "relic", "rune", "snowflake", "star", "key"}
	var anim = motion.NewAnimation(6, true, Decor.Crops("pickup-"+pickupGroups[kind])...)
	var data = &Pickup{
		Object: graphics.NewSprite(x, 0, 1, 0), Z: laneZs[atLane], Kind: kind, Anim: &anim, Lane: atLane, HasShadow: true}
	var collision = Collisions[atLane][0]
	data.SlotUI = -1
	data.Update()
	data.Y = collision.Y - collision.Height/2 - data.Height/2 // taller pickups go below their shadows - pivot bottom

	switch kind {
	case pickup.Coin:
		data.Description = "Gives you 🟨" + Tags[icon.Coin] + "10 coins⬜. Not bad for a single coin, eh?"
		data.Effect = func() { Player.Coins += 10 }
	case pickup.Gem:
		data.Description = "🟩" + Tags[icon.Plus] + "🌗🟩10 " + Tags[icon.Heart] + "health⬜ to all 🟩" +
			Tags[icon.Unit] + "allies⬜."
	case pickup.Crystal:
		data.Description = "🟩" + Tags[icon.Plus] + "🟧10 damage⬜ to all 🟩" + Tags[icon.Unit] + "ally⬜ " +
			Tags[icon.Sword] + "Fighters and " + Tags[icon.Bow] + "Rangers."
	case pickup.Relic:
		data.Description = "🌗🟦" + Tags[icon.Loop] + "Revives⬜ all 🟥" + Tags[icon.Skull] +
			"dead🟩 " + Tags[icon.Unit] + "allies⬜ and gives them 🌗🟩full " + Tags[icon.Heart] + "health⬜."
	case pickup.Rune:
		data.Description = "Prevents 🟥" + Tags[icon.Unit] + "enemies⬜ from appearing for 20s."
	case pickup.Snowflake:
		data.Description = "Prevents all 🟥" + Tags[icon.Unit] + "enemies⬜ from moving. They can still act."
	case pickup.Star:
		data.Description = "Gives you 🟩" + Tags[icon.Glory] + "100 Glory⬜. So glorious!"
	case pickup.Key:
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

		if p.Kind == pickup.Coin {
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
