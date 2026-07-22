package game

// The alive unit in the scene. A unit can be any character, copying its base data at different times,
// then acting upon it and editing it through its behavior (brain function).

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/text"
	"pure-game-kit/packages/utility/time"
)

type Team uint8
type Duty uint8
type Character uint8

type Unit struct {
	graphics.Object
	Stats     Stats
	Character Character
	Team      Team
	Brain     func(self *Unit)
	Anim      *motion.Animation[assets.ImageId]

	SolidsAll, SolidsCenter, SolidsDown,
	SolidsLeft, SolidsRight []geometry.Shape
	CollidableTileIds []uint16

	VelocityX, VelocityY float32
	IsGrounded           bool

	UnitFront, UnitBehind *Unit

	prevX, prevY, currentSpeed float32
}

const TeamAlly, TeamEnemy, TeamNeutral Team = 0, 1, 2
const DutyLow, DutyMiddle, DutyHigh, DutyGarrison Duty = 0, 1, 2, 3
const Gravity = 256

//=================================================================

var Units []*Unit

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

func SpawnUnit(character Character, team Team) {
	var char = Characters[character]
	var anim = motion.NewAnimation(0, false, char.Animations.Idle...)
	var unit = Unit{Object: graphics.NewSprite(0, 0, 1, 0), Character: character, Team: team,
		Brain: char.Brain, Stats: char.Stats, Anim: &anim, SolidsAll: []geometry.Shape{},
		SolidsCenter: []geometry.Shape{}, SolidsDown: []geometry.Shape{}, SolidsLeft: []geometry.Shape{}, SolidsRight: []geometry.Shape{}}

	switch unit.Stats.Duty {
	case DutyLow:
		unit.CollidableTileIds = []uint16{32}
	case DutyMiddle:
		unit.CollidableTileIds = []uint16{32, 16, 17, 18}
	case DutyHigh:
		unit.CollidableTileIds = []uint16{32, 16, 17, 18, 1, 2, 3}
	}

	Units = append(Units, &unit)
}
func UpdateUnits() {
	if Debug {
		for _, u := range Units {
			View.DrawShape(u.X, u.Y, u.Width, u.Height, 0, 0, DebugUnitColor, geometry.Area{})
		}
	}
	for _, u := range Units {
		if Debug {
			var hb = u.Hitbox()
			View.DrawShape(hb.X, hb.Y, hb.Width, hb.Height, 0, hb.Roundness, DebugHitboxColor, geometry.Area{})
		}

		u.applyPhysics()
		u.applyCollisions()
		u.Brain(u)
		u.applyAnimations()

		var curHorSpeed = number.Absolute(u.X-u.prevX) / time.Delta()       // smooth out for FPS dips
		u.currentSpeed = u.currentSpeed + (curHorSpeed-u.currentSpeed)*0.15 // 0.15 = how fast it catches up
		u.prevX, u.prevY = u.X, u.Y
	}
	if Debug {
		for _, u := range Units {
			var x, y = u.AttackPoint()
			View.DrawShape(x, y, 5, 5, 0, 1, DebugAttackColor, geometry.Area{})
		}
	}
}

//=================================================================

func (u *Unit) applyPhysics() {
	u.VelocityY += Gravity * time.Delta()

	if u.IsGrounded && u.Team == TeamAlly {
		u.VelocityX = float32(u.Stats.Speed)
	} else if u.IsGrounded && u.Team == TeamEnemy {
		u.VelocityX = -float32(u.Stats.Speed)
	}
	u.X, u.Y = u.X+u.VelocityX*time.Delta(), u.Y+u.VelocityY*time.Delta()
}
func (u *Unit) applyCollisions() {
	var hb = u.Hitbox()
	var diffX, diffY = u.X - hb.X, u.Y - hb.Y // cache hitbox and obj offset
	var cellX, cellY = CellAtPoint(u.X, u.Y)
	var tileDown = MapLayer.TileAtCell(int(cellX), int(cellY)+1).Id
	var tileLeft = MapLayer.TileAtCell(int(cellX)-1, int(cellY)).Id
	var tileRight = MapLayer.TileAtCell(int(cellX)+1, int(cellY)).Id
	var tileCenter = MapLayer.TileAtCell(int(cellX), int(cellY)).Id

	u.SolidsAll = collection.Clear(u.SolidsAll)
	u.SolidsDown = collection.Clear(u.SolidsDown)
	u.SolidsLeft = collection.Clear(u.SolidsLeft)
	u.SolidsRight = collection.Clear(u.SolidsRight)
	u.SolidsCenter = collection.Clear(u.SolidsCenter)
	if collection.Contains(u.CollidableTileIds, tileDown) {
		u.SolidsDown = Map.TilemapShapesAtCell(int(cellX), int(cellY)+1)
	}
	if collection.Contains(u.CollidableTileIds, tileLeft) {
		u.SolidsLeft = Map.TilemapShapesAtCell(int(cellX)-1, int(cellY))
	}
	if collection.Contains(u.CollidableTileIds, tileRight) {
		u.SolidsRight = Map.TilemapShapesAtCell(int(cellX)+1, int(cellY))
	}
	if collection.Contains(u.CollidableTileIds, tileCenter) {
		u.SolidsCenter = Map.TilemapShapesAtCell(int(cellX), int(cellY))
	}
	u.SolidsAll = collection.Join(u.SolidsAll, u.SolidsCenter, u.SolidsDown, u.SolidsLeft, u.SolidsRight)
	u.IsGrounded = false

	for _, s := range u.SolidsAll {
		if Debug {
			View.DrawShape(s.X, s.Y, s.Width, s.Height, s.Angle, s.Roundness, DebugCollisionColor, geometry.Area{})
		}
		if hb.Overlaps(s) {
			hb = hb.Collide(s)
			u.X, u.Y = hb.X+diffX, hb.Y+diffY
			u.VelocityY = 0
			u.IsGrounded = true
		}
	}

	u.UnitBehind, u.UnitFront = nil, nil
	for _, other := range Units {
		if other == u {
			continue
		}
		var ohb = other.Hitbox()
		if u.Team == TeamAlly && hb.X+hb.Width/2-0.5 > ohb.X-ohb.Width/2+0.5 {
			continue // walking past the proper X positions results in no collision at all
		} else if u.Team == TeamEnemy && hb.X-hb.Width/2+0.5 < ohb.X+ohb.Width/2-0.5 {
			continue // walking past the proper X positions results in no collision at all
		}
		var shouldFightY = number.IsWithin(hb.Y, ohb.Y, max(hb.Height/2, ohb.Height/2))
		if shouldFightY && hb.Overlaps(ohb) { // no collision outside of Y range to fight
			hb = hb.Collide(ohb)
			u.X, u.Y = hb.X+diffX, hb.Y+diffY
			if (u.Team == TeamAlly && u.X < other.X) || (u.Team == TeamEnemy && u.X > other.X) {
				u.UnitFront = other
			} else if (u.Team == TeamAlly && u.X > other.X) || (u.Team == TeamEnemy && u.X < other.X) {
				u.UnitBehind = other
			}
		}
	}
}
func (u *Unit) applyAnimations() {
	if u.IsGrounded {
		if number.IsWithin(u.X-u.prevX, 0, time.Delta()) {
			u.Anim.Frames = Characters[u.Character].Animations.Idle
			u.Anim.IsLooping, u.Anim.FPS = true, 3
		} else {
			u.Anim.Frames = Characters[u.Character].Animations.Walk
			u.Anim.IsLooping, u.Anim.FPS = true, u.currentSpeed*0.25
		}
	}

	var frame = u.Anim.Frame()
	var _, _, w, h = frame.CropArea()

	u.ImageId, u.Width, u.Height = frame, w, h
	if u.Team == TeamEnemy {
		u.Width = -w
	}
	View.DrawObject(&u.Object)
	u.Width = w

	if Debug && u.Object.ContainsPoint(View.MousePosition()) {
		var txt = text.New("Speed: ", number.Round(u.currentSpeed, 2))
		View.DrawText(txt, u.X-u.Width/2, u.Y-100, 20, 0, palette.White, geometry.Area{})
	}
}
