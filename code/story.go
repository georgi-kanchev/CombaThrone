package game

import (
	char "game/code/constants/character"
	"game/code/constants/emotion"
	"game/code/constants/icon"
	"game/code/constants/lane"
	"game/code/constants/pickup"
	"game/code/constants/progress"
	"game/code/constants/screen"
	"game/code/constants/sound"
	"game/code/constants/splash"
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
	"pure-game-kit/packages/utility/time"
)

var Story int

func UpdateStory() {
	var panel = geometry.NewArea(-TileSize*5, -TileSize*3.5, TileSize*7.5, TileSize*2.5)
	if Screen == screen.Scene {
		panel.X, panel.Y = 0, -TileSize/2
	}
	gui.Image(PanelId, panel, false)

	storyTimer += DeltaTimeScaled()

	switch Story {
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
		storySetPanel("⬅(This has to be it...)", char.Adventurer, 0, emotion.Calm, 0)
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
		storyFisherman.X, storyFisherman.Values.MoveSpeed = 0, 30
		Units = append(Units, storyFisherman)

		storySetPanel("When indicated by "+Tags[icon.ArrowRight]+", click on this panel or press any key to continue. ", 0, 0, 0, 0)

		PauseMusic()
		PauseAmbience()
		storyNext()
	case 7:
		storyFisherman.ActTimer = storyFisherman.Values.ActTimer
		if storyFisherman.State == state.Moving {
			storyFisherman.ActTimer = 0
		}

		if storyJustContinued(panel) {
			storyAdventurer = NewUnit(char.Adventurer, team.Ally, lane.Middle)
			storyAdventurer.Values.MoveSpeed = 30
			Units = append(Units, storyAdventurer)
			storySetPanel(". . .", 0, 0, 0, 0)
			storyNext()
		}
	case 8:
		storyFisherman.ActTimer = storyFisherman.Values.ActTimer
		if storyFisherman.State == state.Moving {
			storyFisherman.ActTimer = 0
		}

		if number.IsWithin(storyAdventurer.X, 0, TileSize*2) {
			storySetPanel("⬅G'day, mister! Uhh, wh... what are you doing?", char.Adventurer, char.Fisherman, emotion.Talk, emotion.Calm)
			storyNext()
		}
	case 9:
		storyFisherman.ActTimer, storyAdventurer.Values.DisablesTimer = storyFisherman.Values.ActTimer, 1
		if storyJustContinued(panel) {
			storySetPanel("➡🟨Fishing, obviously?!\nWhat does it look like?", char.Adventurer, char.Fisherman, emotion.Calm, emotion.Talk)
			storyNext()
		}
	case 10:
		storyFisherman.ActTimer, storyAdventurer.Values.DisablesTimer = storyFisherman.Values.ActTimer, 1
		if storyJustContinued(panel) {
			storySetPanel("⬅Ha-ha-ha!", char.Adventurer, char.Fisherman, emotion.Happy, emotion.Calm)
			storyNext()
		}
	case 11:
		storyFisherman.ActTimer, storyAdventurer.Values.DisablesTimer = storyFisherman.Values.ActTimer, 1
		if storyJustContinued(panel) {
			storySetPanel("➡🟨What's so funny?\nWho are you anyway?", char.Adventurer, char.Fisherman, emotion.Happy, emotion.Talk)
			storyNext()
		}
	case 12:
		storyFisherman.ActTimer, storyAdventurer.Values.DisablesTimer = storyFisherman.Values.ActTimer, 1
		if storyJustContinued(panel) {
			storySetPanel("⬅I'm an adventurer by heart and archeologist by duty, sir.",
				char.Adventurer, char.Fisherman, emotion.Talk, emotion.Calm)
			storyNext()
		}
	case 13:
		storyFisherman.ActTimer, storyAdventurer.Values.DisablesTimer = storyFisherman.Values.ActTimer, 1
		if storyJustContinued(panel) {
			storySetPanel("⬅I heard of the recent accident in the mines. Do you happen to know any of the miners?",
				char.Adventurer, char.Fisherman, emotion.Talk, emotion.Calm)
			storyNext()
		}
	case 14:
		storyFisherman.ActTimer, storyAdventurer.Values.DisablesTimer = storyFisherman.Values.ActTimer, 1
		if storyJustContinued(panel) {
			storySetPanel("➡🟨Yes! One of them is my brother. I'm glad that lucky bastard was unscratched. He is my only family.",
				char.Adventurer, char.Fisherman, emotion.Calm, emotion.Talk)
			storyNext()
		}
	case 15:
		storyFisherman.ActTimer, storyAdventurer.Values.DisablesTimer = storyFisherman.Values.ActTimer, 1
		if storyJustContinued(panel) {
			storySetPanel("➡🟨Wait here, I will go bring him. Not getting any bites anyway.",
				char.Adventurer, char.Fisherman, emotion.Calm, emotion.Talk)
			storyNext()
		}
	case 16:
		storyFisherman.ActTimer, storyAdventurer.Values.DisablesTimer = storyFisherman.Values.ActTimer, 1
		if storyJustContinued(panel) {
			storyFisherman.ActTimer = 0
			storySetPanel(". . .", 0, 0, 0, 0)
			storyNext()
		}
	case 17:
		storyFisherman.ActTimer, storyAdventurer.Values.DisablesTimer = storyFisherman.Values.ActTimer, 1
		storyWaitFor(0.65)
	case 18:
		storyFisherman.IsReturning = true
		storyNext()
	case 19:
		storyFisherman.ActTimer, storyAdventurer.Values.DisablesTimer = storyFisherman.Values.ActTimer, 1
		storyWaitFor(4.0)
	case 20:
		storySetPanel(
			"⬅(These villagers must have lost their minds! Fishing on dry land? And I swear, these roads and doors make no sense...)",
			char.Adventurer, 0, emotion.Calm, 0)
		storyNext()
	case 21:
		storyFisherman.ActTimer, storyAdventurer.Values.DisablesTimer = storyFisherman.Values.ActTimer, 1
		if !storyFisherman.IsOutsideOwnBase() {
			storyMiner = NewUnit(char.Miner, team.Enemy, lane.Middle)
			storyMiner.Values.MoveSpeed = 30
			Units = append(Units, storyMiner)
			storyNext()
		}
	case 22:
		storyAdventurer.Values.DisablesTimer = 1
		if number.IsWithin(storyMiner.X, -TileSize, 8) {
			storySetPanel("➡🟨G'day, milady! Don't mind my delusional brother. So you've received our letter.",
				char.Adventurer, char.Miner, emotion.Calm, emotion.Talk)
			storyNext()
		}
	case 23:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("⬅Yes, sir! What happened down there?", char.Adventurer, char.Miner, emotion.Talk, emotion.Calm)
			storyNext()
		}
	case 24:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel(
				"➡🟨Oh, dear! We struck a vein of horrors! Great evil chased us away! Some of us didn't make it out al... alive...",
				char.Adventurer, char.Miner, emotion.Talk, emotion.Sad)
			storyNext()
		}
	case 25:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("⬅This is absolutely terrible! I should investigate and bring peace upon your people!",
				char.Adventurer, char.Miner, emotion.Talk, emotion.Sad)
			storyNext()
		}
	case 26:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("➡🟨Don't be silly, dear! This horrible place is deadly for our strongest men, let alone a lady.",
				char.Adventurer, char.Miner, emotion.Calm, emotion.Sad)
			storyNext()
		}
	case 27:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("➡🟨We need your help to understand more about the red crystals these... creatures were made of.",
				char.Adventurer, char.Miner, emotion.Calm, emotion.Angry)
			storyNext()
		}
	case 28:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("➡🟨We managed to bring this one home. Have a look.", char.Adventurer, char.Miner, emotion.Calm, emotion.Calm)
			var crystal = NewPickup(storyMiner.X-12, pickup.Crystal, lane.MiddleOff)
			Pickups = append(Pickups, crystal)
			Splashes = append(Splashes, NewSplash(splash.Pickup, team.Enemy, crystal.X, crystal.Y, nil))
			storyNext()
		}
	case 29:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("⬅Hmm... looks like a regular boring crystal to me, perhaps a ruby... but not quite the same.",
				char.Adventurer, char.Miner, emotion.Calm, emotion.Calm)
			storyNext()
		}
	case 30:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("⬅Did you bring more samples? I will have to break this one to do some tests on its internal structure.",
				char.Adventurer, char.Miner, emotion.Calm, emotion.Calm)
			storyNext()
		}
	case 31:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("➡🟨Oh, not really. You have to understand, we were busy running for our lives, ma'am!",
				char.Adventurer, char.Miner, emotion.Calm, emotion.Happy)
			storyNext()
		}
	case 32:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
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

var storyAdventurer, storyFisherman, storyMiner *Unit
var storyTimer float32
var storyPortraitLeft, storyPortraitRight assets.ImageId
var storyText text.Dynamic

func storyNext() {
	storyTimer = 0
	Story++
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
	var indicator = geometry.NewArea(panel.X+panel.Width/2-8, panel.Y+panel.Height/2-8, TileSize, TileSize)
	var scale = number.Map(number.Sine(time.Running()*3), -1, 1, 0.3, 0.5)
	gui.Text(Tags[icon.ArrowRight], indicator, geometry.Area{}, ThemeGUI, false, scale)

	if TimeScale == 0 {
		return false
	}

	var result = keyboard.IsAnyKeyJustPressed() && !keyboard.IsKeyJustPressed(key.Escape)
	if panel.ContainsPoint(View.MousePosition()) {
		mouse.SetCursor(cursor.Hand)
		if mouse.IsAnyButtonJustPressed() {
			result = true
		}
	}
	if result {
		PlaySound(sound.Button)
	}
	return result
}
