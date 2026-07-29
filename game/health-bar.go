package game

import (
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/text"
)

var healthBarColors = map[Team]uint{TeamAlly: palette.Green, TeamEnemy: palette.Red, TeamNeutral: palette.Cyan}

type HealthBar struct {
	background, fill, label graphics.Object

	big bool
}

func NewHealthBar(width float32, big bool, team Team) HealthBar {
	var bgr = graphics.NewShapeRectangle(0, 0, width, 4, 0)
	var fill = graphics.NewShapeRectangle(0, 0, 0, 0, 0)
	var label = graphics.NewTextbox(0, 0, 50, 10, 0)
	if !big { // smaller unit bar
		bgr.Height = 2
	}
	bgr.Effects.BorderSize = 0
	fill.Effects.FillColor = healthBarColors[team]
	bgr.Effects.FillColor = color.Darken(fill.Effects.FillColor, 0.9)
	fill.Effects.BorderSize = 0
	label.Effects.TextLineHeight, label.Effects.FillColor = min(8, width/3), 0
	label.Effects.TextAlignX, label.Effects.TextAlignY = 0.5, 1
	return HealthBar{background: bgr, fill: fill, label: label, big: big}
}

func (hb *HealthBar) Update(target geometry.Shape, health, maxHealth int) {
	if health == 0 || maxHealth == 0 {
		return
	}

	var border float32 = 2
	hb.background.X, hb.background.Y = target.X, target.Y-target.Height/2
	if !hb.big {
		border = 1
	}

	hb.fill.Width = number.Map(float32(health), 0, float32(maxHealth), 0, hb.background.Width-border)
	hb.fill.Height = hb.background.Height - border
	hb.fill.X, hb.fill.Y = hb.background.X-hb.background.Width/2+hb.fill.Width/2+border/2, hb.background.Y
	hb.label.Text = text.New(health)
	hb.label.X, hb.label.Y = hb.background.X, hb.background.Y-hb.label.Height/2

	View.DrawObject(&hb.background)
	View.DrawObject(&hb.fill)
	View.DrawObject(&hb.label)
}
