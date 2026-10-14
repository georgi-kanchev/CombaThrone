package game

import (
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/motion/particles"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/random"
)

func ParticlesGlass(p *particles.Particle) (alive bool) {
	if p.Age == 0 {
		p.CustomData["startY"] = p.Y + random.Range[float32](0, 12)
		p.CustomData["height"] = random.Range[float32](0.2, 3)
		p.Scale = random.Range[float32](0.2, 2)
		p.VelocityX = random.Range[float32](-40, 40)
		p.VelocityY = random.Range[float32](-100, -40)
		p.Color = random.Range[uint](200, 255) // only red
	}

	var dt = DeltaTimeScaled()
	var ang = random.Range[float32](0, 360)
	p.Age += dt
	p.VelocityY += Gravity / 2 * dt

	p.X += p.VelocityX * dt
	p.Y += p.VelocityY * dt

	if p.Y > p.CustomData["startY"].(float32) {
		p.Y = p.CustomData["startY"].(float32)
		p.VelocityX, p.VelocityY = 0, 0
		ang = 45
	}

	var alpha = number.Map(p.Age, 1.5, 2.0, 255, 0)
	var col = color.RGBA(byte(p.Color), byte(p.Color), byte(p.Color), byte(number.Limit(alpha, 0, 255)))
	View.DrawShape(p.X, p.Y, p.Scale, p.CustomData["height"].(float32), ang, 0.1, col, geometry.Area{})
	return p.Age < 2.0
}
