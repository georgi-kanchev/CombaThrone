package game

import (
	char "game/code/constants/character"
	"game/code/constants/emotion"
	"game/code/constants/progress"
	"game/code/constants/screen"
	"game/code/constants/sound"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/execution/condition"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/gui"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/cursor"
	"pure-game-kit/packages/utility/text"
)

var EventPortraitLeft, EventPortraitRight assets.ImageId
var EventText text.Dynamic

func UpdateEvents(panel geometry.Area) {
	eventTimer += DeltaTimeScaled()

	switch event {
	case 0:
		eventWaitFor(2)
	case 1:
		var start = storyMapProgressPoints[progress.Start]
		MoveOnMap(start.X+TileSize, start.Y)
		eventNext()
	case 2:
		eventWaitFor(0.7)
	case 3:
		ProgressOnMap()
		eventNext()
		EventPortraitLeft = Characters[char.Explorer].Portraits[emotion.Happy]
		EventText.Set("🔡⬅Finally a place to rest!\nI sure hope they treat outlanders well...")
	case 4:
		for i := progress.Tutorial; i <= progress.Underground; i++ {
			var p = storyMapProgressPoints[i]
			p.Width, p.Height = TileSize*1.4, TileSize*1.4

			var hovered = p.ContainsPoint(View.MousePosition())
			if mouse.IsAnyButtonJustPressed() && hovered && Player.Progress == i && TimeScale > 0 {
				MoveOnMap(p.X, p.Y)
				eventNext()
				PlaySound(sound.Button)
				EventPortraitLeft = Characters[char.Explorer].Portraits[emotion.Calm]
				EventText.Set("⬅...")
			}
		}

		if condition.TrueEvery(2, 0xCFA) {
			var p = storyMapProgressPoints[progress.Tutorial]
			PlayCallForAction(p.X, p.Y)
		}
	case 5:
		eventWaitFor(0.7)
	case 6:
		Screen = screen.Scene
		mouse.SetCursor(cursor.Arrow)
		ScreenFade.GoTo(0, nil, 255).Wait(0.2).GoTo(2.5, nil, 0)
		PauseMusic()
		PauseAmbience()
		eventNext()
	case 7:
		eventWaitFor(10.0)
	case 8:
		Units[0].Values.DisablesTimer = 5
	}

	gui.Image(EventPortraitLeft, geometry.NewArea(panel.X-panel.Width/2+TileSize*1.5, panel.Y, TileSize*2, TileSize*3), false)
	gui.Image(EventPortraitRight, geometry.NewArea(panel.X+panel.Width/2-TileSize-6, panel.Y, -TileSize*2, TileSize*3), false)
	panel.Width, panel.X = panel.Width-TileSize*4.5, panel.X+4
	gui.Text(EventText.Get(), panel, geometry.Area{}, ThemeGUI, false)
}

// private ========================================================

var event int
var eventTimer float32

func eventNext() {
	eventTimer = 0
	event++
}
func eventWaitFor(seconds float32) {
	if eventTimer > seconds {
		eventNext()
	}
}
