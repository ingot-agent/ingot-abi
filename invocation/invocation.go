// Package invocation defines the read-only Runtime Image invocation
// information injected into Components.
//
// The generated runtime owns the process arguments and the operation mode.
// Plugins never parse os.Args directly; they consume this capability to
// decide how to behave, and they cannot change what the runtime was asked
// to do.
//
// Mode is the only handshake between the Builder and the Runtime Image:
// ModeRun is the normal application mode, ModeCheck is the pre-switch
// construction validation mode. In ModeCheck a Component may validate its
// Config and Dependencies, but it must not start an interaction loop or
// hold external resources.
package invocation

// Mode identifies the runtime operation mode.
type Mode uint8

const (
	// ModeRun is the normal application mode: the Runtime Image executes
	// the complete Component Graph and enters the application loop.
	ModeRun Mode = iota + 1
	// ModeCheck is the Builder's construction validation mode. Components
	// validate configuration and dependencies only; they must not start
	// interaction loops or occupy external resources.
	ModeCheck
)

// Invocation exposes immutable invocation metadata of the Runtime Image.
// Implementations are safe for concurrent use.
type Invocation interface {
	// Arguments returns the caller-owned, immutable runtime arguments
	// (excluding ingot-owned flags such as --ingot-check).
	Arguments() []string
	// Mode returns the runtime operation mode.
	Mode() Mode
}
