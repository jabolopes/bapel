package ir

type OpMode byte

const (
	ImmediateMode = OpMode(iota) // Relative to PC.
	VarMode                      // Relative to FP.
	StackMode                    // Relative to SP.
	maxOpMode
)
