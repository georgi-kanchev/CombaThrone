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
)

type Icon uint8
type GUI struct {
	View                             *graphics.View
	Top, ZoneInfo, Coins, UnitsPanel *graphics.Object
	TeamGlory                        [TeamCount]*graphics.Object

	Tooltip        *Tooltip
	HoverHighlight *graphics.Object

	SummonIndex int
}

type Tooltip struct {
	view   *graphics.View
	shape  geometry.Shape
	cursor int
}

const IconHealth, IconCoin, IconGlory, IconDeath, IconCount Icon = 0, 1, 2, 3, 4

var Tags = [IconCount]string{IconHealth: "~", IconCoin: "$", IconGlory: "*"}

func NewUI() *GUI {
	var panelNinePatch = assets.LoadImage9Patch(UserInterface.Crops("panel")[0], 8, 8, 8, 8)
	var view = graphics.NewView(1)
	var top = graphics.NewSprite(0, 0, 1, UserInterface.Crops("hud-top")[0])

	var highlightNinePatch = assets.LoadImage9Patch(UserInterface.Crops("highlight")[0], 4, 4, 4, 4)
	var highlight = graphics.NewSprite(0, 0, 1, highlightNinePatch)

	var unitsPanel = graphics.NewSprite(0, 0, 1, panelNinePatch)
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

	return &GUI{
		View: &view, ZoneInfo: &info, Top: &top, TeamGlory: glory, Coins: &coins, UnitsPanel: &unitsPanel,
		HoverHighlight: &highlight, SummonIndex: -1,
	}
}

//=================================================================

func (u *GUI) UnitPosition(index int) (x, y float32) {
	var width, height = float32(TileSize) + 7, float32(TileSize) + 7
	var offX, offY float32 = width/2 - width*2, height/2 - height
	var i, j = index / 4, index % 4
	return u.UnitsPanel.X + float32(j)*width + offX, u.UnitsPanel.Y + float32(i)*height + offY
}

func (u *GUI) Update() {
	const scale = 1
	var noMask = geometry.Area{}
	var sz float32 = TileSize
	var click = mouse.IsButtonJustPressed(button.Left)
	var drop = u.SummonIndex >= 0 && mouse.IsButtonJustReleased(button.Left)
	u.View.Zoom = scale * 5

	var tx, ty = u.View.PointFromEdge(0.5, 0)
	u.UnitsPanel.Y = ty + u.Top.Height + u.UnitsPanel.Height/2 - sz - 2
	u.View.DrawObject(u.UnitsPanel)

	var lastSummonIndex = u.SummonIndex
	if mouse.IsAnyButtonJustPressed() {
		u.SummonIndex = -1
	}

	var slotImageId = UserInterface.Crops("slot")[0]
	for index, unit := range Player.Units {
		var x, y = u.UnitPosition(index)
		var shape = geometry.NewRoundedRectangle(x, y, sz, sz, 0, 0)
		var hovered = shape.ContainsPoint(u.View.MousePosition())
		u.View.DrawImage(x, y, sz, sz, 0, slotImageId, palette.White, noMask)

		if hovered && lastSummonIndex >= 0 && index != lastSummonIndex {
			u.Highlight(u.View, shape, palette.White)
			if click || drop {
				collection.Swap(Player.Units, lastSummonIndex, index)
				u.SummonIndex = -1
			}
		}
		if unit == nil {
			continue
		}

		var iSz = sz / 2.5
		var icons = UserInterface.Crops("icons")
		var newCursor = cursor.Hand
		if unit.State == StateDecaying || unit.IsSummoned() {
			newCursor = cursor.NotAllowed
		}
		var tint = palette.White
		if index == u.SummonIndex {
			tint = color.RGBA(255, 255, 255, 127)
		}
		u.View.DrawImage(x, y, sz, sz, 0, Characters[unit.Character].Icon, tint, noMask)
		if unit.State == StateDecaying {
			var timerWidth = number.Map(unit.hurtTimer, 0, -unit.Stats.RespawnTimer, sz-4, 0)
			var icon, col = IconDeath, palette.Red
			if unit.Stats.Health > 0 { // got into the enemy base
				icon, col = IconGlory, teamColors[TeamAlly]
			}

			u.View.DrawShape(x, y, sz, sz, 0, 0, color.RGBA(0, 0, 0, 150), noMask)
			u.View.DrawImage(x-sz/2+iSz/2, y+sz/2-iSz/2-3, iSz, iSz, 0, icons[icon], col, noMask)
			u.View.DrawShape(x+sz/2-sz/2, y+sz/2-2, sz-2, 3, 0, 0, palette.Black, noMask)
			u.View.DrawShape(x+timerWidth/2-sz/2+2, y+sz/2-2, timerWidth, 1, 0, 0, palette.White, noMask)
		} else if unit.IsSummoned() {
			var hp, maxHp = float32(unit.Stats.Health), float32(Characters[unit.Character].Stats.Health)
			var hpWidth = number.Map(hp, 0, maxHp, 0, sz-4)
			u.View.DrawShape(x, y, sz, sz, 0, 0, color.RGBA(0, 0, 0, 150), noMask)
			u.View.DrawImage(x-sz/2+iSz/2, y+sz/2-iSz/2-2, iSz, iSz, 0, icons[IconHealth], teamColors[TeamAlly], noMask)
			u.View.DrawShape(x+sz/2-sz/2, y+sz/2-2, sz-2, 3, 0, 0, palette.Black, noMask)
			u.View.DrawShape(x+hpWidth/2-sz/2+2, y+sz/2-2, hpWidth, 1, 0, 0, teamColors[TeamAlly], noMask)
		} else if hovered && click && lastSummonIndex < 0 {
			// Units = append(Units, unit)
			// unit.PrepareSpawn()
			u.SummonIndex = index
		}

		u.TryShowTooltip(u.View, shape, newCursor)

		if u.SummonIndex == index {
			u.Highlight(u.View, shape, palette.White)
		}

		if unit.IsSummoned() && hovered {
			u.Highlight(View, unit.Shape, palette.LightGray)
		} else if unit.IsSummoned() && unit.Shape.ContainsPoint(View.MousePosition()) {
			u.Highlight(u.View, shape, palette.LightGray)
		}
	}

	u.Top.X, u.Top.Y = tx, ty+u.Top.Height/2
	u.View.DrawObject(u.Top)

	u.ZoneInfo.Text = zoneInfos[CurrentZone.kind]
	u.ZoneInfo.X, u.ZoneInfo.Y = u.Top.X, u.Top.Y-sz+4.5
	u.View.DrawObject(u.ZoneInfo)

	u.Coins.X, u.Coins.Y = u.Top.X, u.Top.Y-9
	u.View.DrawObject(u.Coins)
	u.TryShowTooltip(u.View, u.Coins.Shape, cursor.Arrow)

	u.TeamGlory[TeamAlly].X, u.TeamGlory[TeamAlly].Y = u.Top.X-sz*3.5, u.Top.Y-6
	u.TeamGlory[TeamEnemy].X, u.TeamGlory[TeamEnemy].Y = u.Top.X+sz*3.5, u.Top.Y-6
	u.View.DrawObject(u.TeamGlory[TeamAlly])
	u.View.DrawObject(u.TeamGlory[TeamEnemy])
	u.TryShowTooltip(u.View, u.TeamGlory[TeamAlly].Shape, cursor.Arrow)
	u.TryShowTooltip(u.View, u.TeamGlory[TeamEnemy].Shape, cursor.Arrow)

	if u.SummonIndex >= 0 {
		var unitX, unitY = u.UnitPosition(u.SummonIndex)
		var mx, my = u.View.MousePosition()
		mouse.SetCursor(cursor.Move)
		u.DrawPath(unitX, unitY, mx, my)
		u.View.DrawImage(mx, my, sz, sz, 0, Characters[Player.Units[u.SummonIndex].Character].Icon, palette.White, noMask)
	} else if u.Tooltip != nil {
		mouse.SetCursor(u.Tooltip.cursor)
		u.Highlight(u.Tooltip.view, u.Tooltip.shape, highlightCursorColors[u.Tooltip.cursor])
	}
}

func (u *GUI) TryShowTooltip(view *graphics.View, shape geometry.Shape, cursor int) {
	if !shape.ContainsPoint(view.MousePosition()) {
		return
	}
	u.Tooltip = &Tooltip{view: view, shape: shape, cursor: cursor}
}
func (u *GUI) Highlight(view *graphics.View, shape geometry.Shape, color uint) {
	shape.Width += 2
	shape.Height += 2
	u.HoverHighlight.Shape = shape
	u.HoverHighlight.Effects.Tint = color
	view.DrawObject(u.HoverHighlight)
}

func (u *GUI) DrawPath(fromX, fromY, toX, toY float32) {
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
		u.View.DrawShape(line.X, line.Y, width, size, line.Angle, 1, palette.Black, geometry.Area{})
		u.View.DrawShape(line.X, line.Y, width-2, size-2, line.Angle, 1, palette.White, geometry.Area{})
	}
}
