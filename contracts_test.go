package ingotabi_test

import (
	"context"

	ingotabi "github.com/ingot-agent/ingot-abi"
	"github.com/ingot-agent/ingot-abi/invocation"
	"github.com/ingot-agent/ingot-abi/lifecycle"
	"github.com/ingot-agent/ingot-abi/state"
)

// These external-package compile assertions model how an independent Plugin
// consumes the fixed Runtime ABI. They protect the exact nominal shapes
// that Builder-generated code recognizes.

type invocationValue struct{}

func (invocationValue) Arguments() []string { return nil }
func (invocationValue) Mode() invocation.Mode {
	return invocation.ModeCheck
}

var _ invocation.Invocation = invocationValue{}

type controllerValue struct{}

func (controllerValue) RequestShutdown(error) {}

var _ lifecycle.Controller = controllerValue{}

type scopeValue struct{ dir string }

func (s scopeValue) Dir() string { return s.dir }

var _ state.Scope = scopeValue{}

var _ ingotabi.Cleanup = func(context.Context) error { return nil }

var _ = ingotabi.Some(ingotabi.Named[lifecycle.Controller]{
	Name:  "primary",
	Value: controllerValue{},
})

// The fixed Component constructor shape recognized by the Builder:
// func(context.Context, Config, Dependencies) (Exports, ingotabi.Cleanup, error)
type Config struct{}

type Dependencies struct {
	Invocation invocation.Invocation
	Lifecycle  lifecycle.Controller
	State      state.Scope
}

type Exports struct{}

func Component(
	_ context.Context,
	_ Config,
	_ Dependencies,
) (Exports, ingotabi.Cleanup, error) {
	return Exports{}, nil, nil
}
