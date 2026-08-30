package ingotabi

import (
	"context"
	"errors"
	"testing"
)

func TestOptionalDistinguishesAbsentAndZeroValue(t *testing.T) {
	t.Parallel()

	none := None[int]()
	if none.Valid || none.Value != 0 {
		t.Fatalf("None[int]() = %#v", none)
	}

	some := Some(0)
	if !some.Valid || some.Value != 0 {
		t.Fatalf("Some(0) = %#v", some)
	}
}

func TestCheckUniqueNames(t *testing.T) {
	t.Parallel()

	if err := CheckUniqueNames([]Named[int]{
		{Name: "primary", Value: 1},
		{Name: "fallback", Value: 2},
	}); err != nil {
		t.Fatalf("unique names: %v", err)
	}
}

func TestCheckUniqueNamesRejectsEmptyName(t *testing.T) {
	t.Parallel()

	err := CheckUniqueNames([]Named[int]{{Value: 1}})
	if !errors.Is(err, ErrEmptyName) {
		t.Fatalf("empty name error = %v", err)
	}
}

func TestCheckUniqueNamesRejectsDuplicateName(t *testing.T) {
	t.Parallel()

	err := CheckUniqueNames([]Named[int]{
		{Name: "same", Value: 1},
		{Name: "same", Value: 2},
	})
	if !errors.Is(err, ErrDuplicateName) {
		t.Fatalf("duplicate name error = %v", err)
	}
}

func TestCleanupIsCallableFunctionType(t *testing.T) {
	t.Parallel()

	var cleanup Cleanup = func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	}
	if err := cleanup(context.Background()); err != nil {
		t.Fatalf("cleanup error = %v", err)
	}
}

func TestCleanupNilIsValid(t *testing.T) {
	t.Parallel()

	var cleanup Cleanup
	if cleanup != nil {
		t.Fatal("zero Cleanup must be nil")
	}
}
