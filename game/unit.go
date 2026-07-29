package game

// The alive unit in the scene. A unit can be any character, copying its base data at different times,
// then acting upon it and editing it through its behavior (brain function).

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/time"
)

type Unit struct {
	graphics.Object
	Stats     Stats
	Character Character
	Duty      Duty
	Team      Team
	Brain     func(self *Unit)
	Anim      *motion.Animation[assets.ImageId]
	HealthBar HealthBar

	VelocityX, VelocityY float32
	IsGrounded           bool

	UnitFront, UnitBehind *Unit

	prevX, prevY, currentSpeed float32
}

type Team uint8
type Duty uint8
type Character uint8

const TeamAlly, TeamEnemy, TeamNeutral Team = 0, 1, 2
const DutyLower, DutyMiddle, DutyUpper, DutyGarrison Duty = 0, 1, 2, 3
const Gravity = 256

func NewUnit(character Character, team Team, duty Duty) *Unit {
	var char = Characters[character]
	var anim = motion.NewAnimation(0, false, char.Animations.Idle...)
	var unit = Unit{Object: graphics.NewSprite(0, 0, 1, 0), Character: character, Team: team, Duty: duty,
		Brain: char.Brain, Stats: char.Stats, Anim: &anim}

	unit.applyAnimations()

	var lane = Collisions[duty]
	switch duty {
	case DutyLower:
		unit.X, unit.Y = lane[0].X+lane[0].Width/2-48, lane[0].Y-lane[0].Height/2-unit.Height/2
	case DutyMiddle:
		unit.X, unit.Y = lane[0].X+lane[0].Width/2-80, lane[0].Y-lane[0].Height/2-unit.Height/2
	case DutyUpper:
		unit.X, unit.Y = lane[0].X+lane[0].Width/2-112, lane[0].Y-lane[0].Height/2-unit.Height/2
	}
	if team == TeamAlly {
		unit.X = -unit.X
	}

	var hb = unit.Hitbox()
	unit.HealthBar = NewHealthBar(hb.Width-1, false, team)
	return &unit
}

//=================================================================

func (u *Unit) Hitbox() geometry.Shape {
	var char = Characters[u.Character]
	var hitbox = char.Hitbox
	hitbox.X, hitbox.Y = u.X+hitbox.X, u.Y+hitbox.Y
	return hitbox
}
func (u *Unit) AttackPoint() (x, y float32) {
	var hb = u.Hitbox()
	if u.Team == TeamAlly {
		return hb.X + hb.Width, hb.Y
	}
	if u.Team == TeamEnemy {
		return hb.X - hb.Width, hb.Y
	}
	return hb.X, hb.Y
}

//=================================================================

func (u *Unit) Update() {
	u.Mask = Masks[u.Duty] // applied every frame to account for any changes in duty
	u.applyPhysics()
	u.applyCollisions()
	u.Brain(u)
	u.applyAnimations()

	var curHorSpeed = number.Absolute(u.X-u.prevX) / time.Delta()       // smooth out for FPS dips
	u.currentSpeed = u.currentSpeed + (curHorSpeed-u.currentSpeed)*0.15 // 0.15 = how fast it catches up
	u.prevX, u.prevY = u.X, u.Y
}

func (u *Unit) applyPhysics() {
	u.VelocityY += Gravity * time.Delta()

	if u.IsGrounded && u.Team == TeamAlly {
		u.VelocityX = float32(u.Stats.MoveSpeed)
	} else if u.IsGrounded && u.Team == TeamEnemy {
		u.VelocityX = -float32(u.Stats.MoveSpeed)
	}
	u.X, u.Y = u.X+u.VelocityX*time.Delta(), u.Y+u.VelocityY*time.Delta()
}
func (u *Unit) applyCollisions() {
	var hb = u.Hitbox()
	var diffX, diffY = u.X - hb.X, u.Y - hb.Y // cache hitbox and obj offset

	u.IsGrounded = false
	if u.VelocityY > 0 { // collide with ground only when falling down (allows jumping up to a lane/other duty)
		for _, s := range Collisions[u.Duty] {
			if hb.Overlaps(s) {
				hb = hb.Collide(s)
				u.X, u.Y = hb.X+diffX, hb.Y+diffY
				u.VelocityY = 0
				u.IsGrounded = true
			}
		}
	}

	u.UnitBehind, u.UnitFront = nil, nil
	for _, other := range Units {
		var ohb = other.Hitbox()
		if other == u || u.Duty != other.Duty || !hb.Overlaps(ohb) {
			continue
		}
		hb = hb.Collide(ohb)
		u.X, u.Y = hb.X+diffX, hb.Y+diffY
		if (u.Team == TeamAlly && u.X < other.X) || (u.Team == TeamEnemy && u.X > other.X) {
			u.UnitFront = other
		} else if (u.Team == TeamAlly && u.X > other.X) || (u.Team == TeamEnemy && u.X < other.X) {
			u.UnitBehind = other
		}
	}
}
func (u *Unit) applyAnimations() {
	if !u.IsGrounded || u.X == u.prevX {
		u.Anim.Frames = Characters[u.Character].Animations.Idle
		u.Anim.IsLooping, u.Anim.FPS = true, 3
	} else if u.IsGrounded && u.X != u.prevX {
		u.Anim.Frames = Characters[u.Character].Animations.Walk
		u.Anim.IsLooping, u.Anim.FPS = true, u.currentSpeed*0.25
	}

	var frame = u.Anim.Frame()
	var crop = frame.CropArea()

	u.ImageId, u.Width, u.Height = frame, crop.Width, crop.Height
	if u.Team == TeamEnemy {
		u.Width = -crop.Width
	}
	View.DrawObject(&u.Object)
	u.Width = crop.Width
}
