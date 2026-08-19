package game

import (
	"pure-game-kit/packages/execution/condition"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/random"
	"pure-game-kit/packages/utility/time"
)

type EntranceKind uint8
type Entrance struct {
	Tiles     []*graphics.Object
	Kind      EntranceKind
	Lane      Lane
	Team      Team
	HealthBar *HealthBar

	MaxHealth, Health int

	originalTileYs              []float32
	openY, maxOpenY, shakeTimer float32
}

const EntranceNone, EntranceDoor, EntranceShortGate, EntranceTallGate EntranceKind = 0, 1, 2, 3

func NewEntrance(entry EntranceKind, base BaseKind, team Team, lane Lane) *Entrance {
	var data = Entrance{Kind: entry, Team: team, Lane: lane, maxOpenY: TileSize + (TileSize * (float32(entry) - 1))}
	var x, y float32 = -208, 48 // ally upper lane by default

	if entry != EntranceNone {
		data.HealthBar = NewHealthBar(TileSize-2, team, false)
		data.MaxHealth, data.Health = 100, 100
	}

	if lane == LaneMiddle {
		x -= TileSize
		y += TileSize
	}
	if lane == LaneLower {
		x -= TileSize * 2
		y += TileSize * 2
	}

	switch entry {
	case EntranceNone:
		var hole = graphics.NewSprite(x, y, 1, Decor.Crops("hole")[0])
		if base < BaseBarrack {
			hole.X = -Background.Width / 2 // pull back entrance to edge of scene, fighting can happen further back
			hole.Effects.Tint = 0          // and hide the hole
		}

		data.Tiles = []*graphics.Object{&hole}
	case EntranceDoor:
		var door = graphics.NewSprite(x, y, 1, Decor.Crops("door")[1])
		data.Tiles = []*graphics.Object{&door}
	case EntranceShortGate:
		var top0 = graphics.NewSprite(x, y-TileSize, 1, Decor.Crops("gate_top")[0])
		var mid0 = graphics.NewSprite(x, y, 1, Decor.Crops("gate_middle")[0])
		var bot0 = graphics.NewSprite(x, y+TileSize, 1, Decor.Crops("gate_bottom")[0])
		var top1 = graphics.NewSprite(x, y-TileSize, 1, Decor.Crops("gate_top")[1])
		var mid1 = graphics.NewSprite(x, y, 1, Decor.Crops("gate_middle")[1])
		var bot1 = graphics.NewSprite(x, y+TileSize, 1, Decor.Crops("gate_bottom")[1])
		data.Tiles = []*graphics.Object{&top0, &mid0, &bot0, &top1, &mid1, &bot1}
	case EntranceTallGate:
		y -= TileSize / 2
		var top0 = graphics.NewSprite(x, y-TileSize*1.5, 1, Decor.Crops("gate_top")[0])
		var midU0 = graphics.NewSprite(x, y-TileSize*0.5, 1, Decor.Crops("gate_middle")[0])
		var midD0 = graphics.NewSprite(x, y+TileSize*0.5, 1, Decor.Crops("gate_middle")[0])
		var bot0 = graphics.NewSprite(x, y+TileSize*1.5, 1, Decor.Crops("gate_bottom")[0])
		var top1 = graphics.NewSprite(x, y-TileSize*1.5, 1, Decor.Crops("gate_top")[1])
		var midU1 = graphics.NewSprite(x, y-TileSize*0.5, 1, Decor.Crops("gate_middle")[1])
		var midD1 = graphics.NewSprite(x, y+TileSize*0.5, 1, Decor.Crops("gate_middle")[1])
		var bot1 = graphics.NewSprite(x, y+TileSize*1.5, 1, Decor.Crops("gate_bottom")[1])
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

//=================================================================

func (e *Entrance) IsOpen() bool {
	return number.IsWithin(e.openY, e.maxOpenY, 0.1)
}

func (e *Entrance) Update() {
	e.shakeTimer -= DeltaTimeScaled()

	var sensorDistance = float32(TileSize) * 0.75
	if e.Kind == EntranceDoor {
		sensorDistance = TileSize
	}
	var shortestDistance float32 = sensorDistance
	for _, u := range Units {
		var distance = number.Absolute(u.X - e.Tiles[0].X)
		var sameOrOffLane = e.Lane == u.Lane || (e.Lane == u.Lane-1 && u.IsOffLaner())
		if e.Team == u.Team && sameOrOffLane && distance < shortestDistance {
			shortestDistance = distance
		}
	}
	var holdOpenDistance float32 = sensorDistance / 2
	var distance = number.Limit(shortestDistance-holdOpenDistance, 0, sensorDistance-holdOpenDistance)
	var gate = e.Kind == EntranceShortGate || e.Kind == EntranceTallGate
	e.openY = number.Map(distance, sensorDistance-holdOpenDistance, 0, 0, e.maxOpenY)

	var shakeOffsetX, shakeOffsetY float32 = 0, 0
	if e.shakeTimer > 0 {
		shakeOffsetX, shakeOffsetY = random.Range[float32](-3, 3)*e.shakeTimer, random.Range[float32](-3, 3)*e.shakeTimer
	}

	switch e.Kind {
	case EntranceDoor:
		var breakIndex = number.Map(e.Health, 0, e.MaxHealth, 5, 1)
		var isOpen = e.IsOpen()
		if isOpen {
			breakIndex = 0
		}

		if condition.JustTurnedTrue(isOpen, int(e.Tiles[0].X)) {
			PlaySound(AudioDoorOpen)
		}
		if condition.JustTurnedTrue(!isOpen, int(e.Tiles[0].X)*30) && time.Frame() > 1 {
			PlaySound(AudioDoorClose)
		}

		e.Tiles[0].ImageId = Decor.Crops("door")[breakIndex]
	case EntranceShortGate, EntranceTallGate:
		if condition.JustTurnedTrue(e.openY > 0, int(e.Tiles[0].X)) {
			PlaySound(AudioGateOpen)
		}
		if condition.JustTurnedTrue(e.openY == 0, int(e.Tiles[0].X*30)) && time.Frame() > 1 {
			PlaySound(AudioGateClose)
		}

		for i := len(e.Tiles) / 2; i < len(e.Tiles); i++ {
			e.Tiles[i].Y = e.originalTileYs[i] + e.openY
		}
	}

	for i, t := range e.Tiles {
		var lastX, lastY = t.X, t.Y
		if !gate || (gate && i >= len(e.Tiles)/2) {
			t.X, t.Y = t.X+float32(shakeOffsetX), t.Y+float32(shakeOffsetY)
		}
		View.DrawObject(t)
		t.X, t.Y = lastX, lastY
	}
}
func (e *Entrance) TakeDamage(damage int) {
	if e.Health <= 0 {
		return
	}

	e.Health -= damage
	e.shakeTimer = number.Map(float32(damage), 1, 20, 0.1, 0.5)

	if e.Health <= 0 {
		e.HealthBar.FadeOut(1.5)

		switch e.Kind {
		case EntranceDoor:
			PlaySound(AudioDoorCrumble)
		case EntranceShortGate, EntranceTallGate:
			PlaySound(AudioGateCrumble)
		}
	}

	var breakIndex = number.Map(e.Health, 0, e.MaxHealth, 5, 1)
	switch e.Kind { // EntryDoor done in update
	case EntranceShortGate:
		e.Tiles[3].ImageId = Decor.Crops("gate_top")[breakIndex]
		e.Tiles[4].ImageId = Decor.Crops("gate_middle")[breakIndex]
		e.Tiles[5].ImageId = Decor.Crops("gate_bottom")[breakIndex]
	case EntranceTallGate:
		e.Tiles[4].ImageId = Decor.Crops("gate_top")[breakIndex]
		e.Tiles[5].ImageId = Decor.Crops("gate_middle")[breakIndex]
		e.Tiles[6].ImageId = Decor.Crops("gate_middle")[breakIndex]
		e.Tiles[7].ImageId = Decor.Crops("gate_bottom")[breakIndex]
	}
}
