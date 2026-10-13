package lane

type Enum uint8

const ( // (collision layers)
	Lower Enum = iota
	LowerOff
	Middle
	MiddleOff
	Upper
	UpperOff
	Garrison1
	Garrison2
	Garrison3
	GarrisonPlus1
	GarrisonPlus2
	GarrisonPlus3
	Count
	LayerOffset = 32
)
