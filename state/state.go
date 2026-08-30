// Package state defines the plugin-scoped persistent state location owned
// by the generated Runtime Image.
//
// The runtime assigns every Component the state scope of its Plugin
// identity and injects it into Dependencies. Scope carries only the
// persistent location: file access, schema migration and any business
// persistence remain the Plugin's responsibility.
package state

// Scope is the persistent state location assigned to one Plugin.
//
// Dir returns the absolute plugin-scoped state directory. Implementations
// are immutable and safe for concurrent use. Dir never returns an empty
// path.
type Scope interface {
	// Dir returns the absolute plugin-scoped state directory.
	Dir() string
}
