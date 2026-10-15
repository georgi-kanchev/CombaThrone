package zone

type Enum uint8

const (
	Field Enum = iota
	Ruins

	Forest
	Swamp
	Desert

	Docks
	Ocean
	Glacier

	Cave
	Mine
	Hell

	COUNT
	LayerOffset = 10
)
