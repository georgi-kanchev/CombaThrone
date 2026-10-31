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

	loadSound(sound.DoorOpen, 4, "door-open1", "door-open2")
	loadSound(sound.DoorOpen, 4, "door-close1", "door-close2")
	loadSound(sound.DoorOpen, 4, "gate-open1")
	loadSound(sound.DoorOpen, 4, "gate-close1", "gate-close2")
	loadSound(sound.DoorCrumble, 4, "door-crumble1", "door-crumble2", "door-crumble3", "door-crumble4", "door-crumble5")
	loadSound(sound.GateCrumble, 4, "gate-crumble1", "gate-crumble2", "gate-crumble3", "gate-crumble4", "gate-crumble5")

	loadSound(sound.MeleeFlesh, 4, "hit-flesh1", "hit-flesh2", "hit-flesh3", "hit-flesh4", "hit-flesh5", "hit-flesh6")
	loadSound(sound.MeleeWood, 4, "hit-wood1", "hit-wood2")
	loadSound(sound.MeleeFlesh, 4, "hit-metal1", "hit-metal2", "hit-metal3", "hit-metal4")

	loadSound(sound.Bow, 4, "bow")
	loadSound(sound.ArrowGround, 4, "projectile-ground")
	loadSound(sound.ArrowWood, 4, "projectile-wood")
	loadSound(sound.ArrowMetal, 4, "projectile-metal1", "projectile-metal2")
	loadSound(sound.ArrowFlesh, 4, "projectile-flesh1", "projectile-flesh2", "projectile-flesh3", "projectile-flesh4",
		"projectile-flesh5", "projectile-flesh6", "projectile-flesh7")

	loadSound(sound.GloryAlly, 4, "event-positive1", "event-positive2")
	loadSound(sound.GloryEnemy, 4, "event-negative1", "event-negative2")
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
		allMusic[music] = append(allMusic[music], audio.New(assets.LoadMusic("data/music/"+n+".mp3")))
	}
}
