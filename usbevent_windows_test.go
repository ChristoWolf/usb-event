package usbevent_test

import (
	"context"
	"log"
	"testing"
	"time"

	usbevent "github.com/christowolf/usb-event"
	"github.com/christowolf/usb-event/internal/winuser"
)

// Example shows how to use the package.
func Example() {
	n, err := usbevent.Register()
	if err != nil {
		log.Fatalf("registration failed: %v\n", err)
	}
	go func() {
		for e := range n.Channel {
			log.Printf("%+v\n\n", e)
		}
	}()
	n.Run(context.Background())
}

// TestRegister verifies that the Register function
// returns a valid Notifier.
// This test is not parallel because it
// registers a window class.
func TestRegister(t *testing.T) {
	n, err := usbevent.Register()
	if err != nil {
		t.Errorf("registration failed: %v", err)
	}
	if n == nil {
		t.Errorf("got nil notifier")
	}
}

// TestRegisterFailure verifies that the Register function
// returns an error when the window class registration fails.
// This test is not parallel because it
// registers a window class.
func TestRegisterFailure(t *testing.T) {
	// We don't care what happens at first.
	usbevent.Register()
	// But latest by now, registration should fail.
	_, err := usbevent.Register()
	if err == nil {
		t.Errorf("registration succeeded, but should not")
	}
}

// TestRun verifies that the Run function
// can be cancelled/timed out using a context.
func TestRun(t *testing.T) {
	t.Parallel()
	timeout := 100 * time.Millisecond
	n := &usbevent.Notifier{Channel: make(chan winuser.EventInfo)}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	start := time.Now()
	n.Run(ctx)
	if time.Since(start) < timeout {
		t.Errorf("timed out too early")
	}
	cancel()
}
