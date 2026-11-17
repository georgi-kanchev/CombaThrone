package game

import (
	"game/code/constants/progress"
	"game/code/constants/screen"
	"pure-game-kit/packages/execution/condition"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/cursor"
)

var event int
var eventTimer float32

func UpdateEvents() {
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
	case 4:
		for i := progress.Tutorial; i <= progress.Underground; i++ {
			var p = storyMapProgressPoints[i]
			p.Width, p.Height = TileSize*1.4, TileSize*1.4
			var hovered = p.ContainsPoint(View.MousePosition())

			if mouse.IsAnyButtonJustPressed() && hovered && Player.Progress == i {
				MoveOnMap(p.X, p.Y)
				eventNext()
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
	}
}
func eventNext() {
	eventTimer = 0
	event++
}
func eventWaitFor(seconds float32) {
	if eventTimer > seconds {
		eventNext()
	}
}
