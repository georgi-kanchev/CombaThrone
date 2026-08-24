package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/text"
)

type BaseKind uint8
type Garrison uint8
type Base struct {
	Kind          BaseKind
	Garrison      Garrison
	EntranceKinds [3]EntranceKind

	Entrances []*Entrance

	Glory int

	Back, Front, GarrisonBack, GarrisonFront, Flag *graphics.Object

	FlagAnim *motion.Animation[assets.ImageId]

	team      Team
	lastGlory int
}

const (
	BaseNone BaseKind = iota
	BaseCamp
	BaseBarrack
	BaseFort
	BaseFortress
)
const GarrisonNone, Garrison1, Garrison2, Garrison3 Garrison = 0, 1, 2, 3

var Bases [TeamCount]*Base

func NewBase(team Team, kind BaseKind, garrison Garrison, entrances [3]EntranceKind) *Base {
	var b = &Base{team: team, Kind: kind, Garrison: garrison, Glory: baseGlory[kind]}
	b.Entrances = make([]*Entrance, 3)
	b.Entrances[LaneLower/2] = NewEntrance(entrances[LaneLower/2], b.Kind, team, LaneLower)
	b.Entrances[LaneMiddle/2] = NewEntrance(entrances[LaneMiddle/2], b.Kind, team, LaneMiddle)
	b.Entrances[LaneUpper/2] = NewEntrance(entrances[LaneUpper/2], b.Kind, team, LaneUpper)

	var x, y = PointAtCell(14, 5.5)
	if team == TeamAlly {
		x, _ = PointAtCell(3, 5.5)
	}

	var flag = graphics.NewSprite(x, y, 1, 0)
	var anim = motion.NewAnimation[assets.ImageId](2, true)
	b.Flag, b.FlagAnim = &flag, &anim

	switch b.Kind {
	case BaseNone:
		return b
	case BaseCamp:
		var back = graphics.NewTilemap(Layers[BaseCamp-1])
		b.Back = &back
		return b
	case BaseBarrack:
		var barrack = graphics.NewTilemap(Layers[BaseBarrack-1])
		b.Back = &barrack
		return b
	}

	var back, front = graphics.NewTilemap(Layers[BaseFort-1]), graphics.NewTilemap(Layers[BaseFort])
	b.Back, b.Front = &back, &front

	if team == TeamEnemy {
		b.Back.Width *= -1
		b.Front.Width *= -1
	}
	if b.Kind == BaseFort {
		b.Back.Y += TileSize
		b.Front.Y += TileSize

		bringGarrisonLanesDown(team)
	}

	if b.Garrison != GarrisonNone {
		var backIndex, frontIndex = int(BaseFort-1) + int(b.Garrison)*2, int(BaseFort) + int(b.Garrison)*2
		var garBack, garFront = graphics.NewTilemap(Layers[backIndex]), graphics.NewTilemap(Layers[frontIndex])
		b.GarrisonBack, b.GarrisonFront = &garBack, &garFront

		if b.Kind == BaseFort {
			b.GarrisonBack.Y += TileSize
			b.GarrisonFront.Y += TileSize
		}
		if team == TeamEnemy {
			b.GarrisonBack.Width *= -1
			b.GarrisonFront.Width *= -1
		}
	}
	return b
}

//=================================================================

var baseGlory = map[BaseKind]int{BaseNone: 10, BaseCamp: 15, BaseBarrack: 30, BaseFort: 90, BaseFortress: 360}

func (b *Base) UpdateBack() {
	View.DrawObject(b.Back)
	View.DrawObject(b.GarrisonBack)

	var group = "flag-" + zoneNames[CurrentZone.kind]
	var x, y = b.Flag.X, b.Flag.Y
	if b.team == TeamAlly {
		group = "flag-player"
	} else if CurrentZone.kind == ZoneDocks {
		x, y = PointAtCell(13.85, -0.15)
		b.FlagAnim.FPS = 5
	} else {
		x, y = PointAtCell(14, 5.5)
		b.FlagAnim.FPS = 2
	}
	b.Flag.X, b.Flag.Y = x, y
	b.FlagAnim.TimeScale = TimeScale * CurrentZone.WindSpeed
	b.FlagAnim.Frames = Decor.Crops(group)
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

	if b.team == TeamAlly && Player.CoinsJustChanged() {
		UI.Coins.Text = text.New(Player.Coins, "$")
	}
	if b.Glory != b.lastGlory {
		UI.TeamGlory[b.team].Text = text.New(b.Glory, Tags[IconGlory])
		b.lastGlory = b.Glory
	}

	View.DrawObject(b.Front)
	View.DrawObject(b.GarrisonFront)
}
