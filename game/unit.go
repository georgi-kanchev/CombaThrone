package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/random"
)

type Team uint8
type Lane uint8
type Character uint8
type State uint8
type Unit struct {
	graphics.Object
	Stats     Stats
	Character Character
	Lane      Lane
	Team      Team
	Behavior  func(self *Unit)
	Anim      *motion.Animation[assets.ImageId]
	HealthBar *HealthBar
	State     State

	Blood *motion.ParticleSystem

	VelocityX, VelocityY, Z float32

	IsGrounded   bool
	IsReturning  bool // for OffLaners only
	IsAtGarrison bool // cannot act before garrison position - once there, can never move again (only act)

	UnitFront, UnitBehind, ClosestEnemyInRange *Unit

	lastX, lastY, moveSpeedX float32
	actionTimer, hurtTimer   float32 // negative values can be used for "time since last"
}

const ( // states
	StateIdling  State = iota // continuous
	StateWalking              // continuous

	StateHurtStart // single frame
	StateHurting   // continuous

	StateDyingStart // single frame
	StateDying      // continuous
	StateDyingEnd   // single frame
	StateDecaying   // continuous

	StateActionStart      // single frame
	StateActionCharging   // continuous
	StateActionTrigger    // single frame
	StateActionRecovering // continuous
	StateActionEnd        // single frame
)

const TeamAlly, TeamEnemy, TeamNeutral Team = 0, 1, 2

const Gravity, GroundFrictionPercent, DeathFadeOutTime = 256, 15.0, 30.0

var Units []*Unit = make([]*Unit, 0, 16)

func NewUnit(character Character, team Team, lane Lane) *Unit {
	var char = Characters[character]
	var anim = motion.NewAnimation(0, false, char.Animations.Idle...)
	var unit = Unit{Object: graphics.NewSprite(-2000, -2000, 1, 0), Character: character, Team: team, Lane: lane,
		Behavior: char.Behavior, Stats: char.Stats, Anim: &anim, actionTimer: number.NaN(), hurtTimer: number.NaN()}

	unit.Blood = motion.NewParticleSystem(unit.particlesBlood)

	unit.draw() // update frame size

	var col = laneCollisions[lane]
	var laneX, laneY = col[0].X + col[0].Width/2, col[0].Y - col[0].Height/2 - unit.Height/2
	switch lane {
	case LaneLower, LaneLowerOff:
		unit.X, unit.Y = laneX-40, laneY
	case LaneMiddle, LaneMiddleOff:
		unit.X, unit.Y = laneX-72, laneY
	case LaneUpper, LaneUpperOff:
		unit.X, unit.Y = laneX-104, laneY
	case LaneGarrison1, LaneGarrison2, LaneGarrison3, LaneGarrison4, LaneGarrison5:
		unit.X, unit.Y = CurrentZone.Background.Width/2+unit.Width/2, laneY
	case LaneGarrisonPlus1, LaneGarrisonPlus2, LaneGarrisonPlus3, LaneGarrisonPlus4, LaneGarrisonPlus5:
		unit.X, unit.Y = PointAtCell(18, 3)
	}
	if unit.IsOffLaner() {
		unit.Y += TileSize / 2
	}

	if team == TeamAlly {
		if AllyBase.Kind < BaseBarrack {
			unit.X = CurrentZone.Background.Width/2 + unit.Width/2
		}

		unit.X = -unit.X
	} else if EnemyBase.Kind < BaseBarrack {
		unit.X = CurrentZone.Background.Width / 2
	}

	var hb = unit.Hitbox()
	unit.HealthBar = NewHealthBar(hb.Width-1, team, unit.IsOffLaner())
	return &unit
}

//=================================================================

func (u *Unit) Hitbox() geometry.Shape {
	var char = Characters[u.Character]
	var hitbox = char.Hitbox
	hitbox.X, hitbox.Y = u.X+hitbox.X, u.Y+hitbox.Y
	return hitbox
}
func (u *Unit) EnemyEntrance() (canBeActedUpon bool, entrance *Entrance) {
	var e *Entrance
	if u.Team != TeamNeutral && (u.IsLaner() || u.IsOffLaner()) {
		if u.Team == TeamAlly {
			e = EnemyBase.Entrances[u.Lane/2]
		} else {
			e = AllyBase.Entrances[u.Lane/2]
		}
		var actionRange = float32(u.Stats.ActRange) * TileSize
		var melee = u.Stats.ActRange == 1 && number.IsWithin(u.X, e.Tiles[0].X, TileSize/2)
		var ranged = u.Stats.ActRange > 1 && number.Absolute(u.X-e.Tiles[0].X) < actionRange
		canBeActedUpon = !e.IsOpen() && e.Health > 0 && (melee || ranged)
	}
	return canBeActedUpon, e
}
func (u *Unit) IsLaner() bool {
	return u.Lane == LaneLower || u.Lane == LaneMiddle || u.Lane == LaneUpper
}
func (u *Unit) IsOffLaner() bool {
	return u.Lane == LaneLowerOff || u.Lane == LaneMiddleOff || u.Lane == LaneUpperOff
}
func (u *Unit) IsGarrisoner() bool {
	return u.Lane >= LaneGarrison1
}

func (u *Unit) Update() {
	u.Anim.TimeScale = TimeScale
	u.hurtTimer -= DeltaTimeScaled()
	u.actionTimer -= DeltaTimeScaled()
	u.Z = laneZs[u.Lane]

	u.Mask = laneMasks[u.Lane] // applied every frame to account for any changes in lane
	if (u.X < 0 && AllyBase.Kind < BaseBarrack) || (u.X > 0 && EnemyBase.Kind < BaseBarrack) {
		u.Mask = geometry.Area{}
	}

	if TimeScale > 0 {
		u.applyState()
		u.actUponState()
		u.applyPhysics()
		u.applyCollisions()
		u.Behavior(u)
	}
	u.draw()

	if TimeScale > 0 {
		u.Blood.Update()

		var speedX = number.Absolute(u.X-u.lastX) / DeltaTimeScaled() // smooth out for FPS dips
		u.moveSpeedX = u.moveSpeedX + (speedX-u.moveSpeedX)*0.15      // 0.15 = how fast it catches up
		if number.IsNaN(u.moveSpeedX) {
			u.moveSpeedX = 0
		}
		u.lastX, u.lastY = u.X, u.Y
	}
}
func (u *Unit) TakeDamage(damage int) {
	if u.Stats.Health > 0 {
		u.Stats.Health -= damage
		u.hurtTimer = u.Stats.HurtTime
	}
}

// private ========================================================

var laneZs = [LaneCount]float32{
	LaneLower: 0, LaneLowerOff: 0.5, LaneMiddle: 1, LaneMiddleOff: 1.5, LaneUpper: 2, LaneUpperOff: 2.5,
	LaneGarrison1: 0, LaneGarrison2: 0.5, LaneGarrison3: 1, LaneGarrison4: 1.5, LaneGarrison5: 2,
	LaneGarrisonPlus1: 2.5, LaneGarrisonPlus2: 2.5, LaneGarrisonPlus3: 2.5, LaneGarrisonPlus4: 2.5, LaneGarrisonPlus5: 2.5,
}
var laneMasks = map[Lane]geometry.Area{
	LaneLower:     geometry.NewArea(0, 0, 556, 1000),
	LaneLowerOff:  geometry.NewArea(0, 0, 556, 1000),
	LaneMiddle:    geometry.NewArea(0, 0, 492, 1000),
	LaneMiddleOff: geometry.NewArea(0, 0, 492, 1000),
	LaneUpper:     geometry.NewArea(0, 0, 428, 1000),
	LaneUpperOff:  geometry.NewArea(0, 0, 428, 1000),
}
var laneCollisions = map[Lane][]geometry.Shape{}

func (u *Unit) particlesBlood(p *motion.Particle) (alive bool) {
	if p.Age == 0 {
		p.CustomData["offsetY"] = random.Range[float32](-4, 4)
		p.Scale = random.Range[float32](0.2, 3)
		p.VelocityX = random.Range[float32](30, 50)
		p.VelocityY = random.Range[float32](-20, 20)
		p.Color = random.Range[uint](128, 255) // only red

		if u.Team == TeamAlly {
			p.VelocityX *= -1
		}
	}

	var dt = DeltaTimeScaled()
	p.Age += dt
	p.VelocityY += Gravity / 2 * dt

	p.X += p.VelocityX * dt
	p.Y += p.VelocityY * dt

	if p.Y > u.Y+u.Height/2+p.CustomData["offsetY"].(float32) {
		p.Y = u.Y + u.Height/2 + p.CustomData["offsetY"].(float32)
		p.VelocityX, p.VelocityY = 0, 0
	}

	var alpha = number.Map(p.Age, 1.0, 1.5, 255, 0)
	var col = color.RGBA(byte(p.Color), 0, 0, byte(number.Limit(alpha, 0, 255)))
	View.DrawShape(p.X, p.Y, p.Scale, p.Scale, 0, 1, col, u.Mask)
	return p.Age < 1.5
}

func (u *Unit) applyState() {
	var canBeActedUpon, entrance = u.EnemyEntrance()
	var canAct = u.actionTimer < 0 || number.IsNaN(u.actionTimer)
	var enemyEntranceInRange = canBeActedUpon && entrance != nil
	var hasMeleeTarget = u.UnitFront != nil && u.Team != u.UnitFront.Team
	var melee = canAct && (hasMeleeTarget || enemyEntranceInRange) && u.Stats.ActRange == 1

	var closestDistX = number.ValueBiggest[float32]()
	var actionRange = float32(u.Stats.ActRange) * TileSize
	u.ClosestEnemyInRange = nil
	for _, t := range Units {
		if u == t || t.Stats.Health <= 0 || t.IsOffLaner() {
			continue
		}
		var _, myEntrance = t.EnemyEntrance()
		if myEntrance != nil && t.Team == TeamEnemy && myEntrance.Tiles[0].X > t.X {
			continue // enemy is already inside my base - can't shoot through the wall
		} else if myEntrance != nil && t.Team == TeamAlly && myEntrance.Tiles[0].X < t.X {
			continue // enemy is already inside my base - can't shoot through the wall
		}

		var distX = number.Absolute(t.X - u.X)
		var allyEnemy, enemyAlly = u.Team == TeamAlly && t.Team == TeamEnemy, u.Team == TeamEnemy && t.Team == TeamAlly
		var isEnemy = allyEnemy || enemyAlly
		var isInFront = (allyEnemy && u.X < t.X) || (enemyAlly && u.X > t.X)
		var closeEnough = distX < actionRange
		if isInFront && isEnemy && distX < closestDistX && closeEnough {
			closestDistX = distX
			u.ClosestEnemyInRange = t
		}
	}
	var ranged = u.Stats.ActRange > 1 && (u.ClosestEnemyInRange != nil || enemyEntranceInRange)
	var garrisonOrNot = !u.IsGarrisoner() || (u.IsGarrisoner() && u.IsAtGarrison)

	if u.State == StateWalking && u.Stats.Health > 0 && (!u.IsGrounded || u.moveSpeedX < 0.01) {
		u.State = StateIdling
	} else if u.State == StateIdling && u.UnitFront == nil && !canBeActedUpon &&
		u.IsGrounded && u.Stats.Health > 0 && !u.IsAtGarrison {
		u.State = StateWalking
	}

	if u.State == StateActionEnd && u.Stats.Health > 0 {
		u.State = StateIdling
	} else if u.State == StateActionRecovering && u.Anim.IsJustFinished() {
		u.State = StateActionEnd
	} else if u.State == StateActionTrigger {
		u.State = StateActionRecovering
	} else if u.State == StateActionCharging && u.Anim.IsJustFinished() {
		u.State = StateActionTrigger
	} else if u.State == StateActionStart {
		u.State = StateActionCharging
	} else if (u.State == StateIdling || u.State == StateWalking) && melee {
		u.State = StateActionStart
	} else if (u.State == StateIdling || u.State == StateWalking) && ranged && garrisonOrNot {
		if canAct {
			u.State = StateActionStart
		} else if u.Stats.Health > 0 {
			u.State = StateIdling // enemy in range but waiting for action timer (stay in one place, don't keep walking)
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
		u.VelocityX = 0

		var arrived = u.moveSpeedX != 0 && u.moveSpeedX < 0.01
		if arrived && u.IsGarrisoner() {
			u.IsAtGarrison = true
		}
	case StateWalking:
		u.Anim.Frames = Characters[u.Character].Animations.Walk
		u.Anim.IsLooping, u.Anim.FPS = true, u.moveSpeedX*0.25

		var _, e = u.EnemyEntrance()
		switch u.Team {
		case TeamAlly:
			u.VelocityX = float32(u.Stats.MoveSpeed)
			if e != nil && u.IsOffLaner() && u.X > e.Tiles[0].X-TileSize {
				u.IsReturning = true
			}
			if e != nil && u.X > e.Tiles[0].X {
				u.HealthBar.MoveToGlory(2.5)
			}
		case TeamEnemy:
			u.VelocityX = -float32(u.Stats.MoveSpeed)
			if e != nil && u.IsOffLaner() && u.X < e.Tiles[0].X+TileSize {
				u.IsReturning = true
			}
			if e != nil && u.X < e.Tiles[0].X {
				u.HealthBar.MoveToGlory(2.5)
			}
		}
		if u.IsReturning {
			u.VelocityX = -u.VelocityX
		}

		if u.X < -CurrentZone.Background.Width/2-u.Width*2 || u.X > CurrentZone.Background.Width/2+u.Width*2 {
			u.hurtTimer = 0 // no instant delete - to have time to play glory text animation etc
			u.State = StateDecaying
		}

	case StateActionStart:
		u.actionTimer = float32(u.Stats.ActSpeed) / 10
		u.Anim.Frames = Characters[u.Character].Animations.ActionStart
		u.Anim.IsLooping, u.Anim.FPS, u.Anim.Time = false, 8, 0
		u.VelocityX = 0
		PlaySound(Characters[u.Character].Sounds.ActionStart)
	case StateActionTrigger:
		u.Anim.Frames = Characters[u.Character].Animations.ActionEnd
		u.Anim.IsLooping, u.Anim.FPS, u.Anim.Time = false, 8, 0

		var dmg = u.Stats.ActValue
		var canBeActedUpon, e = u.EnemyEntrance()
		if u.Stats.ActRange == 1 {
			if u.UnitFront != nil {
				u.UnitFront.TakeDamage(dmg)
				PlaySound(Characters[u.Character].Sounds.HitFlesh)
			} else if canBeActedUpon && e != nil {
				e.TakeDamage(dmg)
				if e.Health > 0 && e.Kind == EntranceDoor {
					PlaySound(Characters[u.Character].Sounds.HitWood)
				} else if e.Health > 0 && (e.Kind == EntranceShortGate || e.Kind == EntranceTallGate) {
					PlaySound(Characters[u.Character].Sounds.HitMetal)
				}
			}
			break
		}

		var t = u.ClosestEnemyInRange
		if t != nil && !t.IsOffLaner() {
			var prediction = t.VelocityX
			if number.Absolute(t.X-u.X) < TileSize*3 {
				prediction = 0 // target is too close - don't predict movement to not shoot behind self
			}
			var proj = u.NewProjectile(u.X, u.Y, u.Z, t.X+prediction, t.Y+t.Height/2-8, t.Z, dmg, ProjectileArrow, nil)
			Projectiles = append(Projectiles, proj)
			PlaySound(Characters[u.Character].Sounds.ActionTrigger)
		} else if canBeActedUpon && e != nil {
			var x, y = e.Tiles[0].X, e.Tiles[0].Y
			if e.Kind == EntranceTallGate {
				y += TileSize
			}
			Projectiles = append(Projectiles, u.NewProjectile(u.X, u.Y, u.Z, x, y, laneZs[e.Lane], dmg, ProjectileArrow, e))
			PlaySound(Characters[u.Character].Sounds.ActionTrigger)
		}
	case StateActionCharging, StateActionRecovering, StateActionEnd: // empty
	case StateHurtStart:
		u.Anim.Frames = Characters[u.Character].Animations.Hurt
		u.Anim.IsLooping, u.Anim.FPS, u.Anim.Time = false, 5, 0
		u.VelocityX = 0
		u.Blood.EmitFromLine(30, u.X, u.Y-6, u.X, u.Y+6)
	case StateHurting: // empty
	case StateDyingStart:
		u.Anim.Frames = Characters[u.Character].Animations.Die
		u.Anim.IsLooping, u.Anim.FPS, u.Anim.Time = false, 8, 0
		u.HealthBar.FadeOut(1.5)
		u.VelocityX = 0
		u.Blood.EmitFromLine(30, u.X, u.Y-6, u.X, u.Y+6)
	case StateDying, StateDyingEnd: // empty
	case StateDecaying:
		if u.hurtTimer < -DeathFadeOutTime || u.IsGarrisoner() {
			Units = collection.Remove(Units, u)
		} else if u.hurtTimer < 0 {
			u.Effects.Tint = color.RGBA(255, 255, 255, byte(number.Map(u.hurtTimer, 0, -DeathFadeOutTime, 255, 0)))
		}
	}
}

func (u *Unit) applyPhysics() {
	if u.IsAtGarrison {
		return // no physics for garrison position - once there, the unit is locked
	}

	if u.IsGrounded {
		u.VelocityX *= 1.0 - (GroundFrictionPercent/100)*DeltaTimeScaled()
	}
	var canBeActedUpon, entry = u.EnemyEntrance()
	if canBeActedUpon && entry != nil {
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
		for _, s := range laneCollisions[u.Lane] {
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
		var isGarrison = other.IsGarrisoner() || u.IsGarrisoner()
		var isOffLaner = other.IsOffLaner() || u.IsOffLaner()
		if other == u || u.Lane != other.Lane || anyoneDead || isGarrison || isOffLaner || !hb.Overlaps(ohb) {
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

	if u.IsOffLaner() {
		// for _, p := range Pickups {
		// if p {

		// }
		// }
	}
}
func (u *Unit) draw() {
	var frame = u.Anim.Frame()
	var crop = frame.CropArea()

	u.ImageId, u.Width, u.Height = frame, crop.Width, crop.Height

	if u.Stats.Health > 0 && !u.IsGarrisoner() {
		DrawShadow(u.X, u.Z, u.Width*0.6, u.Height*0.1, 0, u.Mask)
	}

	if u.Team == TeamEnemy {
		u.Width = -crop.Width
	}
	if u.IsReturning {
		u.Width = -u.Width
	}
	View.DrawObject(&u.Object)
	u.Width = crop.Width
}
