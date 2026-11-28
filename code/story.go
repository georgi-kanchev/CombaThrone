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
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/text"
	"pure-game-kit/packages/utility/time"
)

const STORY_HOVER, STORY_STATS, STORY_BENCH, STORY_SUMMON = 40, 43, 47, 49

var Story int

func UpdateStory() {

	var panel = geometry.NewArea(-TileSize*5, -TileSize*3.5, TileSize*7.5, TileSize*2.5)
	if Screen == screen.Scene {
		panel.X, panel.Y = 0, -TileSize/2
	}
	if storyText.Get() != "" {
		gui.Image(PanelId, panel, false)
	}

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
		storySetPanel("⬅🩵(This has to be it...)", char.Adventurer, 0, emotion.Calm, 0)
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
			Player.Units[0] = storyAdventurer
			Units = append(Units, storyAdventurer)
			storySetPanel("", 0, 0, 0, 0)
			storyNext()
		}
	case 8:
		storyFisherman.ActTimer = storyFisherman.Values.ActTimer
		if storyFisherman.State == state.Moving {
			storyFisherman.ActTimer = 0
		}

		if number.IsWithin(storyAdventurer.X, 0, TileSize*2) {
			storySetPanel("⬅🩵G'day, mister! Uhh, wh... what are you doing?",
				char.Adventurer, char.Fisherman, emotion.Talk, emotion.Calm)
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
			storySetPanel("⬅🩵Ha-ha-ha!", char.Adventurer, char.Fisherman, emotion.Happy, emotion.Calm)
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
			storySetPanel("⬅🩵I'm an adventurer by heart and archeologist by duty, sir.",
				char.Adventurer, char.Fisherman, emotion.Talk, emotion.Calm)
			storyNext()
		}
	case 13:
		storyFisherman.ActTimer, storyAdventurer.Values.DisablesTimer = storyFisherman.Values.ActTimer, 1
		if storyJustContinued(panel) {
			storySetPanel("⬅🩵I heard of the recent accident in the mines. Do you happen to know any of the miners?",
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
			storySetPanel("", 0, 0, 0, 0)
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
			"⬅🩵(These villagers must have lost their minds! Fishing on dry land? And I swear, these roads and doors make no sense...)",
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
			storySetPanel("⬅🩵Yes, sir! What happened down there?", char.Adventurer, char.Miner, emotion.Talk, emotion.Calm)
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
			storySetPanel("⬅🩵This is absolutely terrible!",
				char.Adventurer, char.Miner, emotion.Talk, emotion.Sad)
			storyNext()
		}
	case 26:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("➡🟨We need your help to understand more about the red crystals these... creatures were made of.",
				char.Adventurer, char.Miner, emotion.Calm, emotion.Angry)
			storyNext()
		}
	case 27:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("➡🟨We managed to bring this one home. Here, have a look.",
				char.Adventurer, char.Miner, emotion.Calm, emotion.Talk)
			var crystal = NewPickup(storyMiner.X-12, pickup.Crystal, lane.MiddleOff)
			Pickups = append(Pickups, crystal)
			Splashes = append(Splashes, NewSplash(splash.Pickup, team.Enemy, crystal.X, crystal.Y, nil))
			storyNext()
		}
	case 28:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("⬅🩵Hmm... looks just like a boring crystal to me, perhaps a ruby... or...",
				char.Adventurer, char.Miner, emotion.Talk, emotion.Calm)
			storyNext()
		}
	case 29:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("⬅🩵Did you bring more samples? I will have to break this one to do some tests on its internal structure.",
				char.Adventurer, char.Miner, emotion.Talk, emotion.Calm)
			storyNext()
		}
	case 30:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("➡🟨Oh, not really. You have to understand - we were busy running for our lives, ma'am!",
				char.Adventurer, char.Miner, emotion.Calm, emotion.Happy)
			storyNext()
		}
	case 31:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("⬅🩵Well, in that case, I have to visit the site to retrieve some more straight from the source.",
				char.Adventurer, char.Miner, emotion.Talk, emotion.Calm)
			storyNext()
		}
	case 32:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("➡🟨Don't be silly, dear! This horrible place is deadly for our strongest men, let alone a lady.",
				char.Adventurer, char.Miner, emotion.Calm, emotion.Talk)
			storyNext()
		}
	case 33:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("⬅🩵You can't scare an advanturer like me! Bring me to that place!",
				char.Adventurer, char.Miner, emotion.Talk, emotion.Calm)
			storyNext()
		}
	case 34:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("➡🟨If I can't change your mind, I can at least help you prepare for what's out there. What can you do?",
				char.Adventurer, char.Miner, emotion.Calm, emotion.Talk)
			storyNext()
		}
	case 35:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("⬅🩵Oh! Oh! I can play this whistle! Listen...", char.Adventurer, char.Miner, emotion.Talk, emotion.Calm)
			PlaySound(sound.Whistle)
			storyNext()
		}
	case 36:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("➡🟨I... I meant for fighting those creatures, dummy. Show me your combat skills.",
				char.Adventurer, char.Miner, emotion.Calm, emotion.Talk)
			storyNext()
		}
	case 37:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storyAdventurer.Values.DisablesTimer = 0
			storySetPanel("", 0, 0, 0, 0)
			storyNext()
		}
	case 38:
		storyMiner.Values.DisablesTimer = 1
		if storyAdventurer.X > storyMiner.X && !number.IsWithin(storyAdventurer.X, storyMiner.X, TileSize) {
			storyNext()
		}
	case 39:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		storySetPanel("Every "+Tags[icon.Unit]+"unit has a skill they can perform. Hover over them "+
			"with your mouse cursor, to understand what they do. Try it!", 0, 0, 0, 0)
		storyNext()
	case STORY_HOVER:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("Click on a "+Tags[icon.Unit]+"unit to keep their hover info visible at all times (indicated by "+
				Tags[icon.Locked]+").",
				0, 0, 0, 0)
			storyNext()
		}
	case 41:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storyAdventurer.IsReturning, storyMiner.IsReturning = true, true
			storySetPanel("⬅🟨Not too bad! That move will definitely help you during combat.",
				char.Miner, char.Adventurer, emotion.Talk, emotion.Calm)
			storyNext()
		}
	case 42:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("⬅🟨Let's test your reactions. THINK FAST!!", char.Miner, char.Adventurer, emotion.Talk, emotion.Calm)
			storyDummy = NewUnit(char.Dummy, team.Enemy, lane.Middle)
			Units = append(Units, storyDummy)
			storyNext()
		}
	case STORY_STATS:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		storyDummy.X = storyAdventurer.X + 12
		storyNext()
	case 44:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyAdventurer.UnitBehind != nil {
			storyAdventurer.AffectHealth(-(storyAdventurer.Health + storyAdventurer.Values.ActPoints), true, storyDummy)
			storyNext()
		}
	case 45:
		storyAdventurer.HurtTimer, storyMiner.Values.DisablesTimer = 0, 1
		if storyJustContinued(panel) {
			storySetPanel("➡🩵Ouch! Where did that come from?! Wait... am I... DEAD??",
				char.Miner, char.Adventurer, emotion.Calm, emotion.Sad)
			storyNext()
		}
	case 46:
		storyAdventurer.HurtTimer, storyMiner.Values.DisablesTimer = 0, 1
		if storyJustContinued(panel) {
			storySetPanel("⬅🟨Danger can come from anywhere at any time, milady. "+
				"Death is part of life. But we live on in other forms... and the beautiful cycle continues...",
				char.Miner, char.Adventurer, emotion.Talk, emotion.Sad)
			storyNext()
		}
	case STORY_BENCH:
		storyAdventurer.HurtTimer, storyMiner.Values.DisablesTimer = 0, 1
		if storyJustContinued(panel) {
			storySetPanel("🟥"+Tags[icon.Skull]+"Dead⬜ characters 🌗🟦"+Tags[icon.Revive]+"revive⬜ after some time, "+
				"indicated above the "+Tags[icon.Unit]+"unit, upon hover and also in their slot in the panel above.", 0, 0, 0, 0)
			storyNext()
		}
	case 48:
		storyMiner.Values.DisablesTimer = 1
		if storyAdventurer.State == state.WaitingToBeSummoned {
			storySetPanel("Once 🌗🟦"+Tags[icon.Revive]+"revived⬜, the character may be summoned again by clicking or dragging, "+
				"then dropping on the desired lane.", 0, 0, 0, 0)
			storyNext()
		}
	case STORY_SUMMON:
		storyMiner.Values.DisablesTimer = 1
		if condition.TrueEvery(2, 0xCFA) {
			if SceneHUD.SummonIndex == 0 {
				var summon = Bases[team.Ally].Entrances[1].SummonIndicator()
				PlayCallForAction(summon.X, summon.Y)
			} else {
				var x, y = SceneHUD.UnitIconPosition(0)
				PlayCallForAction(View.PointFromView(SceneHUD.View, x, y))
			}
		}

		if storyAdventurer.State == state.Summoned {
			storyAdventurer.Values.MoveSpeed = 30
			storySetPanel("", 0, 0, 0, 0)
			storyNext()
		}
	case 50:
		storyMiner.Values.DisablesTimer = 1
		if number.IsWithin(storyAdventurer.X, storyMiner.X, TileSize) {
			storyMiner.IsReturning = false
			storySetPanel("⬅🩵Hey, that wasn't very nice!", char.Adventurer, char.Miner, emotion.Talk, emotion.Calm)
			storyNext()
		}
	case 51:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("➡🟨Learning the hard way is never nice but always necessary, dear. You need to get used to dying - so is the life of the ⬜"+Tags[icon.Defender]+"Defender🟨.", char.Adventurer, char.Miner, emotion.Calm, emotion.Talk)
			storyNext()
		}
	case 52:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("⬅🩵Me? A ⬜"+Tags[icon.Defender]+"Defender🩵? Cool, I like it! What does my role do?",
				char.Adventurer, char.Miner, emotion.Happy, emotion.Calm)
			storyNext()
		}
	case 53:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("➡🟨All ⬜"+Tags[icon.Defender]+"Defenders🟨 participate in the combat directly, alongside the ⬜"+
				Tags[icon.Fighter]+"Fighters🟨 and ⬜"+Tags[icon.Backliner]+"Backliners🟨.",
				char.Adventurer, char.Miner, emotion.Calm, emotion.Talk)
			storyNext()
		}
	case 54:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("➡⬜"+Tags[icon.Defender]+"Defenders🟨 push or slow lanes down. They 🟧"+
				Tags[icon.Defender]+"block🟨 some incoming 🟧"+Tags[icon.Wound]+"damage🟨 from the ⬜"+Tags[icon.Fighter]+
				"Fighters🟨 and ⬜"+Tags[icon.Backliner]+"Backliners🟨 but don't cause any themselves.",
				char.Adventurer, char.Miner, emotion.Calm, emotion.Talk)
			storyNext()
		}
	case 55:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("⬅🩵Awesome! Wait, so if everybody 🌗🟦"+Tags[icon.Revive]+"revives🩵, what is our combat goal then?",
				char.Adventurer, char.Miner, emotion.Talk, emotion.Calm)
			storyNext()
		}
	case 56:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("➡🟨We fight for ⬜"+Tags[icon.Glory]+"Glory🟨, madam! We lose some of it when an 🟥"+
				Tags[icon.Unit]+"enemy🟨 or a stranger like you enters our base. Try it - get in our barrack! We don't bite!",
				char.Adventurer, char.Miner, emotion.Calm, emotion.Happy)
			storyNext()
		}
	case 57:
		storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
		if storyJustContinued(panel) {
			storySetPanel("➡🟨Did I hear \"bite\"? The fish must be calling me!", 0, char.Fisherman, 0, emotion.Happy)
			Units = collection.Remove(Units, storyFisherman)
			storyFisherman = NewUnit(char.Fisherman, team.Enemy, lane.LowerOff)
			storyFisherman.Values.MoveSpeed = 30
			Units = append(Units, storyFisherman)
			storyNext()
		}
	case 58:
		storyMiner.Values.DisablesTimer, storyFisherman.ActTimer = 1, storyFisherman.Values.ActTimer
		storyWaitFor(5.0)
	case 59:
		storyFisherman.ActTimer = 0
		storyNext()
	case 60:
		storyMiner.Values.DisablesTimer, storyFisherman.ActTimer = 1, storyFisherman.Values.ActTimer
		if storyFisherman.State == state.ActPreparing {
			storyNext()
		}
	case 61:
		storyMiner.Values.DisablesTimer, storyFisherman.ActTimer = 1, storyFisherman.Values.ActTimer
		if storyAdventurer.State == state.Decaying {
			storySetPanel("🟩"+Tags[icon.Unit]+"Allies⬜ that reach the 🟥enemy base⬜ will 🌗🟦"+Tags[icon.Revive]+
				"return⬜ after some time, as if they were to 🌗🟦"+Tags[icon.Revive]+"revive⬜. Their remaining 🌗🟩"+
				Tags[icon.Health]+"health⬜ is subtracted from the 🟥enemy ⬜"+Tags[icon.Glory]+"Glory. The first team to have no "+
				Tags[icon.Glory]+"Glory loses the battle.",
				0, 0, 0, 0)
			storyNext()
		}
	case 62:
		storyMiner.Values.DisablesTimer, storyFisherman.ActTimer = 1, storyFisherman.Values.ActTimer
		if storyAdventurer.State == state.WaitingToBeSummoned {
			Player.Units[1] = NewUnit(char.Bowyer, team.Ally, 0)
			Player.Units[2] = NewUnit(char.Smith, team.Ally, 0)
			Player.Units[3] = NewUnit(char.Cook, team.Ally, 0)
			Player.Units[4] = NewUnit(char.Kid, team.Ally, 0)

			storySetPanel("", 0, 0, 0, 0)
			Units = append(Units, storyAdventurer)
			storyAdventurer.Lane = lane.Middle
			storyAdventurer.PrepareSpawn()
			storyNext()
		}
	case 63:
		storyMiner.Values.DisablesTimer, storyFisherman.ActTimer = 1, storyFisherman.Values.ActTimer
		if number.IsWithin(storyAdventurer.X, storyMiner.X, TileSize) {
			storySetPanel("⬅🩵I made some friends inside. They agreed to help me on my upcoming journey you are preparing me for. "+
				"My whistle is our signal!",
				char.Adventurer, char.Miner, emotion.Talk, emotion.Calm)
			storyNext()
		}
	case 64:
		storyMiner.Values.DisablesTimer, storyFisherman.ActTimer = 1, storyFisherman.Values.ActTimer
		storyAdventurer.Values.DisablesTimer = 1
	}

	// storyAdventurer.Values.DisablesTimer, storyMiner.Values.DisablesTimer = 1, 1
	// if storyJustContinued(panel) {
	// 	storySetPanel("There are 6 roles in total:\n"+
	// 		Tags[icon.Fighter]+"Fighter, "+Tags[icon.Backliner]+"Backliner, "+Tags[icon.Defender]+"Defender,\n"+
	// 		Tags[icon.Griefer]+"Griefer, "+Tags[icon.Supplier]+"Supplier & "+Tags[icon.Collector]+"Collector.\n"+
	// 		"The first 3 participate in the combat directly, while the second 3 do so indirectly.\n"+
	// 		"The "+Tags[icon.Unit]+"unit hover info shows their role.", 0, 0, 0, 0)
	// 	storyNext()
	// }

	// 	storyAdventurer.HurtTimer, storyMiner.Values.DisablesTimer = 0, 1
	// 	if storyJustContinued(panel) {
	// 		storyMiner.IsReturning = false
	// 		storySetPanel(Tags[icon.Unit]+"Units that reach the end of the lane and enter the enemy base will "+
	// 			"🌗🟦"+Tags[icon.Revive]+"return⬜ after some time.", 0, 0, 0, 0)
	// 		storyNext()

	if storyPortraitLeft != 0 {
		gui.Image(storyPortraitLeft, geometry.NewArea(panel.X-panel.Width/2+TileSize, panel.Y, TileSize*2, TileSize*3), false)
	}
	if storyPortraitRight != 0 {
		gui.Image(storyPortraitRight, geometry.NewArea(panel.X+panel.Width/2-TileSize, panel.Y, -TileSize*2, TileSize*3), false)
	}
	if storyPortraitLeft != 0 || storyPortraitRight != 0 {
		panel.Width = panel.Width - TileSize*3.5
	} else {
		panel.Width *= 0.9
	}

	gui.Text(storyText.Get(), panel, geometry.Area{}, ThemeGUI, false)
}

// private ========================================================

var storyAdventurer, storyFisherman, storyMiner, storyDummy *Unit
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
