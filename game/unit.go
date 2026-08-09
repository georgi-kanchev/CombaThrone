package game

// The alive unit in the scene. A unit can be any character, copying its base data at different times,
// then acting upon it and editing it through its behavior (brain function).

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/number"
)

type Unit struct {
	graphics.Object
	Stats     Stats
	Character Character
	Lane      Lane
	Team      Team
	Behavior  func(self *Unit)
	Anim      *motion.Animation[assets.ImageId]
	HealthBar HealthBar
	State     State

	VelocityX, VelocityY, Z float32
	IsGrounded              bool

	UnitFront, UnitBehind, ClosestEnemyInRange *Unit

	lastX, lastY, moveSpeedX float32
	attackTimer, hurtTimer   float32 // negative values can be used for "time since last"
}

type Team uint8
type Lane uint8
type Character uint8
type State uint8

const ( // states that end on 'ing' are continuous, single frame otherwise
	StateIdling State = iota
	StateWalking

	StateHurtStart
	StateHurting

	StateDyingStart
	StateDying
	StateDyingEnd
	StateDecaying

	StateAttackStart
	StateAttackCharging
	StateAttackTrigger
	StateAttackRecovering
	StateAttackEnd
)

const TeamAlly, TeamEnemy, TeamNeutral Team = 0, 1, 2
const Gravity, GroundFrictionPercent, DeathFadeOutTime = 256, 10.0, 30.0

func NewUnit(character Character, team Team, lane Lane) *Unit {
	var char = Characters[character]
	var anim = motion.NewAnimation(0, false, char.Animations.Idle...)
	var unit = Unit{Object: graphics.NewSprite(-2000, -2000, 1, 0), Character: character, Team: team, Lane: lane,
		Behavior: char.Brain, Stats: char.Stats, Anim: &anim, attackTimer: number.NaN(), hurtTimer: number.NaN()}

	unit.draw() // update frame size

	var col = Collisions[lane]
	switch lane {
	case LaneLower:
		unit.X, unit.Y = col[0].X+col[0].Width/2-40, col[0].Y-col[0].Height/2-unit.Height/2
	case LaneMiddle:
		unit.X, unit.Y = col[0].X+col[0].Width/2-72, col[0].Y-col[0].Height/2-unit.Height/2
	case LaneUpper:
		unit.X, unit.Y = col[0].X+col[0].Width/2-104, col[0].Y-col[0].Height/2-unit.Height/2
	}
	if team == TeamAlly {
		unit.X = -unit.X
	}

	var hb = unit.Hitbox()
	unit.HealthBar = NewHealthBar(hb.Width-1, team)
	return &unit
}

func (u *Unit) Hitbox() geometry.Shape {
	var char = Characters[u.Character]
	var hitbox = char.Hitbox
	hitbox.X, hitbox.Y = u.X+hitbox.X, u.Y+hitbox.Y
	return hitbox
}
func (u *Unit) EnemyEntrance() (attackable bool, entrance *EntranceData) {
	var e *EntranceData
	if u.Team != TeamNeutral && (u.Lane == LaneUpper || u.Lane == LaneMiddle || u.Lane == LaneLower) {
		e = Entrances[int(3*(1-u.Team))+int(u.Lane)]
		var attackRange = float32(u.Stats.AttackRange) * TileSize
		var melee = u.Stats.AttackRange == 1 && number.IsWithin(u.X, e.Tiles[0].X, TileSize/2)
		var ranged = u.Stats.AttackRange > 1 && number.Absolute(u.X-e.Tiles[0].X) < attackRange
		attackable = !e.IsOpen() && e.Health > 0 && (melee || ranged)
	}
	return attackable, e
}

func (u *Unit) Update() {
	u.Anim.TimeScale = TimeScale
	u.hurtTimer -= DeltaTimeScaled()
	u.attackTimer -= DeltaTimeScaled()
	u.Mask = Masks[u.Lane] // applied every frame to account for any changes in lane

	u.Z = map[Lane]float32{LaneUpper: 2, LaneMiddle: 1, LaneLower: 0}[u.Lane]

	u.applyState()
	u.actUponState()
	u.applyPhysics()
	u.applyCollisions()
	u.Behavior(u)
	u.draw()

	var speedX = number.Absolute(u.X-u.lastX) / DeltaTimeScaled() // smooth out for FPS dips
	u.moveSpeedX = u.moveSpeedX + (speedX-u.moveSpeedX)*0.15      // 0.15 = how fast it catches up
	if number.IsNaN(u.moveSpeedX) {
		u.moveSpeedX = 0
	}
	u.lastX, u.lastY = u.X, u.Y
}
func (u *Unit) TakeDamage(damage int) {
	if u.Stats.Health > 0 {
		u.Stats.Health -= damage
		u.hurtTimer = u.Stats.HurtTime
	}
}

func (u *Unit) applyState() {
	var attackable, entrance = u.EnemyEntrance()
	var canAttack = u.attackTimer < 0 || number.IsNaN(u.attackTimer)
	var enemyEntranceInRange = attackable && entrance != nil
	var hasMeleeTarget = u.UnitFront != nil && u.Team != u.UnitFront.Team
	var melee = canAttack && (hasMeleeTarget || enemyEntranceInRange) && u.Stats.AttackRange == 1

	var closestDistX = number.ValueBiggest[float32]()
	var attackRange = float32(u.Stats.AttackRange) * TileSize
	u.ClosestEnemyInRange = nil
	for _, t := range Units {
		if u == t || t.Stats.Health <= 0 {
			continue
		}
		var distX = number.Absolute(t.X - u.X)
		var allyEnemy, enemyAlly = u.Team == TeamAlly && t.Team == TeamEnemy, u.Team == TeamEnemy && t.Team == TeamAlly
		var isEnemy = allyEnemy || enemyAlly
		var isInFront = (allyEnemy && u.X < t.X) || (enemyAlly && u.X > t.X)
		var closeEnough = distX < attackRange
		if isInFront && isEnemy && distX < closestDistX && closeEnough {
			closestDistX = distX
			u.ClosestEnemyInRange = t
		}
	}
	var ranged = u.Stats.AttackRange > 1 && (u.ClosestEnemyInRange != nil || enemyEntranceInRange)

	if u.State == StateWalking && u.Stats.Health > 0 && (!u.IsGrounded || u.moveSpeedX < 0.01) {
		u.State = StateIdling
	} else if u.State == StateIdling && u.UnitFront == nil && !attackable && u.IsGrounded && u.Stats.Health > 0 {
		u.State = StateWalking
	}

	if u.State == StateAttackEnd && u.Stats.Health > 0 {
		u.State = StateIdling
	} else if u.State == StateAttackRecovering && u.Anim.IsJustFinished() {
		u.State = StateAttackEnd
	} else if u.State == StateAttackTrigger {
		u.State = StateAttackRecovering
	} else if u.State == StateAttackCharging && u.Anim.IsJustFinished() {
		u.State = StateAttackTrigger
	} else if u.State == StateAttackStart {
		u.State = StateAttackCharging
	} else if (u.State == StateIdling || u.State == StateWalking) && melee {
		u.State = StateAttackStart
	} else if (u.State == StateIdling || u.State == StateWalking) && ranged {
		if canAttack {
			u.State = StateAttackStart
		} else if u.Stats.Health > 0 {
			u.State = StateIdling // enemy in range but waiting for attack timer (stay in one place, don't keep walking)
		}
	}

	if u.State == StateHurting && u.Stats.Health <= 0 {
		u.State = StateDyingStart // bug fix for units sometimes staying alive
	} else if u.State == StateHurting && u.hurtTimer < 0 && u.Stats.Health > 0 {
		u.State = StateIdling
	} else if u.State == StateHurtStart {
		u.State = StateHurting
	}
	if u.State != StateDyingStart && u.State != StateDying && u.State != StateDecaying &&
		u.State != StateHurting && u.hurtTimer > 0 {
		u.State = StateHurtStart // can interupt other states
	}

	if u.State == StateDyingEnd {
		u.State = StateDecaying
	}
	if u.State == StateDying && u.Anim.IsJustFinished() {
		u.State = StateDyingEnd
	} else if u.State == StateDyingStart {
		u.State = StateDying
	} else if u.State == StateHurtStart && u.Stats.Health <= 0 {
		u.State = StateDyingStart
	}
}
func (u *Unit) actUponState() {
	switch u.State {
	case StateIdling:
		u.Anim.Frames = Characters[u.Character].Animations.Idle
		u.Anim.IsLooping, u.Anim.FPS = true, 3
	case StateWalking:
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
	case StateAttackTrigger:
		u.Anim.Frames = Characters[u.Character].Animations.AttackEnd
		u.Anim.IsLooping, u.Anim.FPS, u.Anim.Time = false, 8, 0

		var dmg = u.Stats.AttackDamage
		var attackable, entrance = u.EnemyEntrance()
		if u.Stats.AttackRange == 1 {
			if u.UnitFront != nil {
				u.UnitFront.TakeDamage(dmg)
			} else if attackable && entrance != nil {
				entrance.TakeDamage(dmg)
			}
			break
		}
		var t = u.ClosestEnemyInRange
		if t != nil {
			Projectiles = append(Projectiles, NewProjectile(u.X, u.Y, u.Z, t.X, t.Y+t.Height/2-8, t.Z, dmg, entrance))
		} else if attackable && entrance != nil {
			var x, y = entrance.Tiles[0].X, entrance.Tiles[0].Y
			var z = map[Lane]float32{LaneLower: 0, LaneMiddle: 1, LaneUpper: 2}[entrance.Lane]
			Projectiles = append(Projectiles, NewProjectile(u.X, u.Y, u.Z, x, y, z, dmg, entrance))
		}
	case StateAttackCharging, StateAttackRecovering, StateAttackEnd: // empty
	case StateHurtStart:
		u.Anim.Frames = Characters[u.Character].Animations.Hurt
		u.Anim.IsLooping, u.Anim.FPS, u.Anim.Time = false, 5, 0
	case StateHurting: // empty
	case StateDyingStart:
		u.Anim.Frames = Characters[u.Character].Animations.Die
		u.Anim.IsLooping, u.Anim.FPS, u.Anim.Time = false, 8, 0
		u.HealthBar.FadeOut(1.5)
	case StateDying, StateDyingEnd: // empty
	case StateDecaying:
		if u.hurtTimer < -DeathFadeOutTime {
			Units = collection.Remove(Units, u)
		} else if u.hurtTimer < 0 {
			u.Effects.Tint = color.RGBA(255, 255, 255, byte(number.Map(u.hurtTimer, 0, -DeathFadeOutTime, 255, 0)))
		}
	}
}

func (u *Unit) applyPhysics() {
	if u.IsGrounded {
		u.VelocityX *= 1.0 - (GroundFrictionPercent / 100.0)
	}
	var attack, entry = u.EnemyEntrance()
	if attack && entry != nil {
		u.VelocityX = 0
	}

	u.VelocityY += Gravity * DeltaTimeScaled()
	u.X, u.Y = u.X+u.VelocityX*DeltaTimeScaled(), u.Y+u.VelocityY*DeltaTimeScaled()
}
func (u *Unit) applyCollisions() {
	var hb = u.Hitbox()
	var diffX, diffY = u.X - hb.X, u.Y - hb.Y // cache hitbox and obj offset

	u.IsGrounded = false
	if u.VelocityY > 0 { // collide with ground only when falling down (allows jumping up to a lane)
		for _, s := range Collisions[u.Lane] {
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
		if other == u || u.Lane != other.Lane || anyoneDead || !hb.Overlaps(ohb) {
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

	if u.Stats.Health > 0 {
		DrawShadow(u.X, u.Z, u.Width*0.6, u.Height*0.1, 0, u.Mask)
	}

	if u.Team == TeamEnemy {
		u.Width = -crop.Width
	}
	View.DrawObject(&u.Object)
	u.Width = crop.Width
}
