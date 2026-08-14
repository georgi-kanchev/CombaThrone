package game

import (
	"pure-game-kit/packages/graphics"
)

type Base uint8
type Garrison uint8

type BaseData struct {
	Health int

	Back, Front, GarrisonBack, GarrisonFront *graphics.Object

	base     Base
	garrison Garrison
}

const BaseBarrack, BaseFort, BaseFortress Base = 0, 1, 2
const Garrison0, Garrison1, Garrison2, Garrison3 Garrison = 0, 1, 2, 3

func NewBase(base Base, garrison Garrison, ally bool) BaseData {
	var health = map[Base]int{BaseBarrack: 50, BaseFort: 100, BaseFortress: 200}[base]
	var b = BaseData{base: base, garrison: garrison, Health: health}

	if base == BaseBarrack {
		var barrack = graphics.NewTilemap(Layers[0])
		b.Front = &barrack
		return b
	}

	var back, front = graphics.NewTilemap(Layers[1]), graphics.NewTilemap(Layers[2])
	b.Back, b.Front = &back, &front

	if garrison != Garrison0 {
		var indexes = map[Garrison][2]int{Garrison1: {3, 4}, Garrison2: {5, 6}, Garrison3: {7, 8}}[garrison]
		var garBack, garFront = graphics.NewTilemap(Layers[indexes[0]]), graphics.NewTilemap(Layers[indexes[1]])
		b.GarrisonBack, b.GarrisonFront = &garBack, &garFront

		if base == BaseFort {
			b.GarrisonBack.Y += TileSize
			b.GarrisonFront.Y += TileSize
		}
	}

	if base == BaseFort {
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
