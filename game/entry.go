package game

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/number"
)

type Entry uint8

const EntryHole, EntryDoor, EntryShortGate, EntryTallGate Entry = 0, 1, 2, 3

type EntryData struct {
	Tiles     []*graphics.Object
	Entry     Entry
	Duty      Duty
	Team      Team
	HealthBar HealthBar

	MaxHealth, Health int

	originalTileYs  []float32
	openY, maxOpenY float32
}

func NewEntry(entry Entry, team Team, duty Duty) *EntryData {
	var data = EntryData{Entry: entry, Team: team, Duty: duty, maxOpenY: TileSize + (TileSize * (float32(entry) - 1))}
	var x, y float32 = -208, 48 // ally upper lane by default

	if entry != EntryHole {
		data.HealthBar = NewHealthBar(TileSize-2, team)
		data.MaxHealth, data.Health = 100, 100
	}

	if duty == DutyMiddle {
		x -= TileSize
		y += TileSize
	}
	if duty == DutyLower {
		x -= TileSize * 2
		y += TileSize * 2
	}

	switch entry {
	case EntryHole:
		var hole = graphics.NewSprite(x, y, 1, TilesetCrops.Frame("hole", 0))
		data.Tiles = []*graphics.Object{&hole}
	case EntryDoor:
		var door = graphics.NewSprite(x, y, 1, TilesetCrops.Frame("door", 1))
		data.Tiles = []*graphics.Object{&door}
	case EntryShortGate:
		var top0 = graphics.NewSprite(x, y-TileSize, 1, TilesetCrops.Frame("gate-top", 0))
		var mid0 = graphics.NewSprite(x, y, 1, TilesetCrops.Frame("gate-middle", 0))
		var bot0 = graphics.NewSprite(x, y+TileSize, 1, TilesetCrops.Frame("gate-bottom", 0))
		var top1 = graphics.NewSprite(x, y-TileSize, 1, TilesetCrops.Frame("gate-top", 1))
		var mid1 = graphics.NewSprite(x, y, 1, TilesetCrops.Frame("gate-middle", 1))
		var bot1 = graphics.NewSprite(x, y+TileSize, 1, TilesetCrops.Frame("gate-bottom", 1))
		data.Tiles = []*graphics.Object{&top0, &mid0, &bot0, &top1, &mid1, &bot1}
	case EntryTallGate:
		y -= TileSize / 2
		var top0 = graphics.NewSprite(x, y-TileSize*1.5, 1, TilesetCrops.Frame("gate-top", 0))
		var midU0 = graphics.NewSprite(x, y-TileSize*0.5, 1, TilesetCrops.Frame("gate-middle", 0))
		var midD0 = graphics.NewSprite(x, y+TileSize*0.5, 1, TilesetCrops.Frame("gate-middle", 0))
		var bot0 = graphics.NewSprite(x, y+TileSize*1.5, 1, TilesetCrops.Frame("gate-bottom", 0))
		var top1 = graphics.NewSprite(x, y-TileSize*1.5, 1, TilesetCrops.Frame("gate-top", 1))
		var midU1 = graphics.NewSprite(x, y-TileSize*0.5, 1, TilesetCrops.Frame("gate-middle", 1))
		var midD1 = graphics.NewSprite(x, y+TileSize*0.5, 1, TilesetCrops.Frame("gate-middle", 1))
		var bot1 = graphics.NewSprite(x, y+TileSize*1.5, 1, TilesetCrops.Frame("gate-bottom", 1))
		data.Tiles = []*graphics.Object{&top0, &midU0, &midD0, &bot0, &top1, &midU1, &midD1, &bot1}
	}

	for _, t := range data.Tiles {
		data.originalTileYs = append(data.originalTileYs, t.Y)
		if team == TeamEnemy {
			t.X *= -1
			t.Width *= -1
		}
	}
	return &data
}

func (e *EntryData) IsOpen() bool {
	return number.IsWithin(e.openY, e.maxOpenY, 0.1)
}

//=================================================================

func (e *EntryData) Update() {
	var sensorDistance = float32(TileSize) * 0.75
	if e.Entry == EntryDoor {
		sensorDistance = TileSize
	}
	var shortestDistance float32 = sensorDistance
	for _, u := range Units {
		var distance = number.Absolute(u.X - e.Tiles[0].X)
		if e.Team == u.Team && e.Duty == u.Duty && distance < shortestDistance {
			shortestDistance = distance
		}
	}
	var holdOpenDistance float32 = sensorDistance / 2
	var distance = number.Limit(shortestDistance-holdOpenDistance, 0, sensorDistance-holdOpenDistance)
	e.openY = number.Map(distance, sensorDistance-holdOpenDistance, 0, 0, e.maxOpenY)

	switch e.Entry {
	case EntryDoor:
		var breakIndex = number.Map(e.Health, 0, e.MaxHealth, 5, 1)
		if e.IsOpen() {
			breakIndex = 0
		}
		e.Tiles[0].ImageId = TilesetCrops.Frame("door", breakIndex)
	case EntryShortGate, EntryTallGate:

		for i := len(e.Tiles) / 2; i < len(e.Tiles); i++ {
			e.Tiles[i].Y = e.originalTileYs[i] + e.openY
		}
	}

	for _, t := range e.Tiles {
		View.DrawObject(t)
	}
}
func (e *EntryData) TakeDamage(damage int) {
	if e.Health <= 0 {
		return
	}

	e.Health -= damage

	if e.Health <= 0 {
		e.HealthBar.FadeOut(1.5)
	}

	var breakIndex = number.Map(e.Health, 0, e.MaxHealth, 5, 1)
	switch e.Entry { // EntryDoor done in update
	case EntryShortGate:
		e.Tiles[3].ImageId = TilesetCrops.Frame("gate-top", breakIndex)
		e.Tiles[4].ImageId = TilesetCrops.Frame("gate-middle", breakIndex)
		e.Tiles[5].ImageId = TilesetCrops.Frame("gate-bottom", breakIndex)
	case EntryTallGate:
		e.Tiles[4].ImageId = TilesetCrops.Frame("gate-top", breakIndex)
		e.Tiles[5].ImageId = TilesetCrops.Frame("gate-middle", breakIndex)
		e.Tiles[6].ImageId = TilesetCrops.Frame("gate-middle", breakIndex)
		e.Tiles[7].ImageId = TilesetCrops.Frame("gate-bottom", breakIndex)
	}
}

//=================================================================
