package prenn_test

import (
	"context"
	"testing"

	"github.com/dundunlabs/prenn"
)

func TestParam(t *testing.T) {
	params := prenn.Params{
		"foo": "bar",
	}

	ctx := &prenn.Context{
		Context: context.Background(),
		Params:  params,
	}

	if got, want := ctx.Param("foo"), "bar"; got != want {
		t.Errorf("should return %q, got %q", want, got)
	}
}

func TestWithValue(t *testing.T) {
	ctx := &prenn.Context{
		Context: context.Background(),
	}
	ctx.WithValue("foo", "bar")

	if got, want := ctx.Value("foo"), "bar"; got != want {
		t.Errorf("should return %q, got %q", want, got)
	}
}
