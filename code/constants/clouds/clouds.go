package clouds

type Enum uint8

const (
	NONE Enum = iota
	Normal
	Windy
	COUNT
)
