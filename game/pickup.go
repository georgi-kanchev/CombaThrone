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

	Description string
	Lane        Lane
	Kind        PickupKind
	Z           float32

	Anim   *motion.Animation[assets.ImageId]
	Target *Unit

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
	var data = &Pickup{Object: graphics.NewSprite(x, 0, 1, 0), Z: laneZs[lane], Kind: kind, Anim: &anim, Lane: lane}
	var collision = Collisions[lane][0]
	data.SlotUI = -1
	data.Update()
	data.Y = collision.Y - collision.Height/2 - data.Height/2 // taller pickups go below their shadows - pivot bottom

	switch kind {
	case PickupCoin:
		data.Description = "Gives you 🟨" + Tags[IconCoin] + "10 coins⬜.\nNot bad for a single coin, eh?"
		data.Effect = func() { Player.Coins += 10 }
	case PickupGem:
		data.Description = "🟩" + Tags[IconHealth] + "Doubles health⬜ of all your units."
	case PickupCrystal:
		data.Description = "🟥" + Tags[IconMelee] + Tags[IconRanged] + Tags[IconTank] + Tags[IconMage] +
			Tags[IconHealer] + Tags[IconCollector] + Tags[IconSupplier] + Tags[IconTrapper] +
			"\nDoubles action points\n" +
			"⬜of all your units."
	case PickupRelic:
		data.Description = "Revives all 🟥" + Tags[IconDeath] +
			"dead⬜ units and gives them 🟩" + Tags[IconHealth] + "full health⬜."
	case PickupRune:
		data.Description = "Prevents any enemy units from appearing for 20s."
	case PickupSnowflake:
		data.Description = "Prevents all enemy units from moving. They can still act."
	case PickupStar:
		data.Description = "Gives you 🟩" + Tags[IconGlory] + "100 Glory⬜. So glorious!"
	case PickupKey:
		data.Description = "Unlocks a treasure chest for you\n(if owned)."
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
		p.X, p.Y = point.MoveToPointSmooth(p.X, p.Y, x, y, 0.06)
	} else if p.Target == nil {
		DrawShadow(p.X, p.Z-0.1, p.Width*0.6, p.Height*0.15, 0, p.Mask)
	} else {
		var hb = p.Target.Hitbox()
		p.Mask = p.Target.Mask
		p.X, p.Y = point.MoveToPointSmooth(p.X, p.Y, hb.X-hb.Width/2-p.Width/2, hb.Y, 0.4)
	}

	view.DrawObject(&p.Object)
}

func (p *Pickup) DrawTooltip(bench bool) {
	const width, height = 140.0, TileSize + 12
	var shape = GameHUD.ShapeToUI(p.Object.Shape)
	var col, noMask = palette.White, geometry.Area{}
	var icon = UserInterface.Crops("icons-pickup")[p.Kind]
	var x, y = shape.X, shape.Y - shape.Height/2 - height/2
	if bench {
		shape = p.Object.Shape
		x, y = shape.X, shape.Y+shape.Height/2+height/2
	}
	var area = geometry.NewArea(x, y, width, height).Inside(GameHUD.View.Bounds())
	x, y = area.X, area.Y

	GameHUD.Highlight(GameHUD.View, shape, palette.White)
	GameHUD.View.DrawImage(x, y, width, height, 0, PanelNinePatchId, col, noMask)
	GameHUD.View.DrawImage(x+width/2-TileSize/2-6, y, -TileSize, TileSize, 0, icon, col, noMask)

	TooltipLabel.Shape = geometry.NewRectangle(x-TileSize/2, y, width-TileSize-12, height-12, 0)
	TooltipLabel.Text = p.Description
	TooltipLabel.Effects.TextAlignX, TooltipLabel.Effects.TextAlignY = 0.5, 0.5
	GameHUD.View.DrawObject(TooltipLabel)
}
