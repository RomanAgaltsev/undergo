package envdecode

import (
	"strings"
	"testing"
)

func env(pairs map[string]string) func(string) (string, bool) {
	return func(k string) (string, bool) {
		v, ok := pairs[k]
		return v, ok
	}
}

type config struct {
	Host    string `env:"HOST"`
	Port    int    `env:"PORT"`
	Debug   bool   `env:"DEBUG"`
	Size    int64  `env:"SIZE"`
	Token   string `env:"TOKEN,required"`
	Skipped string
	hidden  string //nolint:unused // must be ignored by Decode
}

func TestDecodeFillsTaggedFields(t *testing.T) {
	var c config
	err := Decode(&c, env(map[string]string{
		"HOST": "localhost", "PORT": "8080", "DEBUG": "true",
		"SIZE": "1048576", "TOKEN": "abc",
	}))
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if c.Host != "localhost" || c.Port != 8080 || !c.Debug || c.Size != 1048576 || c.Token != "abc" {
		t.Fatalf("Decode gave %+v", c)
	}
	if c.Skipped != "" {
		t.Errorf("Skipped = %q, want untouched", c.Skipped)
	}
}

func TestDecodeLeavesAbsentOptionalFieldsAlone(t *testing.T) {
	c := config{Host: "kept"}
	if err := Decode(&c, env(map[string]string{"TOKEN": "abc"})); err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if c.Host != "kept" {
		t.Errorf("Host = %q, want %q", c.Host, "kept")
	}
}

func TestDecodeRequiresRequired(t *testing.T) {
	var c config
	err := Decode(&c, env(map[string]string{"HOST": "x"}))
	if err == nil {
		t.Fatal("a missing required variable was accepted")
	}
	if !strings.Contains(err.Error(), "TOKEN") {
		t.Errorf("the error must name the variable, got: %v", err)
	}
}

func TestDecodeRejectsUnparseableValues(t *testing.T) {
	var c config
	if err := Decode(&c, env(map[string]string{"PORT": "eighty", "TOKEN": "abc"})); err == nil {
		t.Fatal("a non-numeric port was accepted")
	}
}

func TestDecodeRejectsBadDestinations(t *testing.T) {
	var c config
	if err := Decode(c, env(nil)); err == nil {
		t.Error("a non-pointer destination was accepted")
	}
	if err := Decode(nil, env(nil)); err == nil {
		t.Error("a nil destination was accepted")
	}
	s := "not a struct"
	if err := Decode(&s, env(nil)); err == nil {
		t.Error("a pointer to a non-struct was accepted")
	}
}
