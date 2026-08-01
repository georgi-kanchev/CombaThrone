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
	Behavior  func(self *Unit)
	Anim      *motion.Animation[assets.ImageId]
	HealthBar HealthBar
	State     State

	VelocityX, VelocityY float32
	IsGrounded           bool

	UnitFront, UnitBehind *Unit

	lastX, lastY, moveSpeedX float32
	attackTimer, hurtTimer   float32 // negative values can be used for "time since last"
}

type Team uint8
type Duty uint8
type Character uint8
type State uint8

const StateIdle, StateWalk, StateAttackStart, StateAttacking, StateAttackEnd, StateHurt, StateDead = 0, 1, 2, 3, 4, 5, 6
const TeamAlly, TeamEnemy, TeamNeutral Team = 0, 1, 2
const DutyLower, DutyMiddle, DutyUpper, DutyGarrison Duty = 0, 1, 2, 3
const Gravity = 256
const GroundFrictionPercent = 5.0

func NewUnit(character Character, team Team, duty Duty) *Unit {
	var char = Characters[character]
	var anim = motion.NewAnimation(0, false, char.Animations.Idle...)
	var unit = Unit{Object: graphics.NewSprite(-2000, -2000, 1, 0), Character: character, Team: team, Duty: duty,
		Behavior: char.Brain, Stats: char.Stats, Anim: &anim, attackTimer: number.NaN(), hurtTimer: number.NaN()}

	unit.draw() // update frame size

	var lane = Collisions[duty]
	switch duty {
	case DutyLower:
		unit.X, unit.Y = lane[0].X+lane[0].Width/2-56, lane[0].Y-lane[0].Height/2-unit.Height/2
	case DutyMiddle:
		unit.X, unit.Y = lane[0].X+lane[0].Width/2-88, lane[0].Y-lane[0].Height/2-unit.Height/2
	case DutyUpper:
		unit.X, unit.Y = lane[0].X+lane[0].Width/2-120, lane[0].Y-lane[0].Height/2-unit.Height/2
	}
	if team == TeamAlly {
		unit.X = -unit.X
	}

	var hb = unit.Hitbox()
	unit.HealthBar = NewHealthBar(hb.Width-1, team)
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
	u.hurtTimer -= time.Delta()
	u.attackTimer -= time.Delta()
	u.Mask = Masks[u.Duty] // applied every frame to account for any changes in duty

	u.applyState()
	u.actUponState()
	u.applyPhysics()
	u.applyCollisions()
	u.Behavior(u)
	u.draw()

	var speedX = number.Absolute(u.X-u.lastX) / time.Delta() // smooth out for FPS dips
	u.moveSpeedX = u.moveSpeedX + (speedX-u.moveSpeedX)*0.15 // 0.15 = how fast it catches up
	u.lastX, u.lastY = u.X, u.Y
}
func (u *Unit) TakeDamage(damage int) {
	if u.Stats.Health <= 0 {
		return
	}

	u.Stats.Health -= damage
	u.hurtTimer = u.Stats.HurtTime

	u.Anim.Frames = Characters[u.Character].Animations.Hurt
	u.Anim.IsLooping, u.Anim.FPS, u.Anim.Time = false, 4, 0

	if u.Stats.Health <= 0 {
		u.Anim.Frames = Characters[u.Character].Animations.Die
		u.Anim.IsLooping, u.Anim.FPS = false, 8
		u.HealthBar.FadeOut(1.5)
	}
}

//=================================================================

func (u *Unit) applyState() {
	if u.State == StateWalk && (!u.IsGrounded || u.moveSpeedX < 0.01) {
		u.State = StateIdle
	}
	if u.State == StateAttackEnd || (u.IsGrounded && u.moveSpeedX > 0.01 && u.UnitFront == nil) {
		u.State = StateWalk
	}
	if u.State == StateAttacking && u.Anim.IsFinished() {
		u.State = StateAttackEnd
	}
	if u.State == StateAttackStart {
		u.State = StateAttacking
	}
	if u.UnitFront != nil && (u.attackTimer < 0 || number.IsNaN(u.attackTimer)) {
		u.State = StateAttackStart
	}
	if u.hurtTimer > 0 {
		u.State = StateHurt
	}
	if u.Stats.Health <= 0 {
		u.State = StateDead
	}
}
func (u *Unit) actUponState() {
	switch u.State {
	case StateIdle:
		u.Anim.Frames = Characters[u.Character].Animations.Idle
		u.Anim.IsLooping, u.Anim.FPS = true, 3
	case StateWalk:
		u.Anim.Frames = Characters[u.Character].Animations.Walk
		u.Anim.IsLooping, u.Anim.FPS = true, u.moveSpeedX*0.25

		switch u.Team {
		case TeamAlly:
			u.VelocityX = float32(u.Stats.MoveSpeed)
		case TeamEnemy:
			u.VelocityX = -float32(u.Stats.MoveSpeed)
		}
	case StateAttackStart:
		u.attackTimer = float32(u.Stats.AttackSpeed) / 10
		u.Anim.Frames = Characters[u.Character].Animations.AttackStart
		u.Anim.IsLooping, u.Anim.FPS, u.Anim.Time = false, 8, 0
	case StateAttacking:
		if u.Anim.IsFinished() {
			u.Anim.Frames = Characters[u.Character].Animations.AttackEnd
			u.Anim.IsLooping, u.Anim.FPS, u.Anim.Time = false, 8, 0
		}
	case StateAttackEnd:
		if u.UnitFront != nil {
			u.UnitFront.TakeDamage(u.Stats.AttackDamage)
		}
	case StateHurt, StateDead: // empty - just prevents running any other code, animations play once upon taking damage
	}
}

func (u *Unit) applyPhysics() {
	if u.IsGrounded {
		u.VelocityX *= 1.0 - (GroundFrictionPercent / 100.0)
	}
	u.VelocityY += Gravity * time.Delta()
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
		var anyoneDead = u.Stats.Health <= 0 || other.Stats.Health <= 0
		if other == u || u.Duty != other.Duty || anyoneDead || !hb.Overlaps(ohb) {
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
func (u *Unit) draw() {
	var frame = u.Anim.Frame()
	var crop = frame.CropArea()

	u.ImageId, u.Width, u.Height = frame, crop.Width, crop.Height
	if u.Team == TeamEnemy {
		u.Width = -crop.Width
	}
	View.DrawObject(&u.Object)
	u.Width = crop.Width
}
