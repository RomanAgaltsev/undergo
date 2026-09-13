package structof

import (
	"encoding/json"
	"strings"
	"testing"
)

type Row struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email,omitempty"`
	secret  string //nolint:unused // present so the test can ask for it by name
	Balance int64  `json:"balance"`
}

func sample() Row {
	return Row{ID: 7, Name: "ada", Email: "ada@example.com", secret: "x", Balance: 42}
}

func TestProjectKeepsNamesTypesAndTags(t *testing.T) {
	got, err := Project(sample(), "ID", "Balance")
	if err != nil {
		t.Fatalf("Project: %v", err)
	}

	encoded, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshalling the projection: %v", err)
	}
	if want := `{"id":7,"balance":42}`; string(encoded) != want {
		t.Errorf("got %s, want %s", encoded, want)
	}
}

func TestProjectKeepsTheRequestedOrder(t *testing.T) {
	got, err := Project(sample(), "Name", "ID")
	if err != nil {
		t.Fatalf("Project: %v", err)
	}

	encoded, _ := json.Marshal(got)
	if want := `{"name":"ada","id":7}`; string(encoded) != want {
		t.Errorf("got %s, want %s — fields must follow the order they were asked for", encoded, want)
	}
}

func TestProjectPreservesTagOptions(t *testing.T) {
	row := sample()
	row.Email = ""

	got, err := Project(row, "Email", "ID")
	if err != nil {
		t.Fatalf("Project: %v", err)
	}

	encoded, _ := json.Marshal(got)
	if want := `{"id":7}`; string(encoded) != want {
		t.Errorf("got %s, want %s — omitempty must survive the copy", encoded, want)
	}
}

func TestProjectAcceptsAPointer(t *testing.T) {
	row := sample()
	got, err := Project(&row, "ID")
	if err != nil {
		t.Fatalf("Project: %v", err)
	}
	if encoded, _ := json.Marshal(got); string(encoded) != `{"id":7}` {
		t.Errorf("got %s, want {\"id\":7}", encoded)
	}
}

func TestProjectNoFields(t *testing.T) {
	got, err := Project(sample())
	if err != nil {
		t.Fatalf("Project with no fields: %v", err)
	}
	if encoded, _ := json.Marshal(got); string(encoded) != `{}` {
		t.Errorf("got %s, want {}", encoded)
	}
}

func TestProjectUnknownFieldIsAnError(t *testing.T) {
	_, err := Project(sample(), "ID", "Nope")
	if err == nil {
		t.Fatal("no error for an unknown field")
	}
	if !strings.Contains(err.Error(), "Nope") {
		t.Errorf("error %q does not name the missing field", err)
	}
}

// TestProjectUnexportedFieldIsAnError is the sharp edge: reflect.StructOf panics
// on an unexported field name rather than returning an error.
func TestProjectUnexportedFieldIsAnError(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Project panicked instead of returning an error: %v", r)
		}
	}()

	_, err := Project(sample(), "secret")
	if err == nil {
		t.Fatal("no error for an unexported field")
	}
}

func TestProjectNonStructIsAnError(t *testing.T) {
	if _, err := Project(42, "ID"); err == nil {
		t.Error("no error for a non-struct argument")
	}
	if _, err := Project("hello"); err == nil {
		t.Error("no error for a string argument")
	}
}
