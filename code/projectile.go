package game

import (
	"game/code/constants/effect"
	"game/code/constants/entrance"
	proj "game/code/constants/projectile"
	"game/code/constants/splash"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/motion/particles"
	"pure-game-kit/packages/utility/angle"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/point"
	"pure-game-kit/packages/utility/random"
)

type Projectile struct {
	graphics.Object
	Value int

	StartX, StartY, StartZ    float32
	TargetX, TargetY, TargetZ float32

	Z, Age, TravelTime, ArcHeight float32

	FadeOutTime float32
	Kind        proj.Enum

	Owner         *Unit
	EnemyEntrance *Entrance
	Particles     *particles.System

	Triggered, FreeRoation bool

	Hidden bool // makes sure the Particles finish off even after the Projectile itself is gone, FadeOutTime removes it later
}
type ProjectileBase struct {
	SplashTop, SplashMiddle, SplashBottom splash.Enum
	FadeOutTime, MissMultiplier           float32
	FreeRotation, PersistOnImpact         bool
	ParticlesFunc                         func(p *particles.Particle) bool
	Effect                                effect.Enum
}

var Projectiles, ProjectilesBehind []*Projectile = make([]*Projectile, 0, 32), make([]*Projectile, 0, 32)
var ProjectilesBase = [proj.COUNT]ProjectileBase{
	proj.Arrow:       {FadeOutTime: 10, PersistOnImpact: true, MissMultiplier: 1.0},
	proj.SmallBullet: {FadeOutTime: 0.1, MissMultiplier: 1.0},
	proj.RedPotion: {FadeOutTime: 10, FreeRotation: true, MissMultiplier: 1.0, ParticlesFunc: ParticlesGlass,
		SplashTop: splash.Fire2, SplashMiddle: splash.Fire3, SplashBottom: splash.Fire2},
	proj.GreenPotion: {FadeOutTime: 10, FreeRotation: true, MissMultiplier: 0.2, ParticlesFunc: ParticlesGlass,
		SplashMiddle: splash.Heal},
	proj.YellowPotion: {FadeOutTime: 10, FreeRotation: true, MissMultiplier: 1.0, ParticlesFunc: ParticlesGlass,
		SplashMiddle: splash.PoisonCloud},
	proj.MediumSpit: {FadeOutTime: 0.1, MissMultiplier: 1.0, SplashMiddle: splash.PoisonSpit, Effect: effect.Shaman},
}

func (u *Unit) NewProjectile(x, y, z, targetX, targetY, targetZ float32, value int, kind proj.Enum, enemyEntrance *Entrance) *Projectile {
	var speed float32 = 100 * Characters[u.Character].Projectile.Speed
	var missMultiplier float32 = 1 / Characters[u.Character].Projectile.Speed
	var randomY = random.Range[float32](-12, 12)
	missMultiplier *= ProjectilesBase[kind].MissMultiplier
	randomY *= missMultiplier
	if enemyEntrance != nil { // should not miss the entrances
		missMultiplier = 0
		targetY += TileSize / 2
	}

	var fadeOutTime, freeRotation = ProjectilesBase[kind].FadeOutTime, ProjectilesBase[kind].FreeRotation
	var dist = point.DistanceToPoint(x, y, targetX, targetY)
	var totalTime = max(dist/speed, 0.01) // prevent division by zero
	var parabola = Characters[u.Character].Projectile.ParabolaMultiplier
	var proj = &Projectile{Owner: u, Kind: kind,
		Object: graphics.NewSprite(x, y, 1, Decor.Crops("projectiles")[kind]),
		StartX: x, StartY: y, StartZ: z, Z: z,
		TargetX:    targetX + random.Range[float32](-12, 12)*missMultiplier,
		TargetY:    targetY + randomY,
		TargetZ:    targetZ + random.Range[float32](-0.35, 0.35)*missMultiplier,
		TravelTime: totalTime, ArcHeight: dist / 3 * parabola, Value: value, EnemyEntrance: enemyEntrance,
		FadeOutTime: fadeOutTime, FreeRoation: freeRotation, Particles: particles.NewSystem(),
	}
	proj.Details.Tint = Characters[u.Character].Projectile.Tint
	return proj
}

//=================================================================

func (p *Projectile) Update() {
	p.Age += DeltaTimeScaled()
	p.Particles.Update()

	var enemyTeam = 1 - p.Owner.Team
	var progress = min(1.0, p.Age/p.TravelTime)
	var groundX = p.StartX + (p.TargetX-p.StartX)*progress
	var groundY = p.StartY + (p.TargetY-p.StartY)*progress
	var heightArc = 4 * p.ArcHeight * progress * (1 - progress)

	p.X = groundX
	p.Y = groundY - heightArc
	p.Z = p.StartZ + (p.TargetZ-p.StartZ)*progress

	var velX = p.TargetX - p.StartX
	var velY = (p.TargetY - p.StartY) - 4*p.ArcHeight*(1-2*progress)

	if p.FreeRoation && !p.Triggered {
		p.Angle += velX * DeltaTimeScaled()
	} else {
		p.Angle = angle.BetweenPoints(0, 0, velX, velY)
	}

	if progress >= 1.0 {
		var e = p.EnemyEntrance
		if !p.Triggered {

			if !ProjectilesBase[p.Kind].PersistOnImpact {
				p.Hidden = true
			}

			var _, aboveGroundY = PointAtCell(0, 6)
			if p.Y > aboveGroundY { // rendering projectiles above walls is fine
				Projectiles = collection.Remove(Projectiles, p)
				ProjectilesBehind = append(ProjectilesBehind, p)
			}

			if e != nil && !e.IsOpen() && e.Health > 0 && number.IsWithin(p.X, e.Tiles[0].X, TileSize/2) {
				e.TakeDamage(p.Value)
				if e.Kind == entrance.ShortGate || e.Kind == entrance.TallGate {
					p.trigger(28)
				} else {
					p.trigger(0)
				}
				if e.Health > 0 { // crumble sound played by entrance itself
					var sounds = Characters[p.Owner.Character].Sounds.HitWood
					if e.Kind == entrance.ShortGate || e.Kind == entrance.TallGate {
						sounds = Characters[p.Owner.Character].Sounds.HitMetal
					}
					PlaySound(sounds)
				}
			} else {
				p.trigger(0)
				PlaySound(Characters[p.Owner.Character].Sounds.HitGround)
			}
		}

		var alpha = min(255, number.Map(p.Age, p.TravelTime, p.TravelTime+p.FadeOutTime, 255, 0))
		var entranceDestroyed = e != nil && (e.IsOpen() || e.OpenY > 0 || e.Health <= 0)
		p.Details.Tint = color.RGBA(255, 255, 255, byte(alpha))

		if p.Age > p.TravelTime+p.FadeOutTime || entranceDestroyed {
			Projectiles = collection.Remove(Projectiles, p)
			ProjectilesBehind = collection.Remove(ProjectilesBehind, p)
		}
		if !p.Hidden {
			View.DrawObject(&p.Object)
		}
		return
	}

	for _, u := range Units {
		if u.Health <= 0 || enemyTeam != u.Team {
			continue
		}

		var hb = u.Hitbox()
		if u.IsGarrisoner() { // hide behind wall
			hb.Height /= 2        // shrink hitbox
			hb.Y -= hb.Height / 2 // move up
		}

		if !p.Triggered && p.Shape.Overlaps(hb) && number.IsWithin(p.Z, u.Z, 0.2) {
			p.trigger(0)
			u.AffectHealth(-p.Value, true, p.Owner)
			u.AddEffect(ProjectilesBase[p.Kind].Effect)
			p.Hidden = true
			PlaySound(Characters[p.Owner.Character].Sounds.HitFlesh)
			return
		}
	}

	var shadowAngle = angle.BetweenPoints(p.StartX, p.StartY, p.TargetX, p.TargetY)
	if p.TargetZ == p.Z { // same lane path shouldn't bend the shadow angle
		shadowAngle = 0
	}
	if !p.Hidden {
		DrawShadow(groundX, p.Z, number.Absolute(p.Width), p.Height*0.25, shadowAngle, p.Mask)
		View.DrawObject(&p.Object)
	}
}

// private ========================================================

func (p *Projectile) trigger(offsetY float32) {
	var base = ProjectilesBase[p.Kind]
	var top, mid, bot = base.SplashTop, base.SplashMiddle, base.SplashBottom
	p.Triggered = true
	p.Particles.EmitFromPoint(50, p.X, p.Y, base.ParticlesFunc)

	if top != splash.NONE {
		Splashes = append(Splashes, NewSplash(top, p.Owner.Team, p.X, p.Y-8+offsetY, p.Owner))
	}
	if mid != splash.NONE {
		Splashes = append(Splashes, NewSplash(mid, p.Owner.Team, p.X, p.Y+offsetY, p.Owner))
	}
	if bot != splash.NONE {
		Splashes = append(Splashes, NewSplash(bot, p.Owner.Team, p.X, p.Y+8+offsetY, p.Owner))
	}
}
