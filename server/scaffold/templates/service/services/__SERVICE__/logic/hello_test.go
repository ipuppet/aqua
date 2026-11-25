package logic

import (
	"context"
	"log/slog"
	"testing"
)

func TestHelloService_SayHello(t *testing.T) {
	svc := NewHelloService(slog.Default())

	got, err := svc.SayHello(context.Background(), "")
	if err != nil {
		t.Fatalf("SayHello() error = %v", err)
	}
	if got == nil || got.Message != "hello, world" {
		t.Fatalf("SayHello() = %#v, want message %q", got, "hello, world")
	}

	got, err = svc.SayHello(context.Background(), "aqua")
	if err != nil {
		t.Fatalf("SayHello() error = %v", err)
	}
	if got == nil || got.Message != "hello, aqua" {
		t.Fatalf("SayHello() = %#v, want message %q", got, "hello, aqua")
	}
}
