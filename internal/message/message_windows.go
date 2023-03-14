package message

import (
	"github.com/christowolf/usb-event/internal/types"
)

const (
	// See https://learn.microsoft.com/en-us/windows/win32/devio/wm-devicechange
	// and https://github.com/lxn/win/blob/a377121e959e22055dd01ed4bb2383e5bd02c238/user32.go#L720.
	WM_DEVICECHANGE = types.DWORD(537)
)

const (
	// See https://learn.microsoft.com/en-us/windows/win32/devio/wm-devicechange.
	DBT_DEVICEARRIVAL = types.DWORD(0x8000)
)
