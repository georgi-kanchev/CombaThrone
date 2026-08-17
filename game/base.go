package game

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/text"
)

type Base uint8
type Garrison uint8
type SaveState struct {
	Base     Base
	Garrison Garrison
	Coins    int
}
type BaseData struct {
	SaveState

	Glory int

	Back, Front, GarrisonBack, GarrisonFront, GloryLabel *graphics.Object

	lastGlory int
}

const BaseCamp, BaseBarrack, BaseFort, BaseFortress Base = 0, 1, 2, 3
const Garrison0, Garrison1, Garrison2, Garrison3 Garrison = 0, 1, 2, 3

func NewBase(saveState SaveState, ally bool) BaseData {
	var glory = map[Base]int{BaseCamp: 15, BaseBarrack: 30, BaseFort: 90, BaseFortress: 360}[saveState.Base]
	var b = BaseData{SaveState: saveState, Glory: glory}

	var ptsX, ptsY = PointAtCell(14, 5.5)
	if ally {
		ptsX, ptsY = PointAtCell(3, 5.5)
	}
	var label = graphics.NewTextbox(ptsX, ptsY, TileSize, TileSize, 0)
	label.Effects.FillColor, label.Effects.TextLineHeight = 0, 11
	label.Effects.TextAlignX, label.Effects.TextAlignY = 0.5, 0.5
	label.Effects.TextShadowColor, label.Effects.TextShadowWeight = 0, 0
	label.Effects.OutlineSize, label.Effects.OutlineColor = 0.4, palette.Black
	b.GloryLabel = &label

	switch b.Base {
	case BaseCamp:
		var back = graphics.NewTilemap(Layers[0])
		b.Back = &back
		return b
	case BaseBarrack:
		var barrack = graphics.NewTilemap(Layers[1])
		b.Back = &barrack
		return b
	}

	var back, front = graphics.NewTilemap(Layers[2]), graphics.NewTilemap(Layers[3])
	b.Back, b.Front = &back, &front

	if b.Garrison != Garrison0 {
		var indexes = map[Garrison][2]int{Garrison1: {4, 5}, Garrison2: {6, 7}, Garrison3: {8, 9}}[b.Garrison]
		var garBack, garFront = graphics.NewTilemap(Layers[indexes[0]]), graphics.NewTilemap(Layers[indexes[1]])
		b.GarrisonBack, b.GarrisonFront = &garBack, &garFront

		if b.Base == BaseFort {
			b.GarrisonBack.Y += TileSize
			b.GarrisonFront.Y += TileSize
		}
	}

	if b.Base == BaseFort {
		b.Back.Y += TileSize
		b.Front.Y += TileSize

		if !ally {
			b.Back.Width *= -1
			b.Front.Width *= -1
		}
		bringGarrisonLanesDown(ally)
	}
	return b
}

func (b *BaseData) UpdateBack() {
	View.DrawObject(b.Back)
	View.DrawObject(b.GarrisonBack)
}
func (b *BaseData) UpdateFront() {
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
