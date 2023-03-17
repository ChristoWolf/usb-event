package winuser_test

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
	"unsafe"

	"github.com/christowolf/usb-event/internal/message"
	"github.com/christowolf/usb-event/internal/types"
	"github.com/christowolf/usb-event/internal/user32"
	"github.com/christowolf/usb-event/internal/winuser"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

// TestNotifierWndProc verifies that the
// WndProc callback function reacts appropriately
// to data-less WM_DEVICECHANGE messages.
func TestNotifierWndProc(t *testing.T) {
	t.Parallel()
	type args struct {
		hwnd   windows.HWND
		msg    types.DWORD
		wParam uintptr
		lParam uintptr
	}
	tests := []struct {
		name     string
		msg      types.DWORD
		wantType int
		args     args
	}{
		{
			"no device change",
			0,
			0,
			args{0, 0, 0, 0},
		},
		{
			"device change, no arrival",
			win.WM_DEVICECHANGE,
			0,
			args{0, win.WM_DEVICECHANGE, 0, 0},
		},
		{
			"device change, arrival, no data",
			win.WM_DEVICECHANGE,
			int(message.DBT_DEVICEARRIVAL),
			args{0, win.WM_DEVICECHANGE, uintptr(message.DBT_DEVICEARRIVAL), 0},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			n := &winuser.Notifier{make(chan winuser.EventInfo)}
			if tt.wantType != 0 {
				go func() {
					got := <-n.Channel
					if got.EventType != winuser.EventType(tt.wantType) {
						t.Errorf("want %v, got: %v", tt.wantType, got.EventType)
					}
				}()
			}
			want := win.DefWindowProc(win.HWND(tt.args.hwnd), uint32(tt.args.msg), tt.args.wParam, tt.args.lParam)
			got := n.WndProc(tt.args.hwnd, tt.args.msg, tt.args.wParam, tt.args.lParam)
			if got != want {
				t.Errorf("want %v, got: %v", want, got)
			}
		})
	}
}

// TestNotifierWndProcDataDeviceArrival verifies that the
// WndProc callback function correctly parses
// data from DBT_DEVICEARRIVAL messages.
func TestNotifierWndProcDataDeviceArrival(t *testing.T) {
	t.Parallel()
	n := &winuser.Notifier{make(chan winuser.EventInfo)}
	type deviceInfo struct {
		size       types.DWORD
		deviceType types.DWORD
		reserved   types.DWORD
		classGuid  windows.GUID
	}
	guid, _ := windows.GUIDFromString("{a5dcbf10-6530-11d2-901f-00c04fb951ed}")
	devInfo := deviceInfo{
		size:       220,
		deviceType: user32.DBT_DEVTYP_DEVICEINTERFACE,
		reserved:   0,
		classGuid:  guid,
	}
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	binary.Write(w, binary.LittleEndian, devInfo)
	w.Flush()
	lParam := uintptr(unsafe.Pointer(&b.Bytes()[0]))
	want := winuser.EventInfo{
		DeviceType: devInfo.deviceType,
		Guid:       guid,
		DeviceName: "",
		EventType:  winuser.EventType(message.DBT_DEVICEARRIVAL),
	}
	go func() {
		got := <-n.Channel
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want %v, got: %v", want, got.EventType)
		}
	}()
	_ = n.WndProc(0, win.WM_DEVICECHANGE, uintptr(message.DBT_DEVICEARRIVAL), lParam)
}
