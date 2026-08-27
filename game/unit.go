package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/cursor"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/random"
)

type Team uint8
type Lane uint8
type State uint8
type Role uint8
type Unit struct {
	graphics.Object
	Health    int
	Stats     Stats
	Character CharacterKind
	Lane      Lane
	Team      Team
	Behavior  func(self *Unit)
	Anim      *motion.Animation[assets.ImageId]
	HealthBar *HealthBar
	State     State

	Blood *motion.ParticleSystem

	VelocityX, VelocityY, Z float32

	IsGrounded, IsAtWall bool
	IsReturning          bool // for OffLaners only

	UnitFront, UnitBehind, ClosestEnemyInRange *Unit

	Carrying *Pickup

	lastX, lastY, moveSpeedX float32
	actionTimer, hurtTimer   float32 // negative values can be used for "time since last"
	lastState                State
}

const ( // states
	StateSummoned            State = iota // single frame
	StateWaitingToBeSummoned              // continuous

	StateIdling  // continuous
	StateWalking // continuous

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

const TeamAlly, TeamEnemy, TeamNeutral, TeamCount Team = 0, 1, 2, 3

const RoleMelee, RoleRanged, RoleTank, RoleMage Role = 0, 1, 2, 3
const RoleHealer, RoleCollector, RoleSupplier, RoleTrapper, RoleCount Role = 4, 5, 6, 7, 8

const Gravity, GroundFrictionPercent, BloodMultiplier = 256.0, 15.0, 40.0

var Units []*Unit = make([]*Unit, 0, 16)
var Collisions = map[Lane][]geometry.Shape{}
var PinnedUnit *Unit

func NewUnit(character CharacterKind, team Team, lane Lane) *Unit {
	var char = Characters[character]
	var anim = motion.NewAnimation(0, false, char.Animations.Idle...)
	var unit = Unit{Object: graphics.NewSprite(-2000, -2000, 1, 0), Character: character, Team: team, Lane: lane,
		Behavior: char.Behavior, Anim: &anim, actionTimer: number.NaN(), hurtTimer: number.NaN()}

	if team == TeamAlly {
		unit.State = StateWaitingToBeSummoned
	}

	unit.Blood = motion.NewParticleSystem(unit.particlesBlood)

	unit.draw() // update frame size
	unit.PrepareSpawn()
	return &unit
}

//=================================================================

func (u *Unit) Hitbox() geometry.Shape {
	var char = Characters[u.Character]
	var hitbox = char.Hitbox
	hitbox.X, hitbox.Y = u.X+hitbox.X, u.Y+hitbox.Y
	return hitbox
}
func (u *Unit) MyEntrance() (entrance *Entrance) {
	return Bases[u.Team].Entrances[u.Lane/2]
}
func (u *Unit) EnemyEntrance() (canBeActedUpon bool, entrance *Entrance) {
	var e *Entrance
	if u.Team != TeamNeutral && (u.IsLaner() || u.IsOffLaner()) {
		e = Bases[1-u.Team].Entrances[u.Lane/2]
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

func (u *Unit) PrepareSpawn() {
	u.actionTimer, u.hurtTimer = number.NaN(), number.NaN()
	u.Effects.Tint = palette.White
	u.Stats = Characters[u.Character].Stats
	u.IsReturning, u.Health = false, u.Stats.MaxHealth
	u.VelocityX, u.VelocityY = 0, 0

	if u.IsGarrisoner() {
		u.Stats.ActRange += 2
	}

	var col = Collisions[u.Lane]
	var laneX, laneY = col[0].X + col[0].Width/2, col[0].Y - col[0].Height/2 - u.Height/2
	switch u.Lane {
	case LaneLower, LaneLowerOff:
		u.X, u.Y = laneX-40, laneY
	case LaneMiddle, LaneMiddleOff:
		u.X, u.Y = laneX-72, laneY
	case LaneUpper, LaneUpperOff:
		u.X, u.Y = laneX-104, laneY
	case LaneGarrison1, LaneGarrison2, LaneGarrison3:
		u.X, u.Y = CurrentZone.Ground.Width/2+u.Width/2, laneY
	case LaneGarrisonPlus1, LaneGarrisonPlus2, LaneGarrisonPlus3:
		u.X, u.Y = PointAtCell(18, 3)
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

	var hb = u.Hitbox()
	u.HealthBar = NewHealthBar(hb.Width-1, u.Team, u.IsOffLaner())
}
func (u *Unit) Update() {
	if u == nil {
		return
	}

	u.actionTimer -= DeltaTimeScaled()
	u.hurtTimer -= DeltaTimeScaled()

	if !u.IsSummoned() && (number.IsNaN(u.hurtTimer) || u.hurtTimer < -float32(u.Stats.RespawnTimer)/10) {
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
		u.applyState()
		u.actUponState()
		u.applyPhysics()
		u.applyCollisions()
		u.Behavior(u)
	}
	u.draw()
	u.Carrying.Update()

	if TimeScale > 0 {
		u.Blood.Update()

		var speedX = number.Absolute(u.X-u.lastX) / DeltaTimeScaled() // smooth out for FPS dips
		u.moveSpeedX = u.moveSpeedX + (speedX-u.moveSpeedX)*0.15      // 0.15 = how fast it catches up
		if number.IsNaN(u.moveSpeedX) || u.moveSpeedX < 0.01 {
			u.moveSpeedX = 0
		}
		u.lastX, u.lastY = u.X, u.Y
	}
	u.lastState = u.State
}
func (u *Unit) DrawTooltip(shape geometry.Shape, bench bool) {
	var tsz float32 = TileSize
	var col, noMask = palette.White, geometry.Area{}
	const width, height = 140.0, 94.0
	var x, y = shape.X, shape.Y - shape.Height/2 - height/2
	if bench {
		x, y = 0, GameHUD.UnitsPanel.Y+GameHUD.UnitsPanel.Height/2+height/2
	}
	var area = geometry.NewArea(x, y, width, height).Inside(GameHUD.View.Bounds())
	x, y = area.X, area.Y

	GameHUD.Highlight(GameHUD.View, shape, palette.White)
	GameHUD.View.DrawImage(x, y, width, height, 0, PanelNinePatchId, col, noMask)

	var char = Characters[u.Character]
	var icon = char.Icon
	GameHUD.View.DrawImage(x+width/2-tsz/2-6, y-height/2+tsz/2+6, tsz, tsz, 0, SlotId, col, noMask)
	GameHUD.View.DrawImage(x+width/2-tsz/2-6, y-height/2+tsz/2+6, -tsz, tsz, 0, icon, col, noMask)

	var b, s = char.Stats, u.Stats
	var health1, health2, speed, val, rng, actionTimer, respawnTimer string
	if u.IsSummoned() {
		if u.Health != b.MaxHealth {
			health1 = TooltipTexts[0].Set(u.Health, "/")
		}
		if s.MaxHealth != b.MaxHealth {
			health2 = statEquation(s.MaxHealth, b.MaxHealth, TooltipTexts[1])
		}
		if s.Speed != b.Speed {
			speed = statEquation(s.Speed, b.Speed, TooltipTexts[2])
		}
		if s.ActValue != b.ActValue {
			val = statEquation(s.ActValue, b.ActValue, TooltipTexts[3])
		}
		if s.ActRange != b.ActRange {
			rng = statEquation(s.ActRange, b.ActRange, TooltipTexts[4])
		}
		if u.actionTimer > 0 {
			actionTimer = TooltipTexts[5].Set(number.Round(u.actionTimer, 1), "/")
		}
	}
	if u.State == StateDecaying && u.hurtTimer > -float32(s.RespawnTimer)/10 {
		respawnTimer = TooltipTexts[6].Set(number.Round(float32(s.RespawnTimer)/10+u.hurtTimer, 1), "/")
	}
	TooltipLabel.Shape = geometry.NewRectangle(x, y, width-16, height-12, 0)
	TooltipLabel.Effects.TextAlignX, TooltipLabel.Effects.TextAlignY = 0, 0.5
	TooltipLabel.Text = TooltipTexts[7].Set(
		"🟩", Tags[IconHealth], health1, b.MaxHealth, " health ", health2, "\n",
		"🟨", Tags[IconMove], s.Speed, " speed ", speed, "\n",
		"🟥", Tags[roleIcons[s.Role]], s.ActValue, " ", char.ActValueName, " ", val, "\n",
		"🟧", Tags[IconRange], s.ActRange, " range ", rng, "\n",
		"🌗🟪", Tags[IconTimer], actionTimer, number.Round(float32(b.ActTime)/10, 1), "s action\n",
		"🌗🟦", Tags[IconLoop], respawnTimer, number.Round(float32(b.RespawnTimer)/10, 1), "s respawn\n",
		"⬜", char.Info,
	)
	GameHUD.View.DrawObject(TooltipLabel)

	TooltipLabel.Effects.TextAlignX, TooltipLabel.Effects.TextAlignY = 1, 1
	TooltipLabel.Text = TooltipTexts[8].Set(s.Name, "\n",
		Tags[roleIcons[s.Role]], roleNames[s.Role], "\n\n\n",
		Tags[IconHome], zoneNames[char.Origin])
	GameHUD.View.DrawObject(TooltipLabel)

	if PinnedUnit == u {
		var pin = UserInterface.Crops("icons-text")[IconLocked]
		var sz float32 = TileSize / 3
		GameHUD.View.DrawImage(x+width/2-4, y-height/2+4, sz, sz, 0, pin, palette.White, geometry.Area{})
	}
}

func (u *Unit) TakeDamage(damage int) {
	if u.Health > 0 {
		u.Health -= damage
		u.hurtTimer = u.Stats.HurtTime
	}
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
var teamColors = [TeamCount]uint{TeamAlly: palette.Green, TeamEnemy: palette.Red, TeamNeutral: palette.Orange}

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
		if u == t || t.Health <= 0 || t.IsOffLaner() {
			continue
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

	if u.State == StateWalking && u.Health > 0 && (!u.IsGrounded || u.moveSpeedX < 0.01) {
		u.State = StateIdling
	} else if u.State == StateIdling && u.UnitFront == nil && !u.IsAtWall && u.IsGrounded && u.Health > 0 && !canBeActedUpon {
		u.State = StateWalking
	}

	if u.State == StateSummoned && u.lastState == StateSummoned {
		u.State = StateWalking // first frame is event, second frame (now) starts walking
	}

	if u.State == StateActionEnd && u.Health > 0 {
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
	} else if (u.State == StateIdling || u.State == StateWalking) && ranged && garrisonOrNot && canShoot {
		if canAct {
			u.State = StateActionStart
		} else if u.Health > 0 && u.IsLaner() { // no shoot-move-shoot-move for laners - but garrisoners should
			u.State = StateIdling // enemy in range but waiting for action timer (stay in one place, don't keep walking)
		}
	}

	if u.State == StateHurting && u.Health <= 0 {
		u.State = StateDyingStart // bug fix for units sometimes staying alive
	} else if u.State == StateHurting && u.hurtTimer < 0 && u.Health > 0 {
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
		u.VelocityX = 0
	case StateWalking:
		u.Anim.Frames = Characters[u.Character].Animations.Walk
		u.Anim.IsLooping, u.Anim.FPS = true, u.moveSpeedX*0.25

		if u.IsLaner() && u.IsInsideEnemyBase(TileSize/1.5) {
			u.HealthBar.MoveToGlory(2.5)
		} else if u.IsOffLaner() && u.IsInsideEnemyBase(-TileSize) {
			u.IsReturning = true
		}
		if (u.IsInsideEnemyBase(0) || !u.IsOutsideOwnBase()) && u == PinnedUnit && !u.IsGarrisoner() {
			PinnedUnit = nil
		}

		u.VelocityX = float32([TeamCount]int{u.Stats.Speed, -u.Stats.Speed}[u.Team])
		if u.IsAtWall { // chill out when next to a wall
			u.VelocityX = 0
		}
		if u.IsReturning {
			u.VelocityX = -u.VelocityX
		}

		if u.X < -CurrentZone.Ground.Width/2-u.Width*2 || u.X > CurrentZone.Ground.Width/2+u.Width*2 {
			u.hurtTimer = 0 // no instant delete - to have time to play glory text animation etc
			u.State = StateDecaying

			if u.IsOffLaner() && u.IsReturning && u.Carrying != nil {
				u.Carrying.Target = nil
				u.Carrying.Mask = geometry.Area{}
				u.Carrying.SlotUI = GameHUD.FreePickupSlot()
				GameHUD.Pickups[u.Carrying.SlotUI] = u.Carrying
				collection.Remove(Pickups, u.Carrying)
				u.Carrying = nil
			}
		}
	case StateActionStart: // random delay to balance same sided units melee VVVVVVV
		u.actionTimer = float32(u.Stats.ActTime)/10 + random.Range[float32](0, 0.1)
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
		var percent = 1 - float32(u.Health)/float32(u.Stats.MaxHealth)
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
		var respawnTimer = -float32(u.Stats.RespawnTimer) / 10
		if u.hurtTimer < respawnTimer || u.IsGarrisoner() {
			Units = collection.Remove(Units, u)
			if PinnedUnit == u {
				PinnedUnit = nil
			}

			if u.Team == TeamAlly {
				u.State = StateWaitingToBeSummoned
				u.PrepareSpawn()
			}
		} else if u.hurtTimer < 0 {
			u.Effects.Tint = color.RGBA(255, 255, 255, byte(number.Map(u.hurtTimer, 0, respawnTimer, 255, 0)))
		}
	}
}

func (u *Unit) applyPhysics() {
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

	u.UnitBehind, u.UnitFront = nil, nil
	for _, other := range Units {
		var ohb = other.Hitbox()
		var anyoneDead = u.Health <= 0 || other.Health <= 0
		var isGarrison = other.IsGarrisoner() || u.IsGarrisoner()
		var isOffLaner = other.IsOffLaner() || u.IsOffLaner()
		if other == u || u.Lane != other.Lane || anyoneDead || isGarrison || isOffLaner || !hb.Overlaps(ohb) ||
			other.IsInsideEnemyBase(TileSize/2) {
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

	if u.IsOffLaner() && GameHUD.FreePickupSlot() >= 0 && u.Carrying == nil {
		for _, p := range Pickups {
			if p != nil && p.Target == nil && p.lane == u.Lane && p.Overlaps(hb) {
				u.Carrying = p
				p.Target = u
				u.IsReturning = true
				collection.Remove(Pickups, p)
			}
		}
	}
}
func (u *Unit) draw() {
	var frame = u.Anim.Frame()
	var crop = frame.CropArea()
	u.ImageId, u.Width, u.Height = frame, crop.Width, crop.Height

	if u.Health > 0 && !u.IsGarrisoner() {
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

	if !u.IsGarrisoner() && (!u.IsOutsideOwnBase() || u.IsInsideEnemyBase(0)) {
		return
	}

	if u.Object.ContainsPoint(View.MousePosition()) {
		mouse.SetCursor(cursor.Hand)

		if mouse.IsAnyButtonJustPressed() {
			PinnedUnit = u
		}
	}
}
