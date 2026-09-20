package suite

import (
	"fmt"
	"os"
	"testing"
)

func TestAlpha(t *testing.T) { fmt.Println("ALPHA RAN") }

// TestBeta stands for any test that reaches code calling os.Exit on a path it
// believes is a clean shutdown.
func TestBeta(t *testing.T) {
	fmt.Println("BETA RAN")
	os.Exit(0)
}

func TestGamma(t *testing.T) { fmt.Println("GAMMA RAN") }
