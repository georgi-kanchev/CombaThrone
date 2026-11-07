package sound

type Enum uint16

const (
	Ambience Enum = iota

	Button
	Play

	DoorOpen
	DoorClose
	DoorCrumble
	GateOpen
	GateClose
	GateCrumble

	GloryAlly
	GloryEnemy

	Bow
	ShotGround
	ShotFlesh
	ShotWood
	ShotMetal

	SwingLight
	SwingHeavy
	HitFlesh
	HitWood
	HitMetal
	Stab

	Drop
	Death

	Stomp
	Heal
)
