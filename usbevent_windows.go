package usbevent

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/christowolf/usb-event/internal/user32"
	"github.com/christowolf/usb-event/internal/winuser"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

// Functions based on https://gist.github.com/nathan-osman/18c2e227ad00a223b61c0b3c16d452c3.
type Notifier struct {
	Channel chan winuser.EventInfo
	Hwnd    windows.HWND
}

func Register() (*Notifier, error) {
	n := Notifier{Channel: make(chan winuser.EventInfo)}
	wn := winuser.Notifier{Channel: n.Channel}
	cb := windows.NewCallback(wn.WndProc)
	// The following is based on
	// https://github.com/hallazzang/go-windows-programming/blob/ff0b400d8c7ba888340412472d92765a8412dc0d/example/gui/basic/main.go#L51.
	inst := win.GetModuleHandle(nil)
	cn, err := syscall.UTF16PtrFromString("usbeventWindow")
	if err != nil {
		return nil, fmt.Errorf("failed to convert window class name to UTF16: %w", err)
	}
	wc := win.WNDCLASSEX{
		HInstance:     inst,
		LpfnWndProc:   cb,
		LpszClassName: cn,
	}
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	if win.RegisterClassEx(&wc) == 0 {
		return nil, fmt.Errorf("failed to register window class: %w", syscall.GetLastError())
	}
	wName, err := syscall.UTF16PtrFromString("usbevent.exe")
	if err != nil {
		return nil, fmt.Errorf("failed to convert window name to UTF16: %w", err)
	}
	wdw := win.CreateWindowEx(
		0,
		wc.LpszClassName,
		wName,
		win.WS_MINIMIZE|win.WS_OVERLAPPEDWINDOW,
		win.CW_USEDEFAULT,
		win.CW_USEDEFAULT,
		100,
		100,
		0,
		0,
		wc.HInstance,
		nil)
	if wdw == 0 {
		return nil, fmt.Errorf("failed to create window: %w", syscall.GetLastError())
	}
	err = user32.RegisterDeviceNotificationW(windows.Handle(wdw), user32.DEVICE_NOTIFY_WINDOW_HANDLE)
	if err != nil {
		return nil, fmt.Errorf("failed device notification registration: %w", err)
	}
	_ = win.ShowWindow(wdw, win.SW_HIDE)
	win.UpdateWindow(wdw)
	n.Hwnd = windows.HWND(wdw)
	return &n, nil
}

func (n *Notifier) Run() {
	for {
		var msg win.MSG
		got := win.GetMessage(&msg, win.HWND(n.Hwnd), 0, 0)
		if got == 0 {
			win.TranslateMessage(&msg)
			win.DispatchMessage(&msg)
		}
	}
}
