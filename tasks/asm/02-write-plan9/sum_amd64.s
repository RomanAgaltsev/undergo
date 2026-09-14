#include "textflag.h"

// func SumInt64(xs []int64) int64
//
// Implement this. The frame size and argument size in the TEXT directive are
// part of the exercise: `go vet` checks them against the Go declaration and
// will tell you when they are wrong.
TEXT ·SumInt64(SB), NOSPLIT, $0-32
	// undergo: implement SumInt64
	MOVQ $0, ret+24(FP)
	RET
