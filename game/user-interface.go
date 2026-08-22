package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
)

type Icon uint8
type GUI struct {
	View                             *graphics.View
	Top, ZoneInfo, Coins, UnitsPanel *graphics.Object
	TeamGlory                        [TeamCount]*graphics.Object
}

const IconGlory, IconHealth, IconCoin Icon = 0, 1, 2

var UITags = [5]string{IconGlory: "*", IconHealth: "~", IconCoin: "$"}

func NewUI() *GUI {
	var panelNinePatch = assets.LoadImage9Patch(UserInterface.Crops("panel")[0], 8, 8, 8, 8)
	var view = graphics.NewView(1)
	var top = graphics.NewSprite(0, 0, 1, UserInterface.Crops("hud-top")[0])

	var unitsPanel = graphics.NewSprite(0, 0, 1, panelNinePatch)
	unitsPanel.Width, unitsPanel.Height = TileSize*5.5, TileSize*3

	var label = graphics.NewTextbox(0, 0, TileSize+6, TileSize/2, 0)
	label.Effects.FillColor, label.Effects.TextLineHeight = 0, 14
	label.Effects.TextAlignX, label.Effects.TextAlignY = 0.5, 0.5
	label.Effects.TextShadowColor, label.Effects.TextShadowWeight = 0, 0
	label.Effects.OutlineSize, label.Effects.OutlineColor = 0.4, palette.Black

	var coins = label
	coins.Effects.TextLineHeight, coins.Effects.Tint = 10, palette.White
	var info = label

	info.Width, info.Effects.TextLineHeight = TileSize*8.5, 8
	info.Effects.TextSymbolGap, info.Effects.TextColor = 10, palette.Gold

	var ally, enemy = label, label
	var glory = [TeamCount]*graphics.Object{&ally, &enemy}

	return &GUI{
		View: &view, ZoneInfo: &info, Top: &top, TeamGlory: glory, Coins: &coins, UnitsPanel: &unitsPanel,
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
	for i := range 2 {
		for j := range 4 {
			var x, y = u.UnitsPanel.X + float32(j)*width + offX, u.UnitsPanel.Y + float32(i)*height + offY
			u.View.DrawImage(x, y, sz, sz, 0, slotImageId, palette.White, geometry.Area{})
		}
	}

	u.Top.X, u.Top.Y = tx, ty+u.Top.Height/2
	u.View.DrawObject(u.Top)

	u.ZoneInfo.Text = zoneInfos[CurrentZone.kind]
	u.ZoneInfo.X, u.ZoneInfo.Y = u.Top.X, u.Top.Y-sz+4.5
	u.View.DrawObject(u.ZoneInfo)

	u.Coins.X, u.Coins.Y = u.Top.X, u.Top.Y-7
	u.View.DrawObject(u.Coins)

	u.TeamGlory[TeamAlly].X, u.TeamGlory[TeamAlly].Y = u.Top.X-sz*3.5+0.5, u.Top.Y-6
	u.TeamGlory[TeamEnemy].X, u.TeamGlory[TeamEnemy].Y = u.Top.X+sz*3.5+0.5, u.Top.Y-6
	u.View.DrawObject(u.TeamGlory[TeamAlly])
	u.View.DrawObject(u.TeamGlory[TeamEnemy])
}
