package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/point"
	"pure-game-kit/packages/utility/random"
)

type Team uint8
type Lane uint8
type State uint8
type Role uint8
type Unit struct {
	graphics.Object
	Health    int
	Values    Values
	Character CharacterKind
	Lane      Lane
	Team      Team
	Anim      *motion.Animation[assets.ImageId]
	HealthBar *HealthBar
	State     State

	Effects map[Effect]Values

	Blood *motion.ParticleSystem

	VelocityX, VelocityY, Z float32

	IsGrounded, IsAtWall bool
	IsReturning          bool // for OffLaners only

	UnitFront, UnitBehind, ClosestEnemyInRange *Unit

	Carrying []*Pickup

	LastX, LastY, MoveSpeedX float32
	ActTimer, HurtTimer      float32 // negative values can be used for "time since last"
	LastState                State
}

const ( // states
	StateSummoned            State = iota // single frame
	StateWaitingToBeSummoned              // continuous

	StateIdling // continuous
	StateMoving // continuous

	StateHurtStart // single frame
	StateHurting   // continuous

	StateDyingStart // single frame
	StateDying      // continuous
	StateDyingEnd   // single frame
	StateDecaying   // continuous

	StateActStart      // single frame
	StateActCharging   // continuous
	StateActTrigger    // single frame
	StateActRecovering // continuous
	StateActEnd        // single frame
)

const TeamAlly, TeamEnemy, TeamCount Team = 0, 1, 2

const RoleFighter, RoleRanger, RoleDefender, RoleGriefer, RoleSupplier, RoleCollector, RoleCount Role = 0, 1, 2, 3, 4, 5, 6

const Gravity, GroundFriction, BloodMultiplier = 256.0, 400.0, 40.0

var Units []*Unit = make([]*Unit, 0, 16)
var Collisions = map[Lane][]geometry.Shape{}
var PinnedUnit *Unit

func NewUnit(character CharacterKind, team Team, lane Lane) *Unit {
	var char = Characters[character]
	var anim = motion.NewAnimation(0, false, char.Animations.Idle...)
	var unit = Unit{Object: graphics.NewSprite(-2000, -2000, 1, 0), Character: character, Team: team, Lane: lane,
		Anim: &anim, ActTimer: number.NaN(), HurtTimer: number.NaN(), LastState: StateWaitingToBeSummoned,
		Effects: make(map[Effect]Values, 4),
	}

	if len(anim.Frames) == 0 {
		anim.Frames = char.Animations.Move
	}

	if team == TeamAlly {
		unit.State = StateWaitingToBeSummoned
	}

	unit.Blood = motion.NewParticleSystem(unit.particlesBlood)

	unit.draw() // update frame size
	unit.PrepareSpawn()
	return &unit
}

//=================================================================

func (u *Unit) Hitbox(additionalWidth ...float32) geometry.Shape {
	var char = Characters[u.Character]
	var hitbox = char.Hitbox
	if u.IsFacingLeft() {
		hitbox.X *= -1
	}
	hitbox.X, hitbox.Y = u.X+hitbox.X, u.Y+hitbox.Y
	if len(additionalWidth) == 1 {
		hitbox.Width += additionalWidth[0]
	}
	return hitbox
}
func (u *Unit) MyEntrance() (entrance *Entrance) {
	return Bases[u.Team].Entrances[u.Lane/2]
}
func (u *Unit) EnemyEntrance() (canBeActedUpon bool, entrance *Entrance) {
	var e *Entrance
	if u.IsLaner() || u.IsOffLaner() {
		e = Bases[1-u.Team].Entrances[u.Lane/2]
		var actRange = float32(u.Values.ActRange) * TileSize
		var melee = u.Values.ActRange == 1 && number.IsWithin(u.X, e.Tiles[0].X, TileSize/2)
		var ranged = u.Values.ActRange > 1 && number.Absolute(u.X-e.Tiles[0].X) < actRange
		canBeActedUpon = !e.IsOpen() && e.Health > 0 && (melee || ranged)
	}
	return canBeActedUpon, e
}
func (u *Unit) IsLaner() bool {
	return u.Lane == LaneLower || u.Lane == LaneMiddle || u.Lane == LaneUpper
}
func (u *Unit) IsOffLaner() bool {
	return u.Values.Role == RoleCollector || u.Values.Role == RoleGriefer || u.Values.Role == RoleSupplier
}
func (u *Unit) IsGarrisoner() bool {
	return u.Lane >= LaneGarrison1
}
func (u *Unit) IsSummoned() bool {
	return u != nil && u.State != StateWaitingToBeSummoned
}
func (u *Unit) IsOutsideOwnBase() bool {
	if u.IsGarrisoner() {
		return false
	}
	var myEntrance = Bases[u.Team].Entrances[u.Lane/2]
	if u.Team == TeamAlly && u.X > myEntrance.Tiles[0].X {
		return true
	} else if u.Team == TeamEnemy && u.X < myEntrance.Tiles[0].X {
		return true
	}
	return false
}
func (u *Unit) IsInsideEnemyBase(offset float32) bool {
	var _, entrance = u.EnemyEntrance()
	if entrance != nil && u.Team == TeamEnemy && u.X < entrance.Tiles[0].X-offset {
		return true
	} else if entrance != nil && u.Team == TeamAlly && u.X > entrance.Tiles[0].X+offset {
		return true
	}
	return false
}
func (u *Unit) IsOnScreen() bool {
	var bounds = View.Bounds()
	bounds.Width -= u.Width / 2
	bounds.Height -= u.Height / 2
	return bounds.ContainsPoint(u.X, u.Y)
}
func (u *Unit) IsAlive() bool {
	return u.Values.MaxHealth == 0 || u.Health > 0
}
func (u *Unit) IsFacingLeft() bool {
	return (u.Team == TeamEnemy && !u.IsReturning) || (u.Team == TeamAlly && u.IsReturning)
}

func (u *Unit) PrepareSpawn() {
	u.ActTimer, u.HurtTimer = number.NaN(), number.NaN()
	u.Details.Tint = palette.White
	u.Values = Characters[u.Character].Values
	u.IsReturning, u.Health = false, u.Values.MaxHealth
	u.VelocityX, u.VelocityY = 0, 0
	clear(u.Effects)

	if u.IsGarrisoner() {
		u.AddEffect(EffectMoreRangeOnGarrison)
	}

	var col = Collisions[u.Lane]
	var hb = u.Hitbox()
	var laneY = col[0].Y - col[0].Height/2 - hb.Height
	switch u.Lane {
	case LaneLower, LaneLowerOff:
		u.X, u.Y = TileSize*9.5, laneY
	case LaneMiddle, LaneMiddleOff:
		u.X, u.Y = TileSize*8.5, laneY
	case LaneUpper, LaneUpperOff:
		u.X, u.Y = TileSize*7.5, laneY
	case LaneGarrison1, LaneGarrison2, LaneGarrison3:
		u.X, u.Y = CurrentZone.Ground.Width/2+u.Width/2, laneY
	case LaneGarrisonPlus1, LaneGarrisonPlus2, LaneGarrisonPlus3:
		u.X, u.Y = PointAtCell(18.5, 3)
	}
	if u.IsOffLaner() {
		u.Y += TileSize / 2
	}

	if u.Team == TeamAlly {
		if Bases[u.Team].Kind < BaseBarrack {
			u.X = CurrentZone.Ground.Width/2 + u.Width/2
		}
		u.X = -u.X
	} else if Bases[u.Team].Kind < BaseBarrack {
		u.X = CurrentZone.Ground.Width / 2
	}

	u.HealthBar = NewHealthBar(hb.Width-1, u.Team, u.IsOffLaner())
}
func (u *Unit) Update() {
	if u == nil {
		return
	}

	u.ActTimer -= DeltaTimeScaled()
	u.HurtTimer -= DeltaTimeScaled()
	u.Values.SleepTimer -= DeltaTimeScaled()

	if !u.IsSummoned() && (number.IsNaN(u.HurtTimer) || u.HurtTimer < -u.Values.ReviveTimer) {
		u.State = StateWaitingToBeSummoned
		return
	}

	u.Anim.TimeScale = TimeScale
	u.Z = laneZs[u.Lane]

	u.Mask = laneMasks[u.Lane] // applied every frame to account for any changes in lane
	if (u.X < 0 && Bases[TeamAlly].Kind < BaseBarrack) || (u.X > 0 && Bases[TeamEnemy].Kind < BaseBarrack) {
		u.Mask = geometry.Area{}
	}

	if TimeScale > 0 {
		u.updateEffects()
		u.applyState()
		var behavior = Behaviors[u.Character]
		if behavior != nil && u.State != StateDecaying {
			behavior(u)
		}
		u.actUponState()
		u.applyPhysics()
		u.applyCollisions()
	}
	u.draw()
	for i, p := range u.Carrying {
		p.Update()
		var hb = u.Hitbox()
		var offsetX = hb.Width/2 + TileSize/3
		if u.IsFacingLeft() {
			offsetX *= -1
		}

		p.Mask = u.Mask
		p.X, p.Y = point.MoveToPointSmooth(p.X, p.Y, hb.X+offsetX, hb.Y-float32(i)*TileSize/2, 0.4)
	}

	if TimeScale > 0 {
		u.Blood.Update()

		var speedX = number.Absolute(u.X-u.LastX) / DeltaTimeScaled() // smooth out for FPS dips
		u.MoveSpeedX = u.MoveSpeedX + (speedX-u.MoveSpeedX)*0.15      // 0.15 = how fast it catches up
		if number.IsNaN(u.MoveSpeedX) || u.MoveSpeedX < 0.01 {
			u.MoveSpeedX = 0
		}
		u.LastX, u.LastY = u.X, u.Y
	}
	u.LastState = u.State
}

func (u *Unit) ProjectHealth(value int) (remainingHealth int) {
	if !u.IsAlive() || u.State == StateDecaying {
		return u.Health
	}
	if u.Values.Role <= RoleDefender && value > 0 {
		return min(u.Health+value, u.Values.MaxHealth)
	}
	if u.Values.Role >= RoleDefender && value < 0 {
		value = min(value+u.Values.ActPoints, 0)
	}
	return u.Health + value
}
func (u *Unit) AffectHealth(value int) {
	var newHealth = u.ProjectHealth(value)
	if newHealth != u.Health {
		u.Health = newHealth

		if value < 0 {
			u.HurtTimer = 0.5
		} else if value > 0 {
			// TODO: particles
		}
	}
}

func (u *Unit) AddEffect(effect Effect) {
	var _, has = u.Effects[effect]
	if !has {
		u.Effects[effect] = Effects[effect]
		u.TickEffect(effect)
	}
}
func (u *Unit) TickEffect(effect Effect) {
	var _, has = u.Effects[effect]
	if !has {
		return
	}
	var values = u.Effects[effect]
	u.Values.MaxHealth += values.MaxHealth
	u.Values.ReviveTimer += values.ReviveTimer
	u.Values.MoveSpeed += values.MoveSpeed
	u.Values.ActPoints += values.ActPoints
	u.Values.ActRange += values.ActRange
	u.Values.ActTimer += values.ActTimer

	u.AffectHealth(-values.EffectTickDamage)
}
func (u *Unit) RemoveEffect(effect Effect) {
	var _, has = u.Effects[effect]
	if !has {
		return
	}

	var values = u.Effects[effect]
	u.Values.MaxHealth -= values.MaxHealth
	u.Values.ReviveTimer -= values.ReviveTimer
	u.Values.MoveSpeed -= values.MoveSpeed
	u.Values.ActPoints -= values.ActPoints
	u.Values.ActRange -= values.ActRange
	u.Values.ActTimer -= values.ActTimer

	delete(u.Effects, effect)
}

// private ========================================================

var laneZs = [LaneCount]float32{
	LaneLower: 0, LaneLowerOff: 0.5, LaneMiddle: 1, LaneMiddleOff: 1.5, LaneUpper: 2, LaneUpperOff: 2.5,
	LaneGarrison1: 0, LaneGarrison2: 1, LaneGarrison3: 2,
	LaneGarrisonPlus1: 2.5, LaneGarrisonPlus2: 2.5, LaneGarrisonPlus3: 2.5,
}
var laneMasks = map[Lane]geometry.Area{
	LaneLower:     geometry.NewArea(0, 0, 556, 1000),
	LaneLowerOff:  geometry.NewArea(0, 0, 556, 1000),
	LaneMiddle:    geometry.NewArea(0, 0, 492, 1000),
	LaneMiddleOff: geometry.NewArea(0, 0, 492, 1000),
	LaneUpper:     geometry.NewArea(0, 0, 428, 1000),
	LaneUpperOff:  geometry.NewArea(0, 0, 428, 1000),
}
var teamColors = [TeamCount]uint{TeamAlly: palette.Green, TeamEnemy: palette.Red}
var effectsToRemove []Effect

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
	var canAct = u.ActTimer < 0 || number.IsNaN(u.ActTimer)
	var enemyEntranceInRange = canBeActedUpon && entrance != nil
	var hasMeleeTarget = u.UnitFront != nil && u.Team != u.UnitFront.Team
	var melee = canAct && (hasMeleeTarget || enemyEntranceInRange) && u.Values.ActRange == 1
	var canMove = !u.IsAtWall && u.IsGrounded

	var closestDistX = number.ValueBiggest[float32]()
	var actRange = float32(u.Values.ActRange) * TileSize
	u.ClosestEnemyInRange = nil
	for _, t := range Units {
		if u == t || !t.IsAlive() || t.IsOffLaner() {
			continue
		}

		var distX = number.Absolute(t.X - u.X)
		var allyEnemy, enemyAlly = u.Team == TeamAlly && t.Team == TeamEnemy, u.Team == TeamEnemy && t.Team == TeamAlly
		var isEnemy = allyEnemy || enemyAlly
		var isInFront = (allyEnemy && u.X < t.X) || (enemyAlly && u.X > t.X)
		var closeEnough = distX < actRange
		if isInFront && isEnemy && distX < closestDistX && closeEnough && !t.IsInsideEnemyBase(TileSize) {
			closestDistX = distX
			u.ClosestEnemyInRange = t
		}
	}
	var ranged = u.Values.ActRange > 1 && (u.ClosestEnemyInRange != nil || enemyEntranceInRange)
	var garrisonOrNot = !u.IsGarrisoner() || (u.IsGarrisoner() && u.IsOnScreen())
	var myEntrance *Entrance
	if u.IsLaner() {
		myEntrance = Bases[u.Team].Entrances[u.Lane/2]
	}
	var sameLaneWithTarget = u.ClosestEnemyInRange != nil && u.Lane == u.ClosestEnemyInRange.Lane
	var openDoorShoot = u.IsLaner() && !u.IsOutsideOwnBase() && myEntrance.IsOpen() && sameLaneWithTarget
	var canShoot = u.IsOutsideOwnBase() || openDoorShoot
	if u.IsGarrisoner() {
		canShoot = true
	}

	if u.State == StateMoving && u.IsAlive() && (!u.IsGrounded || u.MoveSpeedX < 0.01) {
		u.State = StateIdling
	} else if u.State == StateIdling && u.UnitFront == nil && canMove && u.IsAlive() && !canBeActedUpon {
		u.State = StateMoving
	}

	if u.State == StateSummoned && u.LastState == StateSummoned {
		u.State = StateMoving // first frame is event, second frame (now) starts walking
	}

	if u.State == StateActEnd && u.IsAlive() {
		u.State = StateIdling
	} else if u.State == StateActRecovering && u.Anim.IsJustFinished() {
		u.State = StateActEnd
	} else if u.State == StateActTrigger {
		u.State = StateActRecovering
	} else if u.State == StateActCharging && u.Anim.IsJustFinished() {
		u.State = StateActTrigger
	} else if u.State == StateActStart {
		u.State = StateActCharging
	} else if (u.State == StateIdling || u.State == StateMoving) && melee && !u.IsOffLaner() {
		u.State = StateActStart
	} else if (u.State == StateIdling || u.State == StateMoving) && ranged && garrisonOrNot && canShoot && !u.IsOffLaner() {
		if canAct {
			u.State = StateActStart
		} else if u.IsAlive() && u.IsLaner() { // no shoot-move-shoot-move for laners - but garrisoners should
			u.State = StateIdling // enemy in range but waiting for act timer (stay in one place, don't keep walking)
		}
	}

	if u.Values.SleepTimer > 0 {
		u.State = StateIdling
	}

	if u.State == StateHurting && !u.IsAlive() {
		u.State = StateDyingStart // bug fix for units sometimes staying alive
	} else if u.State == StateHurting && u.HurtTimer < 0 && u.IsAlive() {
		u.State = StateIdling
	} else if u.State == StateHurtStart {
		u.State = StateHurting
	}
	if u.State != StateDyingStart && u.State != StateDying && u.State != StateDecaying &&
		u.State != StateHurting && u.HurtTimer > 0 {
		u.State = StateHurtStart // can interupt other states
	}

	if u.State == StateDyingEnd {
		u.State = StateDecaying
	}
	if u.State == StateDying && u.Anim.IsJustFinished() {
		u.State = StateDyingEnd
	} else if u.State == StateDyingStart {
		u.State = StateDying
	} else if u.State == StateHurtStart && u.Health <= 0 {
		u.State = StateDyingStart
	}
}
func (u *Unit) actUponState() {
	switch u.State {
	case StateWaitingToBeSummoned, StateSummoned: // empty
	case StateIdling:
		u.Anim.Frames = Characters[u.Character].Animations.Idle
		u.Anim.IsLooping, u.Anim.FPS = true, 3

		if number.IsNaN(u.Values.SleepTimer) || u.Values.SleepTimer < 0 {
			u.VelocityX = 0
		}
	case StateMoving:
		u.Anim.Frames = Characters[u.Character].Animations.Move
		u.Anim.IsLooping, u.Anim.FPS = true, u.MoveSpeedX*0.25

		if u.IsLaner() && u.IsInsideEnemyBase(TileSize/1.5) {
			u.HealthBar.MoveToGlory(2.5)
		} else if u.IsOffLaner() && u.IsInsideEnemyBase(-TileSize) {
			u.IsReturning = true
		}
		if u == PinnedUnit && !View.IsAreaVisible(u.Bounds()) {
			PinnedUnit = nil
		}

		u.VelocityX = float32([TeamCount]int{u.Values.MoveSpeed, -u.Values.MoveSpeed}[u.Team])
		if u.IsAtWall { // chill out when next to a wall
			u.VelocityX = 0
		}
		if u.IsReturning {
			u.VelocityX = -u.VelocityX
		}

		if u.X < -CurrentZone.Ground.Width/2-u.Width*2 || u.X > CurrentZone.Ground.Width/2+u.Width*2 {
			u.HurtTimer = 0 // no instant delete - to have time to play glory text animation etc
			u.State = StateDecaying
		}
		if u.IsOffLaner() && u.IsReturning && !u.IsOutsideOwnBase() {
			for _, c := range u.Carrying {
				c.Mask = geometry.Area{}
				c.SlotUI = GameHUD.FreePickupSlot()
				if c.Kind == PickupCoin {
					c.SlotUI = len(GameHUD.Pickups)
				} else {
					GameHUD.PickupsCarrying--
				}
				GameHUD.Pickups = collection.Add(GameHUD.Pickups, c)
				Pickups = collection.Remove(Pickups, c)
			}
			u.Carrying = collection.Clear(u.Carrying)
		}
	case StateActStart: // random delay to balance same sided units melee VVVVVVV
		u.ActTimer = u.Values.ActTimer + random.Range[float32](0, 0.1)
		u.Anim.Frames = Characters[u.Character].Animations.Prepare
		u.Anim.IsLooping, u.Anim.FPS, u.Anim.Time = false, 8, 0
		u.VelocityX = 0
		PlaySound(Characters[u.Character].Sounds.ActStart)
	case StateActTrigger:
		u.Anim.Frames = Characters[u.Character].Animations.Recover
		u.Anim.IsLooping, u.Anim.FPS, u.Anim.Time = false, 8, 0

		if u.Values.Role == RoleDefender {
			return // defenders cannot deal damage
		}

		var dmg = max(u.Values.ActPoints, 0)
		var canBeActedUpon, e = u.EnemyEntrance()
		if u.Values.ActRange == 1 {
			if u.UnitFront != nil {
				u.UnitFront.AffectHealth(-dmg)
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

		var projectileKind = Characters[u.Character].Projectile.Kind
		var t = u.ClosestEnemyInRange
		if t != nil && !t.IsOffLaner() {
			var prediction = t.VelocityX
			if number.Absolute(t.X-u.X) < TileSize*3 {
				prediction = 0 // target is too close - don't predict movement to not shoot behind self
			}
			var offsetX = u.Width / 3
			if u.Team == TeamEnemy {
				offsetX = -u.Width / 3
			}
			var proj = u.NewProjectile(u.X+offsetX, u.Y, u.Z, t.X+prediction, t.Y+t.Height/2-8, t.Z, dmg, projectileKind, nil)
			Projectiles = append(Projectiles, proj)
			PlaySound(Characters[u.Character].Sounds.ActTrigger)
		} else if canBeActedUpon && e != nil {
			var x, y = e.Tiles[0].X, e.Tiles[0].Y
			if e.Kind == EntranceTallGate {
				y += TileSize
			}
			Projectiles = append(Projectiles, u.NewProjectile(u.X, u.Y, u.Z, x, y, laneZs[e.Lane], dmg, projectileKind, e))
			PlaySound(Characters[u.Character].Sounds.ActTrigger)
		}
	case StateActCharging, StateActRecovering, StateActEnd: // empty
	case StateHurtStart:
		u.Anim.Frames = Characters[u.Character].Animations.Hurt
		u.Anim.IsLooping, u.Anim.FPS, u.Anim.Time = false, 5, 0
		u.VelocityX = 0
		var percent = 1 - float32(u.Health)/float32(u.Values.MaxHealth)
		u.Blood.EmitFromLine(int(percent*BloodMultiplier), u.X, u.Y-6, u.X, u.Y+6)
	case StateHurting: // empty
	case StateDyingStart:
		u.Anim.Frames = Characters[u.Character].Animations.Die
		u.Anim.IsLooping, u.Anim.FPS, u.Anim.Time = false, 8, 0
		u.HealthBar.FadeOut(1.5)
		u.VelocityX = 0
		u.Blood.EmitFromLine(BloodMultiplier, u.X, u.Y-6, u.X, u.Y+6)
	case StateDying, StateDyingEnd: // empty
	case StateDecaying:
		if u.HurtTimer < -u.Values.ReviveTimer || u.IsGarrisoner() {
			Units = collection.Remove(Units, u)
			if PinnedUnit == u {
				PinnedUnit = nil
			}

			if u.Team == TeamAlly {
				u.State = StateWaitingToBeSummoned
				u.PrepareSpawn()
			}
		} else if u.HurtTimer < 0 {
			u.Details.Tint = color.RGBA(255, 255, 255, byte(number.Map(u.HurtTimer, 0, -u.Values.ReviveTimer, 255, 0)))
		}
	}
}
func (u *Unit) updateEffects() {
	if len(u.Effects) == 0 {
		return
	}

	collection.Clear(effectsToRemove)
	for k, v := range u.Effects {
		v.EffectTimer -= DeltaTimeScaled()
		u.Effects[k] = v

		if v.EffectTimer < 0 {
			effectsToRemove = collection.Add(effectsToRemove, k)
		}
	}
	for _, k := range effectsToRemove {
		u.RemoveEffect(k)
	}
}

func (u *Unit) applyPhysics() {
	if u.IsGrounded {
		u.VelocityX *= 1.0 - (GroundFriction/100)*DeltaTimeScaled()
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

	u.IsGrounded, u.IsAtWall = false, false
	if u.VelocityY > 0 { // collide with ground only when falling down (allows jumping up to a lane)
		for _, s := range Collisions[u.Lane] {
			if hb.Overlaps(s) {
				hb = hb.Collide(s)
				u.X, u.Y = hb.X+diffX, hb.Y+diffY
				u.VelocityY = 0
				u.IsGrounded = true

				if s.Height > 12 { // we have a wall
					u.VelocityX = 0
					u.IsAtWall = true
				}
			}
		}
	}

	if u.IsInsideEnemyBase(TileSize / 2) {
		return
	}

	u.UnitBehind, u.UnitFront = nil, nil
	for _, other := range Units {
		var ohb = other.Hitbox()
		var anyoneDead = !u.IsAlive() || !other.IsAlive()
		var isGarrison = other.IsGarrisoner() || u.IsGarrisoner()
		var intruder = !u.IsOutsideOwnBase() && other.IsInsideEnemyBase(TileSize/2)
		if other == u || u.Lane != other.Lane || anyoneDead || isGarrison || !hb.Overlaps(ohb) || intruder {
			continue
		}
		var opposite = (!u.IsReturning && other.IsReturning) || (u.Team != other.Team)
		if u.IsOffLaner() && other.IsOffLaner() && opposite {
			u.IsReturning = true // face to face, gotta return to base
		}

		hb = hb.Collide(ohb)
		u.X, u.Y = hb.X+diffX, hb.Y+diffY
		if (u.Team == TeamAlly && u.X < other.X) || (u.Team == TeamEnemy && u.X > other.X) {
			u.UnitFront = other
		} else if (u.Team == TeamAlly && u.X > other.X) || (u.Team == TeamEnemy && u.X < other.X) {
			u.UnitBehind = other
		}
	}

	if u.Values.Role != RoleCollector {
		return
	}

	var canCarry = len(u.Carrying) < u.Values.ActPoints
	var emptySlot = GameHUD.FreePickupSlot() >= 0
	var canFitNext = len(GameHUD.Pickups)+GameHUD.PickupsCarrying+1 <= 4
	var allowed = emptySlot && canCarry && canFitNext
	for _, p := range Pickups {
		if p != nil && p.Lane == u.Lane && p.Overlaps(hb) && (allowed || p.Kind == PickupCoin) {
			p.HasShadow = false
			u.Carrying = collection.Add(u.Carrying, p)
			Pickups = collection.Remove(Pickups, p)
			if p.Kind != PickupCoin {
				GameHUD.PickupsCarrying++
			}
			if len(u.Carrying) >= u.Values.ActPoints {
				u.IsReturning = true
			}
		}
	}
}
func (u *Unit) draw() {
	if len(u.Anim.Frames) == 0 {
		return
	}

	var frame = u.Anim.Frame()
	var crop = frame.CropArea()
	u.ImageId, u.Width, u.Height = frame, crop.Width, crop.Height

	if u.IsAlive() && !u.IsGarrisoner() {
		var hb = u.Hitbox()
		DrawShadow(hb.X, u.Z, hb.Width, hb.Height*0.1, 0, u.Mask)
	}

	if u.IsFacingLeft() {
		u.Width = -crop.Width
	}
	View.DrawObject(&u.Object)
	u.Width = crop.Width
}
