// Package lifecycle defines the request interface between Plugins and the
// generated Runtime Image for process shutdown.
//
// lifecycle.Controller is the only process-control surface exposed to
// Plugins. It does not exit the process, it does not allow a Plugin to
// choose the exit code, and it does not expose process context.
//
// Semantics of RequestShutdown:
//
//   - concurrent-safe and non-blocking;
//   - the first request cancels the runtime context;
//   - a nil error expresses a normal completion intent, a non-nil error
//     expresses a process-level failure cause;
//   - non-nil causes received after the first request are aggregated with
//     Cleanup errors instead of being dropped by first-wins competition;
//   - the process exits 0 when the final aggregated result has no error and
//     1 when any error is present.
//
// The runtime owns lifecycle mechanics: context cancellation, cleanup
// ordering and the process result. Plugins only submit intent.
package lifecycle

// Controller accepts runtime shutdown requests. Implementations are safe
// for concurrent use and never block the caller.
type Controller interface {
	// RequestShutdown submits the current completion intent. nil means the
	// component finished normally; a non-nil error is a process-level
	// failure cause.
	RequestShutdown(error)
}
