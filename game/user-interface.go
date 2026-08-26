package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/button"
	"pure-game-kit/packages/input/mouse/cursor"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/easing"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/point"
	"pure-game-kit/packages/utility/text"
)

type Icon uint8
type HUD struct {
	View                             *graphics.View
	Top, ZoneInfo, Coins, UnitsPanel *graphics.Object
	TeamGlory                        [TeamCount]*graphics.Object

	Tooltip        *Tooltip
	HoverHighlight *graphics.Object

	Pickups []*Pickup

	SummonIndex              int
	SummonDragX, SummonDragY float32
}

type Tooltip struct {
	shape geometry.Shape
	draw  func(shape geometry.Shape)
}

const (
	IconHealth Icon = iota
	IconCoin
	IconGlory
	IconTimer
	IconTank
	IconRanged
	IconMelee
	IconMage
	IconRange
	IconDeath
	IconLocked
	IconUnlocked
	IconStory
	IconHealer
	IconCollector
	IconSupplier
	IconTrapper
	IconMove
	IconHome
	IconLoop
	IconCount
)

var Tags = []string{
	IconHealth: "~", IconCoin: "$", IconGlory: "*", IconDeath: "`", IconStory: "@", IconMove: ">", IconRange: "#",
	IconTimer: "^", IconMelee: "&", IconRanged: "<", IconHealer: "{", IconHome: "}", IconLoop: "[",
}

var ThemeUI assets.GUIThemeId
var UserInterface assets.AtlasId
var SlotId, ButtonUpId, ButtonDownId assets.ImageId
var PanelNinePatchId assets.ImageId
var GameHUD *HUD
var TooltipLabel *graphics.Object

func NewHUD() *HUD {
	UserInterface = assets.LoadAtlas(assets.LoadImage("data/user-interface.png"), "data/user-interface.xml")

	for i, t := range Tags {
		if t != "" {
			assets.FontId(0).EmbedImage(text.At(t, 0), UserInterface.Crops("text-icons")[i])
		}
	}

	ThemeUI = assets.LoadGUITheme("data/user-interface-theme.xml")
	SlotId = UserInterface.Crops("slot")[0]

	TooltipLabel = new(graphics.NewTextbox(0, 0, 100, 100, 0))
	TooltipLabel.Effects.FillColor, TooltipLabel.Effects.TextLineHeight = 0, 10
	TooltipLabel.Effects.TextLineGap = -40
	TooltipLabel.Effects.TextAlignX, TooltipLabel.Effects.TextAlignY = 0.5, 0.5

	var btn = UserInterface.Crops("button")
	ButtonUpId = assets.LoadImage9Patch(btn[0], 8, 8, 8, 8)
	ButtonDownId = assets.LoadImage9Patch(btn[1], 8, 8, 8, 8)
	PanelNinePatchId = assets.LoadImage9Patch(UserInterface.Crops("panel")[0], 8, 8, 8, 8)

	var view = graphics.NewView(1)
	var top = graphics.NewSprite(0, 0, 1, UserInterface.Crops("hud-top")[0])

	var highlightNinePatch = assets.LoadImage9Patch(UserInterface.Crops("highlight")[0], 4, 4, 4, 4)
	var highlight = graphics.NewSprite(0, 0, 1, highlightNinePatch)

	var unitsPanel = graphics.NewSprite(0, 0, 1, PanelNinePatchId)
	unitsPanel.Width, unitsPanel.Height = TileSize*5.5, TileSize*3

	var label = graphics.NewTextbox(0, 0, TileSize+6, TileSize/2+2, 0)
	label.Effects.FillColor, label.Effects.TextLineHeight = 0, 14
	label.Effects.TextAlignX, label.Effects.TextAlignY = 0.5, 0.5
	label.Effects.TextShadowColor, label.Effects.TextShadowWeight = 0, 0
	label.Effects.OutlineSize, label.Effects.OutlineColor = 0.4, palette.Black

	var coins = label
	coins.Effects.TextLineHeight, coins.Effects.Tint = 10, palette.White
	coins.Text = "0" + Tags[IconCoin]

	var info = label
	info.Width, info.Effects.TextLineHeight = TileSize*8.5, 8
	info.Effects.TextSymbolGap, info.Effects.TextColor = 10, palette.Gold

	var ally, enemy = label, label
	var glory = [TeamCount]*graphics.Object{&ally, &enemy}

	return &HUD{
		View: &view, ZoneInfo: &info, Top: &top, TeamGlory: glory, Coins: &coins, UnitsPanel: &unitsPanel,
		HoverHighlight: &highlight, SummonIndex: -1, Pickups: make([]*Pickup, 4),
	}
}

//=================================================================

func (h *HUD) UnitIconPosition(index int) (x, y float32) {
	var width, height = float32(TileSize) + 7, float32(TileSize) + 7
	var offX, offY float32 = width/2 - width*2, height/2 - height
	var i, j = index / 4, index % 4
	return h.UnitsPanel.X + float32(j)*width + offX, h.UnitsPanel.Y + float32(i)*height + offY
}
func (h *HUD) FreePickupSlot() int {
	for i, p := range h.Pickups {
		if p == nil {
			return i
		}
	}
	return -1
}
func (h *HUD) PickupSlotPosition(slot int) (x, y float32) {
	switch slot {
	case 0:
		return h.Coins.X - 45, h.Coins.Y + 5
	case 1:
		return h.Coins.X - 28, h.Coins.Y + 3
	case 2:
		return h.Coins.X + 28, h.Coins.Y + 3
	case 3:
		return h.Coins.X + 45, h.Coins.Y + 5
	}
	return number.NaN(), number.NaN()
}

func (h *HUD) UpdateBack() {
	const scale = 0.9
	var sz float32 = TileSize
	h.View.Zoom = scale * 5

	var tx, ty = h.View.PointFromEdge(0.5, 0)
	h.UnitsPanel.Y = ty + h.Top.Height + h.UnitsPanel.Height/2 - sz - 2
	h.View.DrawObject(h.UnitsPanel)

	var lastSummonIndex = h.SummonIndex
	if mouse.IsAnyButtonJustPressed() {
		h.SummonIndex = -1
	}

	h.drawUnitsBench(lastSummonIndex)

	h.Top.X, h.Top.Y = tx, ty+h.Top.Height/2
	h.View.DrawObject(h.Top)

	h.ZoneInfo.Text = zoneInfos[CurrentZone.kind]
	h.ZoneInfo.X, h.ZoneInfo.Y = h.Top.X, h.Top.Y-sz+4.5
	h.View.DrawObject(h.ZoneInfo)

	h.Coins.X, h.Coins.Y = h.Top.X, h.Top.Y-9
	h.View.DrawObject(h.Coins)
	h.TryShowTooltip(h.View, h.Coins.Shape, false, func(shape geometry.Shape) {
		const width, height = 80.0, 25.0
		var col, noMask = palette.White, geometry.Area{}
		h.View.DrawImage(shape.X, shape.Y+shape.Height/2+height/2, width, height, 0, PanelNinePatchId, col, noMask)

		var x, y = shape.X, shape.Y + shape.Height/2 + height/2
		TooltipLabel.Shape = geometry.NewRectangle(x, y, width-16, height, 0)
		TooltipLabel.Effects.TextAlignX, TooltipLabel.Effects.TextAlignY = 0.5, 0.5
		TooltipLabel.Text = text.New("Your 🟨", Tags[IconCoin], "Coins⬜.")
		GameHUD.View.DrawObject(TooltipLabel)
	})

	h.TeamGlory[TeamAlly].X, h.TeamGlory[TeamAlly].Y = h.Top.X-sz*3.5, h.Top.Y-6
	h.TeamGlory[TeamEnemy].X, h.TeamGlory[TeamEnemy].Y = h.Top.X+sz*3.5, h.Top.Y-6
	h.View.DrawObject(h.TeamGlory[TeamAlly])
	h.View.DrawObject(h.TeamGlory[TeamEnemy])
	h.TryShowTooltip(h.View, h.TeamGlory[TeamAlly].Shape, false, func(shape geometry.Shape) {
		const width, height = 80.0, 25.0
		var col, noMask = palette.White, geometry.Area{}
		h.View.DrawImage(shape.X, shape.Y+shape.Height/2+height/2, width, height, 0, PanelNinePatchId, col, noMask)

		var x, y = shape.X, shape.Y + shape.Height/2 + height/2
		TooltipLabel.Shape = geometry.NewRectangle(x, y, width-16, height, 0)
		TooltipLabel.Effects.TextAlignX, TooltipLabel.Effects.TextAlignY = 0.5, 0.5
		TooltipLabel.Text = text.New("Your 🟩", Tags[IconGlory], "Glory⬜.")
		GameHUD.View.DrawObject(TooltipLabel)
	})
	h.TryShowTooltip(h.View, h.TeamGlory[TeamEnemy].Shape, false, func(shape geometry.Shape) {
		const width, height = 130.0, 25.0
		var col, noMask = palette.White, geometry.Area{}
		h.View.DrawImage(shape.X, shape.Y+shape.Height/2+height/2, width, height, 0, PanelNinePatchId, col, noMask)

		var x, y = shape.X, shape.Y + shape.Height/2 + height/2
		TooltipLabel.Shape = geometry.NewRectangle(x, y, width-16, height, 0)
		TooltipLabel.Effects.TextAlignX, TooltipLabel.Effects.TextAlignY = 0.5, 0.5
		TooltipLabel.Text = text.New("The 🟥", Tags[IconGlory], "Glory⬜ of the Enemy.")
		GameHUD.View.DrawObject(TooltipLabel)
	})

	iterateRemovable(&h.Pickups, func(p *Pickup) { p.Update() })

	h.trySummon(lastSummonIndex)
}
func (h *HUD) UpdateFront() {
	if h.SummonIndex < 0 && h.Tooltip != nil {
		h.Tooltip.draw(h.Tooltip.shape)
	}
}

func (h *HUD) TryShowTooltip(view *graphics.View, shape geometry.Shape, force bool, draw func(shape geometry.Shape)) {
	if PinnedUnit != nil && !force {
		return
	}

	if force || shape.ContainsPoint(view.MousePosition()) {
		var x, y = view.PointToView(GameHUD.View, shape.X, shape.Y)
		var width, height = shape.Width, shape.Height
		if view != h.View {
			width *= view.Zoom / GameHUD.View.Zoom
			height *= view.Zoom / GameHUD.View.Zoom
		}

		h.Tooltip = &Tooltip{shape: geometry.NewRectangle(x, y, width, height, 0), draw: draw}
	}
}
func (h *HUD) Highlight(view *graphics.View, shape geometry.Shape, color uint) {
	shape.Width += 2
	shape.Height += 2
	h.HoverHighlight.Shape = shape
	h.HoverHighlight.Effects.Tint = color
	view.DrawObject(h.HoverHighlight)
}

func (h *HUD) DrawPath(view *graphics.View, fromX, fromY, toX, toY float32) {
	const size = 4.0
	var count float32 = point.DistanceToPoint(fromX, fromY, toX, toY) / 15
	for i := 1; i < int(count+2); i++ {
		var x1 = number.Map(easing.Linear(float32(i-1)/count), 0, 1, fromX, toX)
		var y1 = number.Map(easing.QuadIn(float32(i-1)/count), 0, 1, fromY, toY)
		var x2 = number.Map(easing.Linear(float32(i)/count), 0, 1, fromX, toX)
		var y2 = number.Map(easing.QuadIn(float32(i)/count), 0, 1, fromY, toY)
		if i == int(count+1) {
			x2, y2 = toX, toY
		}
		var line = geometry.NewLine(x1, y1, x2, y2, 1)
		var width = max(size, line.Width-size)
		view.DrawShape(line.X, line.Y, width, size, line.Angle, 1, palette.Black, geometry.Area{})
		view.DrawShape(line.X, line.Y, width-2, size-2, line.Angle, 1, palette.White, geometry.Area{})
	}
}

// private ========================================================

func (h *HUD) drawUnitsBench(lastSummonIndex int) {
	var noMask geometry.Area
	var drop = h.SummonIndex >= 0 && mouse.IsButtonJustReleased(button.Left)
	var click = mouse.IsButtonJustPressed(button.Left)
	var sz float32 = TileSize

	for index, unit := range Player.Units {
		var x, y = h.UnitIconPosition(index)
		var unitShape = geometry.NewRoundedRectangle(x, y, sz, sz, 0, 0)
		var hovered = unitShape.ContainsPoint(h.View.MousePosition())
		h.View.DrawImage(x, y, sz, sz, 0, SlotId, palette.White, noMask)

		if hovered && lastSummonIndex >= 0 && index != lastSummonIndex {
			h.Highlight(h.View, unitShape, palette.White)

			if click || drop {
				collection.Swap(Player.Units, lastSummonIndex, index)
				h.SummonIndex = -1
			}
		}
		if unit == nil {
			continue
		}

		var iSz = sz / 2.5
		var icons = UserInterface.Crops("text-icons")
		var tint = palette.White
		if index == h.SummonIndex {
			tint = color.RGBA(255, 255, 255, 127)
		}
		h.View.DrawImage(x, y, sz, sz, 0, Characters[unit.Character].Icon, tint, noMask)
		if unit.State == StateDecaying {
			var timerWidth = number.Map(unit.hurtTimer, 0, -float32(unit.Stats.RespawnTimer)/10, sz-4, 0)
			var icon, col = IconDeath, palette.Red
			if unit.Stats.Health > 0 { // got into the enemy base
				icon, col = IconGlory, teamColors[TeamAlly]
			}

			h.View.DrawShape(x, y, sz, sz, 0, 0, color.RGBA(0, 0, 0, 150), noMask)
			h.View.DrawImage(x-sz/2+iSz/2, y+sz/2-iSz/2-3, iSz, iSz, 0, icons[icon], col, noMask)
			h.View.DrawShape(x+sz/2-sz/2, y+sz/2-2, sz-2, 3, 0, 0, palette.Black, noMask)
			h.View.DrawShape(x+timerWidth/2-sz/2+2, y+sz/2-2, timerWidth, 1, 0, 0, palette.White, noMask)
		} else if unit.IsSummoned() {
			var hp, maxHp = float32(unit.Stats.Health), float32(Characters[unit.Character].Stats.Health)
			var hpWidth = number.Map(hp, 0, maxHp, 0, sz-4)
			h.View.DrawShape(x, y, sz, sz, 0, 0, color.RGBA(0, 0, 0, 150), noMask)
			h.View.DrawImage(x-sz/2+iSz/2, y+sz/2-iSz/2-2, iSz, iSz, 0, icons[IconHealth], teamColors[TeamAlly], noMask)
			h.View.DrawShape(x+sz/2-sz/2, y+sz/2-2, sz-2, 3, 0, 0, palette.Black, noMask)
			h.View.DrawShape(x+hpWidth/2-sz/2+2, y+sz/2-2, hpWidth, 1, 0, 0, teamColors[TeamAlly], noMask)
		} else if hovered && click && lastSummonIndex < 0 {
			h.SummonIndex = index
			h.SummonDragX, h.SummonDragY = h.View.PointToView(View, x, y)
		}

		h.TryShowTooltip(h.View, unitShape, false, func(shape geometry.Shape) {
			const width, height = 200.0, 100.0
			// var col, noMask = palette.White, geometry.Area{}

			mouse.SetCursor(cursor.Hand)
			if unit.State == StateDecaying || unit.IsSummoned() {
				mouse.SetCursor(cursor.NotAllowed)
			}

			unit.drawTooltipInfo(shape, true)

			// GameHUD.Highlight(GameHUD.View, shape, palette.White)
			// h.View.DrawImage(shape.X, shape.Y+shape.Height/2+height/2+4, width, height, 0, PanelNinePatchId, col, noMask)
		})

		if h.SummonIndex == index {
			h.Highlight(h.View, unitShape, palette.White)
		}

		if h.SummonIndex < 0 {
			if unit.IsSummoned() && hovered {
				h.Highlight(View, unit.Shape, palette.LightGray)
			} else if unit.IsSummoned() && unit.Shape.ContainsPoint(View.MousePosition()) {
				h.Highlight(h.View, unitShape, palette.LightGray)
			}
		}
	}
}
func (h *HUD) trySummon(lastSummonIndex int) {
	if lastSummonIndex < 0 {
		return
	}

	var unit = Player.Units[lastSummonIndex]
	if unit == nil {
		return
	}

	var unitX, unitY = h.UnitIconPosition(lastSummonIndex)
	var mx, my = h.View.MousePosition()
	mx, my = h.View.PointToView(View, mx, my)
	unitX, unitY = h.View.PointToView(View, unitX, unitY)

	mouse.SetCursor(cursor.Move)

	var targetX, targetY = mx, my
	var hoverShape geometry.Shape
	var drop, click = mouse.IsButtonJustReleased(button.Left), mouse.IsButtonJustPressed(button.Left)
	var smallestDist = number.Infinity()
	const size = TileSize * 0.85
	var lane Lane

	for _, e := range Bases[TeamAlly].Entrances {
		var shape = e.Shape()
		if Bases[TeamAlly].Kind < BaseBarrack {
			shape.X += TileSize / 1.5
			shape.Y += TileSize
		} else if e.Kind == EntranceTallGate {
			shape.Y += TileSize / 2
		}
		shape.Width, shape.Height = size, size

		var dist = point.DistanceToPoint(shape.X, shape.Y, mx, my)
		if dist < smallestDist && dist < TileSize*1.5 {
			smallestDist = dist
			hoverShape = shape
			lane = e.Lane
		}
		View.DrawShape(shape.X, shape.Y, shape.Width, shape.Height, 0, 0, color.RGBA(0, 255, 0, 80), geometry.Area{})
		h.Highlight(View, shape, palette.Green)
	}

	var from, to = int(LaneGarrison3 + Lane(Bases[TeamAlly].Garrison)), int(LaneGarrison1)
	var takenGarrisons = make([]Lane, 0, 8)
	for _, pu := range Player.Units {
		if pu != nil && pu.IsSummoned() && pu.IsGarrisoner() && pu.Stats.Health > 0 {
			takenGarrisons = append(takenGarrisons, pu.Lane)
		}
	}
	if unit.Stats.ActRange > 1 {
		var garrisonsFull = len(takenGarrisons) == 6
		var garrisonX, garrisonY = PointAtCell(0, 4)
		var shape = geometry.NewRectangle(garrisonX, garrisonY, size, size, 0)

		var dist = point.DistanceToPoint(shape.X, shape.Y, mx, my)
		if !garrisonsFull && dist < smallestDist && dist < TileSize*1.5 {
			smallestDist = dist
			hoverShape = shape
			lane = LaneGarrison1

			for i := from; i >= to; i-- {
				if !collection.Contains(takenGarrisons, Lane(i)) {
					lane = Lane(i)
					break
				}
			}
		}

		var garrColBg = color.RGBA(0, 255, 0, 80)
		var garrCol = palette.Green
		if garrisonsFull {
			garrColBg = color.RGBA(255, 0, 0, 80)
			garrCol = palette.Red
		}
		View.DrawShape(shape.X, shape.Y, shape.Width, shape.Height, 0, 0, garrColBg, geometry.Area{})
		h.Highlight(View, shape, garrCol)
	}

	if hoverShape != (geometry.Shape{}) {
		h.Highlight(View, hoverShape, palette.White)
		targetX, targetY = hoverShape.X, hoverShape.Y

		if click || drop {
			h.SummonIndex = -1
			Units = append(Units, unit)
			unit.Lane = lane

			if unit.IsOffLaner() {
				unit.Lane++
			}
			unit.PrepareSpawn()
		}
	}

	h.SummonDragX, h.SummonDragY = point.MoveToPointSmooth(h.SummonDragX, h.SummonDragY, targetX, targetY, 0.4)

	h.DrawPath(View, unitX, unitY, h.SummonDragX, h.SummonDragY)
	var icon = Characters[unit.Character].Icon
	View.DrawImage(h.SummonDragX, h.SummonDragY, size, size, 0, icon, palette.White, geometry.Area{})
}
