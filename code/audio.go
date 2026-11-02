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
	"pure-game-kit/packages/utility/random"
	"pure-game-kit/packages/utility/text"
)

func InitAudio() {
	loadMusic(music.Ambience, "ambience_field_S")
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
	var variations, has = allSounds[audio]
	if has {
		var sound = random.PickFrom(variations)
		sound.Play()
	}
}
func PlayMusic(audio music.Enum) {
	var variations, has = allMusic[audio]
	for _, m := range allMusic {
		for _, a := range m {
			a.Pause()
		}
	}
	if has {
		var music = random.PickFrom(variations)
		music.Play()
	}
}
func PlayAmbience(zone zone.Enum) {
	allMusic[music.Ambience][zone].Play()
}
func StopMusic() {
	for _, v := range allMusic {
		for _, a := range v {
			a.Pause()
		}
	}
}

func CurrentMusic() *audio.Audio {
	for _, variations := range allMusic {
		for _, v := range variations {
			if v.IsPlaying() {
				return v
			}
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

var allSounds = make(map[sound.Enum][]*audio.Audio)
var allMusic = make(map[music.Enum][]*audio.Audio)
var musicNames = make(map[assets.AudioId]string)

func loadSound(sound sound.Enum, stacks uint8, namesVariants ...string) {
	for _, n := range namesVariants {
		allSounds[sound] = append(allSounds[sound], new(audio.New(assets.LoadSound(stacks, "data/sounds/"+n+".ogg"))))
	}
}
func loadMusic(music music.Enum, namesVariants ...string) {
	for _, n := range namesVariants {
		var id = assets.LoadMusic("data/music/" + n + ".ogg")
		allMusic[music] = append(allMusic[music], new(audio.New(id)))

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
}
