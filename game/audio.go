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
var AudioEventPositive, AudioEventNegative []audio.Audio

func LoadAudio() {
	ambiences[ZoneField] = assets.LoadMusic("data/music/ambience-plains.mp3")
	AudioAmbience = audio.New(ambiences[ZoneField])

	AudioDoorOpen = append(AudioDoorOpen, audio.New(assets.LoadSound("data/sounds/door-open1.mp3", 4)))
	AudioDoorOpen = append(AudioDoorOpen, audio.New(assets.LoadSound("data/sounds/door-open2.mp3", 4)))
	AudioDoorClose = append(AudioDoorClose, audio.New(assets.LoadSound("data/sounds/door-close1.mp3", 4)))
	AudioDoorClose = append(AudioDoorClose, audio.New(assets.LoadSound("data/sounds/door-close2.mp3", 4)))
	AudioGateOpen = append(AudioGateOpen, audio.New(assets.LoadSound("data/sounds/gate-open1.mp3", 4)))
	AudioGateClose = append(AudioGateClose, audio.New(assets.LoadSound("data/sounds/gate-close1.mp3", 4)))
	AudioGateClose = append(AudioGateClose, audio.New(assets.LoadSound("data/sounds/gate-close2.mp3", 4)))
	for i := 1; i < 5; i++ {
		var door = assets.LoadSound(text.New("data/sounds/door-crumble", i, ".mp3"), 4)
		var gate = assets.LoadSound(text.New("data/sounds/gate-crumble", i, ".mp3"), 4)
		AudioDoorCrumble = append(AudioDoorCrumble, audio.New(door))
		AudioGateCrumble = append(AudioGateCrumble, audio.New(gate))
	}

	AudioHitWood = append(AudioHitWood, audio.New(assets.LoadSound("data/sounds/hit-wood1.mp3", 4)))
	AudioHitWood = append(AudioHitWood, audio.New(assets.LoadSound("data/sounds/hit-wood2.mp3", 4)))
	for i := 1; i < 6; i++ {
		AudioHitFlesh = append(AudioHitFlesh, audio.New(assets.LoadSound(text.New("data/sounds/hit-flesh", i, ".mp3"), 4)))
	}
	for i := 1; i < 4; i++ {
		AudioHitMetal = append(AudioHitMetal, audio.New(assets.LoadSound(text.New("data/sounds/hit-metal", i, ".mp3"), 4)))
	}

	AudioBow = append(AudioBow, audio.New(assets.LoadSound("data/sounds/bow.mp3", 4)))
	AudioProjectileGround = append(AudioProjectileGround, audio.New(assets.LoadSound("data/sounds/projectile-ground.mp3", 4)))
	AudioProjectileWood = append(AudioProjectileWood, audio.New(assets.LoadSound("data/sounds/projectile-wood.mp3", 4)))
	AudioProjectileMetal = append(AudioProjectileMetal, audio.New(assets.LoadSound("data/sounds/projectile-metal1.mp3", 4)))
	AudioProjectileMetal = append(AudioProjectileMetal, audio.New(assets.LoadSound("data/sounds/projectile-metal2.mp3", 4)))
	for i := 1; i <= 7; i++ {
		var asset = assets.LoadSound(text.New("data/sounds/projectile-flesh", i, ".mp3"), 4)
		AudioProjectileFlesh = append(AudioProjectileFlesh, audio.New(asset))
	}

	AudioEventPositive = append(AudioEventPositive, audio.New(assets.LoadSound("data/sounds/event-positive1.mp3", 4)))
	AudioEventPositive = append(AudioEventPositive, audio.New(assets.LoadSound("data/sounds/event-positive2.mp3", 4)))
	AudioEventNegative = append(AudioEventNegative, audio.New(assets.LoadSound("data/sounds/event-negative1.mp3", 4)))
	AudioEventNegative = append(AudioEventNegative, audio.New(assets.LoadSound("data/sounds/event-negative2.mp3", 4)))
}

//=================================================================

func UpdateAudio() {
	if AudioAmbience.IsJustFinished() {
		AudioAmbience.Play()
	}
}

func PlaySound(variations []audio.Audio) {
	if len(variations) == 0 {
		return
	}
	var sound = random.PickFrom(variations)
	sound.Play()
}
func PlayAmbience(zone ZoneKind) {
	AudioAmbience.AssetId = ambiences[zone]
	AudioAmbience.ApplyProperties()
	AudioAmbience.Play()
}

// private ========================================================

var ambiences [ZoneCount]assets.AudioId
