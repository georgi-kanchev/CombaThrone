package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/button"
	"pure-game-kit/packages/input/mouse/cursor"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
)

type Icon uint8
type GUI struct {
	View                             *graphics.View
	Top, ZoneInfo, Coins, UnitsPanel *graphics.Object
	TeamGlory                        [TeamCount]*graphics.Object

	Tooltip        *Tooltip
	HoverHighlight *graphics.Object
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
		HoverHighlight: &highlight,
	}
}

//=================================================================

func (u *GUI) Update() {
	const scale = 1
	var sz float32 = TileSize
	u.View.Zoom = scale * 5

	var tx, ty = u.View.PointFromEdge(0.5, 0)
	u.UnitsPanel.Y = ty + u.Top.Height + u.UnitsPanel.Height/2 - sz - 2
	u.View.DrawObject(u.UnitsPanel)

	var slotImageId = UserInterface.Crops("slot")[0]
	var width, height = sz + 7, sz + 7
	var offX, offY float32 = width/2 - width*2, height/2 - height
	for index, unit := range Player.Units {
		var i, j = index / 4, index % 4
		var x, y = u.UnitsPanel.X + float32(j)*width + offX, u.UnitsPanel.Y + float32(i)*height + offY
		u.View.DrawImage(x, y, sz, sz, 0, slotImageId, palette.White, geometry.Area{})

		if unit != nil {
			var shape = geometry.NewRoundedRectangle(x, y, sz, sz, 0, 0)
			var icons = UserInterface.Crops("icons")
			var newCursor = cursor.Hand
			if unit.State == StateDecaying || unit.IsSummoned() {
				newCursor = cursor.NotAllowed
			}
			u.View.DrawImage(x, y, sz, sz, 0, Characters[unit.Character].Icon, palette.White, geometry.Area{})
			if unit.State == StateDecaying {
				var timerWidth = number.Map(unit.hurtTimer, 0, -unit.Stats.RespawnTimer, sz-4, 0)
				var icon, col = IconDeath, palette.Red
				if unit.Stats.Health > 0 { // got into the enemy base
					icon, col = IconGlory, palette.Green
				}

				u.View.DrawShape(x, y, sz, sz, 0, 0, color.RGBA(0, 0, 0, 127), geometry.Area{})
				u.View.DrawImage(x, y, sz/2, sz/2, 0, icons[icon], col, geometry.Area{})
				u.View.DrawShape(x+TileSize/2-sz/2, y+sz/2-2, TileSize-2, 3, 0, 0, palette.Black, geometry.Area{})
				u.View.DrawShape(x+timerWidth/2-sz/2+2, y+sz/2-2, timerWidth, 1, 0, 0, col, geometry.Area{})
			} else if unit.IsSummoned() {
				u.View.DrawShape(x, y, sz, sz, 0, 0, color.RGBA(0, 0, 0, 127), geometry.Area{})
				u.View.DrawImage(x, y, sz/2, sz/2, 0, icons[IconHealth], palette.Green, geometry.Area{})
			} else if shape.ContainsPoint(u.View.MousePosition()) && mouse.IsButtonJustPressed(button.Left) {
				Units = append(Units, unit)
				unit.PrepareSpawn()
			}

			u.TryShowTooltip(u.View, shape, newCursor)
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

	if u.Tooltip != nil {
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
