package winuser

import (
	"bytes"
	"encoding/binary"
	"reflect"
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
// For now this only supports DBT_DEVICEARRIVAL.
func (n *Notifier) WndProc(hwnd windows.HWND, msg types.DWORD, wParam, lParam uintptr) uintptr {
	switch msg {
	case message.WM_DEVICECHANGE:
		switch wParam {
		case uintptr(message.DBT_DEVICEARRIVAL):
			dType, guid, name, err := readDeviceInfo(lParam)
			if err != nil {
				panic(err)
			}
			n.Channel <- EventInfo{dType, guid, name, Arrival}
		}
		// TODO: https://gist.github.com/nathan-osman/18c2e227ad00a223b61c0b3c16d452c3
	}
	return win.DefWindowProc(win.HWND(hwnd), uint32(msg), wParam, lParam)
}

// readDeviceInfo parses binary data into a readable form.
// Based on https://github.com/unreality/nCryptAgent/blob/eecebcab1e366420f6479090b5cfa803f3979f57/deviceevents/events.go#L63.
func readDeviceInfo(pDevInfo uintptr) (types.DWORD, windows.GUID, string, error) {
	var devInfo deviceInfo
	var devInfoBytes []byte
	// Do some pointer arithmetic to align the struct.
	s1 := (*reflect.SliceHeader)(unsafe.Pointer(&devInfoBytes))
	s1.Data = pDevInfo
	s1.Len = int(uint32(unsafe.Sizeof(devInfo)))
	s1.Cap = s1.Len
	// Read the binary data.
	r := bytes.NewReader(devInfoBytes)
	// Windows is little endian.
	err := binary.Read(r, binary.LittleEndian, &devInfo.size)
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
	var devName []byte
	s2 := (*reflect.SliceHeader)(unsafe.Pointer(&devName))
	s2.Data = pDevInfo + unsafe.Sizeof(devInfo)
	s2.Len = int(uint32(devInfo.size) - uint32(unsafe.Sizeof(devInfo)))
	s2.Cap = int(uint32(devInfo.size) - uint32(unsafe.Sizeof(devInfo)))
	name := string(devName)
	// Return all relevant data.
	return devInfo.deviceType, devInfo.classGuid, name, nil
}
