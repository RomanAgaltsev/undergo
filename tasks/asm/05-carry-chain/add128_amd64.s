#include "textflag.h"

// func Add128(loA, hiA, loB, hiB uint64) (lo, hi uint64)
//
// Four uint64 arguments (32 bytes) and two uint64 results (16 bytes) make the
// argument area 48 bytes, with the results at +32 and +40.
//
// This body adds the low words and ignores the carry out of that addition, so
// the high word is always zero. Replace it.
TEXT ·Add128(SB), NOSPLIT, $0-48
	MOVQ loA+0(FP), AX
	MOVQ loB+16(FP), BX
	ADDQ BX, AX
	MOVQ AX, lo+32(FP)
	MOVQ $0, hi+40(FP)
	RET
