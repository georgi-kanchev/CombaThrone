package game

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/text"
)

type BaseKind uint8
type Garrison uint8
type SaveState struct {
	Kind          BaseKind
	Garrison      Garrison
	EntranceKinds [3]EntranceKind
	Gold          int
}
type Base struct {
	SaveState

	Entrances [3]*Entrance

	Glory int

	Back, Front, GarrisonBack, GarrisonFront, GloryLabel *graphics.Object

	lastGlory int
}

const BaseNone, BaseCamp, BaseBarrack, BaseFort, BaseFortress BaseKind = 0, 1, 2, 3, 4
const GarrisonNone, Garrison1, Garrison2, Garrison3 Garrison = 0, 1, 2, 3

var AllyBase, EnemyBase Base

func NewBase(team Team, saveState SaveState) Base {
	var b = Base{SaveState: saveState, Glory: baseGlory[saveState.Kind]}
	b.Entrances[LaneLower/2] = NewEntrance(saveState.EntranceKinds[LaneLower/2], b.Kind, team, LaneLower)
	b.Entrances[LaneMiddle/2] = NewEntrance(saveState.EntranceKinds[LaneMiddle/2], b.Kind, team, LaneMiddle)
	b.Entrances[LaneUpper/2] = NewEntrance(saveState.EntranceKinds[LaneUpper/2], b.Kind, team, LaneUpper)

	var ptsX, ptsY = PointAtCell(14, 5.5)
	if team == TeamAlly {
		ptsX, ptsY = PointAtCell(3, 5.5)
	}
	var label = graphics.NewTextbox(ptsX+Flags.X, ptsY+Flags.Y, TileSize, TileSize, 0)
	label.Effects.FillColor, label.Effects.TextLineHeight = 0, 11
	label.Effects.TextAlignX, label.Effects.TextAlignY = 0.5, 0.5
	label.Effects.TextShadowColor, label.Effects.TextShadowWeight = 0, 0
	label.Effects.OutlineSize, label.Effects.OutlineColor = 0.4, palette.Black
	b.GloryLabel = &label

	switch b.Kind {
	case BaseNone:
		return b
	case BaseCamp:
		var back = graphics.NewTilemap(Layers[LayerCamp])
		b.Back = &back
		return b
	case BaseBarrack:
		var barrack = graphics.NewTilemap(Layers[LayerBarrack])
		b.Back = &barrack
		return b
	}

	var back, front = graphics.NewTilemap(Layers[LayerFortBack0]), graphics.NewTilemap(Layers[LayerFortFront0])
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
		var backIndex = int(LayerFortBack0) + int(b.Garrison)*2
		var frontIndex = int(LayerFortFront0) + int(b.Garrison)*2
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
}
func (b *Base) UpdateFront() {
	b.Glory = max(b.Glory, 0)

	if b.Glory == 0 {
		//TimeScale = 0 // game over
	}

	if b.Glory != b.lastGlory {
		b.GloryLabel.Text = text.New(b.Glory, Tags[IconGlory])
	}
	b.lastGlory = b.Glory

	View.DrawObject(b.Front)
	View.DrawObject(b.GarrisonFront)
	View.DrawObject(b.GloryLabel)
}
