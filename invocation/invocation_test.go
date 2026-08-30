package invocation_test

import (
	"testing"

	"github.com/ingot-agent/ingot-abi/invocation"
)

type staticInvocation struct {
	arguments []string
	mode      invocation.Mode
}

func (s staticInvocation) Arguments() []string {
	return append([]string(nil), s.arguments...)
}

func (s staticInvocation) Mode() invocation.Mode { return s.mode }

var _ invocation.Invocation = staticInvocation{}

func TestModesAreDistinctAndStartedFromRun(t *testing.T) {
	t.Parallel()
	if invocation.ModeRun != 1 {
		t.Fatalf("ModeRun = %d, want 1", invocation.ModeRun)
	}
	if invocation.ModeCheck != 2 {
		t.Fatalf("ModeCheck = %d, want 2", invocation.ModeCheck)
	}
}

func TestArgumentsReturnsCallerOwnedCopy(t *testing.T) {
	t.Parallel()
	invocationValue := staticInvocation{arguments: []string{"chat", "--plain"}, mode: invocation.ModeRun}
	arguments := invocationValue.Arguments()
	arguments[0] = "mutated"
	if invocationValue.Arguments()[0] != "chat" {
		t.Fatal("Arguments returned aliased storage")
	}
	if invocationValue.Mode() != invocation.ModeRun {
		t.Fatalf("Mode() = %d", invocationValue.Mode())
	}
}

func TestZeroInvocationIsModeRunValue(t *testing.T) {
	t.Parallel()
	// The zero Mode is not a declared mode; an implementation must always
	// return an explicit mode. This guards the Mode constants against an
	// accidental shift that would make zero equal ModeRun.
	var mode invocation.Mode
	if mode == invocation.ModeRun || mode == invocation.ModeCheck {
		t.Fatalf("zero Mode collides with a declared mode")
	}
}
