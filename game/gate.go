package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
)

type Entry uint8

const EntryHole, EntryDoor, EntryShortGate, EntryTallGate Entry = 0, 1, 2, 3

type EntryData struct {
	Tiles    []*graphics.Object
	Material Entry
	Duty     Duty
	Team     Team
	Health   int
}

func NewGate(material Entry, team Team, duty Duty) *EntryData {
	var gate = EntryData{Material: material, Team: team, Duty: duty, Health: 100}
	var x, y float32 = -208, 48 // ally upper by default

	if duty == DutyMiddle {
		x -= TileSize
		y += TileSize
	}
	if duty == DutyLower {
		x -= TileSize * 2
		y += TileSize * 2
	}

	switch material {
	case EntryHole:
		var hole = graphics.NewSprite(x, y, 1, assets.ImageId(GatesData.Frame("hole", 0)))
		gate.Tiles = []*graphics.Object{&hole}
	case EntryDoor:
		var door = graphics.NewSprite(x, y, 1, assets.ImageId(GatesData.Frame("door", 1)))
		gate.Tiles = []*graphics.Object{&door}
	case EntryShortGate:
		var top0 = graphics.NewSprite(x, y-TileSize, 1, assets.ImageId(GatesData.Frame("gate-top", 0)))
		var mid0 = graphics.NewSprite(x, y, 1, assets.ImageId(GatesData.Frame("gate-middle", 0)))
		var bot0 = graphics.NewSprite(x, y+TileSize, 1, assets.ImageId(GatesData.Frame("gate-bottom", 0)))
		var top1 = graphics.NewSprite(x, y-TileSize, 1, assets.ImageId(GatesData.Frame("gate-top", 1)))
		var mid1 = graphics.NewSprite(x, y, 1, assets.ImageId(GatesData.Frame("gate-middle", 1)))
		var bot1 = graphics.NewSprite(x, y+TileSize, 1, assets.ImageId(GatesData.Frame("gate-bottom", 1)))
		gate.Tiles = []*graphics.Object{&top0, &mid0, &bot0, &top1, &mid1, &bot1}
	case EntryTallGate:
		y -= TileSize / 2
		var top0 = graphics.NewSprite(x, y-TileSize*1.5, 1, assets.ImageId(GatesData.Frame("gate-top", 0)))
		var midU0 = graphics.NewSprite(x, y-TileSize*0.5, 1, assets.ImageId(GatesData.Frame("gate-middle", 0)))
		var midD0 = graphics.NewSprite(x, y+TileSize*0.5, 1, assets.ImageId(GatesData.Frame("gate-middle", 0)))
		var bot0 = graphics.NewSprite(x, y+TileSize*1.5, 1, assets.ImageId(GatesData.Frame("gate-bottom", 0)))
		var top1 = graphics.NewSprite(x, y-TileSize*1.5, 1, assets.ImageId(GatesData.Frame("gate-top", 1)))
		var midU1 = graphics.NewSprite(x, y-TileSize*0.5, 1, assets.ImageId(GatesData.Frame("gate-middle", 1)))
		var midD1 = graphics.NewSprite(x, y+TileSize*0.5, 1, assets.ImageId(GatesData.Frame("gate-middle", 1)))
		var bot1 = graphics.NewSprite(x, y+TileSize*1.5, 1, assets.ImageId(GatesData.Frame("gate-bottom", 1)))
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

func (g *EntryData) Update() {
	for _, v := range g.Tiles {
		View.DrawObject(v)
	}
}
