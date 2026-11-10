package game

import (
	"game/code/constants/icon"
	"game/code/constants/music"
	"game/code/constants/popup"
	"game/code/constants/screen"
	"game/code/constants/sound"
	"game/code/constants/zone"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/audio"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/gui"
	"pure-game-kit/packages/utility/file"
	"pure-game-kit/packages/utility/text"
)

func InitAudio() {
	loadAmbience(zone.Field, "field_S")
	loadAmbience(zone.Ruins, "ruins_S")
	loadAmbience(zone.Forest, "forest_S")
	loadAmbience(zone.Swamp, "swamp_S")
	loadAmbience(zone.Desert, "desert_S")
	loadAmbience(zone.Docks, "docks_S")
	loadAmbience(zone.Ocean, "ocean_PM")
	loadAmbience(zone.Glacier, "glacier_S")
	loadAmbience(zone.Cave, "cave_S")
	loadAmbience(zone.Mine, "mine_JFB")
	loadAmbience(zone.Hell, "hell_PM")

	loadMusic(music.StoryMap, "map_Childhood_Home_JS", "map_Pleasant_Peasants_JS")
	loadMusic(music.SceneEarly, "battle1_Duel_For_Honor_JS", "battle1_Funeral_March_JS", "battle1_Greatest_Nation_JS",
		"battle1_Nothing_Can_Stop_You_JS", "battle1_Overwhelming_Numbers_JS", "battle1_Protect_The_Realm_JS",
		"battle1_Shattered_Lands_JS", "battle1_The_Mad_King_JS", "battle1_Tournament_Time_JS")
	loadMusic(music.SceneLate, "battle2_Escape_Room_JS", "battle2_Hidden_Stairways_JS", "battle2_Maze_Trap_JS",
		"battle2_Sacred_Catacombs_JS")

	loadMusic(music.Menu, "menu_Homeland_Return_JS", "menu_Little_Village_JS", "menu_Old_Palace_JS")

	loadSound(sound.Button, 1, "ui_button")

	loadSound(sound.Play, 1, "play")

	loadSound(sound.DoorOpen, 4, "door_open")
	loadSound(sound.DoorClose, 4, "door_close")
	loadSound(sound.DoorCrumble, 4, "door_crumble")
	loadSound(sound.GateOpen, 4, "gate_open")
	loadSound(sound.GateClose, 4, "gate_close")
	loadSound(sound.GateCrumble, 4, "gate_crumble")

	loadSound(sound.GloryAlly, 4, "event_positive")
	loadSound(sound.GloryEnemy, 4, "event_negative")

	loadSound(sound.Drop, 4, "drop")
	loadSound(sound.Death, 4, "death")
	loadSound(sound.Revive, 4, "revive")

	loadSound(sound.SwingLight, 4, "swing_light")
	loadSound(sound.SwingHeavy, 4, "swing_heavy")
	loadSound(sound.HitFlesh, 4, "hit_flesh")
	loadSound(sound.HitWood, 4, "hit_wood")
	loadSound(sound.HitMetal, 4, "hit_metal")
	loadSound(sound.Stab, 4, "stab")
	loadSound(sound.Stomp, 2, "stomp")

	loadSound(sound.Bow, 4, "bow")
	loadSound(sound.ShotGround, 4, "shot_ground")
	loadSound(sound.ShotWood, 4, "shot_wood")
	loadSound(sound.ShotMetal, 4, "shot_metal")
	loadSound(sound.ShotFlesh, 4, "shot_flesh")

	loadSound(sound.Heal, 4, "heal")
}

//=================================================================

func UpdateMusic() {
	var curMusic = CurrentMusic()
	if Popup != popup.NONE && curMusic != nil {
		ShowCurrentMusic()
	}
	if curMusic == nil || curMusic.IsJustFinished() {
		if Screen == screen.Title {
			PlayMusic(music.Menu)
		} else if Screen == screen.StoryMap {
			PlayMusic(music.StoryMap)
		} else if Screen == screen.Scene {
			PlayMusic(music.SceneEarly)
		}
	}
}

func PlaySound(audio sound.Enum) {
	allSounds[audio].Play()
}
func PlayMusic(audio music.Enum) {
	PauseMusic()
	allMusic[audio].Play()
}
func PlayAmbience(zone zone.Enum) {
	PauseAmbience()
	allAmbiences[zone].Play()
}
func PauseAmbience() {
	for _, a := range allAmbiences {
		a.Pause()
	}
}
func PauseMusic() {
	for _, a := range allMusic {
		a.Pause()
	}
}
func PauseSounds() {
	for _, a := range allSounds {
		a.Pause()
	}
}
func StopSounds() {
	for _, a := range allSounds {
		a.Stop()
	}
}
func ResumeSounds() {
	for _, a := range allSounds {
		if a.IsPaused() {
			a.Play()
		}
	}
}

func CurrentMusic() *audio.Audio {
	for _, a := range allMusic {
		if a.IsPlaying() || a.IsJustFinished() {
			return a
		}
	}
	return nil
}
func CurrentAmbience() *audio.Audio {
	for _, a := range allAmbiences {
		if a.IsPlaying() {
			return a
		}
	}
	return nil
}

func ShowCurrentMusic() {
	if Settings.VolumeMusic == 0 {
		return
	}

	var curMusic = CurrentMusic()
	if curMusic != nil {
		gui.Text(musicNames[curMusic.CurrentAssetId()], gui.AreaHUD(1.0, 0.0, 400, 15), geometry.Area{}, ThemeGUI, false)
	}
}

func TryPlayButtonSound() {
	if gui.IsJustPressed() {
		PlaySound(sound.Button)
	}
}

// private ========================================================

var allSounds = make(map[sound.Enum]*audio.Audio)
var allMusic = make(map[music.Enum]*audio.Audio)
var allAmbiences = make(map[zone.Enum]*audio.Audio)
var musicNames = make(map[assets.AudioId]string)

func loadSound(sound sound.Enum, stacks uint8, baseName string) {
	var ids []assets.AudioId
	for i := 1; i <= 9; i++ {
		ids = tryLoadSound(ids, "data/sounds/"+baseName+text.New(i), stacks, "PM", "JFB", "S", "O", "JS")
	}
	allSounds[sound] = new(audio.New(assets.LoadAudioVariations(ids...)))
}
func loadMusic(music music.Enum, namesVariants ...string) {
	var ids = make([]assets.AudioId, len(namesVariants))
	for i, n := range namesVariants {
		var id = assets.LoadMusic("data/music/" + n + ".ogg")
		ids[i] = id

		if text.ContainsOneOf(n, "ambience") {
			n = ""
		} else {
			n = text.Remove(n, "menu_")
			n = text.Remove(n, "battle1_")
			n = text.Remove(n, "battle2_")
			n = text.Remove(n, "map_")
			n = text.Replace(n, "_", " ")
			n = text.Replace(n, "JS", "")
			n = "➡" + Tags[icon.Music] + n
		}
		musicNames[id] = n
	}
	allMusic[music] = new(audio.New(assets.LoadAudioVariations(ids...)))
}
func loadAmbience(zone zone.Enum, name string) {
	allAmbiences[zone] = new(audio.New(assets.LoadMusic("data/ambiences/" + name + ".ogg")))
}

func tryLoadSound(ids []assets.AudioId, path string, stacks uint8, suffixes ...string) []assets.AudioId {
	for _, s := range suffixes {
		var p = path + "_" + s + ".ogg"
		if file.Exists(p) {
			ids = append(ids, assets.LoadSound(stacks, p))
		}
	}
	return ids
}
