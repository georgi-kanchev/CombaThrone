package game

import (
	"game/code/constants/effect"
	"game/code/constants/entrance"
	"game/code/constants/splash"
	"game/code/constants/team"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/direction"
	"pure-game-kit/packages/utility/number"
)

type Splash struct {
	graphics.Object
	Kind   splash.Enum
	Anim   *motion.Animation[assets.ImageId]
	Values SplashValues
	Team   team.Enum
	Owner  *Unit

	VelocityX, VelocityY float32
}
type SplashValues struct {
	AnimName                string
	Damage                  int
	Loop, SingleTick        bool
	Timer, FadeOutTime, FPS float32
	TickTimer               float32 // ticks every 1 second, always ticks once upon creation (first update) since it starts at 0
	Tick                    func(self *Splash)
	Update                  func(self *Splash)
}

var SplashBases = map[splash.Enum]SplashValues{
	splash.Dust:   {AnimName: "dust", FPS: 10},
	splash.Pickup: {AnimName: "pickup", FPS: 10},
	splash.Fire3: {AnimName: "fire-3", FPS: 10, Loop: true, Timer: 5, Damage: 3, FadeOutTime: 0.25,
		Tick: splashTickFire},
	splash.Fire2: {AnimName: "fire-2", FPS: 10, Loop: true, Timer: 5, Damage: 2, FadeOutTime: 0.25,
		Tick: splashTickFire},
	splash.Fire1: {AnimName: "fire-1", FPS: 10, Loop: true, Timer: 5, Damage: 1, FadeOutTime: 0.25,
		Tick: splashTickFire},
	splash.Heal: {AnimName: "heal", FPS: 10, Loop: true, Timer: 5, Damage: -2, FadeOutTime: 0.25,
		Tick: splashTickHeal},
	splash.PoisonCloud: {AnimName: "poison-cloud", FPS: 10, Loop: true, Timer: 5, Damage: 2, FadeOutTime: 0.25,
		Tick: splashTickPoisonCloud},
	splash.Whirlwind: {AnimName: "whirlwind", FPS: 12, Loop: true, Timer: 20, Damage: 6, FadeOutTime: 0.5,
		SingleTick: true, Tick: splashTickWhirlwind, Update: splashUpdateWhirlwind},
	splash.WaterShield: {AnimName: "water-shield", FPS: 12, Loop: true, Timer: 6, FadeOutTime: 0.5,
		Update: splashUpdateWatershield},
	splash.LightningZap: {AnimName: "lightning-zap", FPS: 12, Damage: 6, SingleTick: true,
		Tick: splashTickLightingZap},
	splash.PoisonSpit: {AnimName: "poison-spit", FPS: 16},
}
var Splashes = make([]*Splash, 0, 8)

func NewSplash(kind splash.Enum, team team.Enum, x, y float32, owner *Unit) *Splash {
	var frames = SplashAtlas.Crops(SplashBases[kind].AnimName)
	var anim = motion.NewAnimation[assets.ImageId](SplashBases[kind].FPS, SplashBases[kind].Loop, frames...)
	var spr = graphics.NewSprite(x, y, 0, 0)
	var result = &Splash{Object: spr, Kind: kind, Team: team, Anim: &anim, Values: SplashBases[kind], Owner: owner}
	result.Update() // update width & height to be ready for single tick splashes
	return result
}

//=================================================================

func (s *Splash) Update() {
	s.Values.Timer -= DeltaTimeScaled()
	s.Values.TickTimer -= DeltaTimeScaled()
	s.X += s.VelocityX * DeltaTimeScaled()
	s.Y += s.VelocityY * DeltaTimeScaled()

	var t, bt, f = s.Values.Timer, SplashBases[s.Kind].Timer, SplashBases[s.Kind].FadeOutTime
	if number.IsBetween(t, bt-f, bt, false, false) {
		s.Details.Tint = color.RGBA(255, 255, 255, byte(number.Map(t, bt-f, bt, 255, 0)))
	} else if number.IsBetween(t, 0, f, false, false) {
		s.Details.Tint = color.RGBA(255, 255, 255, byte(number.Map(t, 0, f, 0, 255)))
	}

	s.Anim.FPS = SplashBases[s.Kind].FPS * TimeScale
	var frame = s.Anim.Frame()
	var crop = frame.CropArea()
	s.ImageId, s.Width, s.Height = frame, crop.Width, crop.Height

	if s.Values.Update != nil {
		s.Values.Update(s)
	}

	if s.Values.TickTimer < 0 {
		s.Values.TickTimer = 1

		if s.Values.Tick != nil {
			s.Values.Tick(s)
		}

		if s.Values.SingleTick {
			s.Values.TickTimer = number.Infinity()
		}
	}
	View.DrawObject(&s.Object)

	if (s.Values.Loop && s.Values.Timer < 0) || (!s.Values.Loop && s.Anim.IsJustFinished()) {
		Splashes = collection.Remove(Splashes, s)
	}
}

// private ========================================================

func (s *Splash) overlapsUnit(unit *Unit) bool {
	var hb = unit.Hitbox()
	return s.ContainsPoint(hb.X, hb.Y+hb.Height/2) || s.ContainsPoint(hb.X, hb.Y)
}

func splashTickFire(self *Splash) {
	self.tryDamageEnemy()
	var entrances = Bases[1-self.Team].Entrances
	for _, e := range entrances {
		if e.Kind == entrance.Door && self.Overlaps(e.Tiles[0].Shape) {
			e.TakeDamage(self.Values.Damage)
		}
	}
}
func splashTickPoisonCloud(self *Splash) {
	self.tryDamageEnemy()
	var entrances = Bases[1-self.Team].Entrances
	for _, e := range entrances {
		var contains = false
		for _, t := range e.Tiles {
			if self.Overlaps(t.Shape) {
				contains = true
				break
			}
		}
		if e.Kind > entrance.Door && contains {
			e.TakeDamage(self.Values.Damage)
		}
	}
}
func splashTickHeal(self *Splash) {
	for _, u := range Units {
		if u.Team == self.Team && self.overlapsUnit(u) && u.IsAlive() {
			u.AffectHealth(-self.Values.Damage, false, self.Owner)
		}
	}
}
func splashTickWhirlwind(self *Splash) {
	self.Y -= 12

	var smallestDist float32 = 99999
	var closestUnit *Unit
	for _, u := range Units {
		var ally = self.Team == team.Ally && self.X < u.X
		var enemy = self.Team == team.Enemy && self.X > u.X
		var dist = number.Absolute(u.X - self.X)
		if u.Team != self.Team && u.IsGrounded && (ally || enemy) && dist < smallestDist {
			smallestDist = dist
			closestUnit = u
		}
	}

	const speed = 100.0
	if closestUnit != nil {
		var dirX, dirY = direction.BetweenPoints(self.X, self.Y, closestUnit.X, closestUnit.Y)
		self.VelocityX, self.VelocityY = dirX*speed, dirY*speed
	} else {
		self.VelocityX = speed
		if self.Team == team.Enemy {
			self.VelocityX *= -1
		}
	}
}
func splashUpdateWhirlwind(self *Splash) {
	for _, u := range Units {
		if u.Team != self.Team && self.overlapsUnit(u) && u.IsGrounded && u.IsAlive() {
			u.VelocityY = -220
			u.AffectHealth(-self.Values.Damage, false, self.Owner)
		}
	}
}
func splashUpdateWatershield(self *Splash) {
	if self.Owner != nil {
		var hb = self.Owner.Hitbox()
		self.X, self.Y = hb.X, hb.Y

		if self.Values.Timer < 0.2 {
			self.Owner.RemoveEffect(effect.Aquaweaver)
		} else {
			self.Owner.AddEffect(effect.Aquaweaver)
		}
	}
}
func splashTickLightingZap(self *Splash) {
	var offsetX float32 = TileSize * 1.55
	if self.Team == team.Enemy {
		offsetX *= -1
	}
	self.X += offsetX
	self.Y -= 6
	self.tryDamageEnemy()
}

func (s *Splash) tryDamageEnemy() {
	for _, u := range Units {
		if u.Team != s.Team && s.overlapsUnit(u) && u.IsAlive() {
			u.AffectHealth(-s.Values.Damage, false, s.Owner)
		}
	}
}
