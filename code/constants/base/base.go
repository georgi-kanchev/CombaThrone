package base

type Enum uint8

const (
	NONE Enum = iota
	Camp
	Barrack
	Fort
	Fortress
	COUNT
)
