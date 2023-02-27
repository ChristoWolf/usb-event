package user32_test

import (
	"syscall"
	"testing"

	"github.com/christowolf/usb-event/internal/user32"
	"golang.org/x/sys/windows"
)

// TestRegisterDeviceNotificationW verifies that
// USB device events can be registered to by valid handles.
func TestRegisterDeviceNotificationW(t *testing.T) {
	t.Parallel()
	// We will use a valid window handle.
	hType := user32.DEVICE_NOTIFY_WINDOW_HANDLE
	handle, err := syscall.GetCurrentProcess()
	if err != nil {
		t.Fatalf("test needs handle from current process: %v", err)
	}
	err = user32.RegisterDeviceNotificationW(windows.Handle(handle), hType)
	if err != nil {
		t.Errorf("want: nil, got: %v", err)
	}
}

// TestRegisterDeviceNotificationW verifies that
// an error is returned for invalid handles.
func TestRegisterDeviceNotificationWInvalidHandle(t *testing.T) {
	t.Parallel()
	// We will use a window handle for a service,
	// which should not work.
	hType := user32.DEVICE_NOTIFY_HANDLE_TYPE(0x00000001) // service handle
	handle, err := syscall.GetCurrentProcess()
	if err != nil {
		t.Fatalf("test needs handle from current process: %v", err)
	}
	err = user32.RegisterDeviceNotificationW(windows.Handle(handle), hType)
	if err == nil {
		t.Error("want: non-nil error, got: nil")
	}
	t.Logf("success, got error: %v", err)
}
