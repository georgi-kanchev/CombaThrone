package entrance

type Enum uint8

const (
	NONE Enum = iota
	Door
	ShortGate
	TallGate
	COUNT
)
