package game

// The alive unit in the scene. A unit can be any character, copying its base data at different times,
// then acting upon it and editing it through its behavior (brain function).

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/number"
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
	Duty      Duty
	Brain     func(self *Unit)
	Anim      *motion.Animation[assets.ImageId]

	Collisions, Center, Down, Left, Right []geometry.Shape
	CollidableTileIds                     []uint16

	VelocityX, VelocityY float32
	IsGrounded           bool

	prevX, prevY float32
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

//=================================================================

func SpawnUnit(character Character, duty Duty, team Team) {
	var char = Characters[character]
	var anim = motion.NewAnimation(0, false, char.Animations.Idle...)
	var unit = Unit{Object: graphics.NewSprite(0, 0, 1, 0), Character: character, Team: team, Duty: duty,
		Brain: char.Brain, Stats: char.Stats, Anim: &anim, Collisions: []geometry.Shape{},
		Center: []geometry.Shape{}, Down: []geometry.Shape{}, Left: []geometry.Shape{}, Right: []geometry.Shape{}}

	switch duty {
	case DutyLow:
		unit.CollidableTileIds = []uint16{16, 19, 21}
	case DutyMiddle:
		unit.CollidableTileIds = []uint16{1, 2, 3, 16, 19, 21}
	case DutyGarrison:
		unit.CollidableTileIds = []uint16{4, 5, 6, 19, 20, 21}
	}

	Units = append(Units, &unit)
}
func UpdateUnits() {
	for _, u := range Units {
		u.applyPhysics()
		u.applyCollisions()
		u.Brain(u)
		u.applyAnimations()
		u.prevX, u.prevY = u.X, u.Y
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
	u.X += u.VelocityX * time.Delta()
	u.Y += u.VelocityY * time.Delta()

	if Debug {
		View.DrawShape(u.X, u.Y, u.Width, u.Height, 0, 0, DebugUnitColor, geometry.Area{})
	}
}
func (u *Unit) applyCollisions() {
	var hb = u.Hitbox()
	var cellX, cellY = CellAtPoint(u.X, u.Y)
	var diffX, diffY = u.X - hb.X, u.Y - hb.Y // cache hitbox and obj offset
	var tileDown = Layers[LayerMap].TileAtCell(int(cellX), int(cellY)+1).Id
	var tileLeft = Layers[LayerMap].TileAtCell(int(cellX)-1, int(cellY)).Id
	var tileRight = Layers[LayerMap].TileAtCell(int(cellX)+1, int(cellY)).Id
	var tileCenter = Layers[LayerMap].TileAtCell(int(cellX), int(cellY)).Id

	u.Collisions = collection.Clear(u.Collisions)
	u.Down = collection.Clear(u.Down)
	u.Left = collection.Clear(u.Left)
	u.Right = collection.Clear(u.Right)
	u.Center = collection.Clear(u.Center)
	if collection.Contains(u.CollidableTileIds, tileDown) {
		u.Down = Tilemaps[LayerMap].TilemapShapesAtCell(int(cellX), int(cellY)+1)
	}
	if collection.Contains(u.CollidableTileIds, tileLeft) {
		u.Left = Tilemaps[LayerMap].TilemapShapesAtCell(int(cellX)-1, int(cellY))
	}
	if collection.Contains(u.CollidableTileIds, tileRight) {
		u.Right = Tilemaps[LayerMap].TilemapShapesAtCell(int(cellX)+1, int(cellY))
	}
	if collection.Contains(u.CollidableTileIds, tileCenter) {
		u.Center = Tilemaps[LayerMap].TilemapShapesAtCell(int(cellX), int(cellY))
	}
	u.Collisions = collection.Join(u.Collisions, u.Center, u.Down, u.Left, u.Right)
	u.IsGrounded = false

	for _, s := range u.Collisions {
		if Debug {
			View.DrawShape(s.X, s.Y, s.Width, s.Height, s.Angle, s.Roundness, DebugCollisionColor, geometry.Area{})
		}
		if hb.Overlaps(s) {
			hb = hb.Collide(s)                // move hitbox
			u.X, u.Y = hb.X+diffX, hb.Y+diffY // move object according to collision + original hitbox offset
			u.VelocityY = 0
			u.IsGrounded = true
		}
	}
	if Debug {
		View.DrawShape(hb.X, hb.Y, hb.Width, hb.Height, 0, hb.Roundness, DebugHitboxColor, geometry.Area{})
	}
}
func (u *Unit) applyAnimations() {
	if u.IsGrounded {
		if number.IsWithin(u.X-u.prevX, 0, 0.1) {
			u.Anim.Frames = Characters[u.Character].Animations.Idle
			u.Anim.IsLooping, u.Anim.FPS = true, 3
		} else {
			u.Anim.Frames = Characters[u.Character].Animations.Walk
			u.Anim.IsLooping, u.Anim.FPS = true, number.Absolute(u.X-u.prevX)*15
		}
	}

	var frame = u.Anim.Frame()
	var _, _, w, h = frame.CropArea()

	u.ImageId, u.Width, u.Height = frame, w, h
	if u.Team == TeamEnemy {
		u.Width = -w
	}
	View.DrawObject(&u.Object)
}
