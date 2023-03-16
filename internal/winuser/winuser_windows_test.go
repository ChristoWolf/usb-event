package winuser_test

import (
	"testing"

	"github.com/christowolf/usb-event/internal/message"
	"github.com/christowolf/usb-event/internal/types"
	"github.com/christowolf/usb-event/internal/winuser"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

// TestNotifierWndProc verifies that the
// WndProc callback function TODO
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
