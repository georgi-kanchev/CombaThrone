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

	owner         *Unit
	enemyEntrance *Entrance
	trigger       bool
}

const ProjectileFadeOutTime float32 = 10

func (u *Unit) NewProjectile(x, y, z, targetX, targetY, targetZ float32, dmg int, enemyEntrance *Entrance) *Projectile {
	const speed float32 = 100
	var accuracyMultiplier float32 = 1
	if enemyEntrance != nil { // cannot and should not miss the entrances
		accuracyMultiplier = 0
		targetY += TileSize / 2
	}
	var dist = point.DistanceToPoint(x, y, targetX, targetY)
	var totalTime = max(dist/speed, 0.01) // prevent division by zero
	var proj = &Projectile{owner: u,
		Object: graphics.NewSprite(x, y, 1, DecorCrops.Crops("projectiles")[0]),
		StartX: x, StartY: y, StartZ: z, Z: z,
		TargetX:    targetX + random.Range[float32](-12, 12)*accuracyMultiplier,
		TargetY:    targetY + random.Range[float32](-12, 12),
		TargetZ:    targetZ + random.Range[float32](-0.35, 0.35)*accuracyMultiplier,
		TravelTime: totalTime, ArcHeight: dist / 3, Damage: dmg, enemyEntrance: enemyEntrance,
	}
	return proj
}

func (p *Projectile) Update() {
	p.Age += DeltaTimeScaled()

	var enemyTeam = 1 - p.owner.Team
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

	if progress >= 1.0 {
		var e = p.enemyEntrance
		if !p.trigger {
			p.trigger = true
			Projectiles = collection.Remove(Projectiles, p)
			ProjectilesBehind = append(ProjectilesBehind, p)

			if e != nil && !e.IsOpen() && e.Health > 0 && number.IsWithin(p.X, e.Tiles[0].X, TileSize/2) {
				e.TakeDamage(p.Damage)
				if e.Health > 0 { // crumble sound played by entrance itself
					var sounds = Characters[p.owner.Character].Sounds.HitWood
					if e.Kind == EntranceShortGate || e.Kind == EntranceTallGate {
						sounds = Characters[p.owner.Character].Sounds.HitMetal
					}
					PlaySound(sounds)
				}
			} else {
				PlaySound(Characters[p.owner.Character].Sounds.HitGround)
			}
		}

		var alpha = min(255, number.Map(p.Age, p.TravelTime, p.TravelTime+ProjectileFadeOutTime, 255, 0))
		var entranceDestroyed = e != nil && (e.IsOpen() || e.openY > 0 || e.Health <= 0)
		p.Effects.Tint = color.RGBA(255, 255, 255, byte(alpha))

		if p.Age > p.TravelTime+ProjectileFadeOutTime || entranceDestroyed {
			Projectiles = collection.Remove(Projectiles, p)
			ProjectilesBehind = collection.Remove(ProjectilesBehind, p)
		}
		View.DrawObject(&p.Object)
		return
	}

	for _, u := range Units {
		if u.Stats.Health <= 0 || enemyTeam != u.Team {
			continue
		}

		var hb = u.Hitbox()
		if u.Lane > LaneUpper { // is garrison - hide behind wall (shrink hitbox & move up)
			hb.Height /= 2
			hb.Y -= hb.Height / 2
		}

		if p.Shape.Overlaps(hb) && number.IsWithin(p.Z, u.Z, 0.2) {
			p.trigger = true
			u.TakeDamage(p.Damage)
			Projectiles = collection.Remove(Projectiles, p)
			ProjectilesBehind = collection.Remove(ProjectilesBehind, p)
			PlaySound(Characters[p.owner.Character].Sounds.HitFlesh)
			return
		}
	}

	var shadowAngle = angle.BetweenPoints(p.StartX, p.StartY, p.TargetX, p.TargetY)
	if p.TargetZ == p.Z { // same lane path shouldn't bend the shadow angle
		shadowAngle = 0
	}
	DrawShadow(groundX, p.Z, number.Absolute(p.Width), p.Height*0.25, shadowAngle, p.Mask)
	View.DrawObject(&p.Object)
}
