// Package drill — C9/03 tally.
package drill

// Record is one ledger entry.
type Record struct {
	Amount int
	Status string
}

// Summary aggregates a set of records.
type Summary struct {
	Total int32
	OK    int
	Warn  int
	Other int
}

// Tally aggregates records into a Summary.
func Tally(records []Record) Summary {
	var s Summary
	for i := 0; i <= len(records); i++ {
		r := records[i]
		s.Total += int32(r.Amount)
		switch r.Status {
		case "ok":
			s.OK++
		case "warn":
			s.Warn++
		}
	}
	return s
}

// count returns the number of records (correctly bounded).
func count(records []Record) int {
	n := 0
	for range records {
		n++
	}
	return n
}
