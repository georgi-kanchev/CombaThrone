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
	Background, Fill, Damage, Label, Glory *graphics.Object

	Timer, Duration          float32
	LastValue                int
	Team                     Team
	ToGlory, SubtractedGlory bool
	StartX, StartY           float32
}

func NewHealthBar(width float32, team Team, offLane bool) *HealthBar {
	var height float32 = 6
	if offLane {
		height = 3
	}
	var bgr = graphics.NewShapeRectangle(0, 0, width, height, 0)
	var fill = graphics.NewShapeRectangle(0, 0, 0, 0, 0)
	var data = HealthBar{Team: team, Background: &bgr, Fill: &fill}
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

		data.Label, data.Glory, data.Damage = &label, &glory, &dmg
	}

	return &data
}

//=================================================================

func (hb *HealthBar) FadeOut(duration float32) {
	hb.Duration, hb.Timer = duration, duration
}
func (hb *HealthBar) MoveToGlory(duration float32) {
	if hb.ToGlory {
		return
	}
	hb.ToGlory = true
	hb.Duration, hb.Timer = duration, duration
	hb.StartX, hb.StartY = hb.Label.X, hb.Label.Y
	hb.Glory.Text = text.New(Tags[IconGlory], -hb.LastValue)

	if hb.Team == TeamAlly {
		PlaySound(AudioEventPositive)
	} else {
		PlaySound(AudioEventNegative)
	}
}

func (hb *HealthBar) Update(target geometry.Shape, health, maxHealth int, mask geometry.Area) {
	if hb == nil {
		return
	}

	if hb.Timer > 0 {
		var alpha = uint8(number.Map(hb.Timer, hb.Duration, 0, 255, 0))
		hb.Timer -= DeltaTimeScaled()
		hb.Background.Effects.Tint = color.RGBA(255, 255, 255, alpha)
		hb.Fill.Effects.Tint = color.RGBA(255, 255, 255, alpha)
		hb.Damage.Effects.Tint = color.RGBA(255, 255, 255, alpha)
		hb.Label.Effects.Tint = color.RGBA(255, 255, 255, alpha)
	} else if hb.Timer < 0 {
		if hb.ToGlory && !hb.SubtractedGlory {
			hb.SubtractedGlory = true
			Bases[1-hb.Team].Glory -= hb.LastValue
		}
		return
	}

	var border float32 = 2
	hb.Background.X, hb.Background.Y = target.X, target.Y-target.Height/2

	hb.Fill.Width = number.Map(float32(max(health, 0)), 0, float32(maxHealth), 0, hb.Background.Width-border)
	hb.Fill.Height = hb.Background.Height - border
	hb.Fill.X, hb.Fill.Y = hb.Background.X-hb.Background.Width/2+hb.Fill.Width/2+border/2, hb.Background.Y

	hb.Background.Mask, hb.Fill.Mask = mask, mask
	View.DrawObject(hb.Background)

	if hb.Damage != nil {
		const dmgSpeed = 0.4
		hb.Damage.Width += (hb.Fill.Width - hb.Damage.Width) * float32(1.0-number.Power(0.01, dmgSpeed*DeltaTimeScaled()))
		if number.Absolute(hb.Fill.Width-hb.Damage.Width) < 0.1 {
			hb.Damage.Width = hb.Fill.Width
		}
		hb.Damage.Height = hb.Fill.Height
		hb.Damage.X, hb.Damage.Y = hb.Background.X-hb.Background.Width/2+hb.Damage.Width/2+border/2, hb.Background.Y

		if hb.Team == TeamAlly {
			hb.Damage.Effects.FillColor = palette.Red
		} else {
			hb.Damage.Effects.FillColor = palette.Orange
		}
		hb.Damage.Mask = mask
		View.DrawObject(hb.Damage)
	}

	View.DrawObject(hb.Fill)

	if hb.Label != nil {
		hb.Label.X, hb.Label.Y = hb.Background.X, hb.Background.Y
		if hb.LastValue != health {
			hb.Label.Text = text.New(health)
		}
		hb.LastValue = health
		hb.Label.Mask = mask
		View.DrawObject(hb.Label)
	}

	if hb.ToGlory && hb.Team != TeamNeutral {
		var progress = number.Limit(number.Map(hb.Timer, hb.Duration, 0, 0, 1), 0, 1)
		var targetX, targetY = View.PointFromView(GameHUD.View, GameHUD.TeamGlory[1-hb.Team].X, GameHUD.TeamGlory[1-hb.Team].Y)
		hb.Glory.X = number.Map(easing.CircOut(progress), 0, 1, hb.StartX, targetX)
		hb.Glory.Y = number.Map(easing.CubicIn(progress), 0, 1, hb.StartY, targetY)

		const riseEnd, fallStart = 0.5, 0.75
		if progress < riseEnd { // smooth rise in
			var t = number.Map(progress, 0, riseEnd, 0, 1)
			var alpha = byte(number.Map(easing.CubicOut(t), 0, 1, 0, 255))
			var r, g, b, _ = color.Channels(hb.Glory.Effects.Tint)
			hb.Glory.Effects.TextLineHeight = number.Map(easing.CubicOut(t), 0, 1, 4, 16)
			hb.Glory.Effects.Tint = color.RGBA(r, g, b, alpha)
		} else if progress < fallStart { // breathing hang
			var t = number.Map(progress, riseEnd, fallStart, 0, 1)
			var microFloat = number.Sine(t*3.14159) * 2.0 // adds pulsing 16 -> 18 -> 16...
			hb.Glory.Effects.TextLineHeight = 16 + microFloat

		} else { // faster fall
			var t = number.Map(progress, fallStart, 1, 0, 1)
			var alpha = byte(number.Map(easing.CubicIn(t), 0, 1, 255, 0))
			var r, g, b, _ = color.Channels(hb.Glory.Effects.Tint)
			hb.Glory.Effects.TextLineHeight = number.Map(easing.CubicIn(t), 0, 1, 16, 12)
			hb.Glory.Effects.Tint = color.RGBA(r, g, b, alpha)
		}

		View.DrawObject(hb.Glory)
	}
}
