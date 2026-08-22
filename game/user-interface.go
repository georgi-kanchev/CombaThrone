package game

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
)

type Icon uint8
type HUD struct {
	Top, ZoneInfo, Coins *graphics.Object
	TeamGlory            [2]*graphics.Object
}

const IconGlory, IconHealth, IconCoin Icon = 0, 1, 2

var Tags = [5]string{IconGlory: "*", IconHealth: "~", IconCoin: "$"}

func NewHUD() *HUD {
	var top = graphics.NewSprite(0, 0, 1, UserInterface.Crops("hud-top")[0])
	var label = graphics.NewTextbox(0, 0, TileSize+6, TileSize/2, 0)
	label.Effects.FillColor, label.Effects.TextLineHeight = 0, 14
	label.Effects.TextAlignX, label.Effects.TextAlignY = 0.5, 0.5
	label.Effects.TextShadowColor, label.Effects.TextShadowWeight = 0, 0
	label.Effects.OutlineSize, label.Effects.OutlineColor = 0.4, palette.Black
	var info, ally, enemy, coins = label, label, label, label
	coins.Effects.TextLineHeight, coins.Effects.Tint = 10, palette.White
	info.Width, info.Effects.TextLineHeight = TileSize*8.5, 8
	info.Effects.TextSymbolGap, info.Effects.TextColor = 10, palette.Gold
	return &HUD{ZoneInfo: &info, Top: &top, TeamGlory: [2]*graphics.Object{&ally, &enemy}, Coins: &coins}
}

//=================================================================

func (h *HUD) Update() {
	var tx, ty = View.PointFromEdge(0.5, 0)
	h.Top.X, h.Top.Y = tx, ty+h.Top.Height/2
	View.DrawObject(h.Top)

	h.ZoneInfo.Text = zoneInfos[CurrentZone.kind]
	h.ZoneInfo.X, h.ZoneInfo.Y = h.Top.X, h.Top.Y-TileSize+4.5
	View.DrawObject(h.ZoneInfo)

	h.Coins.X, h.Coins.Y = h.Top.X, h.Top.Y-7
	View.DrawObject(h.Coins)

	h.TeamGlory[TeamAlly].X, h.TeamGlory[TeamAlly].Y = h.Top.X-TileSize*3.5+0.5, h.Top.Y-6
	h.TeamGlory[TeamEnemy].X, h.TeamGlory[TeamEnemy].Y = h.Top.X+TileSize*3.5+0.5, h.Top.Y-6
	View.DrawObject(h.TeamGlory[TeamAlly])
	View.DrawObject(h.TeamGlory[TeamEnemy])
}
