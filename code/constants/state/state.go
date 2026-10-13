package state

type Enum uint8

const ( // states
	Summoned            Enum = iota // single frame
	WaitingToBeSummoned             // continuous

	Idling // continuous
	Moving // continuous

	HurtStart // single frame
	Hurting   // continuous

	DyingStart // single frame
	Dying      // continuous
	DyingEnd   // single frame
	Decaying   // continuous

	ActStart      // single frame
	ActCharging   // continuous
	ActTrigger    // single frame
	ActRecovering // continuous
	ActEnd        // single frame
)
