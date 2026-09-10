# undergo

**Go katas beneath the surface.** Most Go exercises drill idiomatic surface Go.
These ask what a struct actually weighs, how many instantiations the compiler
emitted, whether that value escaped, and which interleavings the memory model
permits.

## Start in 60 seconds

```
git clone https://github.com/RomanAgaltsev/undergo
cd undergo
go run ./cmd/undergo list
go run ./cmd/undergo start layout/01-struct-padding
```

No install, no account, no network after the clone — dependencies are vendored.

## The five modes

| Mode | You produce | How it is graded |
|---|---|---|
| **build** | an implementation | frozen tests |
| **predict** | answers in `prediction.yaml` | the harness measures reality and diffs your guess |
| **optimize** | a faster or smaller implementation | benchmark against the baseline, same machine, same run |
| **review** | findings | you, against the sealed key |
| **design** | a design document | you, against the sealed rubric |

`predict` stores no answer anywhere. The truth is computed when you verify, so
the task cannot be cheated and cannot rot when Go changes. Grading tells you
*which* slots are wrong, never what the right value was.

## Commands

```
undergo list [--track T] [--mode M] [--unsolved]
undergo show <id>
undergo start <id>            copy the task into work/ (gitignored)
undergo verify <id>           grade your work
undergo hint <id>             a nudge — always available
undergo reveal <id>           the solution; refuses while failing, --stuck overrides
undergo progress              your record
undergo doctor                what this machine can grade
```

## About the sealed solutions

Every task ships its reference solution as `solution.sealed` — a gzipped,
base64'd blob. **This is obfuscation, not secrecy**: `base64 -d | tar xz` opens
it, and that is fine. The seal exists so you cannot read the answer by accident,
so `grep` and GitHub code search cannot spoil six tasks at once, and so opening
one is a deliberate act. `undergo hint` is always free.

Your work lives in `work/`, which is gitignored. This repo is the reference
copy; nobody's answers are committed to it.

## Contributing

New tasks are welcome — see [CONTRIBUTING.md](CONTRIBUTING.md) and the candidate
pool in [ROADMAP.md](ROADMAP.md).

MIT licensed.
