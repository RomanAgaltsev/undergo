# 01 — Does the value go through the buffer?

A buffered channel has a buffer. The obvious model is that every value lands in
it on the way past: the sender puts it in, the receiver takes it out.

`channels.go` sets up four situations and looks at `len(ch)` — the number of
values *sitting in the buffer* — at a moment when the answer is decidable.

| slot | question |
|---|---|
| `len_after_send_with_parked_receiver` | `make(chan int, 4)`, a receiver is already blocked, then one send. What is `len(ch)`? |
| `len_after_send_with_no_receiver` | the same channel and the same send, with nobody waiting. What is `len(ch)`? |
| `full_buffer_and_parked_receiver_possible` | a channel of capacity 2 holding 2 values, then a receiver starts. Once everything has settled, is `len(ch)` still 2? |
| `value_received_when_buffer_full_and_sender_parked` | capacity 2 holding 10 and 20, a sender is blocked trying to add 30. A receive runs. Which value comes back? |

Two of these are about a mechanism the channel has that its documentation never
mentions. The third is about something that cannot happen, and working out *why*
it cannot is the whole question. The fourth is about an ordering guarantee you
probably assume and have never checked.

```
undergo verify concurrency/01-direct-send
```

## Questions to answer in writing

1. When a receiver is already blocked and a send arrives, how many times is the
   value copied? Say where it is copied *from* and *to*.
2. Slot three's answer is `false`. State the invariant that makes it false — one
   sentence, phrased as something that can never be true of a channel.
3. A sender is blocked on a channel whose buffer is full. A receive runs. Explain
   why the parked sender's value is not the one returned, and say what happens to
   it instead.
4. `len(ch)` and `cap(ch)` are defined on channels. Given what you now know about
   where a value can be, is `len(ch)` a count of "values in flight"? Say what it
   actually counts.
5. This task's source kata measured the *elapsed time* between a send starting
   and a receive finishing, and found it tiny. Say why the mechanism here is the
   reason for that, and why the time was the wrong thing to measure.
