package pickup

type Enum uint8

const (
	Coin Enum = iota
	Gem
	Crystal
	Relic
	Rune
	Snowflake
	Star
	Key
	COUNT
)
