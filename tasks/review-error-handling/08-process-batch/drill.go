// Package drill — C3/08 process-batch.
package drill

import (
	"errors"
	"fmt"
	"log"
)

// ErrSkip marks an item that could not be processed.
var ErrSkip = errors.New("skip")

func handle(item int) error {
	if item < 0 {
		return fmt.Errorf("negative item %d: %w", item, ErrSkip)
	}
	return nil
}

// ProcessBatch processes each item and returns an error if the batch failed.
func ProcessBatch(items []int) (err error) {
	for _, item := range items {
		if err := handle(item); err != nil {
			log.Printf("item %d failed: %v", item, err)
			continue
		}
	}
	if len(items) == 0 {
		return fmt.Errorf("empty batch: %v", ErrSkip)
	}
	return err
}
