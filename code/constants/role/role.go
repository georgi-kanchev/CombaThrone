package role

type Enum uint8

const (
	Fighter Enum = iota
	Backliner
	Defender
	Griefer
	Supplier
	Collector
	COUNT
)
