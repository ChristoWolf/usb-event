package user32_test

import (
	"fmt"
	"reflect"
	"syscall"
	"testing"

	"github.com/christowolf/usb-event/internal/user32"
)

func TestRegisterDeviceNotificationW(t *testing.T) {
	t.Parallel()
	handle, err := syscall.GetCurrentProcess()
	if err != nil {
		fmt.Println(err)
	}
	if got := user32.RegisterDeviceNotificationW(
		user32.HANDLE(handle),
		,
		tt.args.Flags); !reflect.DeepEqual(got, tt.want) {
		t.Errorf("RegisterDeviceNotificationW() = %v, want %v", got, tt.want)
	}
}
