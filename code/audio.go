package game

import (
	"game/code/constants/icon"
	"game/code/constants/music"
	"game/code/constants/sound"
	"game/code/constants/zone"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/audio"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/gui"
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

	loadMusic(music.Menu, "menu_Homeland_Return_JS", "menu_Little_Village_JS", "menu_Old_Palace_JS")

	loadSound(sound.DoorOpen, 4, "door_open1_PM", "door_open2_PM", "door_open3_PM")
	loadSound(sound.DoorClose, 4, "door_close_S")
	loadSound(sound.DoorCrumble, 4, "door_crumble1_O", "door_crumble2_O")
	loadSound(sound.GateOpen, 4, "gate_open1_S", "gate_open2_S", "gate_open3_S")
	loadSound(sound.GateClose, 4, "gate_close1_S", "gate_close2_S", "gate_close3_S")
	loadSound(sound.GateCrumble, 4, "gate_crumble1_O", "gate_crumble2_O")

	loadSound(sound.Drop, 4, "drop1_S")
	loadSound(sound.Death, 4, "death1_JFB", "death2_JFB")

	loadSound(sound.HitFlesh, 4, "hit_flesh1_PM", "hit_flesh2_PM", "hit_flesh3_PM")
	loadSound(sound.HitWood, 4, "hit_wood1_S", "hit_wood2_S")
	loadSound(sound.HitMetal, 4, "hit_metal1_PM", "hit_metal2_S")
	loadSound(sound.Stab, 4, "stab1_PM", "stab2_PM")

	loadSound(sound.Bow, 4, "bow1_S")
	loadSound(sound.ShotGround, 4, "proj_ground1_PM")
	loadSound(sound.ShotWood, 4, "proj_wood1_S")
	loadSound(sound.ShotMetal, 4, "proj_metal1_PM", "proj_metal2_PM")
	loadSound(sound.ShotFlesh, 4, "proj_flesh1_PM", "proj_flesh2_PM", "proj_flesh3_PM", "proj_flesh4_PM")

	loadSound(sound.GloryAlly, 4, "event_positive1_JFB", "event_positive2_JFB")
	loadSound(sound.GloryEnemy, 4, "event_negative1_JFB", "event_negative2_JFB")
}

//=================================================================

func PlaySound(audio sound.Enum) {
	allSounds[audio].Play()
}
func PlayMusic(audio music.Enum) {
	for _, a := range allMusic {
		a.Stop()
	}
	allMusic[audio].Play()
}
func PlayAmbience(zone zone.Enum) {
	for _, a := range allAmbiences {
		a.Stop()
	}
	allAmbiences[zone].Play()
}
func StopMusic() {
	for _, a := range allMusic {
		a.Stop()
	}
}

func CurrentMusic() *audio.Audio {
	for _, a := range allMusic {
		if a.IsPlaying() {
			return a
		}
	}
	return nil
}

func ShowCurrentMusic() {
	var curMusic = CurrentMusic()
	if curMusic != nil {
		gui.Text(musicNames[curMusic.AssetId], gui.AreaHUD(1.0, 0.0, 400, 15), geometry.Area{}, ThemeGUI, false)
	}
}

// private ========================================================

var allSounds = make(map[sound.Enum]*audio.Audio)
var allMusic = make(map[music.Enum]*audio.Audio)
var allAmbiences = make(map[zone.Enum]*audio.Audio)
var musicNames = make(map[assets.AudioId]string)

func loadSound(sound sound.Enum, stacks uint8, namesVariants ...string) {
	var ids = make([]assets.AudioId, len(namesVariants))
	for i, n := range namesVariants {
		var id = assets.LoadSound(stacks, "data/sounds/"+n+".ogg")
		ids[i] = id
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
			n = text.Replace(n, "_", " ")
			n = text.Replace(n, "JS", "(by Joel Steudler)")
			n = "➡" + Tags[icon.Music] + n + " "
		}
		musicNames[id] = n
	}
	allMusic[music] = new(audio.New(assets.LoadAudioVariations(ids...)))
}
func loadAmbience(zone zone.Enum, name string) {
	allAmbiences[zone] = new(audio.New(assets.LoadMusic("data/ambiences/" + name + ".ogg")))
}
