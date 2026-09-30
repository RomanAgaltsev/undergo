package drill

import "testing"

func TestQueuesAgree(t *testing.T) {
	for _, q := range []Queue{NewSliceQueue(), NewRingQueue(2)} {
		for i := range 5 {
			q.Push(i)
		}
		for want := range 5 {
			if got, ok := q.Pop(); !ok || got != want {
				t.Fatalf("%T: Pop = %d, %v; want %d, true", q, got, ok, want)
			}
		}
		if _, ok := q.Pop(); ok {
			t.Fatalf("%T: Pop on an empty queue reported a value", q)
		}
	}
}

func benchmarkQueue(b *testing.B, q Queue) {
	b.ReportAllocs()
	for b.Loop() {
		q.Push(1)
		q.Pop()
	}
}

func BenchmarkSliceQueue(b *testing.B) { benchmarkQueue(b, NewSliceQueue()) }

func BenchmarkRingQueue(b *testing.B) { benchmarkQueue(b, NewSliceQueue()) }
