package role

type Enum uint8

const (
	Fighter Enum = iota
	Ranger
	Defender
	Griefer
	Supplier
	Collector
	COUNT
)
