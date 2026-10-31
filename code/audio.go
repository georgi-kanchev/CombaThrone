package game

import (
	"game/code/constants/music"
	"game/code/constants/sound"
	"game/code/constants/zone"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/audio"
	"pure-game-kit/packages/utility/random"
)

func InitAudio() {
	loadMusic(music.Ambience, "ambience_field_S")

	loadSound(sound.DoorOpen, 4, "door_open1_PM", "door_open2_PM", "door_open3_PM")
	loadSound(sound.DoorOpen, 4, "door_close_S")
	loadSound(sound.DoorOpen, 4, "gate_open1_S", "gate_open2_S", "gate_open3_S")
	loadSound(sound.DoorOpen, 4, "gate_close1_S", "gate_close2_S", "gate_close3_S")
	loadSound(sound.DoorCrumble, 4, "door_crumble_O")
	loadSound(sound.GateCrumble, 4, "gate_crumble_O")

	loadSound(sound.MeleeFlesh, 4, "hit_flesh1", "hit_flesh2", "hit_flesh3")
	loadSound(sound.MeleeWood, 4, "hit_wood1", "hit_wood2")
	loadSound(sound.MeleeFlesh, 4, "hit_metal1")

	loadSound(sound.Bow, 4, "bow_S")
	loadSound(sound.ArrowGround, 4, "arrow_ground1")
	loadSound(sound.ArrowWood, 4, "arrow_wood1")
	loadSound(sound.ArrowMetal, 4, "arrow_metal1", "arrow_metal2")
	loadSound(sound.ArrowFlesh, 4, "arrow_flesh1", "arrow_flesh2", "arrow_flesh3", "arrow_flesh4")

	loadSound(sound.GloryAlly, 4, "event_positive1", "event_positive2")
	loadSound(sound.GloryEnemy, 4, "event_negative1", "event_negative2")
}

//=================================================================

func UpdateMusic() {
	for _, m := range allMusic {
		for _, a := range m {
			if a.IsJustFinished() {
				a.Play()
			}
		}
	}
}

func PlaySound(audio sound.Enum) {
	var variations, has = allSounds[audio]
	if has {
		var sound = random.PickFrom(variations)
		sound.Play()
	}
}
func PlayAmbience(zone zone.Enum) {
	allMusic[music.Ambience][zone].Play()
}

// private ========================================================

var allSounds = make(map[sound.Enum][]audio.Audio)
var allMusic = make(map[music.Enum][]audio.Audio)

func loadSound(sound sound.Enum, stacks uint8, namesVariants ...string) {
	for _, n := range namesVariants {
		allSounds[sound] = append(allSounds[sound], audio.New(assets.LoadSound("data/sounds/"+n+".mp3", stacks)))
	}
}
func loadMusic(music music.Enum, namesVariants ...string) {
	for _, n := range namesVariants {
		var path = "data/music/" + n + ".mp3"
		allMusic[music] = append(allMusic[music], audio.New(assets.LoadMusic(path)))
	}
}
