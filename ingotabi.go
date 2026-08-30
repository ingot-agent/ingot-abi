// Package ingotabi contains the fixed Component ABI shared by ingot
// Plugins and generated Runtime Images.
//
// This package is the only supported root contract between a Plugin
// Component and the Runtime Image produced by the ingot Builder. It is not
// a general-purpose SDK and it is never user-configurable: the Builder
// pins its module path and exact version, recognizes exactly these types
// when validating Component constructors and graph wiring, and generates
// the host implementations for every other capability.
//
// Plugins import this package for four things:
//
//   - Cleanup, the per-component resource release function returned by the
//     every Component constructor;
//   - Optional[T] and None/Some, the OPTIONAL dependency wrapper used in a
//     Component's Dependencies and Exports;
//   - Named[T] and CheckUniqueNames, the MANY runtime-instance name wrapper
//     used to expose named capability collections;
//   - the sentinel errors returned by CheckUniqueNames.
//
// Everything else a Plugin needs (model, tool, session, prompt, agent,
// interaction and the other replaceable agent capabilities) lives in
// ordinary domain Contract Modules such as github.com/ingot-agent/sdk.
package ingotabi

import (
	"context"
	"errors"
	"fmt"
)

// Cleanup releases the resources owned by one Component instance.
//
// A Cleanup must observe ctx.Done and return promptly when the context is
// cancelled. Generated wiring calls Component cleanups synchronously in the
// reverse of Component creation order, with each call bounded by the
// runtime's cleanup timeout. Cleanup is optional: a constructor may return
// a nil Cleanup when the Component owns no resources.
//
// Cleanup does not express process results; use lifecycle.Controller for
// shutdown requests.
type Cleanup func(context.Context) error

// Optional represents an OPTIONAL Component dependency. Valid distinguishes
// a provided zero value from an absent value.
type Optional[T any] struct {
	Value T
	Valid bool
}

// None returns an Optional with no value.
func None[T any]() Optional[T] {
	return Optional[T]{}
}

// Some returns an Optional containing value.
func Some[T any](value T) Optional[T] {
	return Optional[T]{Value: value, Valid: true}
}

// Named gives a runtime capability instance a stable name. The name is
// scoped to the collection in which the value is exported.
type Named[T any] struct {
	Name  string
	Value T
}

var (
	// ErrEmptyName indicates that a Named value has no runtime instance name.
	ErrEmptyName = errors.New("empty capability name")
	// ErrDuplicateName indicates that two Named values in one collection use
	// the same runtime instance name.
	ErrDuplicateName = errors.New("duplicate capability name")
)

// CheckUniqueNames verifies that every item has a non-empty, unique name.
// It returns an error wrapping ErrEmptyName for an empty name and
// ErrDuplicateName for a duplicate.
func CheckUniqueNames[T any](items []Named[T]) error {
	seen := make(map[string]int, len(items))
	for i, item := range items {
		if item.Name == "" {
			return fmt.Errorf("named capability at index %d: %w", i, ErrEmptyName)
		}
		if previous, ok := seen[item.Name]; ok {
			return fmt.Errorf(
				"named capability %q at index %d duplicates index %d: %w",
				item.Name,
				i,
				previous,
				ErrDuplicateName,
			)
		}
		seen[item.Name] = i
	}
	return nil
}
