package game

import (
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/text"
	"pure-game-kit/packages/utility/time"
)

var healthBarColors = [3]uint{TeamAlly: palette.Green, TeamEnemy: palette.Red, TeamNeutral: palette.Cyan}

type HealthBar struct {
	background, fill, label graphics.Object

	fadeOutTimer, fadeOutDuration float32
}

func NewHealthBar(width float32, team Team) HealthBar {
	var bgr = graphics.NewShapeRectangle(0, 0, width, 6, 0)
	var fill = graphics.NewShapeRectangle(0, 0, 0, 0, 0)
	var label = graphics.NewTextbox(0, 0, 50, 6, 0)
	bgr.Effects.BorderSize = 0
	fill.Effects.FillColor = healthBarColors[team]
	bgr.Effects.FillColor = color.Darken(fill.Effects.FillColor, 0.9)
	fill.Effects.BorderSize = 0
	label.Effects.TextLineHeight, label.Effects.FillColor = 6, 0
	label.Effects.TextWeight = 0.15
	label.Effects.TextShadowColor = 0
	label.Effects.OutlineSize = 0.7
	label.Effects.TextSymbolGap = 20
	label.Effects.OutlineColor = palette.Black
	label.Effects.TextAlignX, label.Effects.TextAlignY = 0.5, 0.5
	return HealthBar{background: bgr, fill: fill, label: label}
}
func (hb *HealthBar) FadeOut(duration float32) {
	hb.fadeOutTimer = duration
	hb.fadeOutDuration = duration
}

func (hb *HealthBar) Update(target geometry.Shape, health, maxHealth int, mask geometry.Area) {
	if hb.fadeOutTimer > 0 {
		var alpha = uint8(number.Map(hb.fadeOutTimer, hb.fadeOutDuration, 0, 255, 0))
		hb.fadeOutTimer -= time.Delta()
		hb.background.Effects.Tint = color.RGBA(255, 255, 255, alpha)
		hb.fill.Effects.Tint = color.RGBA(255, 255, 255, alpha)
		hb.label.Effects.Tint = color.RGBA(255, 255, 255, alpha)
	} else if hb.fadeOutTimer < 0 {
		return
	}

	var border float32 = 2
	hb.background.X, hb.background.Y = target.X, target.Y-target.Height/2

	hb.fill.Width = number.Map(float32(max(health, 0)), 0, float32(maxHealth), 0, hb.background.Width-border)
	hb.fill.Height = hb.background.Height - border
	hb.fill.X, hb.fill.Y = hb.background.X-hb.background.Width/2+hb.fill.Width/2+border/2, hb.background.Y
	hb.label.Text = text.New(health)
	hb.label.X, hb.label.Y = hb.background.X, hb.background.Y

	hb.background.Mask, hb.fill.Mask, hb.label.Mask = mask, mask, mask
	View.DrawObject(&hb.background)
	View.DrawObject(&hb.fill)
	View.DrawObject(&hb.label)
}
