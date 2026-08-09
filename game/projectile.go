package game

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/angle"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/point"
	"pure-game-kit/packages/utility/random"
)

type Projectile struct {
	graphics.Object
	Damage int

	StartX, StartY, StartZ    float32
	TargetX, TargetY, TargetZ float32

	Z, Age, TravelTime, ArcHeight float32

	enemyEntrance *EntranceData
	trigger       bool
}

const ProjectileFadeOutTime float32 = 10

func NewProjectile(x, y, z, targetX, targetY, targetZ float32, damage int, targetEntrance *EntranceData) *Projectile {
	const speed float32 = 100
	var dist = point.DistanceToPoint(x, y, targetX, targetY)
	var totalTime = max(dist/speed, 0.01) // prevent division by zero
	var proj = &Projectile{
		Object: graphics.NewSprite(x, y, 1, TilesetCrops.Crops("projectiles")[0]),
		StartX: x, StartY: y, StartZ: z, Z: z,
		TargetX: targetX, TargetY: targetY + random.Range[float32](-6, 6), TargetZ: targetZ,
		TravelTime: totalTime, ArcHeight: 40.0, Damage: damage, enemyEntrance: targetEntrance,
	}
	return proj
}

func (p *Projectile) Update() {
	p.Age += DeltaTimeScaled()

	var progress = min(1.0, p.Age/p.TravelTime)
	var groundX = p.StartX + (p.TargetX-p.StartX)*progress
	var groundY = p.StartY + (p.TargetY-p.StartY)*progress
	var heightArc = 4 * p.ArcHeight * progress * (1 - progress)

	p.X = groundX
	p.Y = groundY - heightArc
	p.Z = p.StartZ + (p.TargetZ-p.StartZ)*progress

	var velX = p.TargetX - p.StartX
	var velY = (p.TargetY - p.StartY) - 4*p.ArcHeight*(1-2*progress)

	p.Angle = angle.BetweenPoints(0, 0, velX, velY)

	if progress > 0.2 && progress < 1.0 {
		for _, u := range Units {
			if u.Stats.Health <= 0 {
				continue
			}

			var hb = u.Hitbox()
			if p.Shape.Overlaps(hb) && number.IsWithin(p.Z, u.Z, 0.2) {
				p.trigger = true
				u.TakeDamage(p.Damage)
				Projectiles = collection.Remove(Projectiles, p)
				return
			}
		}
	}

	if progress < 1.0 {
		var shadowAngle = p.Angle
		if p.TargetZ == p.Z { // same lane path shouldn't bend the shadow angle
			shadowAngle = 0
		}
		DrawShadow(groundX, p.Z, number.Absolute(p.Width), p.Height*0.25, shadowAngle, p.Mask)
	} else {
		var e = p.enemyEntrance
		if !p.trigger && e != nil && !e.IsOpen() && e.Health > 0 && number.IsWithin(p.X, e.Tiles[0].X, TileSize/2) {
			p.trigger = true
			e.TakeDamage(p.Damage)
		}

		var alpha = min(255, number.Map(p.Age, p.TravelTime, p.TravelTime+ProjectileFadeOutTime, 255, 0))
		var entranceDestroyed = e != nil && (e.IsOpen() || e.openY > 0 || e.Health <= 0)
		p.Effects.Tint = color.RGBA(255, 255, 255, byte(alpha))

		if p.Age > p.TravelTime+ProjectileFadeOutTime || entranceDestroyed {
			Projectiles = collection.Remove(Projectiles, p)
		}
	}
	View.DrawObject(&p.Object)
}
