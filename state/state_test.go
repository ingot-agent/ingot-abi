package state_test

import (
	"testing"

	"github.com/ingot-agent/ingot-abi/state"
)

type fixedScope string

func (s fixedScope) Dir() string { return string(s) }

var _ state.Scope = fixedScope("")

func TestScopeReturnsLocation(t *testing.T) {
	t.Parallel()
	value := fixedScope("/var/lib/ingot/state/example.com/plugin")
	if value.Dir() != "/var/lib/ingot/state/example.com/plugin" {
		t.Fatalf("Dir() = %q", value.Dir())
	}
}

func TestScopeIsConcurrentSafeToRead(t *testing.T) {
	t.Parallel()
	value := fixedScope("/var/lib/ingot/state/example.com/plugin")
	done := make(chan struct{})
	for index := 0; index < 16; index++ {
		go func() {
			defer func() { done <- struct{}{} }()
			if value.Dir() == "" {
				t.Error("Dir() returned an empty location")
			}
		}()
	}
	for index := 0; index < 16; index++ {
		<-done
	}
}
