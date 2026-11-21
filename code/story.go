package game

import (
	char "game/code/constants/character"
	"game/code/constants/emotion"
	"game/code/constants/lane"
	"game/code/constants/progress"
	"game/code/constants/screen"
	"game/code/constants/sound"
	"game/code/constants/state"
	"game/code/constants/team"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/execution/condition"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/gui"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/cursor"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/text"
)

func UpdateStory() {
	var panel = geometry.NewArea(-TileSize*5, -TileSize*3.5, TileSize*7.5, TileSize*2.5)
	if Screen == screen.Scene {
		panel.X, panel.Y = 0, -TileSize/2
	}
	gui.Image(PanelId, panel, false)

	storyTimer += DeltaTimeScaled()

	switch story {
	case 0:
		storyWaitFor(2)
	case 1:
		var start = storyMapProgressPoints[progress.Start]
		MoveOnMap(start.X+TileSize, start.Y)
		storyNext()
	case 2:
		storyWaitFor(0.7)
	case 3:
		ProgressOnMap()
		storyNext()
		storySetPanel("⬅I sure hope they treat outlanders well...", char.Explorer, 0, emotion.Happy, 0)
	case 4:
		for i := progress.Tutorial; i <= progress.Underground; i++ {
			var p = storyMapProgressPoints[i]
			p.Width, p.Height = TileSize*1.4, TileSize*1.4

			var hovered = p.ContainsPoint(View.MousePosition())
			if mouse.IsAnyButtonJustPressed() && hovered && Player.Progress == i && TimeScale > 0 {
				MoveOnMap(p.X, p.Y)
				storyNext()
				PlaySound(sound.Button)
			}
		}

		if condition.TrueEvery(2, 0xCFA) {
			var p = storyMapProgressPoints[progress.Tutorial]
			PlayCallForAction(p.X, p.Y)
		}
	case 5:
		storyWaitFor(0.7)
	case 6:
		Screen = screen.Scene
		mouse.SetCursor(cursor.Arrow)
		ScreenFade.GoTo(0, nil, 255).Wait(0.2).GoTo(2.5, nil, 0)

		storyFisherman = NewUnit(char.Fisherman, team.Enemy, lane.MiddleOff)
		storyFisherman.X = 0
		Units = append(Units, storyFisherman)

		storySetPanel("You can click on this panel or press any key\nto continue the story.\nTry it! :)", 0, 0, 0, 0)

		PauseMusic()
		PauseAmbience()
		storyNext()
	case 7:
		storyFisherman.ActTimer = storyFisherman.Values.ActTimer
		if storyFisherman.State == state.Moving {
			storyFisherman.ActTimer = 0
		}

		if storyJustContinued(panel) {
			storyExplorer = NewUnit(char.Explorer, team.Ally, lane.Middle)
			Units = append(Units, storyExplorer)
			storySetPanel("...", 0, 0, 0, 0)
			storyNext()
		}
	case 8:
		storyFisherman.ActTimer = storyFisherman.Values.ActTimer
		if storyFisherman.State == state.Moving {
			storyFisherman.ActTimer = 0
		}

		if number.IsWithin(storyExplorer.X, 0, TileSize*2) {
			storySetPanel("⬅G'day, mister! Uhh, wh... what are you doing?", char.Explorer, char.Fisherman, emotion.Talk, emotion.Calm)
			storyNext()
		}
	case 9:
		storyFisherman.ActTimer, storyExplorer.Values.DisablesTimer = storyFisherman.Values.ActTimer, 1

		if storyJustContinued(panel) {
			storySetPanel("➡🟨Fishing, obviously?!\nWhat does it look like?", char.Explorer, char.Fisherman, emotion.Calm, emotion.Talk)
			storyNext()
		}
	case 10:
		storyFisherman.ActTimer, storyExplorer.Values.DisablesTimer = storyFisherman.Values.ActTimer, 1

		if storyJustContinued(panel) {
			storySetPanel("⬅Ha-ha-ha!", char.Explorer, char.Fisherman, emotion.Happy, emotion.Calm)
			storyNext()
		}
	case 11:
		storyFisherman.ActTimer, storyExplorer.Values.DisablesTimer = storyFisherman.Values.ActTimer, 1

		if storyJustContinued(panel) {
			storySetPanel("➡🟨What?\nWho are you anyway?", char.Explorer, char.Fisherman, emotion.Happy, emotion.Talk)
			storyNext()
		}
	case 12:
		storyFisherman.ActTimer, storyExplorer.Values.DisablesTimer = storyFisherman.Values.ActTimer, 1
	}

	if storyPortraitLeft != 0 {
		gui.Image(storyPortraitLeft, geometry.NewArea(panel.X-panel.Width/2+TileSize, panel.Y, TileSize*2, TileSize*3), false)
	}
	if storyPortraitRight != 0 {
		gui.Image(storyPortraitRight, geometry.NewArea(panel.X+panel.Width/2-TileSize, panel.Y, -TileSize*2, TileSize*3), false)
	}
	panel.Width, panel.X = panel.Width-TileSize*3.5, panel.X
	gui.Text(storyText.Get(), panel, geometry.Area{}, ThemeGUI, false)
}

// private ========================================================

var storyExplorer, storyFisherman *Unit
var story int
var storyTimer float32
var storyPortraitLeft, storyPortraitRight assets.ImageId
var storyText text.Dynamic

func storyNext() {
	storyTimer = 0
	story++
}
func storyWaitFor(seconds float32) {
	if storyTimer > seconds {
		storyNext()
	}
}
func storySetPanel(text string, charLeft, charRight char.Enum, emotionLeft, emotionRight emotion.Enum) {
	var left, right assets.ImageId = 0, 0
	if charLeft != char.NONE && emotionLeft != emotion.NONE {
		left = Characters[charLeft].Portraits[emotionLeft]
	}
	if charRight != char.NONE && emotionRight != emotion.NONE {
		right = Characters[charRight].Portraits[emotionRight]
	}
	storyPortraitLeft, storyPortraitRight = left, right
	storyText.Set(text)
}
func storyJustContinued(panel geometry.Area) bool {
	if panel.ContainsPoint(View.MousePosition()) {
		mouse.SetCursor(cursor.Hand)
		return mouse.IsAnyButtonJustPressed()
	}
	return keyboard.IsAnyKeyJustPressed() && !keyboard.IsKeyJustPressed(key.Escape)
}
