package game

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/audio"
	"pure-game-kit/packages/utility/random"
	"pure-game-kit/packages/utility/text"
)

var AudioAmbience audio.Audio
var AudioDoorCrumble, AudioDoorOpen, AudioDoorClose, AudioGateCrumble, AudioGateOpen, AudioGateClose []audio.Audio
var AudioBow, AudioProjectileGround, AudioProjectileWood, AudioProjectileMetal, AudioProjectileFlesh []audio.Audio
var AudioHitFlesh, AudioHitWood, AudioHitMetal []audio.Audio

func LoadAudio() {
	ambiences[EnvironmentPlains] = assets.LoadMusic("data/audio/ambience-plains.mp3")
	AudioAmbience = audio.New(ambiences[EnvironmentPlains])

	AudioDoorOpen = append(AudioDoorOpen, audio.New(assets.LoadSound("data/audio/door-open1.mp3", 4)))
	AudioDoorOpen = append(AudioDoorOpen, audio.New(assets.LoadSound("data/audio/door-open2.mp3", 4)))
	AudioDoorClose = append(AudioDoorClose, audio.New(assets.LoadSound("data/audio/door-close1.mp3", 4)))
	AudioDoorClose = append(AudioDoorClose, audio.New(assets.LoadSound("data/audio/door-close2.mp3", 4)))
	AudioGateOpen = append(AudioGateOpen, audio.New(assets.LoadSound("data/audio/gate-open1.mp3", 4)))
	AudioGateClose = append(AudioGateClose, audio.New(assets.LoadSound("data/audio/gate-close1.mp3", 4)))
	AudioGateClose = append(AudioGateClose, audio.New(assets.LoadSound("data/audio/gate-close2.mp3", 4)))
	for i := 1; i < 5; i++ {
		var door = assets.LoadSound(text.New("data/audio/door-crumble", i, ".mp3"), 4)
		var gate = assets.LoadSound(text.New("data/audio/gate-crumble", i, ".mp3"), 4)
		AudioDoorCrumble = append(AudioDoorCrumble, audio.New(door))
		AudioGateCrumble = append(AudioGateCrumble, audio.New(gate))
	}

	AudioHitWood = append(AudioHitWood, audio.New(assets.LoadSound("data/audio/hit-wood1.mp3", 4)))
	AudioHitWood = append(AudioHitWood, audio.New(assets.LoadSound("data/audio/hit-wood2.mp3", 4)))
	for i := 1; i < 6; i++ {
		AudioHitFlesh = append(AudioHitFlesh, audio.New(assets.LoadSound(text.New("data/audio/hit-flesh", i, ".mp3"), 4)))
	}
	for i := 1; i < 4; i++ {
		AudioHitMetal = append(AudioHitMetal, audio.New(assets.LoadSound(text.New("data/audio/hit-metal", i, ".mp3"), 4)))
	}

	AudioBow = append(AudioBow, audio.New(assets.LoadSound("data/audio/bow.mp3", 4)))
	AudioProjectileGround = append(AudioProjectileGround, audio.New(assets.LoadSound("data/audio/projectile-ground.mp3", 4)))
	AudioProjectileWood = append(AudioProjectileWood, audio.New(assets.LoadSound("data/audio/projectile-wood.mp3", 4)))
	AudioProjectileMetal = append(AudioProjectileMetal, audio.New(assets.LoadSound("data/audio/projectile-metal1.mp3", 4)))
	AudioProjectileMetal = append(AudioProjectileMetal, audio.New(assets.LoadSound("data/audio/projectile-metal2.mp3", 4)))
	for i := 1; i <= 7; i++ {
		var asset = assets.LoadSound(text.New("data/audio/projectile-flesh", i, ".mp3"), 4)
		AudioProjectileFlesh = append(AudioProjectileFlesh, audio.New(asset))
	}
}

func UpdateAudio() {
	if AudioAmbience.IsJustFinished() {
		AudioAmbience.Play()
	}
}

func PlaySound(variations []audio.Audio) {
	var sound = random.PickFrom(variations)
	sound.Play()
}
func PlayAmbience(environment Environment) {
	AudioAmbience.AssetId = ambiences[environment]
	AudioAmbience.ApplyProperties()
	AudioAmbience.Play()
}

// private ========================================================

var ambiences [1]assets.AudioId // index is Environment
