#include "textflag.h"

// func SumFree(p *[4]int64) int64
//
// One pointer argument (8 bytes) and one int64 result (8 bytes) makes the
// argument area 16 bytes: $0-16. Count it rather than trusting vet — vet
// checks a named offset exactly and does not check the frame size at all.
TEXT ·SumFree(SB), NOSPLIT, $0-16
	MOVQ p+0(FP), AX
	MOVQ 0(AX), BX
	ADDQ 8(AX), BX
	ADDQ 16(AX), BX
	ADDQ 24(AX), BX
	MOVQ BX, ret+8(FP)
	RET

// func SumKept(p *[4]int64) int64
//
// Byte-identical to SumFree. The only difference is in the Go declaration.
TEXT ·SumKept(SB), NOSPLIT, $0-16
	MOVQ p+0(FP), AX
	MOVQ 0(AX), BX
	ADDQ 8(AX), BX
	ADDQ 16(AX), BX
	ADDQ 24(AX), BX
	MOVQ BX, ret+8(FP)
	RET
