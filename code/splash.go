package game

import (
	"game/code/constants/entrance"
	"game/code/constants/splash"
	"game/code/constants/team"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/number"
)

type Splash struct {
	graphics.Object
	Kind   splash.Enum
	Anim   *motion.Animation[assets.ImageId]
	Values SplashValues
	Team   team.Enum
}
type SplashValues struct {
	AnimName         string
	Damage           int
	Loop, SingleTick bool
	Timer, FadeOutTime, FPS,
	VelocityX, VelocityY float32
	TickTimer float32 // ticks every 1 second, always ticks once upon creation (first update) since it starts at 0
	Tick      func(self *Splash)
}

var SplashBases = map[splash.Enum]SplashValues{
	splash.Dust:   SplashValues{AnimName: "dust", FPS: 10},
	splash.Pickup: SplashValues{AnimName: "pickup", FPS: 10},
	splash.Fire3: SplashValues{AnimName: "fire-3", FPS: 10, Loop: true, Timer: 5, Damage: 3, FadeOutTime: 0.25,
		Tick: splashTickFire},
	splash.Fire2: SplashValues{AnimName: "fire-2", FPS: 10, Loop: true, Timer: 5, Damage: 2, FadeOutTime: 0.25,
		Tick: splashTickFire},
	splash.Fire1: SplashValues{AnimName: "fire-1", FPS: 10, Loop: true, Timer: 5, Damage: 1, FadeOutTime: 0.25,
		Tick: splashTickFire},
	splash.Heal: SplashValues{AnimName: "heal", FPS: 10, Loop: true, Timer: 5, Damage: -2, FadeOutTime: 0.25,
		Tick: splashTickHeal},
}
var Splashes = make([]*Splash, 0, 8)

func NewSplash(kind splash.Enum, team team.Enum, x, y float32) *Splash {
	var frames = SplashAtlas.Crops(SplashBases[kind].AnimName)
	var anim = motion.NewAnimation[assets.ImageId](SplashBases[kind].FPS, SplashBases[kind].Loop, frames...)
	return &Splash{Object: graphics.NewSprite(x, y, 0, 0), Kind: kind, Team: team, Anim: &anim, Values: SplashBases[kind]}
}

//=================================================================

func (s *Splash) Update() {
	s.Values.Timer -= DeltaTimeScaled()
	s.Values.TickTimer -= DeltaTimeScaled()

	if s.Values.TickTimer < 0 {
		s.Values.TickTimer = 1

		if s.Values.Tick != nil {
			s.Values.Tick(s)
		}

		if s.Values.SingleTick {
			s.Values.TickTimer = number.Infinity()
		}
	}

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
	View.DrawObject(&s.Object)

	if (s.Values.Loop && s.Values.Timer < 0) || (!s.Values.Loop && s.Anim.IsJustFinished()) {
		Splashes = collection.Remove(Splashes, s)
	}
}

// private ========================================================

func splashTickFire(self *Splash) {
	for _, u := range Units {
		var hb = u.Hitbox()
		if u.Team != self.Team && self.ContainsPoint(hb.X, hb.Y+hb.Height/2) {
			u.AffectHealth(-self.Values.Damage, false)
		}
	}
	var entrances = Bases[1-self.Team].Entrances
	for _, e := range entrances {
		if e.Kind == entrance.Door && self.ContainsPoint(e.Tiles[0].X, e.Tiles[0].Y) {
			e.TakeDamage(self.Values.Damage)
		}
	}
}
func splashTickHeal(self *Splash) {
	for _, u := range Units {
		var hb = u.Hitbox()
		if u.Team == self.Team && self.ContainsPoint(hb.X, hb.Y+hb.Height/2) {
			u.AffectHealth(-self.Values.Damage, false)
		}
	}
}
