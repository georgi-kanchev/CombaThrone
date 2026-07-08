package game

// The alive unit in the scene. A unit can be any character, copying its base data at different times,
// then acting upon it and editing it through its behavior (brain function).

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/motion"
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
}

const TeamAlly, TeamEnemy, TeamNeutral Team = 0, 1, 2
const DutyWalkStraight, DutyUseStairs, DutyStayGarrison Duty = 0, 1, 2

//=================================================================

var Units []*Unit

func (u *Unit) Hitbox() geometry.Area {
	var scale = SceneScale()
	var char = Characters[u.Character]
	var hitbox = char.Hitbox
	hitbox.X, hitbox.Y = u.X+hitbox.X*scale, u.Y+hitbox.Y*scale
	hitbox.Width, hitbox.Height = hitbox.Width*scale, hitbox.Height*scale
	return hitbox
}

func (u *Unit) PlayIdle() {
	u.Anim.Frames, u.Anim.Time, u.Anim.IsLooping, u.Anim.FPS = Characters[u.Character].Animations.Idle, 0, true, 3
}
func (u *Unit) PlayWalk() {
	u.Anim.Frames, u.Anim.Time, u.Anim.IsLooping, u.Anim.FPS = Characters[u.Character].Animations.Walk, 0, true, 6
}
func (u *Unit) PlayAttack() {
	u.Anim.Frames, u.Anim.Time, u.Anim.IsLooping, u.Anim.FPS = Characters[u.Character].Animations.Attack, 0, false, 7
}
func (u *Unit) PlayHurt() {
	u.Anim.Frames, u.Anim.Time, u.Anim.IsLooping, u.Anim.FPS = Characters[u.Character].Animations.Hurt, 0, false, 4
}
func (u *Unit) PlayDie() {
	u.Anim.Frames, u.Anim.Time, u.Anim.IsLooping, u.Anim.FPS = Characters[u.Character].Animations.Die, 0, false, 4
}

//=================================================================

func SpawnUnit(character Character, duty Duty, team Team) {
	var char = Characters[character]
	var anim = motion.NewAnimation(0, false, char.Animations.Idle...)
	var unit = Unit{Object: graphics.NewSprite(0, 580, 1, 0), Character: character, Team: team, Duty: duty,
		Brain: char.Brain, Stats: char.Stats, Anim: &anim}
	Units = append(Units, &unit)
}

func UpdateUnits() {
	for _, u := range Units {
		u.Brain(u)

		if keyboard.IsKeyJustPressed(key.A) {
			u.PlayWalk()
		} else if keyboard.IsKeyPressed(key.D) {
			u.PlayWalk()
		} else if keyboard.IsKeyJustPressed(key.Space) {
			u.PlayAttack()
		}

		if u.Anim.IsFinished() {
			u.PlayIdle()
		}

		var frame, scale = u.Anim.Frame(), SceneScale()
		var _, _, w, h = frame.CropArea()

		u.ImageId, u.Width, u.Height = frame, w*scale, h*scale

		if Debug {
			var hb = u.Hitbox()
			View.DrawShape(u.X, u.Y, u.Width, u.Height, 0, 0, DebugUnitColor, geometry.Area{})
			View.DrawShape(hb.X, hb.Y, hb.Width, hb.Height, 0, 0, DebugHitboxColor, geometry.Area{})
		}
		View.DrawObject(&u.Object)
	}
}
