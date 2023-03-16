package winuser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/christowolf/usb-event/internal/message"
	"github.com/christowolf/usb-event/internal/types"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

const (
	Arrival = EventType(message.DBT_DEVICEARRIVAL)
)

type deviceInfo struct {
	size       types.DWORD
	deviceType types.DWORD
	reserved   types.DWORD
	classGuid  windows.GUID
}

type EventType int

type EventInfo struct {
	DeviceType types.DWORD
	Guid       windows.GUID
	DeviceName string
	EventType  EventType
}

type Notifier struct {
	Channel chan EventInfo
}

// WndProc realizes the WNDPROC callback function,
// see https://learn.microsoft.com/en-us/windows/win32/api/winuser/nc-winuser-wndproc.
// For now, this only supports DBT_DEVICEARRIVAL.
func (n *Notifier) WndProc(hwnd windows.HWND, msg types.DWORD, wParam, lParam uintptr) uintptr {
	switch msg {
	case message.WM_DEVICECHANGE:
		switch wParam {
		case uintptr(message.DBT_DEVICEARRIVAL):
			dType, guid, name, err := readDeviceInfo(lParam)
			if err != nil {
				name = fmt.Sprintf("error: failed to read device information: %s", err)
			}
			n.Channel <- EventInfo{dType, guid, name, Arrival}
		}
		// TODO: https://gist.github.com/nathan-osman/18c2e227ad00a223b61c0b3c16d452c3
	}
	return win.DefWindowProc(win.HWND(hwnd), uint32(msg), wParam, lParam)
}

// readDeviceInfo parses binary data into a readable form.
// Based on https://github.com/unreality/nCryptAgent/blob/eecebcab1e366420f6479090b5cfa803f3979f57/deviceevents/events.go#L63.
func readDeviceInfo(pDevInfo uintptr) (dtype types.DWORD, guid windows.GUID, name string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	var devInfo deviceInfo
	var devInfoBytes []byte
	// Do some pointer arithmetic to align the struct.
	b := (*reflect.SliceHeader)(unsafe.Pointer(&devInfoBytes))
	b.Data = pDevInfo
	b.Len = int(uint32(unsafe.Sizeof(devInfo)))
	b.Cap = b.Len
	// Read the binary data.
	r := bytes.NewReader(devInfoBytes)
	// Windows is little endian.
	err = binary.Read(r, binary.LittleEndian, &devInfo.size)
	if err != nil {
		return types.DWORD(0), windows.GUID{}, "", err
	}
	err = binary.Read(r, binary.LittleEndian, &devInfo.deviceType)
	if err != nil {
		return types.DWORD(0), windows.GUID{}, "", err
	}
	err = binary.Read(r, binary.LittleEndian, &devInfo.reserved)
	if err != nil {
		return types.DWORD(0), windows.GUID{}, "", err
	}
	err = binary.Read(r, binary.LittleEndian, &devInfo.classGuid)
	if err != nil {
		return types.DWORD(0), windows.GUID{}, "", err
	}
	// Read the device name.
	// Based on (6) in https://pkg.go.dev/unsafe#Pointer.
	var devName string
	s := (*reflect.StringHeader)(unsafe.Pointer(&devName))
	s.Data = pDevInfo + unsafe.Sizeof(devInfo)
	read := uint32(devInfo.size)
	size := uint32(unsafe.Sizeof(devInfo))
	if read < size {
		err = fmt.Errorf("read %d bytes, expected at least %d", read, size)
		return types.DWORD(0), windows.GUID{}, "", err
	}
	s.Len = int(read - size)
	name = strings.ReplaceAll(devName, "\x00", "")
	name = strings.Replace(name, `\\?\`, "", 1)
	name = strings.Replace(name, "#", `\`, 2)
	name = strings.Split(name, "#")[0]
	// Return all relevant data.
	return devInfo.deviceType, devInfo.classGuid, name, nil
}
