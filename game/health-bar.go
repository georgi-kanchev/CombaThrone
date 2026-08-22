package game

import (
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/easing"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/text"
)

type HealthBar struct {
	background, fill, damage, label, glory *graphics.Object

	timer, duration          float32
	lastValue                int
	team                     Team
	toGlory, subtractedGlory bool
	startX, startY           float32
}

func NewHealthBar(width float32, team Team, offLane bool) *HealthBar {
	var height float32 = 6
	if offLane {
		height = 3
	}
	var bgr = graphics.NewShapeRectangle(0, 0, width, height, 0)
	var fill = graphics.NewShapeRectangle(0, 0, 0, 0, 0)
	var data = HealthBar{team: team, background: &bgr, fill: &fill}
	bgr.Effects.FillColor, bgr.Effects.BorderSize = color.Darken(fill.Effects.FillColor, 0.9), 0
	fill.Effects.FillColor, fill.Effects.BorderSize = teamColors[team], 0

	if !offLane {
		var label = graphics.NewTextbox(0, 0, 50, 24, 0)
		label.Effects.TextLineHeight, label.Effects.FillColor = 6, 0
		label.Effects.TextWeight, label.Effects.TextSymbolGap = 0.15, 20
		label.Effects.TextShadowColor, label.Effects.TextShadowWeight = 0, 0
		label.Effects.OutlineSize, label.Effects.OutlineColor = 0.65, palette.Black
		label.Effects.TextAlignX, label.Effects.TextAlignY = 0.5, 0.5
		var glory, dmg = label, fill // copy
		glory.Effects.TextLineHeight, glory.Effects.OutlineSize = 10, 0.45
		glory.Height = 20
		glory.Effects.Tint = teamColors[team]

		data.label, data.glory, data.damage = &label, &glory, &dmg
	}

	return &data
}

//=================================================================

func (hb *HealthBar) FadeOut(duration float32) {
	hb.duration, hb.timer = duration, duration
}
func (hb *HealthBar) MoveToGlory(duration float32) {
	if hb.toGlory {
		return
	}
	hb.toGlory = true
	hb.duration, hb.timer = duration, duration
	hb.startX, hb.startY = hb.label.X, hb.label.Y
	hb.glory.Text = text.New(-hb.lastValue, UITags[IconGlory])

	if hb.team == TeamAlly {
		PlaySound(AudioEventPositive)
	} else {
		PlaySound(AudioEventNegative)
	}
}

func (hb *HealthBar) Update(target geometry.Shape, health, maxHealth int, mask geometry.Area) {
	if hb == nil {
		return
	}

	if hb.timer > 0 {
		var alpha = uint8(number.Map(hb.timer, hb.duration, 0, 255, 0))
		hb.timer -= DeltaTimeScaled()
		hb.background.Effects.Tint = color.RGBA(255, 255, 255, alpha)
		hb.fill.Effects.Tint = color.RGBA(255, 255, 255, alpha)
		hb.damage.Effects.Tint = color.RGBA(255, 255, 255, alpha)
		hb.label.Effects.Tint = color.RGBA(255, 255, 255, alpha)
	} else if hb.timer < 0 {
		if hb.toGlory && !hb.subtractedGlory {
			hb.subtractedGlory = true

			if hb.team == TeamAlly {
				EnemyBase.Glory -= hb.lastValue
			} else {
				AllyBase.Glory -= hb.lastValue
			}
		}
		return
	}

	var border float32 = 2
	hb.background.X, hb.background.Y = target.X, target.Y-target.Height/2

	hb.fill.Width = number.Map(float32(max(health, 0)), 0, float32(maxHealth), 0, hb.background.Width-border)
	hb.fill.Height = hb.background.Height - border
	hb.fill.X, hb.fill.Y = hb.background.X-hb.background.Width/2+hb.fill.Width/2+border/2, hb.background.Y

	hb.background.Mask, hb.fill.Mask = mask, mask
	View.DrawObject(hb.background)

	if hb.damage != nil {
		const dmgSpeed = 0.4
		hb.damage.Width += (hb.fill.Width - hb.damage.Width) * float32(1.0-number.Power(0.01, dmgSpeed*DeltaTimeScaled()))
		if number.Absolute(hb.fill.Width-hb.damage.Width) < 0.1 {
			hb.damage.Width = hb.fill.Width
		}
		hb.damage.Height = hb.fill.Height
		hb.damage.X, hb.damage.Y = hb.background.X-hb.background.Width/2+hb.damage.Width/2+border/2, hb.background.Y

		if hb.team == TeamAlly {
			hb.damage.Effects.FillColor = palette.Red
		} else {
			hb.damage.Effects.FillColor = palette.Orange
		}
		hb.damage.Mask = mask
		View.DrawObject(hb.damage)
	}

	View.DrawObject(hb.fill)

	if hb.label != nil {
		hb.label.X, hb.label.Y = hb.background.X, hb.background.Y
		if hb.lastValue != health {
			hb.label.Text = text.New(health)
		}
		hb.lastValue = health
		hb.label.Mask = mask
		View.DrawObject(hb.label)
	}

	if hb.toGlory && hb.team != TeamNeutral {
		var progress = number.Limit(number.Map(hb.timer, hb.duration, 0, 0, 1), 0, 1)
		var targetX, targetY = View.PointFromView(UI.View, UI.TeamGlory[1-hb.team].X, UI.TeamGlory[1-hb.team].Y)
		hb.glory.X = number.Map(easing.CircOut(progress), 0, 1, hb.startX, targetX)
		hb.glory.Y = number.Map(easing.CubicIn(progress), 0, 1, hb.startY, targetY)

		const riseEnd, fallStart = 0.5, 0.75
		if progress < riseEnd { // smooth rise in
			var t = number.Map(progress, 0, riseEnd, 0, 1)
			var alpha = byte(number.Map(easing.CubicOut(t), 0, 1, 0, 255))
			var r, g, b, _ = color.Channels(hb.glory.Effects.Tint)
			hb.glory.Effects.TextLineHeight = number.Map(easing.CubicOut(t), 0, 1, 4, 16)
			hb.glory.Effects.Tint = color.RGBA(r, g, b, alpha)
		} else if progress < fallStart { // breathing hang
			var t = number.Map(progress, riseEnd, fallStart, 0, 1)
			var microFloat = number.Sine(t*3.14159) * 2.0 // adds pulsing 16 -> 18 -> 16...
			hb.glory.Effects.TextLineHeight = 16 + microFloat

		} else { // faster fall
			var t = number.Map(progress, fallStart, 1, 0, 1)
			var alpha = byte(number.Map(easing.CubicIn(t), 0, 1, 255, 0))
			var r, g, b, _ = color.Channels(hb.glory.Effects.Tint)
			hb.glory.Effects.TextLineHeight = number.Map(easing.CubicIn(t), 0, 1, 16, 12)
			hb.glory.Effects.Tint = color.RGBA(r, g, b, alpha)
		}

		View.DrawObject(hb.glory)
	}
}
