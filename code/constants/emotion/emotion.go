package emotion

type Enum uint8

const (
	NONE Enum = iota
	Calm
	Talk
	Happy
	Sad
	Angry
	COUNT
)
