package sound

type Enum uint16

const (
	Ambience Enum = iota
	DoorOpen
	DoorClose
	DoorCrumble
	GateOpen
	GateClose
	GateCrumble
	Bow
	ArrowGround
	ArrowFlesh
	ArrowWood
	ArrowMetal
	MeleeFlesh
	MeleeWood
	MeleeMetal
	GloryAlly
	GloryEnemy
)
