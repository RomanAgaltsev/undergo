package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestCIVerifyNeverPrintsSolutionSource(t *testing.T) {
	e := sealedRepo(t)
	var out, errOut bytes.Buffer
	e.Out, e.Err = &out, &errOut

	// The fixture's sealed solution is not a compilable task, so this is
	// expected to report a failure — what matters is that it stays quiet
	// about the contents.
	_ = CIVerify(e, nil)

	combined := out.String() + errOut.String()
	if strings.Contains(combined, "package padding // solved") {
		t.Fatal("ci-verify leaked the reference solution into its output")
	}
	if !strings.Contains(combined, "layout/01-struct-padding") {
		t.Errorf("ci-verify did not report on the task:\n%s", combined)
	}
}
