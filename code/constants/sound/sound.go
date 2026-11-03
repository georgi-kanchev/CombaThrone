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
	ShotGround
	ShotFlesh
	ShotWood
	ShotMetal
	Swing
	HitFlesh
	HitWood
	HitMetal
	Stab
	Drop
	Death
	GloryAlly
	GloryEnemy
)
