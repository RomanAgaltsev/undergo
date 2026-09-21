# 10 — The zero that was not zero

IEEE-754 has two zeros. `math.Copysign(0, -1)` is the negative one: its bits
are `0x8000000000000000`, and `nz == 0` is nevertheless `true`.

So "is this zero?" has two defensible answers — the bits say no, the comparison
says yes. Go 1.22 changed which one `reflect.Value.IsZero` gives. The question
here is what Go 1.21 gave, which is not one answer.

## The program

```go
type Small struct{ F float64 }

type Big struct {
	F   float64
	Pad [1024]byte
}

func main() {
	nz := math.Copysign(0, -1)
	small := Small{F: nz}
	big := Big{F: nz}
	fmt.Printf("%v|%v|%v|%v|%v",
		reflect.ValueOf(nz).IsZero(),
		reflect.ValueOf(small).IsZero(),
		reflect.ValueOf(big).IsZero(),
		reflect.ValueOf([2]float64{nz, nz}).IsZero(),
		reflect.ValueOf(small).Field(0).IsZero())
}
```

`Small` and `Big` hold the same negative zero. They differ only in padding the
program never reads.

Five readings, run under two toolchains. The go line is `go 1.21` in every run
and does not vary — only the compiler does.

| slot | value under test |
|---|---|
| `scalar_is_zero_…` | the `float64` itself |
| `small_struct_is_zero_…` | `Small{F: nz}` |
| `big_struct_is_zero_…` | `Big{F: nz}` |
| `array_is_zero_…` | `[2]float64{nz, nz}` |
| `field_is_zero_…` | `Small{F: nz}` reached with `.Field(0)` |

Each has an `_on_local` and an `_under_go1_21` form, so there are ten slots.
Predict `true` or `false` for each.

```
undergo verify versions/10-negative-zero
```

**This task fetches a Go toolchain the first time it runs.** If it cannot be
fetched the task skips and says so.

## Questions to answer in writing

1. Before running anything: you have one negative zero. Write down what you
   expect all five readings to be on an old toolchain, and say which rule you
   used. Then run it. How many did you get?
2. `Small` and `Big` contain the same value. If their answers differ, what
   could `reflect` possibly be looking at other than the value? Find the
   mechanism, and find the number that decides it.
3. On one of these toolchains, a value is zero and the float inside it is not.
   Which of those two answers is wrong? Defend your choice — and then say what
   it would have cost to "fix" the other one instead.
4. This was a bug in a function whose whole job is one boolean. What kind of
   test would have caught it, and why do you think none did for so long?
