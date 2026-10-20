package game

import (
	"game/code/constants/base"
	"game/code/constants/entrance"
	"game/code/constants/garrison"
	"game/code/constants/icon"
	"game/code/constants/lane"
	"game/code/constants/team"
	"game/code/constants/zone"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/text"
)

type Base struct {
	Kind     base.Enum
	Team     team.Enum
	Garrison garrison.Enum

	EntranceKinds [3]entrance.Enum
	Entrances     []*Entrance

	Glory, LastGlory int

	Back, Front, GarrisonBack, GarrisonFront, Flag *graphics.Object

	FlagAnim *motion.Animation[assets.ImageId]
}

var Bases [team.COUNT]*Base

func NewBase(forTeam team.Enum, kind base.Enum, garr garrison.Enum, entranceLow, entranceMid, entranceUp entrance.Enum) *Base {
	var b = &Base{Team: forTeam, Kind: kind, Garrison: garr, Glory: baseGlory[kind]}
	b.Entrances = make([]*Entrance, 3)
	b.Entrances[lane.Lower/2] = NewEntrance(entranceLow, b.Kind, forTeam, lane.Lower)
	b.Entrances[lane.Middle/2] = NewEntrance(entranceMid, b.Kind, forTeam, lane.Middle)
	b.Entrances[lane.Upper/2] = NewEntrance(entranceUp, b.Kind, forTeam, lane.Upper)

	var x, y = PointAtCell(14, 5.5)
	if forTeam == team.Ally {
		x, y = PointAtCell(3, 6)
	}

	var flag = graphics.NewSprite(x, y, 1, 0)
	var anim = motion.NewAnimation[assets.ImageId](2, true)
	b.Flag, b.FlagAnim = &flag, &anim

	switch b.Kind {
	case base.NONE:
		return b
	case base.Camp:
		var back = graphics.NewTilemap(Layers[base.Camp-1])
		b.Back = &back
		if forTeam == team.Enemy {
			b.Back.Width *= -1
		}
		return b
	case base.Barrack:
		var barrack = graphics.NewTilemap(Layers[base.Barrack-1])
		b.Back = &barrack
		if forTeam == team.Enemy {
			b.Back.Width *= -1
		}
		return b
	}

	var back, front = graphics.NewTilemap(Layers[base.Fort-1]), graphics.NewTilemap(Layers[base.Fort])
	b.Back, b.Front = &back, &front

	if forTeam == team.Enemy {
		b.Back.Width *= -1
		b.Front.Width *= -1
	}
	if b.Kind == base.Fort {
		b.Back.Y += TileSize
		b.Front.Y += TileSize

		bringGarrisonLanesDown(forTeam)
	}

	if b.Garrison != garrison.NONE {
		var backIndex, frontIndex = int(base.Fort-1) + int(b.Garrison)*2, int(base.Fort) + int(b.Garrison)*2
		var garBack, garFront = graphics.NewTilemap(Layers[backIndex]), graphics.NewTilemap(Layers[frontIndex])
		b.GarrisonBack, b.GarrisonFront = &garBack, &garFront

		if b.Kind == base.Fort {
			b.GarrisonBack.Y += TileSize
			b.GarrisonFront.Y += TileSize
		}
		if forTeam == team.Enemy {
			b.GarrisonBack.Width *= -1
			b.GarrisonFront.Width *= -1
		}
	}
	return b
}

//=================================================================

var baseGlory = map[base.Enum]int{base.NONE: 10, base.Camp: 15, base.Barrack: 30, base.Fort: 90, base.Fortress: 360}

func (b *Base) UpdateBack() {
	View.DrawObject(b.Back)
	View.DrawObject(b.GarrisonBack)

	var x, y = b.Flag.X, b.Flag.Y
	b.FlagAnim.Frames = CurrentZone.FlagFrames
	if b.Team == team.Ally {
		b.FlagAnim.Frames = Decor.Crops("flag-player")
		b.FlagAnim.FPS = 6
	} else {
		if CurrentZone.Kind == zone.Docks {
			x, y = PointAtCell(13.85, 0.31)
			b.FlagAnim.FPS = 5
		} else {
			x, y = PointAtCell(14, 5.5)
			b.FlagAnim.FPS = 2
		}
	}
	b.Flag.X, b.Flag.Y = x, y
	b.FlagAnim.TimeScale = TimeScale * CurrentZone.WindSpeed

	var frame = b.FlagAnim.Frame()
	var crop = frame.CropArea()
	b.Flag.ImageId, b.Flag.Width, b.Flag.Height = frame, crop.Width, crop.Height
	View.DrawObject(b.Flag)
}
func (b *Base) UpdateFront() {
	b.Glory = max(b.Glory, 0)

	if b.Glory == 0 {
		//TimeScale = 0 // game over
	}

	if b.Team == team.Ally && Player.CoinsJustChanged() {
		SceneHUD.Coins.Text = text.New(Player.Coins, "$")
	}
	if b.Glory != b.LastGlory {
		SceneHUD.TeamGlory[b.Team].Text = text.New(Tags[icon.Glory], b.Glory)
		b.LastGlory = b.Glory
	}

	View.DrawObject(b.Front)
	View.DrawObject(b.GarrisonFront)
}
