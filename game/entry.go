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
}

func NewEntry(entry Entry, team Team, duty Duty) *EntryData {
	var gate = EntryData{Entry: entry, Team: team, Duty: duty}
	var x, y float32 = -208, 48 // ally upper lane by default

	if entry != EntryHole {
		gate.HealthBar = NewHealthBar(TileSize-2, true, team)
		gate.MaxHealth, gate.Health = 100, 100
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
		gate.Tiles = []*graphics.Object{&hole}
	case EntryDoor:
		var door = graphics.NewSprite(x, y, 1, TilesetCrops.Frame("door", 1))
		gate.Tiles = []*graphics.Object{&door}
	case EntryShortGate:
		var top0 = graphics.NewSprite(x, y-TileSize, 1, TilesetCrops.Frame("gate-top", 0))
		var mid0 = graphics.NewSprite(x, y, 1, TilesetCrops.Frame("gate-middle", 0))
		var bot0 = graphics.NewSprite(x, y+TileSize, 1, TilesetCrops.Frame("gate-bottom", 0))
		var top1 = graphics.NewSprite(x, y-TileSize, 1, TilesetCrops.Frame("gate-top", 1))
		var mid1 = graphics.NewSprite(x, y, 1, TilesetCrops.Frame("gate-middle", 1))
		var bot1 = graphics.NewSprite(x, y+TileSize, 1, TilesetCrops.Frame("gate-bottom", 1))
		gate.Tiles = []*graphics.Object{&top0, &mid0, &bot0, &top1, &mid1, &bot1}
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
		gate.Tiles = []*graphics.Object{&top0, &midU0, &midD0, &bot0, &top1, &midU1, &midD1, &bot1}
	}
	if team == TeamEnemy {
		for _, o := range gate.Tiles {
			o.X *= -1
			o.Width *= -1
		}
	}
	return &gate
}

func (g *EntryData) ApplyDamage(damage int) {
	if g.Health <= 0 {
		return
	}

	g.Health -= damage

	var breakIndex = number.Map(g.Health, 0, g.MaxHealth, 5, 1)
	switch g.Entry {
	case EntryDoor:
		g.Tiles[0].ImageId = TilesetCrops.Frame("door", breakIndex)
	case EntryShortGate:
		g.Tiles[3].ImageId = TilesetCrops.Frame("gate-top", breakIndex)
		g.Tiles[4].ImageId = TilesetCrops.Frame("gate-middle", breakIndex)
		g.Tiles[5].ImageId = TilesetCrops.Frame("gate-bottom", breakIndex)
	case EntryTallGate:
		g.Tiles[4].ImageId = TilesetCrops.Frame("gate-top", breakIndex)
		g.Tiles[5].ImageId = TilesetCrops.Frame("gate-middle", breakIndex)
		g.Tiles[6].ImageId = TilesetCrops.Frame("gate-middle", breakIndex)
		g.Tiles[7].ImageId = TilesetCrops.Frame("gate-bottom", breakIndex)
	}
}

func (g *EntryData) Update() {
	for _, v := range g.Tiles {
		View.DrawObject(v)
	}
}
