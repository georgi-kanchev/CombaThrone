package clouds

type Enum uint8

const (
	None Enum = iota
	Normal
	Windy
	Count
)
